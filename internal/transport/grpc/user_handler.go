package grpc

import (
	"context"
	"log"

	"course-platform/internal/domain/user/service"
	"course-platform/internal/shared/pb/userpb"
)

// UserHandler 用户gRPC处理器
// 处理来自API Gateway的gRPC请求，调用用户服务完成业务逻辑
type UserHandler struct {
	userpb.UnimplementedUserServiceServer
	userService service.UserServiceInterface
}

// NewUserHandler 创建用户gRPC处理器实例
func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register 处理用户注册gRPC请求
func (h *UserHandler) Register(ctx context.Context, req *userpb.RegisterRequest) (*userpb.RegisterResponse, error) {
	log.Printf("🔍 gRPC: 收到注册请求 - 用户名: %s, 邮箱: %s", req.Username, req.Email)

	// 使用完整参数注册用户
	user, err := h.userService.Register(req.Username, req.Email, req.Password, req.Nickname)
	if err != nil {
		log.Printf("❌ gRPC: 注册失败 - %v", err)
		return &userpb.RegisterResponse{
			Code:    400,
			Message: err.Error(),
			User:    nil,
		}, nil
	}

	// 转换为protobuf用户对象
	pbUser := &userpb.User{
		Id:        uint32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 注册成功 - 用户ID: %d", user.ID)
	return &userpb.RegisterResponse{
		Code:    200,
		Message: "注册成功",
		User:    pbUser,
	}, nil
}

// Login 处理用户登录gRPC请求
func (h *UserHandler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	log.Printf("🔍 gRPC: 收到登录请求 - 用户名: %s", req.Username)

	// 尝试用用户名或邮箱登录
	token, user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		log.Printf("❌ gRPC: 登录失败 - %v", err)
		return &userpb.LoginResponse{
			Code:    401,
			Message: err.Error(),
			Token:   "",
			User:    nil,
		}, nil
	}

	// 转换为protobuf用户对象
	pbUser := &userpb.User{
		Id:        uint32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 登录成功 - 用户ID: %d, Token长度: %d", user.ID, len(token))
	return &userpb.LoginResponse{
		Code:    200,
		Message: "登录成功",
		Token:   token,
		User:    pbUser,
	}, nil
}

// GetUser 处理获取用户gRPC请求
func (h *UserHandler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	log.Printf("🔍 gRPC: 收到获取用户请求 - 用户名: %s", req.Username)

	// 通过用户名获取用户
	user, err := h.userService.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("❌ gRPC: 获取用户失败 - %v", err)
		return &userpb.GetUserResponse{
			Code:    404,
			Message: err.Error(),
			User:    nil,
		}, nil
	}

	// 转换为protobuf用户对象
	pbUser := &userpb.User{
		Id:        uint32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 获取用户成功 - 用户ID: %d", user.ID)
	return &userpb.GetUserResponse{
		Code:    200,
		Message: "获取成功",
		User:    pbUser,
	}, nil
}

// GetUserByID 处理根据用户ID获取用户gRPC请求
func (h *UserHandler) GetUserByID(ctx context.Context, req *userpb.GetUserByIDRequest) (*userpb.GetUserByIDResponse, error) {
	log.Printf("🔍 gRPC: 收到通过ID获取用户请求 - 用户ID: %d", req.UserId)

	user, err := h.userService.GetUserByID(uint(req.UserId))
	if err != nil {
		log.Printf("❌ gRPC: 通过ID获取用户失败 - %v", err)
		return &userpb.GetUserByIDResponse{
			Code:    404,
			Message: err.Error(),
			User:    nil,
		}, nil
	}

	// 转换为protobuf用户对象
	pbUser := &userpb.User{
		Id:        uint32(user.ID),
		Username:  user.Username,
		Email:     user.Email,
		Nickname:  user.Nickname,
		Avatar:    user.AvatarURL,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 通过ID获取用户成功 - 用户ID: %d", user.ID)
	return &userpb.GetUserByIDResponse{
		Code:    200,
		Message: "获取成功",
		User:    pbUser,
	}, nil
}

// TODO: 以下方法需要在user.proto中添加相应的消息定义后才能实现

// GetMe 处理获取当前用户信息gRPC请求 (暂未实现)
// func (h *UserHandler) GetMe(ctx context.Context, req *userpb.GetMeRequest) (*userpb.GetMeResponse, error) {
//   // 需要在user.proto中添加GetMeRequest和GetMeResponse定义
// }

// UpdateProfile 处理更新用户资料gRPC请求 (暂未实现)
// func (h *UserHandler) UpdateProfile(ctx context.Context, req *userpb.UpdateProfileRequest) (*userpb.UpdateProfileResponse, error) {
//   // 需要在user.proto中添加UpdateProfileRequest和UpdateProfileResponse定义
// }

// ChangePassword 处理修改密码gRPC请求 (暂未实现)
// func (h *UserHandler) ChangePassword(ctx context.Context, req *userpb.ChangePasswordRequest) (*userpb.ChangePasswordResponse, error) {
//   // 需要在user.proto中添加ChangePasswordRequest和ChangePasswordResponse定义
// }

// GetUserByUsername 处理根据用户名获取用户gRPC请求 (暂未实现)
// func (h *UserHandler) GetUserByUsername(ctx context.Context, req *userpb.GetUserByUsernameRequest) (*userpb.GetUserByUsernameResponse, error) {
//   // 需要在user.proto中添加GetUserByUsernameRequest和GetUserByUsernameResponse定义
// }
