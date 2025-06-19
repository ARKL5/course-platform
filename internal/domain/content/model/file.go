package model

import (
	"time"

	userModel "course-platform/internal/domain/user/model"

	"gorm.io/gorm"
)

// FileInfo 文件信息模型
// 用于存储上传文件的元数据信息
type FileInfo struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      `json:"created_at"`           // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 软删除时间

	// 文件信息字段
	FileID   string `gorm:"uniqueIndex;not null;size:64" json:"file_id"` // 文件唯一标识
	FileName string `gorm:"not null;size:255" json:"file_name"`          // 原始文件名
	FileURL  string `gorm:"not null;size:500" json:"file_url"`           // 文件访问URL
	FileType string `gorm:"not null;size:50" json:"file_type"`           // 文件类型 (image/avatar/document)
	FileSize int64  `gorm:"not null" json:"file_size"`                   // 文件大小(字节)
	MimeType string `gorm:"size:100" json:"mime_type"`                   // MIME类型

	// 关联字段
	UserID uint           `gorm:"not null;index" json:"user_id"` // 上传用户ID
	User   userModel.User `gorm:"foreignKey:UserID" json:"-"`    // 关联用户

	// 状态字段
	Status string `gorm:"size:20;default:'active'" json:"status"` // 文件状态 (active/deleted)
}

// TableName 指定表名
func (FileInfo) TableName() string {
	return "files"
}

// BeforeCreate 创建前钩子
func (f *FileInfo) BeforeCreate(tx *gorm.DB) error {
	// 设置默认状态
	if f.Status == "" {
		f.Status = "active"
	}
	return nil
}

// IsImage 判断是否为图片文件
func (f *FileInfo) IsImage() bool {
	return f.FileType == "image" || f.FileType == "avatar"
}

// File 课程文件模型 (用于内容服务)
type File struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FileName   string    `gorm:"size:255;not null" json:"file_name"` // 文件名
	FilePath   string    `gorm:"size:500;not null" json:"file_path"` // 文件路径
	FileURL    string    `gorm:"size:500" json:"file_url"`           // 文件访问URL
	FileSize   int64     `gorm:"not null" json:"file_size"`          // 文件大小
	FileType   string    `gorm:"size:50;not null" json:"file_type"`  // 文件类型 (image, video, document, etc.)
	CourseID   uint      `gorm:"not null;index" json:"course_id"`    // 关联课程ID
	UploaderID uint      `gorm:"not null;index" json:"uploader_id"`  // 上传者ID
	UploadTime time.Time `gorm:"autoCreateTime" json:"upload_time"`  // 上传时间
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (File) TableName() string {
	return "course_files"
}

// FileFilter 文件过滤器
type FileFilter struct {
	CourseID   uint   `json:"course_id"`
	FileType   string `json:"file_type"`
	UploaderID uint   `json:"uploader_id"`
	Page       int    `json:"page"`
	PageSize   int    `json:"page_size"`
}

// FileListCache 文件列表缓存结构
type FileListCache struct {
	Files []File `json:"files"`
	Total int64  `json:"total"`
}
