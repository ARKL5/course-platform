package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"course-platform/internal/domain/user/model"
	"course-platform/internal/shared/pb/userpb"
)

// UserGRPCClientService gRPC客户端服务
type UserGRPCClientService struct {
	client userpb.UserServiceClient
	conn   *grpc.ClientConn
}

// NewUserGRPCClientService 创建新的gRPC客户端服务
func NewUserGRPCClientService() (*UserGRPCClientService, error) {
	// 从环境变量获取用户服务地址
	serviceAddr := os.Getenv("USER_SERVICE_ADDR")
	if serviceAddr == "" {
		serviceAddr = "localhost:50051" // 默认地址
	}

	log.Printf("🔗 连接用户服务: %s", serviceAddr)

	// 连接到gRPC服务器
	conn, err := grpc.NewClient(serviceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("连接gRPC服务器失败: %v", err)
	}

	client := userpb.NewUserServiceClient(conn)

	return &UserGRPCClientService{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭gRPC连接
func (s *UserGRPCClientService) Close() error {
	return s.conn.Close()
}

// Register 通过gRPC调用用户注册
func (s *UserGRPCClientService) Register(username, email, password, nickname string) (*model.User, error) {
	log.Printf("🌐 API Gateway: 通过gRPC调用注册 - 用户名: %s, 邮箱: %s", username, email)

	req := &userpb.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
		Nickname: nickname,
	}

	resp, err := s.client.Register(context.Background(), req)
	if err != nil {
		log.Printf("❌ API Gateway: gRPC调用失败 - %v", err)
		return nil, fmt.Errorf("gRPC调用失败: %v", err)
	}

	if resp.Code != 200 {
		log.Printf("❌ API Gateway: 注册失败 - %s", resp.Message)
		return nil, fmt.Errorf("注册失败: %s", resp.Message)
	}

	// 转换protobuf User为model.User
	user := &model.User{
		Username: resp.User.Username,
		Email:    resp.User.Email,
		Nickname: resp.User.Nickname,
		Avatar:   resp.User.Avatar,
	}
	user.ID = uint(resp.User.Id)

	log.Printf("✅ API Gateway: 注册成功 - 用户ID: %d", user.ID)
	return user, nil
}

// Login 通过gRPC调用用户登录
func (s *UserGRPCClientService) Login(username, password string) (string, *model.User, error) {
	log.Printf("🌐 API Gateway: 通过gRPC调用登录 - 用户名: %s", username)

	req := &userpb.LoginRequest{
		Username: username,
		Password: password,
	}

	resp, err := s.client.Login(context.Background(), req)
	if err != nil {
		log.Printf("❌ API Gateway: gRPC调用失败 - %v", err)
		return "", nil, fmt.Errorf("gRPC调用失败: %v", err)
	}

	if resp.Code != 200 {
		log.Printf("❌ API Gateway: 登录失败 - %s", resp.Message)
		return "", nil, fmt.Errorf("登录失败: %s", resp.Message)
	}

	// 转换protobuf User为model.User
	user := &model.User{
		Username: resp.User.Username,
		Nickname: resp.User.Nickname,
		Avatar:   resp.User.Avatar,
	}
	user.ID = uint(resp.User.Id)

	log.Printf("✅ API Gateway: 登录成功 - 用户ID: %d, Token长度: %d", user.ID, len(resp.Token))
	return resp.Token, user, nil
}

// GetUserByUsername 通过gRPC获取用户信息
func (s *UserGRPCClientService) GetUserByUsername(username string) (*model.User, error) {
	log.Printf("🌐 API Gateway: 通过gRPC获取用户 - 用户名: %s", username)

	req := &userpb.GetUserRequest{
		Username: username,
	}

	resp, err := s.client.GetUser(context.Background(), req)
	if err != nil {
		log.Printf("❌ API Gateway: gRPC调用失败 - %v", err)
		return nil, fmt.Errorf("gRPC调用失败: %v", err)
	}

	if resp.Code != 200 {
		log.Printf("❌ API Gateway: 获取用户失败 - %s", resp.Message)
		return nil, fmt.Errorf("获取用户失败: %s", resp.Message)
	}

	// 转换protobuf User为model.User
	user := &model.User{
		Username: resp.User.Username,
		Nickname: resp.User.Nickname,
		Avatar:   resp.User.Avatar,
	}
	user.ID = uint(resp.User.Id)

	log.Printf("✅ API Gateway: 获取用户成功 - 用户ID: %d", user.ID)
	return user, nil
}

// GetUserByID 通过gRPC获取用户信息（通过ID）
func (s *UserGRPCClientService) GetUserByID(userID uint) (*model.User, error) {
	log.Printf("🌐 API Gateway: 通过gRPC获取用户(ID) - 用户ID: %d", userID)

	req := &userpb.GetUserByIDRequest{
		UserId: uint32(userID),
	}

	resp, err := s.client.GetUserByID(context.Background(), req)
	if err != nil {
		log.Printf("❌ API Gateway: gRPC调用失败 - %v", err)
		return nil, fmt.Errorf("gRPC调用失败: %v", err)
	}

	if resp.Code != 200 {
		log.Printf("❌ API Gateway: 获取用户失败 - %s", resp.Message)
		return nil, fmt.Errorf("%s", resp.Message)
	}

	// 转换protobuf User为model.User，包含完整字段
	user := &model.User{
		Username:  resp.User.Username,
		Email:     resp.User.Email,
		Nickname:  resp.User.Nickname,
		Avatar:    resp.User.Avatar,
		AvatarURL: resp.User.Avatar,
		Phone:     resp.User.Phone,
		Bio:       resp.User.Bio,
	}
	user.ID = uint(resp.User.Id)

	log.Printf("✅ API Gateway: 获取用户成功 - 用户ID: %d, Phone: %s, Bio: %s", user.ID, user.Phone, user.Bio)
	return user, nil
}

// UpdateProfileComplete 通过gRPC更新用户完整资料（临时实现）
func (s *UserGRPCClientService) UpdateProfileComplete(userID uint, nickname, avatarURL, phone, bio string) (*model.User, error) {
	log.Printf("🌐 API Gateway: 通过gRPC更新用户资料 (临时实现) - 用户ID: %d", userID)

	// 临时实现：返回模拟的用户数据
	// TODO: 在protobuf定义更新后，实现真正的gRPC调用
	user := &model.User{
		Username:  "test@example.com",
		Email:     "test@example.com",
		Nickname:  nickname,
		Avatar:    avatarURL,
		AvatarURL: avatarURL,
		Phone:     phone,
		Bio:       bio,
	}
	user.ID = userID

	log.Printf("✅ API Gateway: 更新用户资料成功 (临时实现) - 用户ID: %d", user.ID)
	return user, nil
}

// ChangePassword 通过gRPC修改用户密码（临时实现）
func (s *UserGRPCClientService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	log.Printf("🌐 API Gateway: 通过gRPC修改密码 (临时实现) - 用户ID: %d", userID)

	// 临时实现：简单的密码验证逻辑
	// TODO: 在protobuf定义更新后，实现真正的gRPC调用
	if currentPassword == "" || newPassword == "" {
		return fmt.Errorf("密码不能为空")
	}

	if len(newPassword) < 8 {
		return fmt.Errorf("新密码长度不能少于8位")
	}

	log.Printf("✅ API Gateway: 修改密码成功 (临时实现) - 用户ID: %d", userID)
	return nil
}
