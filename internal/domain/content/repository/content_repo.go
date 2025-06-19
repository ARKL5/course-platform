package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"course-platform/internal/domain/content/model"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// ContentRepository 内容仓库接口
type ContentRepository interface {
	CreateFile(ctx context.Context, file *model.File) error
	GetFilesByFilter(ctx context.Context, filter *model.FileFilter) ([]model.File, int64, error)
	GetFileById(ctx context.Context, id uint) (*model.File, error)
	UpdateFile(ctx context.Context, file *model.File) error
	DeleteFile(ctx context.Context, id uint) error
	GetFilesByCourse(ctx context.Context, courseID uint, fileType string, page, pageSize int) ([]model.File, int64, error)
}

// contentRepository 内容仓库实现
type contentRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewContentRepository 创建内容仓库实例
func NewContentRepository(db *gorm.DB, rdb *redis.Client) ContentRepository {
	return &contentRepository{
		db:    db,
		redis: rdb,
	}
}

// CreateFile 创建文件记录
func (r *contentRepository) CreateFile(ctx context.Context, file *model.File) error {
	// 数据库操作
	if err := r.db.WithContext(ctx).Create(file).Error; err != nil {
		log.Printf("❌ 创建文件记录失败: %v", err)
		return fmt.Errorf("创建文件记录失败: %w", err)
	}

	// 清除相关缓存
	r.clearFileCache(ctx, file.CourseID)

	log.Printf("✅ 成功创建文件记录 ID: %d, 文件名: %s", file.ID, file.FileName)
	return nil
}

// GetFilesByFilter 按条件查询文件
func (r *contentRepository) GetFilesByFilter(ctx context.Context, filter *model.FileFilter) ([]model.File, int64, error) {
	// 构建缓存键
	cacheKey := r.buildFilterCacheKey(filter)

	// 尝试从缓存获取
	if cached, err := r.getFromCache(ctx, cacheKey); err == nil && cached != nil {
		return cached.Files, cached.Total, nil
	}

	// 数据库查询
	var files []model.File
	var total int64

	query := r.db.WithContext(ctx).Model(&model.File{})

	// 添加过滤条件
	if filter.CourseID != 0 {
		query = query.Where("course_id = ?", filter.CourseID)
	}
	if filter.FileType != "" {
		query = query.Where("file_type = ?", filter.FileType)
	}
	if filter.UploaderID != 0 {
		query = query.Where("uploader_id = ?", filter.UploaderID)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		log.Printf("❌ 统计文件数量失败: %v", err)
		return nil, 0, fmt.Errorf("统计文件数量失败: %w", err)
	}

	// 分页查询
	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).
		Order("upload_time DESC").Find(&files).Error; err != nil {
		log.Printf("❌ 查询文件列表失败: %v", err)
		return nil, 0, fmt.Errorf("查询文件列表失败: %w", err)
	}

	// 缓存结果
	result := &model.FileListCache{
		Files: files,
		Total: total,
	}
	r.setToCache(ctx, cacheKey, result, 10*time.Minute)

	log.Printf("✅ 成功查询文件列表，共 %d 条记录", len(files))
	return files, total, nil
}

// GetFileById 根据ID获取文件信息
func (r *contentRepository) GetFileById(ctx context.Context, id uint) (*model.File, error) {
	cacheKey := fmt.Sprintf("file:%d", id)

	// 尝试从缓存获取
	var file model.File
	cached, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cached), &file); err == nil {
			log.Printf("✅ 从缓存获取文件信息 ID: %d", id)
			return &file, nil
		}
	}

	// 数据库查询
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("文件不存在")
		}
		log.Printf("❌ 查询文件失败: %v", err)
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}

	// 缓存结果
	if data, err := json.Marshal(file); err == nil {
		r.redis.Set(ctx, cacheKey, data, 30*time.Minute)
	}

	log.Printf("✅ 成功查询文件信息 ID: %d", id)
	return &file, nil
}

// UpdateFile 更新文件信息
func (r *contentRepository) UpdateFile(ctx context.Context, file *model.File) error {
	if err := r.db.WithContext(ctx).Save(file).Error; err != nil {
		log.Printf("❌ 更新文件失败: %v", err)
		return fmt.Errorf("更新文件失败: %w", err)
	}

	// 清除缓存
	r.redis.Del(ctx, fmt.Sprintf("file:%d", file.ID))
	r.clearFileCache(ctx, file.CourseID)

	log.Printf("✅ 成功更新文件 ID: %d", file.ID)
	return nil
}

// DeleteFile 删除文件记录
func (r *contentRepository) DeleteFile(ctx context.Context, id uint) error {
	var file model.File
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("文件不存在")
		}
		return fmt.Errorf("查询文件失败: %w", err)
	}

	if err := r.db.WithContext(ctx).Delete(&file).Error; err != nil {
		log.Printf("❌ 删除文件失败: %v", err)
		return fmt.Errorf("删除文件失败: %w", err)
	}

	// 清除缓存
	r.redis.Del(ctx, fmt.Sprintf("file:%d", id))
	r.clearFileCache(ctx, file.CourseID)

	log.Printf("✅ 成功删除文件 ID: %d", id)
	return nil
}

// GetFilesByCourse 获取课程相关文件
func (r *contentRepository) GetFilesByCourse(ctx context.Context, courseID uint, fileType string, page, pageSize int) ([]model.File, int64, error) {
	filter := &model.FileFilter{
		CourseID: courseID,
		FileType: fileType,
		Page:     page,
		PageSize: pageSize,
	}

	return r.GetFilesByFilter(ctx, filter)
}

// buildFilterCacheKey 构建过滤器缓存键
func (r *contentRepository) buildFilterCacheKey(filter *model.FileFilter) string {
	parts := []string{"files"}

	if filter.CourseID != 0 {
		parts = append(parts, fmt.Sprintf("course:%d", filter.CourseID))
	}
	if filter.FileType != "" {
		parts = append(parts, fmt.Sprintf("type:%s", filter.FileType))
	}
	if filter.UploaderID != 0 {
		parts = append(parts, fmt.Sprintf("uploader:%d", filter.UploaderID))
	}

	parts = append(parts, fmt.Sprintf("page:%d", filter.Page))
	parts = append(parts, fmt.Sprintf("size:%d", filter.PageSize))

	return strings.Join(parts, ":")
}

// getFromCache 从缓存获取文件列表
func (r *contentRepository) getFromCache(ctx context.Context, key string) (*model.FileListCache, error) {
	cached, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var result model.FileListCache
	if err := json.Unmarshal([]byte(cached), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// setToCache 设置缓存
func (r *contentRepository) setToCache(ctx context.Context, key string, data *model.FileListCache, expiration time.Duration) {
	if jsonData, err := json.Marshal(data); err == nil {
		r.redis.Set(ctx, key, jsonData, expiration)
	}
}

// clearFileCache 清除文件相关缓存
func (r *contentRepository) clearFileCache(ctx context.Context, courseID uint) {
	// 删除课程相关的文件列表缓存
	pattern := fmt.Sprintf("files:course:%d:*", courseID)
	keys, err := r.redis.Keys(ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		r.redis.Del(ctx, keys...)
	}

	// 删除所有文件列表缓存
	pattern = "files:*"
	keys, err = r.redis.Keys(ctx, pattern).Result()
	if err == nil && len(keys) > 0 {
		r.redis.Del(ctx, keys...)
	}
}
