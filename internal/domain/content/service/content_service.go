package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"course-platform/internal/domain/content/model"
	"course-platform/internal/domain/content/repository"
)

// ContentService 内容服务接口
type ContentService interface {
	UploadFile(ctx context.Context, req *UploadFileRequest) (*model.File, error)
	GetFiles(ctx context.Context, filter *model.FileFilter) ([]model.File, int64, error)
	GetFileById(ctx context.Context, id uint) (*model.File, error)
	DeleteFile(ctx context.Context, fileID, userID uint) error
	GetFilesByCourse(ctx context.Context, courseID uint, fileType string, page, pageSize int) ([]model.File, int64, error)
}

// UploadFileRequest 文件上传请求
type UploadFileRequest struct {
	FileHeader *multipart.FileHeader // HTTP文件上传
	FileData   []byte                // gRPC字节数据
	FileName   string                // 文件名
	FileType   string                // 文件类型
	CourseID   uint                  // 课程ID
	UploaderID uint                  // 上传者ID
}

// contentService 内容服务实现
type contentService struct {
	repo      repository.ContentRepository
	uploadDir string // 文件上传目录
	baseURL   string // 文件访问基础URL
}

// NewContentService 创建内容服务实例
func NewContentService(repo repository.ContentRepository, uploadDir, baseURL string) ContentService {
	return &contentService{
		repo:      repo,
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// UploadFile 上传文件
func (s *contentService) UploadFile(ctx context.Context, req *UploadFileRequest) (*model.File, error) {
	// 验证请求
	if err := s.validateUploadRequest(req); err != nil {
		return nil, err
	}

	// 确保上传目录存在
	if err := s.ensureUploadDir(); err != nil {
		return nil, fmt.Errorf("创建上传目录失败: %w", err)
	}

	// 处理文件数据
	var fileData []byte
	var err error
	var fileSize int64

	if req.FileHeader != nil {
		// HTTP multipart文件上传
		fileData, fileSize, err = s.readMultipartFile(req.FileHeader)
		if err != nil {
			return nil, fmt.Errorf("读取上传文件失败: %w", err)
		}
		if req.FileName == "" {
			req.FileName = req.FileHeader.Filename
		}
	} else if req.FileData != nil {
		// gRPC字节数据上传
		fileData = req.FileData
		fileSize = int64(len(req.FileData))
	} else {
		return nil, fmt.Errorf("未提供有效的文件数据")
	}

	// 生成文件路径和URL
	filePath, fileURL, err := s.generateFilePath(req.FileName, req.CourseID)
	if err != nil {
		return nil, fmt.Errorf("生成文件路径失败: %w", err)
	}

	// 保存文件到磁盘
	if err := s.saveFileToDisk(filePath, fileData); err != nil {
		return nil, fmt.Errorf("保存文件失败: %w", err)
	}

	// 创建文件记录
	file := &model.File{
		FileName:   req.FileName,
		FilePath:   filePath,
		FileURL:    fileURL,
		FileSize:   fileSize,
		FileType:   req.FileType,
		CourseID:   req.CourseID,
		UploaderID: req.UploaderID,
		UploadTime: time.Now(),
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		// 如果数据库插入失败，删除已保存的文件
		os.Remove(filePath)
		return nil, fmt.Errorf("保存文件记录失败: %w", err)
	}

	log.Printf("✅ 成功上传文件: %s, 大小: %d 字节", req.FileName, fileSize)
	return file, nil
}

// GetFiles 获取文件列表
func (s *contentService) GetFiles(ctx context.Context, filter *model.FileFilter) ([]model.File, int64, error) {
	// 设置默认分页参数
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100 // 限制每页最大数量
	}

	files, total, err := s.repo.GetFilesByFilter(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("查询文件列表失败: %w", err)
	}

	log.Printf("✅ 查询文件列表成功，共 %d 条记录", len(files))
	return files, total, nil
}

// GetFileById 根据ID获取文件信息
func (s *contentService) GetFileById(ctx context.Context, id uint) (*model.File, error) {
	file, err := s.repo.GetFileById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询文件失败: %w", err)
	}

	return file, nil
}

// DeleteFile 删除文件
func (s *contentService) DeleteFile(ctx context.Context, fileID, userID uint) error {
	// 获取文件信息
	file, err := s.repo.GetFileById(ctx, fileID)
	if err != nil {
		return fmt.Errorf("查询文件失败: %w", err)
	}

	// 验证权限（只有上传者可以删除文件）
	if file.UploaderID != userID {
		return fmt.Errorf("无权限删除此文件")
	}

	// 删除数据库记录
	if err := s.repo.DeleteFile(ctx, fileID); err != nil {
		return fmt.Errorf("删除文件记录失败: %w", err)
	}

	// 删除磁盘文件
	if err := os.Remove(file.FilePath); err != nil {
		log.Printf("⚠️ 删除磁盘文件失败: %v", err)
		// 不返回错误，因为数据库记录已删除
	}

	log.Printf("✅ 成功删除文件: %s", file.FileName)
	return nil
}

// GetFilesByCourse 获取课程相关文件
func (s *contentService) GetFilesByCourse(ctx context.Context, courseID uint, fileType string, page, pageSize int) ([]model.File, int64, error) {
	return s.repo.GetFilesByCourse(ctx, courseID, fileType, page, pageSize)
}

// validateUploadRequest 验证上传请求
func (s *contentService) validateUploadRequest(req *UploadFileRequest) error {
	if req == nil {
		return fmt.Errorf("上传请求不能为空")
	}

	if req.FileName == "" {
		return fmt.Errorf("文件名不能为空")
	}

	if req.FileType == "" {
		return fmt.Errorf("文件类型不能为空")
	}

	// 特殊情况：头像上传时course_id可以为0
	if req.CourseID == 0 && req.FileType != "image" {
		return fmt.Errorf("课程ID不能为空")
	}

	if req.UploaderID == 0 {
		return fmt.Errorf("上传者ID不能为空")
	}

	// 验证文件类型
	allowedTypes := []string{"image", "video", "document", "audio", "other"}
	isValidType := false
	for _, t := range allowedTypes {
		if req.FileType == t {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return fmt.Errorf("不支持的文件类型: %s", req.FileType)
	}

	return nil
}

// ensureUploadDir 确保上传目录存在
func (s *contentService) ensureUploadDir() error {
	if err := os.MkdirAll(s.uploadDir, 0755); err != nil {
		return err
	}
	return nil
}

// readMultipartFile 读取multipart文件
func (s *contentService) readMultipartFile(fileHeader *multipart.FileHeader) ([]byte, int64, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, 0, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, 0, err
	}

	return data, fileHeader.Size, nil
}

// generateFilePath 生成文件路径和URL
func (s *contentService) generateFilePath(fileName string, courseID uint) (string, string, error) {
	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	if ext == "" {
		ext = ".unknown"
	}

	// 生成唯一文件名（使用时间戳+MD5）
	timestamp := time.Now().Format("20060102150405")
	hash := md5.Sum([]byte(fmt.Sprintf("%s_%d_%d", fileName, courseID, time.Now().UnixNano())))
	uniqueName := fmt.Sprintf("%s_%x%s", timestamp, hash, ext)

	// 按课程ID创建子目录，特殊处理头像上传
	var subDir string
	if courseID == 0 {
		subDir = "avatars" // 头像专用目录
	} else {
		subDir = fmt.Sprintf("course_%d", courseID)
	}
	dirPath := filepath.Join(s.uploadDir, subDir)

	// 确保子目录存在
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", "", err
	}

	// 完整文件路径
	filePath := filepath.Join(dirPath, uniqueName)

	// 文件访问URL
	fileURL := fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.baseURL, "/"), subDir, uniqueName)

	return filePath, fileURL, nil
}

// saveFileToDisk 保存文件到磁盘
func (s *contentService) saveFileToDisk(filePath string, data []byte) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}
