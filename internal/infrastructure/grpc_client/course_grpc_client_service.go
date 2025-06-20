package service

import (
	"context"
	"fmt"
	"log"

	"course-platform/internal/shared/pb/coursepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// CourseGRPCClientService 课程gRPC客户端服务
// API Gateway通过此服务与课程微服务通信
type CourseGRPCClientService struct {
	client coursepb.CourseServiceClient
	conn   *grpc.ClientConn
}

// NewCourseGRPCClientService 创建课程gRPC客户端服务
func NewCourseGRPCClientService(courseServiceAddr string) (*CourseGRPCClientService, error) {
	log.Printf("🔍 正在连接课程微服务: %s", courseServiceAddr)

	// 创建gRPC连接
	conn, err := grpc.NewClient(courseServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("❌ 连接课程微服务失败: %v", err)
		return nil, fmt.Errorf("连接课程微服务失败: %w", err)
	}

	// 创建客户端
	client := coursepb.NewCourseServiceClient(conn)

	log.Printf("✅ 成功连接课程微服务")
	return &CourseGRPCClientService{
		client: client,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (s *CourseGRPCClientService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

// CreateCourse 创建课程
func (s *CourseGRPCClientService) CreateCourse(ctx context.Context, title, description string, instructorID, categoryID uint, price float32, coverImage string) (*coursepb.CreateCourseResponse, error) {
	log.Printf("🔍 gRPC Client: 创建课程 - 标题: %s", title)

	req := &coursepb.CreateCourseRequest{
		Title:        title,
		Description:  description,
		InstructorId: uint32(instructorID),
		CategoryId:   uint32(categoryID),
		Price:        price,
		CoverImage:   coverImage,
	}

	resp, err := s.client.CreateCourse(ctx, req)
	if err != nil {
		log.Printf("❌ gRPC Client: 创建课程失败 - %v", err)
		return nil, fmt.Errorf("创建课程失败: %w", err)
	}

	// 检查响应和课程数据
	if resp == nil {
		log.Printf("❌ gRPC Client: 收到空响应")
		return nil, fmt.Errorf("收到空响应")
	}

	if resp.Code != 200 {
		log.Printf("❌ gRPC Client: 创建课程失败 - Code: %d, Message: %s", resp.Code, resp.Message)
		return nil, fmt.Errorf("创建课程失败: %s", resp.Message)
	}

	if resp.Course == nil {
		log.Printf("❌ gRPC Client: 响应中课程数据为空")
		return nil, fmt.Errorf("响应中课程数据为空")
	}

	log.Printf("✅ gRPC Client: 创建课程成功 - 课程ID: %d", resp.Course.Id)
	return resp, nil
}

// GetCourses 获取课程列表
func (s *CourseGRPCClientService) GetCourses(ctx context.Context, page, pageSize uint, categoryID uint) (*coursepb.GetCoursesResponse, error) {
	log.Printf("🔍 gRPC Client: 获取课程列表 - 页码: %d, 页大小: %d", page, pageSize)

	req := &coursepb.GetCoursesRequest{
		Page:       uint32(page),
		PageSize:   uint32(pageSize),
		CategoryId: uint32(categoryID),
	}

	resp, err := s.client.GetCourses(ctx, req)
	if err != nil {
		log.Printf("❌ gRPC Client: 获取课程列表失败 - %v", err)
		return nil, fmt.Errorf("获取课程列表失败: %w", err)
	}

	log.Printf("✅ gRPC Client: 获取课程列表成功 - 数量: %d", len(resp.Courses))
	return resp, nil
}

// GetCourse 获取单个课程
func (s *CourseGRPCClientService) GetCourse(ctx context.Context, courseID uint) (*coursepb.GetCourseResponse, error) {
	log.Printf("🔍 gRPC Client: 获取课程详情 - 课程ID: %d", courseID)

	req := &coursepb.GetCourseRequest{
		CourseId: uint32(courseID),
	}

	resp, err := s.client.GetCourse(ctx, req)
	if err != nil {
		log.Printf("❌ gRPC Client: 获取课程详情失败 - %v", err)
		return nil, fmt.Errorf("获取课程详情失败: %w", err)
	}

	log.Printf("✅ gRPC Client: 获取课程详情成功 - 课程ID: %d", resp.Course.Id)
	return resp, nil
}

// PublishCourse 发布课程
func (s *CourseGRPCClientService) PublishCourse(ctx context.Context, courseID uint) (*coursepb.PublishCourseResponse, error) {
	log.Printf("🔍 gRPC Client: 发布课程 - 课程ID: %d", courseID)

	req := &coursepb.PublishCourseRequest{
		CourseId: uint32(courseID),
	}

	resp, err := s.client.PublishCourse(ctx, req)
	if err != nil {
		log.Printf("❌ gRPC Client: 发布课程失败 - %v", err)
		return nil, fmt.Errorf("发布课程失败: %w", err)
	}

	log.Printf("✅ gRPC Client: 发布课程成功 - 课程ID: %d", resp.Course.Id)
	return resp, nil
}

// UpdateCourse 更新课程
func (s *CourseGRPCClientService) UpdateCourse(ctx context.Context, courseID uint, title, description string, categoryID uint, price float32, coverImage string) (*coursepb.UpdateCourseResponse, error) {
	log.Printf("🔍 gRPC Client: 更新课程 - 课程ID: %d, 标题: %s", courseID, title)

	req := &coursepb.UpdateCourseRequest{
		CourseId:    uint32(courseID),
		Title:       title,
		Description: description,
		CategoryId:  uint32(categoryID),
		Price:       price,
		CoverImage:  coverImage,
	}

	resp, err := s.client.UpdateCourse(ctx, req)
	if err != nil {
		log.Printf("❌ gRPC Client: 更新课程失败 - %v", err)
		return nil, fmt.Errorf("更新课程失败: %w", err)
	}

	// 检查响应
	if resp == nil {
		log.Printf("❌ gRPC Client: 收到空响应")
		return nil, fmt.Errorf("收到空响应")
	}

	if resp.Code != 200 {
		log.Printf("❌ gRPC Client: 更新课程失败 - Code: %d, Message: %s", resp.Code, resp.Message)
		return nil, fmt.Errorf("更新课程失败: %s", resp.Message)
	}

	if resp.Course == nil {
		log.Printf("❌ gRPC Client: 响应中课程数据为空")
		return nil, fmt.Errorf("响应中课程数据为空")
	}

	log.Printf("✅ gRPC Client: 更新课程成功 - 课程ID: %d", resp.Course.Id)
	return resp, nil
}
