package service

import (
	"fmt"
	"time"

	"course-platform/internal/domain/user/model"
	"course-platform/internal/domain/user/repository"
	"course-platform/internal/shared/utils"

	"github.com/golang-jwt/jwt/v5"
)

// UserServiceInterface 用户服务接口
// 定义用户业务逻辑的标准方法
type UserServiceInterface interface {
	// 核心业务方法
	Register(username, email, password, nickname string) (*model.User, error)
	Login(identifier, password string) (string, *model.User, error) // 支持用户名或邮箱登录
	GetUserByID(userID uint) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	UpdateProfile(userID uint, nickname, avatarURL string) (*model.User, error)
	UpdateProfileComplete(userID uint, nickname, avatarURL, phone, bio string) (*model.User, error)
	ChangePassword(userID uint, oldPassword, newPassword string) error

	// JWT相关方法
	GenerateToken(userID uint) (string, error)
	ValidateToken(tokenString string) (uint, error)
}

// UserService 用户服务实现
type UserService struct {
	userRepo  repository.UserRepositoryInterface // 用户仓储接口
	jwtSecret string                             // JWT密钥
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepositoryInterface) UserServiceInterface {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: "course-platform-secret-key-2024", // 实际项目中应从配置文件读取
	}
}

// Register 用户注册
// 处理用户注册业务逻辑，包括参数验证、密码加密、用户创建
func (s *UserService) Register(username, email, password, nickname string) (*model.User, error) {
	// 1. 参数验证
	if err := s.validateRegisterParams(email, password); err != nil {
		return nil, err
	}

	// 2. 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("检查邮箱是否存在失败: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("邮箱 %s 已被注册", email)
	}

	// 3. 密码加密
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	// 4. 创建用户对象
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		Nickname:     nickname,
	}

	// 5. 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return user, nil
}

// Login 用户登录
// 处理用户登录业务逻辑，包括身份验证、JWT生成
// 支持用户名或邮箱登录
func (s *UserService) Login(identifier, password string) (string, *model.User, error) {
	// 1. 参数验证
	if identifier == "" || password == "" {
		return "", nil, fmt.Errorf("用户名/邮箱和密码不能为空")
	}

	// 2. 根据标识符查找用户（先尝试用户名，再尝试邮箱）
	var user *model.User
	var err error

	// 检查是否是邮箱格式
	if contains(identifier, "@") && contains(identifier, ".") {
		// 尝试用邮箱查找
		user, err = s.userRepo.GetByEmail(identifier)
		if err != nil {
			// 邮箱查找失败，再尝试用户名查找
			user, err = s.userRepo.GetByUsername(identifier)
		}
	} else {
		// 尝试用用户名查找
		user, err = s.userRepo.GetByUsername(identifier)
		if err != nil {
			// 用户名查找失败，再尝试邮箱查找（防止用户名包含@符号）
			user, err = s.userRepo.GetByEmail(identifier)
		}
	}

	if err != nil {
		return "", nil, fmt.Errorf("用户不存在或密码错误")
	}

	// 3. 验证密码
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, fmt.Errorf("用户不存在或密码错误")
	}

	// 4. 生成JWT令牌
	token, err := s.GenerateToken(user.ID)
	if err != nil {
		return "", nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	return token, user, nil
}

// GetUserByID 根据ID获取用户信息
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

// GetUserByEmail 根据邮箱获取用户信息
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	return s.userRepo.GetByEmail(email)
}

// GetUserByUsername 根据用户名获取用户信息
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}

// UpdateProfile 更新用户资料
// 允许用户更新昵称和头像
func (s *UserService) UpdateProfile(userID uint, nickname, avatarURL string) (*model.User, error) {
	// 1. 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 2. 更新字段（只更新非空字段）
	updated := false
	if nickname != "" && nickname != user.Nickname {
		user.Nickname = nickname
		updated = true
	}
	if avatarURL != "" && avatarURL != user.AvatarURL {
		user.AvatarURL = avatarURL
		updated = true
	}

	// 3. 如果没有更新，直接返回
	if !updated {
		return user, nil
	}

	// 4. 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户资料失败: %w", err)
	}

	return user, nil
}

// UpdateProfileComplete 完整更新用户资料
// 支持更新昵称、头像、手机号和个人简介
func (s *UserService) UpdateProfileComplete(userID uint, nickname, avatarURL, phone, bio string) (*model.User, error) {
	// 1. 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在: %w", err)
	}

	// 2. 更新字段（只更新非空字段或显式设置的字段）
	updated := false

	if nickname != "" && nickname != user.Nickname {
		user.Nickname = nickname
		updated = true
	}

	if avatarURL != "" && avatarURL != user.AvatarURL {
		user.AvatarURL = avatarURL
		user.Avatar = avatarURL // 保持兼容性
		updated = true
	}

	if phone != user.Phone { // 允许设置为空字符串
		user.Phone = phone
		updated = true
	}

	if bio != user.Bio { // 允许设置为空字符串
		user.Bio = bio
		updated = true
	}

	// 3. 如果没有更新，直接返回
	if !updated {
		return user, nil
	}

	// 4. 保存更新
	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("更新用户资料失败: %w", err)
	}

	return user, nil
}

// ChangePassword 修改密码
// 验证旧密码并设置新密码
func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) error {
	// 1. 参数验证
	if oldPassword == "" || newPassword == "" {
		return fmt.Errorf("旧密码和新密码不能为空")
	}

	if len(newPassword) < 8 {
		return fmt.Errorf("新密码长度不能少于8位")
	}

	// 2. 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}

	// 3. 验证旧密码
	if !utils.CheckPasswordHash(oldPassword, user.PasswordHash) {
		return fmt.Errorf("旧密码不正确")
	}

	// 4. 加密新密码
	hashedNewPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("新密码加密失败: %w", err)
	}

	// 5. 更新密码
	user.PasswordHash = hashedNewPassword
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// GenerateToken 生成JWT令牌
func (s *UserService) GenerateToken(userID uint) (string, error) {
	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天过期
		"iat":     time.Now().Unix(),
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("签名JWT令牌失败: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func (s *UserService) ValidateToken(tokenString string) (uint, error) {
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("解析JWT令牌失败: %w", err)
	}

	// 验证并提取claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["user_id"].(float64); ok {
			return uint(userID), nil
		}
		return 0, fmt.Errorf("JWT令牌中缺少用户ID")
	}

	return 0, fmt.Errorf("无效的JWT令牌")
}

// validateRegisterParams 验证注册参数
func (s *UserService) validateRegisterParams(email, password string) error {
	if email == "" {
		return fmt.Errorf("邮箱不能为空")
	}

	// 简单的邮箱格式验证
	if len(email) < 5 || !contains(email, "@") || !contains(email, ".") {
		return fmt.Errorf("邮箱格式不正确")
	}

	if password == "" {
		return fmt.Errorf("密码不能为空")
	}

	if len(password) < 8 {
		return fmt.Errorf("密码长度不能少于8位")
	}

	if len(password) > 30 {
		return fmt.Errorf("密码长度不能超过30位")
	}

	return nil
}

// contains 检查字符串是否包含子字符串
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
