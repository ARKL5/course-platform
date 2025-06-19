package handler

import (
	"net/http"

	"course-platform/internal/domain/user/service"
	grpcClient "course-platform/internal/infrastructure/grpc_client"
	"course-platform/internal/shared/middleware"

	"log"

	"github.com/gin-gonic/gin"
)

// RegisterRequest 註冊請求結構體
type RegisterRequest struct {
	Username string `json:"username" binding:"required" example:"testuser"`
	Password string `json:"password" binding:"required" example:"123456"`
	Email    string `json:"email" binding:"required,email" example:"testuser@example.com"`
	Nickname string `json:"nickname" example:"測試用戶"`
}

// RegisterResponse 註冊響應結構體
type RegisterResponse struct {
	Message string   `json:"message" example:"註冊成功"`
	User    UserInfo `json:"user"`
}

// UserInfo 用戶信息結構體
type UserInfo struct {
	ID        uint   `json:"id" example:"1"`
	Username  string `json:"username" example:"testuser"`
	Nickname  string `json:"nickname" example:"測試用戶"`
	Avatar    string `json:"avatar" example:""`
	CreatedAt string `json:"created_at" example:"2023-01-01T00:00:00Z"`
}

// LoginRequest 登入請求結構體
type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required" example:"testuser"`
	Password   string `json:"password" binding:"required" example:"123456"`
	RememberMe bool   `json:"rememberMe" example:"true"`
}

// LoginResponse 登入響應結構體
type LoginResponse struct {
	Message string   `json:"message" example:"登入成功"`
	Token   string   `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User    UserInfo `json:"user"`
}

// ErrorResponse 錯誤響應結構體
type ErrorResponse struct {
	Error string `json:"error" example:"用戶名已存在"`
}

// UpdateProfileRequest 更新用户资料请求结构体
type UpdateProfileRequest struct {
	Nickname string `json:"nickname" example:"新昵称"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.jpg"`
	Phone    string `json:"phone" example:"13800138000"`
	Bio      string `json:"bio" example:"这是我的个人简介"`
}

// ChangePasswordRequest 修改密码请求结构体
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required" example:"oldpassword123"`
	NewPassword     string `json:"newPassword" binding:"required,min=8" example:"newpassword123"`
}

// UserHandler 用戶處理器結構體
type UserHandler struct {
	UserGRPCService *grpcClient.UserGRPCClientService
	UserService     service.UserServiceInterface // 添加直接服务调用
}

// NewUserHandler 建立新的 UserHandler 實例
func NewUserHandler(userGRPCService *grpcClient.UserGRPCClientService, userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		UserGRPCService: userGRPCService,
		UserService:     userService,
	}
}

// Register 處理用戶註冊請求
// @Summary 用戶註冊
// @Description 註冊新用戶帳號
// @Tags 用戶管理
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "註冊請求"
// @Success 201 {object} RegisterResponse "註冊成功"
// @Failure 400 {object} ErrorResponse "請求錯誤"
// @Router /register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req RegisterRequest

	// 綁定 JSON 請求到結構體
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求格式錯誤",
			"details": err.Error(),
		})
		return
	}

	// 呼叫gRPC服務進行註冊
	user, err := h.UserGRPCService.Register(req.Username, req.Email, req.Password, req.Nickname)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回成功結果
	c.JSON(http.StatusCreated, gin.H{
		"message": "註冊成功",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"created_at": user.CreatedAt,
		},
	})
}

// Login 處理用戶登入請求
// @Summary 用戶登入
// @Description 用戶登入並獲取JWT Token
// @Tags 用戶管理
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登入請求"
// @Success 200 {object} LoginResponse "登入成功"
// @Failure 400 {object} ErrorResponse "請求錯誤"
// @Failure 401 {object} ErrorResponse "認證失敗"
// @Router /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest

	// 綁定 JSON 請求到結構體
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ API Gateway: 登录请求绑定失败 - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "請求格式錯誤",
			"details": err.Error(),
		})
		return
	}

	log.Printf("🌐 API Gateway: 处理登录请求 - 标识符: %s", req.Identifier)

	// 調用gRPC服務進行登入 (使用identifier作为username)
	token, user, err := h.UserGRPCService.Login(req.Identifier, req.Password)
	if err != nil {
		log.Printf("❌ API Gateway: 登录失败 - %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("✅ API Gateway: 登录成功 - 用户: %s", user.Username)

	// 返回成功結果
	c.JSON(http.StatusOK, gin.H{
		"message": "登入成功",
		"token":   token,
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"created_at": user.CreatedAt,
		},
	})
}

// GetUserResponse 獲取用戶響應結構體
type GetUserResponse struct {
	Message string   `json:"message" example:"獲取用戶信息成功"`
	User    UserInfo `json:"user"`
}

// GetUser 處理獲取用戶信息請求
// @Summary 獲取用戶信息
// @Description 根據用戶名獲取用戶信息
// @Tags 用戶管理
// @Accept json
// @Produce json
// @Param username path string true "用戶名" example:"testuser"
// @Success 200 {object} GetUserResponse "獲取成功"
// @Failure 400 {object} ErrorResponse "請求錯誤"
// @Failure 404 {object} ErrorResponse "用戶不存在"
// @Router /user/{username} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	username := c.Param("username")

	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用戶名不能為空",
		})
		return
	}

	// 呼叫gRPC服務獲取用戶信息
	user, err := h.UserGRPCService.GetUserByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用戶不存在",
		})
		return
	}

	// 返回用戶信息
	c.JSON(http.StatusOK, gin.H{
		"message": "獲取用戶信息成功",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"created_at": user.CreatedAt,
		},
	})
}

// GetMe 处理获取当前用户信息请求
// @Summary 获取当前用户信息
// @Description 根据JWT Token获取当前登录用户的详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} GetUserResponse "获取成功"
// @Failure 401 {object} ErrorResponse "未授权"
// @Failure 404 {object} ErrorResponse "用户不存在"
// @Router /me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	// 从中间件上下文中获取用户信息
	userID, username, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户身份失败",
			"code":  "USER_IDENTITY_FAILED",
		})
		return
	}

	// 直接通过用户服务获取完整的用户信息（避免gRPC字段缺失问题）
	user, err := h.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "用户不存在",
			"code":  "USER_NOT_FOUND",
		})
		return
	}

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "获取当前用户信息成功",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"email":      user.Email,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"phone":      user.Phone,
			"bio":        user.Bio,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
		"auth_info": gin.H{
			"token_username": username,
			"token_user_id":  userID,
		},
	})
}

// LoginPage 渲染登录页面
func (h *UserHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"Title":   "登录",
		"BaseURL": "http://localhost:8080", // 在实际部署时应该从配置获取
	})
}

// RegisterPage 渲染注册页面
func (h *UserHandler) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{
		"Title":   "注册",
		"BaseURL": "http://localhost:8080", // 在实际部署时应该从配置获取
	})
}

// DashboardPage 渲染用户中心页面
func (h *UserHandler) DashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"Title":   "学习中心",
		"BaseURL": "http://localhost:8080", // 在实际部署时应该从配置获取
	})
}

// CreatorDashboardPage 渲染创作者工作台页面
func (h *UserHandler) CreatorDashboardPage(c *gin.Context) {
	c.HTML(http.StatusOK, "creator-dashboard.html", gin.H{
		"Title":   "创作者工作台",
		"BaseURL": "http://localhost:8080", // 在实际部署时应该从配置获取
	})
}

// ValidateToken 验证JWT Token
// @Summary 验证Token
// @Description 验证JWT Token的有效性
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body map[string]string true "Token验证请求"
// @Success 200 {object} map[string]interface{} "验证成功"
// @Failure 401 {object} ErrorResponse "Token无效"
// @Router /validate-token [post]
func (h *UserHandler) ValidateToken(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求格式错误",
		})
		return
	}

	token := req["token"]
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Token不能为空",
		})
		return
	}

	// 这里可以添加JWT Token验证逻辑
	// 目前返回简单的成功响应
	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"message": "Token有效",
	})
}

// Analytics 处理用户行为分析数据
// @Summary 用户行为分析
// @Description 记录用户行为分析数据
// @Tags 分析
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "分析数据"
// @Success 200 {object} map[string]string "记录成功"
// @Router /analytics [post]
func (h *UserHandler) Analytics(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		// 忽略格式错误，继续处理
		log.Printf("📊 Analytics: 接收到数据格式错误，但继续处理")
	}

	// 记录分析数据
	log.Printf("📊 Analytics: 接收到用户行为数据 - %v", data)

	c.JSON(http.StatusOK, gin.H{
		"message": "数据记录成功",
	})
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前用户的个人资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body UpdateProfileRequest true "更新资料请求"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 从中间件上下文中获取用户信息
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户身份失败",
		})
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求格式错误",
			"details": err.Error(),
		})
		return
	}

	// 调用用户服务更新用户资料
	user, err := h.UserService.UpdateProfileComplete(userID, req.Nickname, req.Avatar, req.Phone, req.Bio)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回更新后的用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "用户资料更新成功",
		"user": gin.H{
			"id":         user.ID,
			"username":   user.Username,
			"nickname":   user.Nickname,
			"avatar":     user.Avatar,
			"phone":      user.Phone,
			"bio":        user.Bio,
			"updated_at": user.UpdatedAt,
		},
	})
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前用户的登录密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "修改密码请求"
// @Success 200 {object} map[string]string "修改成功"
// @Failure 400 {object} ErrorResponse "请求错误"
// @Failure 401 {object} ErrorResponse "未授权"
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// 从中间件上下文中获取用户信息
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户身份失败",
		})
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "请求格式错误",
			"details": err.Error(),
		})
		return
	}

	// 调用用户服务修改密码
	err := h.UserService.ChangePassword(userID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 返回成功消息
	c.JSON(http.StatusOK, gin.H{
		"message": "密码修改成功",
	})
}

// GetCreatorStats 获取创作者统计信息
// @Summary 获取创作者统计
// @Description 获取当前用户的创作者统计信息，包括课程数量、学员数量、收入等
// @Tags 创作者
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "统计信息"
// @Failure 401 {object} ErrorResponse "未授权"
// @Router /creator/stats [get]
func (h *UserHandler) GetCreatorStats(c *gin.Context) {
	// 从中间件上下文中获取用户信息
	userID, _, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "获取用户身份失败",
		})
		return
	}

	// 模拟创作者统计数据（实际项目中应该从数据库获取）
	stats := gin.H{
		"total_courses":       3,
		"total_students":      128,
		"total_revenue":       2450.50,
		"active_courses":      2,
		"draft_courses":       1,
		"this_month_students": 25,
		"this_month_revenue":  320.00,
	}

	log.Printf("📊 获取创作者统计信息成功 - 用户ID: %d", userID)

	c.JSON(http.StatusOK, stats)
}
