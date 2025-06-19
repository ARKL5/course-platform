package grpc

import (
	"context"
	"log"

	"course-platform/internal/domain/course/service"
	"course-platform/internal/shared/pb/coursepb"
)

// CourseHandler 课程gRPC处理器
// 处理来自API Gateway的课程相关gRPC请求，调用课程服务完成业务逻辑
type CourseHandler struct {
	coursepb.UnimplementedCourseServiceServer
	courseService service.CourseServiceInterface
}

// NewCourseHandler 创建课程gRPC处理器实例
func NewCourseHandler(courseService service.CourseServiceInterface) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

// CreateCourse 处理创建课程gRPC请求
func (h *CourseHandler) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CreateCourseResponse, error) {
	log.Printf("🔍 gRPC: 收到创建课程请求 - 标题: %s", req.Title)

	// 调用课程服务创建课程
	course, err := h.courseService.CreateCourse(
		req.Title,
		req.Description,
		uint(req.InstructorId),
		uint(req.CategoryId),
		req.Price,
		req.CoverImage,
	)
	if err != nil {
		log.Printf("❌ gRPC: 创建课程失败 - %v", err)
		return &coursepb.CreateCourseResponse{
			Code:    400,
			Message: err.Error(),
			Course:  nil,
		}, nil
	}

	// 转换为protobuf课程对象
	pbCourse := &coursepb.Course{
		Id:           uint32(course.ID),
		Title:        course.Title,
		Description:  course.Description,
		InstructorId: uint32(course.InstructorID),
		CategoryId:   uint32(course.CategoryID),
		Price:        course.Price,
		CoverImage:   course.CoverImage,
		Status:       course.Status,
		CreatedAt:    course.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    course.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 创建课程成功 - 课程ID: %d", course.ID)
	return &coursepb.CreateCourseResponse{
		Code:    200,
		Message: "课程创建成功",
		Course:  pbCourse,
	}, nil
}

// GetCourses 处理获取课程列表gRPC请求
func (h *CourseHandler) GetCourses(ctx context.Context, req *coursepb.GetCoursesRequest) (*coursepb.GetCoursesResponse, error) {
	log.Printf("🔍 gRPC: 收到获取课程列表请求 - 页码: %d, 页大小: %d", req.Page, req.PageSize)

	// 调用课程服务获取课程列表
	courses, total, err := h.courseService.GetCoursesList(
		uint(req.Page),
		uint(req.PageSize),
		uint(req.CategoryId),
	)
	if err != nil {
		log.Printf("❌ gRPC: 获取课程列表失败 - %v", err)
		return &coursepb.GetCoursesResponse{
			Code:    400,
			Message: err.Error(),
			Courses: nil,
			Total:   0,
		}, nil
	}

	// 转换为protobuf课程列表
	var pbCourses []*coursepb.Course
	for _, course := range courses {
		pbCourse := &coursepb.Course{
			Id:           uint32(course.ID),
			Title:        course.Title,
			Description:  course.Description,
			InstructorId: uint32(course.InstructorID),
			CategoryId:   uint32(course.CategoryID),
			Price:        course.Price,
			CoverImage:   course.CoverImage,
			Status:       course.Status,
			CreatedAt:    course.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:    course.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		pbCourses = append(pbCourses, pbCourse)
	}

	log.Printf("✅ gRPC: 获取课程列表成功 - 数量: %d, 总数: %d", len(courses), total)
	return &coursepb.GetCoursesResponse{
		Code:    200,
		Message: "获取成功",
		Courses: pbCourses,
		Total:   uint32(total),
	}, nil
}

// GetCourse 处理获取单个课程gRPC请求
func (h *CourseHandler) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.GetCourseResponse, error) {
	log.Printf("🔍 gRPC: 收到获取课程请求 - 课程ID: %d", req.CourseId)

	// 调用课程服务获取课程
	course, err := h.courseService.GetCourseByID(uint(req.CourseId))
	if err != nil {
		log.Printf("❌ gRPC: 获取课程失败 - %v", err)
		return &coursepb.GetCourseResponse{
			Code:    404,
			Message: err.Error(),
			Course:  nil,
		}, nil
	}

	// 转换为protobuf课程对象
	pbCourse := &coursepb.Course{
		Id:           uint32(course.ID),
		Title:        course.Title,
		Description:  course.Description,
		InstructorId: uint32(course.InstructorID),
		CategoryId:   uint32(course.CategoryID),
		Price:        course.Price,
		CoverImage:   course.CoverImage,
		Status:       course.Status,
		CreatedAt:    course.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    course.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 获取课程成功 - 课程ID: %d", course.ID)
	return &coursepb.GetCourseResponse{
		Code:    200,
		Message: "获取成功",
		Course:  pbCourse,
	}, nil
}

// mustEmbedUnimplementedCourseServiceServer 实现接口要求
func (h *CourseHandler) mustEmbedUnimplementedCourseServiceServer() {}

// PublishCourse 处理发布课程gRPC请求
func (h *CourseHandler) PublishCourse(ctx context.Context, req *coursepb.PublishCourseRequest) (*coursepb.PublishCourseResponse, error) {
	log.Printf("🔍 gRPC: 收到发布课程请求 - 课程ID: %d", req.CourseId)

	// 调用课程服务发布课程
	err := h.courseService.PublishCourse(uint(req.CourseId))
	if err != nil {
		log.Printf("❌ gRPC: 发布课程失败 - %v", err)
		return &coursepb.PublishCourseResponse{
			Code:    400,
			Message: err.Error(),
			Course:  nil,
		}, nil
	}

	// 获取发布后的课程信息
	course, err := h.courseService.GetCourseByID(uint(req.CourseId))
	if err != nil {
		log.Printf("❌ gRPC: 获取发布后的课程失败 - %v", err)
		return &coursepb.PublishCourseResponse{
			Code:    400,
			Message: err.Error(),
			Course:  nil,
		}, nil
	}

	// 转换为protobuf课程对象
	pbCourse := &coursepb.Course{
		Id:           uint32(course.ID),
		Title:        course.Title,
		Description:  course.Description,
		InstructorId: uint32(course.InstructorID),
		CategoryId:   uint32(course.CategoryID),
		Price:        course.Price,
		CoverImage:   course.CoverImage,
		Status:       course.Status,
		CreatedAt:    course.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    course.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 发布课程成功 - 课程ID: %d", course.ID)
	return &coursepb.PublishCourseResponse{
		Code:    200,
		Message: "课程发布成功",
		Course:  pbCourse,
	}, nil
}

// UpdateCourse 处理更新课程gRPC请求
func (h *CourseHandler) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.UpdateCourseResponse, error) {
	log.Printf("🔍 gRPC: 收到更新课程请求 - 课程ID: %d", req.CourseId)

	// 调用课程服务更新课程
	course, err := h.courseService.UpdateCourse(
		uint(req.CourseId),
		req.Title,
		req.Description,
		uint(req.CategoryId),
		req.Price,
		req.CoverImage,
	)
	if err != nil {
		log.Printf("❌ gRPC: 更新课程失败 - %v", err)
		return &coursepb.UpdateCourseResponse{
			Code:    400,
			Message: err.Error(),
			Course:  nil,
		}, nil
	}

	// 转换为protobuf课程对象
	pbCourse := &coursepb.Course{
		Id:           uint32(course.ID),
		Title:        course.Title,
		Description:  course.Description,
		InstructorId: uint32(course.InstructorID),
		CategoryId:   uint32(course.CategoryID),
		Price:        course.Price,
		CoverImage:   course.CoverImage,
		Status:       course.Status,
		CreatedAt:    course.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    course.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	log.Printf("✅ gRPC: 更新课程成功 - 课程ID: %d", course.ID)
	return &coursepb.UpdateCourseResponse{
		Code:    200,
		Message: "课程更新成功",
		Course:  pbCourse,
	}, nil
}
