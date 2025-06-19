package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"course-platform/internal/domain/user/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

// UserRepositoryInterface 用户仓储接口
// 定义用户数据访问的标准方法
type UserRepositoryInterface interface {
	// 基础CRUD操作
	Create(user *model.User) error
	GetByID(id uint) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uint) error

	// 业务查询方法
	ExistsByEmail(email string) (bool, error)
	GetUserList(offset, limit int) ([]*model.User, int64, error)

	// 缓存相关方法
	SetUserCache(user *model.User) error
	GetUserFromCache(email string) (*model.User, error)
	DeleteUserCache(email string) error
}

// UserRepository 用户仓储实现
type UserRepository struct {
	db    *gorm.DB      // 数据库连接
	redis *redis.Client // Redis缓存客户端
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB, redis *redis.Client) UserRepositoryInterface {
	return &UserRepository{
		db:    db,
		redis: redis,
	}
}

// Create 创建新用户
func (r *UserRepository) Create(user *model.User) error {
	// 验证邮箱格式（基本验证）
	if user.Email == "" && user.Username == "" {
		return errors.New("邮箱或用户名不能为空")
	}

	// 确保兼容性：如果只有username，复制到email
	if user.Email == "" && user.Username != "" {
		user.Email = user.Username
	}
	if user.Username == "" && user.Email != "" {
		user.Username = user.Email
	}

	log.Printf("🔍 Repository: 创建用户 - 邮箱: %s, 用户名: %s", user.Email, user.Username)

	// 检查邮箱是否已存在
	exists, err := r.ExistsByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("检查邮箱存在性失败: %w", err)
	}
	if exists {
		return errors.New("邮箱已存在")
	}

	// 检查用户名是否已存在（兼容性检查）
	if user.Username != user.Email {
		exists, err = r.ExistsByUsername(user.Username)
		if err != nil {
			return fmt.Errorf("检查用户名存在性失败: %w", err)
		}
		if exists {
			return errors.New("用户名已存在")
		}
	}

	// 创建用户
	if err := r.db.Create(user).Error; err != nil {
		log.Printf("❌ Repository: 创建用户失败 - %v", err)
		return fmt.Errorf("创建用户失败: %w", err)
	}

	// 缓存用户信息
	r.cacheUser(user)

	log.Printf("✅ Repository: 用户创建成功 - ID: %d", user.ID)
	return nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	log.Printf("🔍 Repository: 根据ID获取用户 - %d", id)

	// 先尝试从缓存获取
	if user := r.getUserFromCache(fmt.Sprintf("id:%d", id)); user != nil {
		log.Printf("✅ Repository: 从缓存获取用户成功")
		return user, nil
	}

	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Repository: 用户不存在 - ID: %d", id)
			return nil, errors.New("用户不存在")
		}
		log.Printf("❌ Repository: 查询用户失败 - %v", err)
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 缓存用户信息
	r.cacheUser(&user)

	log.Printf("✅ Repository: 获取用户成功 - ID: %d", user.ID)
	return &user, nil
}

// GetByUsername 根据用户名获取用户（兼容性方法）
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	log.Printf("🔍 Repository: 根据用户名获取用户 - %s", username)

	// 先尝试从缓存获取
	if user := r.getUserFromCache(fmt.Sprintf("username:%s", username)); user != nil {
		log.Printf("✅ Repository: 从缓存获取用户成功")
		return user, nil
	}

	var user model.User

	// 按用户名查找
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Repository: 用户不存在 - %s", username)
			return nil, errors.New("用户不存在")
		}
		log.Printf("❌ Repository: 查询用户失败 - %v", err)
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 缓存用户信息
	r.cacheUser(&user)

	log.Printf("✅ Repository: 获取用户成功 - ID: %d", user.ID)
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	log.Printf("🔍 Repository: 根据邮箱获取用户 - %s", email)

	// 先尝试从缓存获取
	if user := r.getUserFromCache(fmt.Sprintf("email:%s", email)); user != nil {
		log.Printf("✅ Repository: 从缓存获取用户成功")
		return user, nil
	}

	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("❌ Repository: 用户不存在 - %s", email)
			return nil, errors.New("用户不存在")
		}
		log.Printf("❌ Repository: 查询用户失败 - %v", err)
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 缓存用户信息
	r.cacheUser(&user)

	log.Printf("✅ Repository: 获取用户成功 - ID: %d", user.ID)
	return &user, nil
}

// Update 更新用户信息
func (r *UserRepository) Update(user *model.User) error {
	log.Printf("🔍 Repository: 更新用户 - ID: %d", user.ID)

	if err := r.db.Save(user).Error; err != nil {
		log.Printf("❌ Repository: 更新用户失败 - %v", err)
		return fmt.Errorf("更新用户失败: %w", err)
	}

	// 清除相关缓存
	r.clearUserCache(user)
	// 重新缓存更新后的数据
	r.cacheUser(user)

	log.Printf("✅ Repository: 用户更新成功 - ID: %d", user.ID)
	return nil
}

// Delete 软删除用户
func (r *UserRepository) Delete(id uint) error {
	log.Printf("🔍 Repository: 删除用户 - ID: %d", id)

	// 先获取用户信息，用于清除缓存
	user, err := r.GetByID(id)
	if err != nil {
		return err
	}

	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		log.Printf("❌ Repository: 删除用户失败 - %v", err)
		return fmt.Errorf("删除用户失败: %w", err)
	}

	// 清除缓存
	r.clearUserCache(user)

	log.Printf("✅ Repository: 用户删除成功 - ID: %d", id)
	return nil
}

// ExistsByEmail 检查邮箱是否存在
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查邮箱存在性失败: %w", err)
	}
	return count > 0, nil
}

// GetUserList 获取用户列表（分页）
func (r *UserRepository) GetUserList(offset, limit int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 获取总数
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户总数失败: %w", err)
	}

	// 获取用户列表
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("获取用户列表失败: %w", err)
	}

	return users, total, nil
}

// SetUserCache 设置用户缓存
func (r *UserRepository) SetUserCache(user *model.User) error {
	if r.redis == nil {
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf("user:email:%s", user.Email)

	// 缓存1小时
	return r.redis.Set(ctx, key, user, time.Hour).Err()
}

// GetUserFromCache 从缓存获取用户
func (r *UserRepository) GetUserFromCache(email string) (*model.User, error) {
	if r.redis == nil {
		return nil, fmt.Errorf("Redis未连接")
	}

	ctx := context.Background()
	key := fmt.Sprintf("user:email:%s", email)

	var user model.User
	err := r.redis.Get(ctx, key).Scan(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// DeleteUserCache 删除用户缓存
func (r *UserRepository) DeleteUserCache(email string) error {
	if r.redis == nil {
		return nil
	}

	ctx := context.Background()
	key := fmt.Sprintf("user:email:%s", email)

	return r.redis.Del(ctx, key).Err()
}

// ExistsByUsername 检查用户名是否存在
func (r *UserRepository) ExistsByUsername(username string) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("检查用户名存在性失败: %w", err)
	}
	return count > 0, nil
}

// 缓存相关方法

// cacheUser 缓存用户信息
func (r *UserRepository) cacheUser(user *model.User) {
	if r.redis == nil {
		return
	}

	ctx := context.Background()
	expiration := 30 * time.Minute // 缓存30分钟

	// 序列化用户数据（这里简化处理）
	userKey := fmt.Sprintf("user:id:%d", user.ID)
	emailKey := fmt.Sprintf("user:email:%s", user.Email)
	usernameKey := fmt.Sprintf("user:username:%s", user.Username)

	// 缓存用户ID映射
	r.redis.Set(ctx, userKey, user.ID, expiration)
	r.redis.Set(ctx, emailKey, user.ID, expiration)
	if user.Username != user.Email {
		r.redis.Set(ctx, usernameKey, user.ID, expiration)
	}
}

// getUserFromCache 从缓存获取用户
func (r *UserRepository) getUserFromCache(key string) *model.User {
	if r.redis == nil {
		return nil
	}

	ctx := context.Background()

	// 这里简化实现，实际应该缓存完整用户对象
	userIDStr, err := r.redis.Get(ctx, fmt.Sprintf("user:%s", key)).Result()
	if err != nil {
		return nil
	}

	// 由于简化实现，这里返回nil，让程序去数据库查询
	_ = userIDStr
	return nil
}

// clearUserCache 清除用户缓存
func (r *UserRepository) clearUserCache(user *model.User) {
	if r.redis == nil {
		return
	}

	ctx := context.Background()

	// 清除所有相关缓存
	r.redis.Del(ctx, fmt.Sprintf("user:id:%d", user.ID))
	r.redis.Del(ctx, fmt.Sprintf("user:email:%s", user.Email))
	r.redis.Del(ctx, fmt.Sprintf("user:username:%s", user.Username))
}
