package repository

import (
	"errors"
	"fmt"
	"log"
	"time"

	"course-platform/internal/domain/course/model"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

// CourseRepositoryInterface 课程仓储接口
type CourseRepositoryInterface interface {
	Create(course *model.Course) error
	GetByID(id uint) (*model.Course, error)
	GetList(page, pageSize uint, categoryID uint) ([]*model.Course, uint, error)
	Update(course *model.Course) error
	Delete(id uint) error
	ExistsByTitle(title string) (bool, error)
	GetByInstructorID(instructorID uint) ([]*model.Course, error)
}

// CourseRepository 课程仓储实现
type CourseRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

// NewCourseRepository 创建课程仓储实例
func NewCourseRepository(db *gorm.DB, redis *redis.Client) CourseRepositoryInterface {
	return &CourseRepository{
		db:    db,
		redis: redis,
	}
}

// Create 创建课程
func (r *CourseRepository) Create(course *model.Course) error {
	log.Printf("🔍 Repository: 创建课程 - 标题: %s", course.Title)

	// 验证必需字段
	if course.Title == "" {
		return errors.New("课程标题不能为空")
	}
	if course.InstructorID == 0 {
		return errors.New("讲师ID不能为空")
	}

	// 检查课程标题是否已存在
	exists, err := r.ExistsByTitle(course.Title)
	if err != nil {
		return fmt.Errorf("检查课程标题存在性失败: %w", err)
	}
	if exists {
		return errors.New("课程标题已存在")
	}

	// 设置默认状态
	if course.Status == "" {
		course.Status = "draft" // 草稿状态
	}

	// 创建课程
	if err := r.db.Create(course).Error; err != nil {
		log.Printf("❌ Repository: 创建课程失败 - %v", err)
		return fmt.Errorf("创建课程失败: %w", err)
	}

	// 缓存课程信息
	r.cacheCourse(course)

	log.Printf("✅ Repository: 课程创建成功 - ID: %d", course.ID)
	return nil
}

// GetByID 根据ID获取课程
func (r *CourseRepository) GetByID(id uint) (*model.Course, error) {
	log.Printf("🔍 Repository: 根据ID获取课程 - %d", id)

	// 先尝试从缓存获取
	if course := r.getCourseFromCache(fmt.Sprintf("id:%d", id)); course != nil {
		log.Printf("✅ Repository: 从缓存获取课程成功")
		return course, nil
	}

	var course model.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Repository: 课程不存在 - ID: %d", id)
			return nil, errors.New("课程不存在")
		}
		log.Printf("❌ Repository: 查询课程失败 - %v", err)
		return nil, fmt.Errorf("查询课程失败: %w", err)
	}

	// 缓存课程信息
	r.cacheCourse(&course)

	log.Printf("✅ Repository: 获取课程成功 - ID: %d", course.ID)
	return &course, nil
}

// GetList 获取课程列表（分页）
func (r *CourseRepository) GetList(page, pageSize uint, categoryID uint) ([]*model.Course, uint, error) {
	log.Printf("🔍 Repository: 获取课程列表 - 页码: %d, 页大小: %d, 分类ID: %d", page, pageSize, categoryID)

	// 设置默认值
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	// 构建查询
	query := r.db.Model(&model.Course{})

	// 根据分类ID过滤
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 只查询已发布的课程
	query = query.Where("status = ?", "published")

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("❌ Repository: 获取课程总数失败 - %v", err)
		return nil, 0, fmt.Errorf("获取课程总数失败: %w", err)
	}

	// 分页查询
	var courses []*model.Course
	offset := (page - 1) * pageSize
	err := query.Offset(int(offset)).Limit(int(pageSize)).Order("created_at DESC").Find(&courses).Error
	if err != nil {
		log.Printf("❌ Repository: 获取课程列表失败 - %v", err)
		return nil, 0, fmt.Errorf("获取课程列表失败: %w", err)
	}

	log.Printf("✅ Repository: 获取课程列表成功 - 数量: %d, 总数: %d", len(courses), total)
	return courses, uint(total), nil
}

// Update 更新课程
func (r *CourseRepository) Update(course *model.Course) error {
	log.Printf("🔍 Repository: 更新课程 - ID: %d", course.ID)

	if err := r.db.Save(course).Error; err != nil {
		log.Printf("❌ Repository: 更新课程失败 - %v", err)
		return fmt.Errorf("更新课程失败: %w", err)
	}

	// 清除相关缓存
	r.clearCourseCache(course)
	// 重新缓存更新后的数据
	r.cacheCourse(course)

	log.Printf("✅ Repository: 课程更新成功 - ID: %d", course.ID)
	return nil
}

// Delete 删除课程（软删除）
func (r *CourseRepository) Delete(id uint) error {
	log.Printf("🔍 Repository: 删除课程 - ID: %d", id)

	// 先获取课程信息，用于清除缓存
	course, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(&model.Course{}, id).Error; err != nil {
		log.Printf("❌ Repository: 删除课程失败 - %v", err)
		return fmt.Errorf("删除课程失败: %w", err)
	}

	// 清除缓存
	r.clearCourseCache(course)

	log.Printf("✅ Repository: 课程删除成功 - ID: %d", id)
	return nil
}

// ExistsByTitle 检查课程标题是否存在
func (r *CourseRepository) ExistsByTitle(title string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Course{}).Where("title = ?", title).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查课程标题存在性失败: %w", err)
	}
	return count > 0, nil
}

// GetByInstructorID 根据讲师ID获取课程列表
func (r *CourseRepository) GetByInstructorID(instructorID uint) ([]*model.Course, error) {
	log.Printf("🔍 Repository: 根据讲师ID获取课程 - 讲师ID: %d", instructorID)

	var courses []*model.Course
	err := r.db.Where("instructor_id = ?", instructorID).Order("created_at DESC").Find(&courses).Error
	if err != nil {
		log.Printf("❌ Repository: 根据讲师ID获取课程失败 - %v", err)
		return nil, fmt.Errorf("获取讲师课程失败: %w", err)
	}

	log.Printf("✅ Repository: 根据讲师ID获取课程成功 - 数量: %d", len(courses))
	return courses, nil
}

// 缓存相关方法

// cacheCourse 缓存课程信息
func (r *CourseRepository) cacheCourse(course *model.Course) {
	if r.redis == nil {
		return
	}

	ctx := context.Background()
	expiration := 30 * time.Minute // 缓存30分钟

	// 缓存课程ID映射
	courseKey := fmt.Sprintf("course:id:%d", course.ID)
	titleKey := fmt.Sprintf("course:title:%s", course.Title)

	r.redis.Set(ctx, courseKey, course.ID, expiration)
	r.redis.Set(ctx, titleKey, course.ID, expiration)
}

// getCourseFromCache 从缓存获取课程
func (r *CourseRepository) getCourseFromCache(key string) *model.Course {
	if r.redis == nil {
		return nil
	}

	ctx := context.Background()

	// 这里简化实现，实际应该缓存完整课程对象
	courseIDStr, err := r.redis.Get(ctx, fmt.Sprintf("course:%s", key)).Result()
	if err != nil {
		return nil
	}

	// 由于简化实现，这里返回nil，让程序去数据库查询
	_ = courseIDStr
	return nil
}

// clearCourseCache 清除课程缓存
func (r *CourseRepository) clearCourseCache(course *model.Course) {
	if r.redis == nil {
		return
	}

	ctx := context.Background()

	// 清除所有相关缓存
	r.redis.Del(ctx, fmt.Sprintf("course:id:%d", course.ID))
	r.redis.Del(ctx, fmt.Sprintf("course:title:%s", course.Title))
}
