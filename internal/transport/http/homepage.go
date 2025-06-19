package router

import (
	"log"
	"net/http"

	grpcClient "course-platform/internal/infrastructure/grpc_client"
	"course-platform/internal/shared/pb/coursepb"
	"course-platform/internal/shared/utils"

	"github.com/gin-gonic/gin"
)

// HomepageHandler 首页处理器
type HomepageHandler struct {
	courseService *grpcClient.CourseGRPCClientService
}

// NewHomepageHandler 创建首页处理器
func NewHomepageHandler(courseService *grpcClient.CourseGRPCClientService) *HomepageHandler {
	return &HomepageHandler{
		courseService: courseService,
	}
}

// HandleHomepage 处理首页请求
func (h *HomepageHandler) HandleHomepage(c *gin.Context) {
	hotCourses := h.getHotCourses(c)
	continueCourses := h.getContinueCourses(c)

	c.HTML(http.StatusOK, "index.html", gin.H{
		"SiteName":        "Course Platform",
		"HotCourses":      hotCourses,
		"ContinueCourses": continueCourses,
	})
}

// getHotCourses 获取热门课程数据
func (h *HomepageHandler) getHotCourses(c *gin.Context) []gin.H {
	// 从数据库获取热门课程
	coursesResp, err := h.courseService.GetCourses(c, 1, 8, 0) // 获取前8门课程，不限制分类

	if err != nil {
		log.Printf("❌ 获取课程列表失败: %v", err)
		log.Printf("🔄 使用备用模拟数据")
		return utils.GetMockCourses()
	}

	// 成功获取课程数据，转换为前端需要的格式
	log.Printf("✅ 成功获取课程数据，课程数量: %d", len(coursesResp.Courses))
	return h.convertCoursesToDisplay(coursesResp.Courses)
}

// convertCoursesToDisplay 转换课程数据为显示格式
func (h *HomepageHandler) convertCoursesToDisplay(courses []*coursepb.Course) []gin.H {
	// 预定义的展示数据
	teacherNames := []string{"张三", "李四", "王五", "赵六", "李明", "陈小红", "刘博士", "周工"}
	categories := []string{"编程开发", "架构设计", "运维部署", "容器编排", "前端开发", "人工智能", "数据库"}
	ratings := []float32{4.9, 4.8, 4.7, 4.9, 4.8, 4.7, 4.9, 4.6}
	studentCounts := []int{1250, 850, 1580, 720, 950, 1100, 680, 890}

	var hotCourses []gin.H
	for i, course := range courses {
		displayData := h.getDisplayDataForCourse(i, teacherNames, categories, ratings, studentCounts)

		hotCourses = append(hotCourses, gin.H{
			"Title":        course.Title,
			"TeacherName":  displayData.TeacherName,
			"CoverImage":   course.CoverImage,
			"Rating":       displayData.Rating,
			"StudentCount": displayData.StudentCount,
			"ID":           course.Id,
			"Price":        course.Price,
			"Category":     displayData.Category,
		})
	}

	return hotCourses
}

// CourseDisplayData 课程显示数据
type CourseDisplayData struct {
	TeacherName  string
	Category     string
	Rating       float32
	StudentCount int
}

// getDisplayDataForCourse 获取课程的显示数据
func (h *HomepageHandler) getDisplayDataForCourse(index int, teachers []string, categories []string, ratings []float32, counts []int) CourseDisplayData {
	if index < len(teachers) {
		return CourseDisplayData{
			TeacherName:  teachers[index],
			Category:     categories[index%len(categories)],
			Rating:       ratings[index%len(ratings)],
			StudentCount: counts[index%len(counts)],
		}
	}

	// 默认数据
	return CourseDisplayData{
		TeacherName:  "专业讲师",
		Category:     "技术课程",
		Rating:       4.8,
		StudentCount: 1000,
	}
}

// getContinueCourses 获取继续学习课程
func (h *HomepageHandler) getContinueCourses(c *gin.Context) []gin.H {
	// 检查是否有JWT token来判断用户状态
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil // 未登录用户返回nil，触发空状态
	}

	// 模拟有学习记录的用户
	return []gin.H{
		{
			"Title":       "React高级开发技巧",
			"TeacherName": "李明",
			"CoverImage":  "/static/images/react-course.jpg",
			"Duration":    "08:45",
			"Progress":    "Lesson 3 of 12",
		},
	}
}
