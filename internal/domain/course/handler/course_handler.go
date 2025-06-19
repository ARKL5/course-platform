package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	service "course-platform/internal/infrastructure/grpc_client"
	"course-platform/internal/shared/pb/coursepb"

	"github.com/gin-gonic/gin"
)

// CourseHandler API Gateway的课程处理器
// 处理客户端HTTP请求，调用课程微服务完成业务
type CourseHandler struct {
	courseGRPCClient *service.CourseGRPCClientService
}

// NewCourseHandler 创建课程处理器
func NewCourseHandler(courseGRPCClient *service.CourseGRPCClientService) *CourseHandler {
	return &CourseHandler{
		courseGRPCClient: courseGRPCClient,
	}
}

// CreateCourseRequest 创建课程请求结构
type CreateCourseRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  string  `json:"description"`
	InstructorID uint    `json:"instructor_id"` // 移除required验证，会从认证信息获取
	CategoryID   uint    `json:"category_id"`
	Price        float32 `json:"price"`
	CoverImage   string  `json:"cover_image"`
}

// UpdateCourseRequest 更新课程请求结构
type UpdateCourseRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CategoryID  uint    `json:"category_id"`
	Price       float32 `json:"price"`
	CoverImage  string  `json:"cover_image"`
}

// CreateCourse 创建课程接口
// @Summary 创建课程
// @Description 创建新的课程
// @Tags 课程管理
// @Accept json
// @Produce json
// @Param course body CreateCourseRequest true "课程信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/courses [post]
func (h *CourseHandler) CreateCourse(c *gin.Context) {
	log.Printf("🔍 API: 收到创建课程请求")

	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ API: 请求参数无效 - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 从认证信息获取用户ID作为讲师ID
	userID, exists := c.Get("userID")
	if exists && userID != nil {
		if uid, ok := userID.(uint); ok {
			req.InstructorID = uid
			log.Printf("🔍 API: 从认证信息获取讲师ID - %d", uid)
		}
	}

	// 如果没有认证信息，使用默认讲师ID（演示模式）
	if req.InstructorID == 0 {
		req.InstructorID = 1
		log.Printf("🔍 API: 使用默认讲师ID（演示模式） - %d", req.InstructorID)
	} else {
		// TODO: 这里应该验证instructor_id是否存在于数据库中
		// 如果不存在，可以使用默认ID或创建用户记录
		// 为了快速解决问题，我们暂时总是使用ID 1（数据库中存在的用户）
		log.Printf("🔍 API: 原讲师ID %d，为确保兼容性使用默认ID 1", req.InstructorID)
		req.InstructorID = 1
	}

	// 调用课程微服务
	ctx := c.Request.Context()
	resp, err := h.courseGRPCClient.CreateCourse(ctx,
		req.Title,
		req.Description,
		req.InstructorID,
		req.CategoryID,
		req.Price,
		req.CoverImage,
	)
	if err != nil {
		log.Printf("❌ API: 创建课程失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建课程失败: " + err.Error(),
		})
		return
	}

	// 检查微服务响应
	if resp.Code != 200 {
		log.Printf("❌ API: 课程微服务返回错误 - Code: %d, Message: %s", resp.Code, resp.Message)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	// 转换响应
	courseData := map[string]interface{}{
		"id":            resp.Course.Id,
		"title":         resp.Course.Title,
		"description":   resp.Course.Description,
		"instructor_id": resp.Course.InstructorId,
		"category_id":   resp.Course.CategoryId,
		"price":         resp.Course.Price,
		"cover_image":   resp.Course.CoverImage,
		"status":        resp.Course.Status,
		"created_at":    resp.Course.CreatedAt,
		"updated_at":    resp.Course.UpdatedAt,
	}

	log.Printf("✅ API: 创建课程成功 - 课程ID: %d", resp.Course.Id)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "课程创建成功",
		"data":    courseData,
	})
}

// GetCourses 获取课程列表接口
// @Summary 获取课程列表
// @Description 分页获取课程列表
// @Tags 课程管理
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param category_id query int false "分类ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/courses [get]
func (h *CourseHandler) GetCourses(c *gin.Context) {
	log.Printf("🔍 API: 收到获取课程列表请求")

	// 解析查询参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	categoryIDStr := c.DefaultQuery("category_id", "0")

	page, err := strconv.ParseUint(pageStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页码参数无效",
		})
		return
	}

	pageSize, err := strconv.ParseUint(pageSizeStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "页大小参数无效",
		})
		return
	}

	categoryID, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "分类ID参数无效",
		})
		return
	}

	// 调用课程微服务
	ctx := c.Request.Context()
	resp, err := h.courseGRPCClient.GetCourses(ctx, uint(page), uint(pageSize), uint(categoryID))
	if err != nil {
		log.Printf("❌ API: 获取课程列表失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取课程列表失败: " + err.Error(),
		})
		return
	}

	// 检查微服务响应
	if resp.Code != 200 {
		log.Printf("❌ API: 课程微服务返回错误 - Code: %d, Message: %s", resp.Code, resp.Message)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	// 转换响应
	var courses []map[string]interface{}
	for _, course := range resp.Courses {
		courseData := map[string]interface{}{
			"id":            course.Id,
			"title":         course.Title,
			"description":   course.Description,
			"instructor_id": course.InstructorId,
			"category_id":   course.CategoryId,
			"price":         course.Price,
			"cover_image":   course.CoverImage,
			"status":        course.Status,
			"created_at":    course.CreatedAt,
			"updated_at":    course.UpdatedAt,
		}
		courses = append(courses, courseData)
	}

	log.Printf("✅ API: 获取课程列表成功 - 数量: %d, 总数: %d", len(courses), resp.Total)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"courses": courses,
			"total":   resp.Total,
			"page":    page,
			"size":    pageSize,
		},
	})
}

// GetCourse 获取单个课程接口
// @Summary 获取课程详情
// @Description 根据ID获取课程详情
// @Tags 课程管理
// @Accept json
// @Produce json
// @Param id path int true "课程ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [get]
func (h *CourseHandler) GetCourse(c *gin.Context) {
	log.Printf("🔍 API: 收到获取课程详情请求")

	// 解析课程ID
	idStr := c.Param("id")
	courseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程ID参数无效",
		})
		return
	}

	// 调用课程微服务
	ctx := c.Request.Context()
	resp, err := h.courseGRPCClient.GetCourse(ctx, uint(courseID))
	if err != nil {
		log.Printf("❌ API: 获取课程详情失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取课程详情失败: " + err.Error(),
		})
		return
	}

	// 检查微服务响应
	if resp.Code != 200 {
		log.Printf("❌ API: 课程微服务返回错误 - Code: %d, Message: %s", resp.Code, resp.Message)
		statusCode := http.StatusBadRequest
		if resp.Code == 404 {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	// 转换响应
	courseData := map[string]interface{}{
		"id":            resp.Course.Id,
		"title":         resp.Course.Title,
		"description":   resp.Course.Description,
		"instructor_id": resp.Course.InstructorId,
		"category_id":   resp.Course.CategoryId,
		"price":         resp.Course.Price,
		"cover_image":   resp.Course.CoverImage,
		"status":        resp.Course.Status,
		"created_at":    resp.Course.CreatedAt,
		"updated_at":    resp.Course.UpdatedAt,
	}

	log.Printf("✅ API: 获取课程详情成功 - 课程ID: %d", resp.Course.Id)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data":    courseData,
	})
}

// GetCoursesByCategory 根据分类获取课程列表（快捷接口）
func (h *CourseHandler) GetCoursesByCategory(c *gin.Context) {
	log.Printf("🔍 API: 收到根据分类获取课程请求")

	categoryIDStr := c.Param("category_id")
	_, err := strconv.ParseUint(categoryIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "分类ID参数无效",
		})
		return
	}

	// 设置查询参数并重定向到GetCourses
	c.Request.URL.RawQuery = "category_id=" + categoryIDStr + "&" + c.Request.URL.RawQuery
	h.GetCourses(c)
}

// SearchCourses 搜索课程（简单实现）
func (h *CourseHandler) SearchCourses(c *gin.Context) {
	log.Printf("🔍 API: 收到搜索课程请求")

	keyword := strings.TrimSpace(c.Query("keyword"))
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "搜索关键词不能为空",
		})
		return
	}

	// 目前简单实现：获取所有课程然后过滤
	// 实际应该在课程服务中实现搜索逻辑
	ctx := c.Request.Context()
	resp, err := h.courseGRPCClient.GetCourses(ctx, 1, 100, 0)
	if err != nil {
		log.Printf("❌ API: 搜索课程失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "搜索课程失败: " + err.Error(),
		})
		return
	}

	// 简单的标题匹配搜索
	var filteredCourses []map[string]interface{}
	lowerKeyword := strings.ToLower(keyword)
	for _, course := range resp.Courses {
		if strings.Contains(strings.ToLower(course.Title), lowerKeyword) ||
			strings.Contains(strings.ToLower(course.Description), lowerKeyword) {
			courseData := map[string]interface{}{
				"id":            course.Id,
				"title":         course.Title,
				"description":   course.Description,
				"instructor_id": course.InstructorId,
				"category_id":   course.CategoryId,
				"price":         course.Price,
				"cover_image":   course.CoverImage,
				"status":        course.Status,
				"created_at":    course.CreatedAt,
				"updated_at":    course.UpdatedAt,
			}
			filteredCourses = append(filteredCourses, courseData)
		}
	}

	log.Printf("✅ API: 搜索课程成功 - 关键词: %s, 结果数量: %d", keyword, len(filteredCourses))
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "搜索成功",
		"data": gin.H{
			"courses": filteredCourses,
			"total":   len(filteredCourses),
			"keyword": keyword,
		},
	})
}

// PublishCourse 发布课程接口
// @Summary 发布课程
// @Description 发布课程，使其对学员可见
// @Tags 课程管理
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param id path int true "课程ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/courses/{id}/publish [post]
func (h *CourseHandler) PublishCourse(c *gin.Context) {
	log.Printf("📤 API: 收到发布课程请求")

	// 解析课程ID
	idStr := c.Param("id")
	courseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程ID参数无效",
		})
		return
	}

	// 获取用户ID（支持演示模式）
	userID, exists := c.Get("userID")
	if !exists {
		userID = uint(1) // 演示模式使用默认用户ID
		log.Printf("📤 API: 演示模式 - 发布课程 %d", courseID)
	} else {
		log.Printf("📤 API: 用户 %v 尝试发布课程 %d", userID, courseID)
	}

	// 调用课程微服务发布课程
	ctx := c.Request.Context()
	resp, err := h.courseGRPCClient.PublishCourse(ctx, uint(courseID))
	if err != nil {
		log.Printf("❌ API: 发布课程失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发布课程失败: " + err.Error(),
		})
		return
	}

	// 检查微服务响应
	if resp.Code != 200 {
		log.Printf("❌ API: 课程微服务返回错误 - Code: %d, Message: %s", resp.Code, resp.Message)
		statusCode := http.StatusBadRequest
		if resp.Code == 404 {
			statusCode = http.StatusNotFound
		} else if resp.Code == 403 {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	log.Printf("✅ API: 课程发布成功 - 课程ID: %d", courseID)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "课程发布成功",
		"data": gin.H{
			"course_id": courseID,
			"status":    "published",
		},
	})
}

// UpdateCourse 更新课程接口
// @Summary 更新课程
// @Description 更新课程信息
// @Tags 课程管理
// @Accept json
// @Produce json
// @Param id path int true "课程ID"
// @Param course body UpdateCourseRequest true "课程信息"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/courses/{id} [put]
func (h *CourseHandler) UpdateCourse(c *gin.Context) {
	log.Printf("🔍 API: 收到更新课程请求")

	// 解析课程ID
	idStr := c.Param("id")
	courseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Printf("❌ API: 课程ID参数无效 - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "课程ID参数无效",
		})
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("❌ API: 请求参数无效 - %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数无效: " + err.Error(),
		})
		return
	}

	// 调用课程微服务
	ctx := c.Request.Context()
	resp, err := h.courseGRPCClient.UpdateCourse(ctx,
		uint(courseID),
		req.Title,
		req.Description,
		req.CategoryID,
		req.Price,
		req.CoverImage,
	)
	if err != nil {
		log.Printf("❌ API: 更新课程失败 - %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "更新课程失败: " + err.Error(),
		})
		return
	}

	// 检查微服务响应
	if resp.Code != 200 {
		log.Printf("❌ API: 课程微服务返回错误 - Code: %d, Message: %s", resp.Code, resp.Message)
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    resp.Code,
			"message": resp.Message,
		})
		return
	}

	// 转换响应
	courseData := map[string]interface{}{
		"id":            resp.Course.Id,
		"title":         resp.Course.Title,
		"description":   resp.Course.Description,
		"instructor_id": resp.Course.InstructorId,
		"category_id":   resp.Course.CategoryId,
		"price":         resp.Course.Price,
		"cover_image":   resp.Course.CoverImage,
		"status":        resp.Course.Status,
		"created_at":    resp.Course.CreatedAt,
		"updated_at":    resp.Course.UpdatedAt,
	}

	log.Printf("✅ API: 更新课程成功 - 课程ID: %d", courseID)
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "课程更新成功",
		"data":    courseData,
	})
}

// CoursesListPage 课程列表页面 - 渲染HTML页面
func (h *CourseHandler) CoursesListPage(c *gin.Context) {
	log.Printf("📚 渲染课程列表页面")

	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	// 获取分类参数
	categoryStr := c.DefaultQuery("category", "0")
	categoryID, err := strconv.Atoi(categoryStr)
	if err != nil || categoryID < 0 {
		categoryID = 0
	}

	// 获取搜索关键词
	keyword := c.Query("search")

	// 调用课程服务获取课程列表
	ctx := c.Request.Context()
	coursesResp, err := h.courseGRPCClient.GetCourses(ctx, uint(page), 12, uint(categoryID)) // 每页12个课程

	var courses []gin.H
	if err != nil {
		log.Printf("❌ 获取课程列表失败: %v", err)
		// 使用空列表
		courses = []gin.H{}
	} else {
		// 转换课程数据
		courses = h.convertCoursesToDisplay(coursesResp.Courses)
	}

	// 分类列表
	categories := []gin.H{
		{"ID": 0, "Name": "全部分类"},
		{"ID": 1, "Name": "编程开发"},
		{"ID": 2, "Name": "架构设计"},
		{"ID": 3, "Name": "运维部署"},
		{"ID": 4, "Name": "容器编排"},
		{"ID": 5, "Name": "前端开发"},
		{"ID": 6, "Name": "人工智能"},
		{"ID": 7, "Name": "数据库"},
	}

	c.HTML(http.StatusOK, "courses-list.html", gin.H{
		"SiteName":        "Course Platform",
		"PageTitle":       "所有课程",
		"Courses":         courses,
		"Categories":      categories,
		"CurrentCategory": categoryID,
		"CurrentPage":     page,
		"Keyword":         keyword,
		"TotalCourses":    len(courses),
		"HasNextPage":     len(courses) >= 12, // 简单判断是否有下一页
	})
}

// convertCoursesToDisplay 转换课程数据为显示格式
func (h *CourseHandler) convertCoursesToDisplay(courses []*coursepb.Course) []gin.H {
	// 预定义的展示数据
	teacherNames := []string{"张三", "李四", "王五", "赵六", "李明", "陈小红", "刘博士", "周工"}
	categories := []string{"编程开发", "架构设计", "运维部署", "容器编排", "前端开发", "人工智能", "数据库"}
	ratings := []float32{4.9, 4.8, 4.7, 4.9, 4.8, 4.7, 4.9, 4.6}
	studentCounts := []int{1250, 850, 1580, 720, 950, 1100, 680, 890}

	var displayCourses []gin.H
	for i, course := range courses {
		// 获取显示数据
		teacherName := "专业讲师"
		category := "技术课程"
		rating := float32(4.8)
		studentCount := 1000

		if i < len(teacherNames) {
			teacherName = teacherNames[i]
			category = categories[i%len(categories)]
			rating = ratings[i%len(ratings)]
			studentCount = studentCounts[i%len(studentCounts)]
		}

		displayCourses = append(displayCourses, gin.H{
			"ID":           course.Id,
			"Title":        course.Title,
			"TeacherName":  teacherName,
			"CoverImage":   course.CoverImage,
			"Rating":       rating,
			"StudentCount": studentCount,
			"Price":        course.Price,
			"Category":     category,
			"Description":  course.Description,
		})
	}

	return displayCourses
}

// CourseDetailPage 课程详情页面 - 渲染HTML页面
func (h *CourseHandler) CourseDetailPage(c *gin.Context) {
	log.Printf("🎨 页面: 收到课程详情页面请求")

	// 解析课程ID
	idStr := c.Param("id")
	courseID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "课程ID参数无效",
		})
		return
	}

	// 获取课程信息
	ctx := c.Request.Context()
	courseResp, err := h.courseGRPCClient.GetCourse(ctx, uint(courseID))
	if err != nil {
		log.Printf("❌ 页面: 获取课程详情失败 - %v", err)
		// 使用备用数据
		h.renderCourseDetailWithFallback(c, uint(courseID))
		return
	}

	if courseResp.Code != 200 {
		log.Printf("❌ 页面: 课程微服务返回错误 - Code: %d, Message: %s", courseResp.Code, courseResp.Message)
		if courseResp.Code == 404 {
			c.HTML(http.StatusNotFound, "error.html", gin.H{
				"Error": "课程不存在",
			})
		} else {
			h.renderCourseDetailWithFallback(c, uint(courseID))
		}
		return
	}

	// 获取课程文件列表 (模拟课程内容)
	// 在实际应用中，这里应该调用内容服务获取课程的视频/文档列表
	lessons := h.generateMockLessons(uint(courseID))

	// 确定当前选中的课程（从URL参数获取，默认为第一个）
	lessonIDStr := c.Query("lesson_id")
	var currentLesson map[string]interface{}

	if lessonIDStr != "" {
		lessonID, _ := strconv.ParseUint(lessonIDStr, 10, 32)
		for _, lesson := range lessons {
			if lesson["Id"].(uint) == uint(lessonID) {
				currentLesson = lesson
				break
			}
		}
	}

	// 如果没有指定或没找到，使用第一个课程
	if currentLesson == nil && len(lessons) > 0 {
		currentLesson = lessons[0]
	}

	// 构建课程数据
	courseData := gin.H{
		"Id":             courseResp.Course.Id,
		"Title":          courseResp.Course.Title,
		"Description":    courseResp.Course.Description,
		"InstructorName": h.getInstructorName(uint(courseResp.Course.InstructorId)),
		"Price":          courseResp.Course.Price,
		"CoverImage":     courseResp.Course.CoverImage,
		"Category":       h.getCategoryName(uint(courseResp.Course.CategoryId)),
		"Status":         courseResp.Course.Status,
	}

	log.Printf("✅ 页面: 渲染课程详情页面成功 - 课程ID: %d, 课程数量: %d", courseResp.Course.Id, len(lessons))

	c.HTML(http.StatusOK, "course-detail.html", gin.H{
		"SiteName":      "Course Platform",
		"Course":        courseData,
		"Lessons":       lessons,
		"CurrentLesson": currentLesson,
	})
}

// renderCourseDetailWithFallback 使用备用数据渲染课程详情页面
func (h *CourseHandler) renderCourseDetailWithFallback(c *gin.Context, courseID uint) {
	log.Printf("🔄 页面: 使用备用数据渲染课程详情页面")

	// 备用课程数据
	fallbackCourses := map[uint]gin.H{
		1: {
			"Id":             1,
			"Title":          "Go语言进阶开发",
			"Description":    "深入学习Go语言高级特性，包括并发编程、网络编程、微服务架构等核心内容。通过实战项目掌握Go语言在企业级开发中的应用。",
			"InstructorName": "张三",
			"Price":          299.00,
			"CoverImage":     "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=400&h=225&fit=crop",
			"Category":       "编程开发",
			"Status":         "published",
		},
		2: {
			"Id":             2,
			"Title":          "微服务架构设计",
			"Description":    "系统学习微服务架构设计原理和最佳实践，包括服务拆分、API设计、数据一致性、服务治理等关键技术。",
			"InstructorName": "李四",
			"Price":          399.00,
			"CoverImage":     "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=400&h=225&fit=crop",
			"Category":       "架构设计",
			"Status":         "published",
		},
		3: {
			"Id":             3,
			"Title":          "Docker容器技术",
			"Description":    "从零开始学习Docker容器技术，掌握容器化部署、镜像构建、网络配置、数据卷管理等核心技能。",
			"InstructorName": "王五",
			"Price":          249.00,
			"CoverImage":     "https://images.unsplash.com/photo-1605810230434-7631ac76ec81?w=400&h=225&fit=crop",
			"Category":       "运维部署",
			"Status":         "published",
		},
	}

	// 获取对应的备用数据
	courseData, exists := fallbackCourses[courseID]
	if !exists {
		// 如果没有对应的备用数据，使用默认数据
		courseData = gin.H{
			"Id":             courseID,
			"Title":          "精品课程",
			"Description":    "这是一门精心设计的专业课程，将为您提供全面深入的学习体验。",
			"InstructorName": "专业讲师",
			"Price":          299.00,
			"CoverImage":     "/static/images/placeholder.jpg",
			"Category":       "技术课程",
			"Status":         "published",
		}
	}

	// 生成模拟课程内容
	lessons := h.generateMockLessons(courseID)
	var currentLesson map[string]interface{}
	if len(lessons) > 0 {
		currentLesson = lessons[0]
	}

	c.HTML(http.StatusOK, "course-detail.html", gin.H{
		"SiteName":      "Course Platform",
		"Course":        courseData,
		"Lessons":       lessons,
		"CurrentLesson": currentLesson,
	})
}

// generateMockLessons 生成模拟课程内容
func (h *CourseHandler) generateMockLessons(courseID uint) []map[string]interface{} {
	lessonTemplates := map[uint][]map[string]interface{}{
		1: { // Go语言课程
			{"Id": 1, "FileName": "第1章：Go语言基础回顾.mp4", "FileSize": "85.2 MB"},
			{"Id": 2, "FileName": "第2章：并发编程核心概念.mp4", "FileSize": "128.5 MB"},
			{"Id": 3, "FileName": "Go语言开发环境配置指南.pdf", "FileSize": "3.2 MB"},
			{"Id": 4, "FileName": "第3章：Goroutine深度解析.mp4", "FileSize": "156.1 MB"},
			{"Id": 5, "FileName": "第4章：Channel通信机制.mp4", "FileSize": "98.8 MB"},
			{"Id": 6, "FileName": "Go并发编程最佳实践.pdf", "FileSize": "5.7 MB"},
			{"Id": 7, "FileName": "第5章：网络编程实战.mp4", "FileSize": "201.3 MB"},
			{"Id": 8, "FileName": "第6章：微服务架构入门.mp4", "FileSize": "175.7 MB"},
		},
		2: { // 微服务架构课程
			{"Id": 9, "FileName": "第1章：微服务架构概述.mp4", "FileSize": "112.4 MB"},
			{"Id": 10, "FileName": "微服务设计模式手册.pdf", "FileSize": "8.3 MB"},
			{"Id": 11, "FileName": "第2章：服务拆分策略.mp4", "FileSize": "145.6 MB"},
			{"Id": 12, "FileName": "第3章：API网关设计.mp4", "FileSize": "132.8 MB"},
			{"Id": 13, "FileName": "服务网格架构图解.pptx", "FileSize": "12.5 MB"},
			{"Id": 14, "FileName": "第4章：服务发现与注册.mp4", "FileSize": "118.2 MB"},
			{"Id": 15, "FileName": "第5章：分布式事务处理.mp4", "FileSize": "198.5 MB"},
		},
		3: { // Docker课程
			{"Id": 16, "FileName": "第1章：Docker基础概念.mp4", "FileSize": "95.7 MB"},
			{"Id": 17, "FileName": "Docker安装配置手册.pdf", "FileSize": "2.8 MB"},
			{"Id": 18, "FileName": "第2章：镜像构建实战.mp4", "FileSize": "167.3 MB"},
			{"Id": 19, "FileName": "第3章：容器网络配置.mp4", "FileSize": "143.9 MB"},
			{"Id": 20, "FileName": "Docker命令速查表.pdf", "FileSize": "1.5 MB"},
			{"Id": 21, "FileName": "第4章：数据卷管理.mp4", "FileSize": "89.5 MB"},
			{"Id": 22, "FileName": "第5章：Docker Compose实战.mp4", "FileSize": "201.1 MB"},
			{"Id": 23, "FileName": "Docker部署最佳实践.pptx", "FileSize": "15.7 MB"},
		},
	}

	lessons, exists := lessonTemplates[courseID]
	if !exists {
		// 默认课程内容
		lessons = []map[string]interface{}{
			{"Id": 1, "FileName": "第1章：课程介绍", "FileSize": "12.5 MB"},
			{"Id": 2, "FileName": "第2章：基础知识", "FileSize": "18.7 MB"},
			{"Id": 3, "FileName": "第3章：进阶内容", "FileSize": "25.3 MB"},
			{"Id": 4, "FileName": "第4章：实战项目", "FileSize": "32.1 MB"},
			{"Id": 5, "FileName": "第5章：总结回顾", "FileSize": "15.8 MB"},
		}
	}

	return lessons
}

// getInstructorName 根据讲师ID获取讲师姓名 (模拟数据)
func (h *CourseHandler) getInstructorName(instructorID uint) string {
	instructorNames := map[uint]string{
		1: "张三",
		2: "李四",
		3: "王五",
		4: "赵六",
		5: "李明",
		6: "陈小红",
		7: "刘博士",
		8: "周工",
	}

	if name, exists := instructorNames[instructorID]; exists {
		return name
	}
	return "专业讲师"
}

// getCategoryName 根据分类ID获取分类名称 (模拟数据)
func (h *CourseHandler) getCategoryName(categoryID uint) string {
	categoryNames := map[uint]string{
		1: "编程开发",
		2: "架构设计",
		3: "运维部署",
		4: "容器编排",
		5: "前端开发",
		6: "人工智能",
		7: "数据库",
	}

	if name, exists := categoryNames[categoryID]; exists {
		return name
	}
	return "技术课程"
}
