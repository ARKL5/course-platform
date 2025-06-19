package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
// 遵循swagger.yaml中的用户字段定义，同时保持向后兼容
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"` // 主键ID
	CreatedAt time.Time      `json:"created_at"`           // 创建时间
	UpdatedAt time.Time      `json:"updated_at"`           // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`       // 软删除时间

	// 兼容性字段 - 保持旧接口正常工作
	Username     string `gorm:"uniqueIndex;not null;size:100" json:"username"` // 用户名（暂时兼容）
	PasswordHash string `gorm:"not null;size:255" json:"-"`                    // 密码哈希值，不返回给前端

	// 新字段 - 基于swagger.yaml设计
	Email     string `gorm:"uniqueIndex;not null;size:100" json:"email"` // 邮箱地址，唯一
	Nickname  string `gorm:"size:100" json:"nickname"`                   // 用户昵称
	AvatarURL string `gorm:"size:500" json:"avatar_url"`                 // 头像URL
	Avatar    string `gorm:"size:500" json:"avatar"`                     // 头像URL（兼容字段）
	Phone     string `gorm:"size:20" json:"phone"`                       // 手机号
	Bio       string `gorm:"size:500" json:"bio"`                        // 个人简介
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// BeforeCreate 创建前钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 如果昵称为空，设置默认昵称
	if u.Nickname == "" {
		u.Nickname = "新用户"
	}

	// 同步Email和Username字段以保持兼容性
	if u.Email != "" && u.Username == "" {
		u.Username = u.Email
	} else if u.Username != "" && u.Email == "" {
		u.Email = u.Username
	}

	// 同步头像字段
	if u.AvatarURL != "" && u.Avatar == "" {
		u.Avatar = u.AvatarURL
	} else if u.Avatar != "" && u.AvatarURL == "" {
		u.AvatarURL = u.Avatar
	}

	return nil
}

// BeforeUpdate 更新前钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	// 同步Email和Username字段以保持兼容性
	if u.Email != "" && u.Username == "" {
		u.Username = u.Email
	} else if u.Username != "" && u.Email == "" {
		u.Email = u.Username
	}

	// 同步头像字段
	if u.AvatarURL != "" && u.Avatar == "" {
		u.Avatar = u.AvatarURL
	} else if u.Avatar != "" && u.AvatarURL == "" {
		u.AvatarURL = u.Avatar
	}

	return nil
}
