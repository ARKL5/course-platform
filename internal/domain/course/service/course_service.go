package service

import (
	"errors"
	"log"
	"strings"

	"course-platform/internal/domain/course/model"
	"course-platform/internal/domain/course/repository"
	userRepository "course-platform/internal/domain/user/repository"
)

// CourseServiceInterface 课程服务接口
type CourseServiceInterface interface {
	CreateCourse(title, description string, instructorID, categoryID uint, price float32, coverImage string) (*model.Course, error)
	GetCourseByID(id uint) (*model.Course, error)
	GetCoursesList(page, pageSize uint, categoryID uint) ([]*model.Course, uint, error)
	UpdateCourse(id uint, title, description string, categoryID uint, price float32, coverImage string) (*model.Course, error)
	DeleteCourse(id uint) error
	PublishCourse(id uint) error
	GetCoursesByInstructor(instructorID uint) ([]*model.Course, error)
}

// CourseService 课程服务实现
type CourseService struct {
	courseRepo repository.CourseRepositoryInterface
	userRepo   userRepository.UserRepositoryInterface
}

// NewCourseService 创建课程服务实例
func NewCourseService(courseRepo repository.CourseRepositoryInterface, userRepo userRepository.UserRepositoryInterface) CourseServiceInterface {
	return &CourseService{
		courseRepo: courseRepo,
		userRepo:   userRepo,
	}
}

// CreateCourse 创建课程
func (s *CourseService) CreateCourse(title, description string, instructorID, categoryID uint, price float32, coverImage string) (*model.Course, error) {
	log.Printf("🔍 Service: 创建课程 - 标题: %s, 讲师ID: %d", title, instructorID)

	// 验证输入参数
	if err := s.validateCourseInput(title, description, instructorID); err != nil {
		return nil, err
	}

	// 验证讲师是否存在
	instructor, err := s.userRepo.GetByID(instructorID)
	var teacherName string
	if err != nil {
		log.Printf("⚠️ Service: 讲师不存在 - ID: %d，使用默认信息", instructorID)
		// 为了兼容性，如果讲师不存在，使用默认讲师信息
		teacherName = "默认讲师"
	} else {
		teacherName = instructor.Nickname
		if teacherName == "" {
			teacherName = instructor.Username
		}
	}

	// 创建课程对象
	course := &model.Course{
		Title:        title,
		Description:  description,
		InstructorID: instructorID,
		CategoryID:   categoryID,
		Price:        price,
		CoverImage:   coverImage,
		Status:       "draft",     // 默认为草稿状态
		TeacherName:  teacherName, // 设置讲师姓名以保持兼容性
	}

	// 创建课程
	if err := s.courseRepo.Create(course); err != nil {
		log.Printf("❌ Service: 创建课程失败 - %v", err)
		return nil, err
	}

	log.Printf("✅ Service: 课程创建成功 - ID: %d", course.ID)
	return course, nil
}

// GetCourseByID 根据ID获取课程
func (s *CourseService) GetCourseByID(id uint) (*model.Course, error) {
	log.Printf("🔍 Service: 获取课程 - ID: %d", id)

	if id == 0 {
		return nil, errors.New("课程ID不能为空")
	}

	course, err := s.courseRepo.GetByID(id)
	if err != nil {
		log.Printf("❌ Service: 获取课程失败 - %v", err)
		return nil, err
	}

	log.Printf("✅ Service: 获取课程成功 - ID: %d", course.ID)
	return course, nil
}

// GetCoursesList 获取课程列表
func (s *CourseService) GetCoursesList(page, pageSize uint, categoryID uint) ([]*model.Course, uint, error) {
	log.Printf("🔍 Service: 获取课程列表 - 页码: %d, 页大小: %d, 分类ID: %d", page, pageSize, categoryID)

	// 设置默认分页参数
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // 限制最大页大小
	}

	courses, total, err := s.courseRepo.GetList(page, pageSize, categoryID)
	if err != nil {
		log.Printf("❌ Service: 获取课程列表失败 - %v", err)
		return nil, 0, err
	}

	log.Printf("✅ Service: 获取课程列表成功 - 数量: %d, 总数: %d", len(courses), total)
	return courses, total, nil
}

// UpdateCourse 更新课程
func (s *CourseService) UpdateCourse(id uint, title, description string, categoryID uint, price float32, coverImage string) (*model.Course, error) {
	log.Printf("🔍 Service: 更新课程 - ID: %d", id)

	if id == 0 {
		return nil, errors.New("课程ID不能为空")
	}

	// 获取现有课程
	course, err := s.courseRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 验证输入参数
	if err := s.validateCourseUpdateInput(title, description); err != nil {
		return nil, err
	}

	// 更新字段
	if title != "" {
		course.Title = title
	}
	if description != "" {
		course.Description = description
	}
	if categoryID > 0 {
		course.CategoryID = categoryID
	}
	if price >= 0 {
		course.Price = price
	}
	if coverImage != "" {
		course.CoverImage = coverImage
	}

	// 保存更新
	if err := s.courseRepo.Update(course); err != nil {
		log.Printf("❌ Service: 更新课程失败 - %v", err)
		return nil, err
	}

	log.Printf("✅ Service: 课程更新成功 - ID: %d", course.ID)
	return course, nil
}

// DeleteCourse 删除课程
func (s *CourseService) DeleteCourse(id uint) error {
	log.Printf("🔍 Service: 删除课程 - ID: %d", id)

	if id == 0 {
		return errors.New("课程ID不能为空")
	}

	// 检查课程是否存在
	_, err := s.courseRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除课程
	if err := s.courseRepo.Delete(id); err != nil {
		log.Printf("❌ Service: 删除课程失败 - %v", err)
		return err
	}

	log.Printf("✅ Service: 课程删除成功 - ID: %d", id)
	return nil
}

// PublishCourse 发布课程
func (s *CourseService) PublishCourse(id uint) error {
	log.Printf("🔍 Service: 发布课程 - ID: %d", id)

	if id == 0 {
		return errors.New("课程ID不能为空")
	}

	// 获取课程
	course, err := s.courseRepo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查课程是否可以发布
	if err := s.validateCourseForPublish(course); err != nil {
		return err
	}

	// 更新状态为已发布
	course.Status = "published"

	if err := s.courseRepo.Update(course); err != nil {
		log.Printf("❌ Service: 发布课程失败 - %v", err)
		return err
	}

	log.Printf("✅ Service: 课程发布成功 - ID: %d", id)
	return nil
}

// GetCoursesByInstructor 获取讲师的课程列表
func (s *CourseService) GetCoursesByInstructor(instructorID uint) ([]*model.Course, error) {
	log.Printf("🔍 Service: 获取讲师课程 - 讲师ID: %d", instructorID)

	if instructorID == 0 {
		return nil, errors.New("讲师ID不能为空")
	}

	// 验证讲师是否存在
	_, err := s.userRepo.GetByID(instructorID)
	if err != nil {
		log.Printf("❌ Service: 讲师不存在 - ID: %d", instructorID)
		return nil, errors.New("讲师不存在")
	}

	courses, err := s.courseRepo.GetByInstructorID(instructorID)
	if err != nil {
		log.Printf("❌ Service: 获取讲师课程失败 - %v", err)
		return nil, err
	}

	log.Printf("✅ Service: 获取讲师课程成功 - 数量: %d", len(courses))
	return courses, nil
}

// 私有验证方法

// validateCourseInput 验证课程创建输入
func (s *CourseService) validateCourseInput(title, description string, instructorID uint) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("课程标题不能为空")
	}

	if len(title) > 200 {
		return errors.New("课程标题不能超过200个字符")
	}

	if len(description) > 5000 {
		return errors.New("课程描述不能超过5000个字符")
	}

	if instructorID == 0 {
		return errors.New("讲师ID不能为空")
	}

	return nil
}

// validateCourseUpdateInput 验证课程更新输入
func (s *CourseService) validateCourseUpdateInput(title, description string) error {
	if title != "" {
		if len(title) > 200 {
			return errors.New("课程标题不能超过200个字符")
		}
	}

	if description != "" {
		if len(description) > 5000 {
			return errors.New("课程描述不能超过5000个字符")
		}
	}

	return nil
}

// validateCourseForPublish 验证课程是否可以发布
func (s *CourseService) validateCourseForPublish(course *model.Course) error {
	if course.Status == "published" {
		return errors.New("课程已经发布")
	}

	if strings.TrimSpace(course.Title) == "" {
		return errors.New("课程标题不能为空")
	}

	if strings.TrimSpace(course.Description) == "" {
		return errors.New("课程描述不能为空")
	}

	if course.InstructorID == 0 {
		return errors.New("课程必须有讲师")
	}

	return nil
}
