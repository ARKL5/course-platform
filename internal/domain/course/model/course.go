package model

import (
	"time"

	"gorm.io/gorm"
)

// Course 课程模型
// 遵循swagger.yaml中的课程字段定义
type Course struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      `json:"created_at"`           // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 软删除时间

	// 核心字段 - 基于swagger.yaml和protobuf设计
	Title        string  `gorm:"not null;size:200" json:"title"`        // 课程标题
	Description  string  `gorm:"type:text" json:"description"`          // 课程描述
	InstructorID uint    `gorm:"not null;index" json:"instructor_id"`   // 讲师ID（用户ID）
	CategoryID   uint    `gorm:"index" json:"category_id"`              // 分类ID
	Price        float32 `gorm:"not null;default:0" json:"price"`       // 课程价格
	CoverImage   string  `gorm:"size:500" json:"cover_image"`           // 课程封面图片URL
	Status       string  `gorm:"size:20;default:'draft'" json:"status"` // 课程状态

	// 兼容性字段 - 保持向后兼容
	TeacherName   string `gorm:"size:100" json:"teacher_name"`    // 讲师姓名（兼容旧版本）
	Category      string `gorm:"size:50;index" json:"category"`   // 课程分类名称（兼容）
	CoverImageURL string `gorm:"size:500" json:"cover_image_url"` // 封面图片URL（兼容）

	// 扩展字段
	StudentCount int     `gorm:"default:0" json:"student_count"` // 学生数量
	Rating       float32 `gorm:"default:0" json:"rating"`        // 课程评分
	ViewCount    int     `gorm:"default:0" json:"view_count"`    // 浏览次数
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// BeforeCreate 创建前钩子
func (c *Course) BeforeCreate(tx *gorm.DB) error {
	// 设置默认状态
	if c.Status == "" {
		c.Status = "draft" // 默认为草稿状态
	}

	// 同步封面图片字段以保持兼容性
	if c.CoverImage != "" && c.CoverImageURL == "" {
		c.CoverImageURL = c.CoverImage
	} else if c.CoverImageURL != "" && c.CoverImage == "" {
		c.CoverImage = c.CoverImageURL
	}

	return nil
}

// BeforeUpdate 更新前钩子
func (c *Course) BeforeUpdate(tx *gorm.DB) error {
	// 同步封面图片字段以保持兼容性
	if c.CoverImage != "" && c.CoverImageURL == "" {
		c.CoverImageURL = c.CoverImage
	} else if c.CoverImageURL != "" && c.CoverImage == "" {
		c.CoverImage = c.CoverImageURL
	}

	return nil
}

// IsPublished 检查课程是否已发布
func (c *Course) IsPublished() bool {
	return c.Status == "published"
}

// IsDraft 检查课程是否为草稿状态
func (c *Course) IsDraft() bool {
	return c.Status == "draft"
}
