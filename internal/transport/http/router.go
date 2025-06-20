package router

import (
	"log"

	_ "course-platform/docs"
	"course-platform/internal/configs"
	contentHandler "course-platform/internal/domain/content/handler"
	courseHandler "course-platform/internal/domain/course/handler"
	userHandler "course-platform/internal/domain/user/handler"
	"course-platform/internal/domain/user/repository"
	"course-platform/internal/domain/user/service"
	grpcClient "course-platform/internal/infrastructure/grpc_client"
	"course-platform/internal/shared/middleware"
	templatefuncs "course-platform/internal/shared/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// SetupRouter 设置路由和所有依赖注入
func SetupRouter(db *gorm.DB, rdb *redis.Client) *gin.Engine {
	// 初始化 Gin 引擎
	r := gin.Default()

	// 设置模板和静态文件
	setupTemplatesAndStatic(r)

	// 初始化服务
	services := initializeServices(db, rdb)

	// 初始化处理器
	handlers := initializeHandlers(services)

	// 设置路由
	setupAllRoutes(r, handlers, services.CourseGRPCService)

	return r
}

// setupTemplatesAndStatic 设置模板和静态文件
func setupTemplatesAndStatic(r *gin.Engine) {
	// 设置模板函数
	r.SetFuncMap(templatefuncs.GetTemplateFunctions())

	// 获取静态文件配置
	staticConfig := configs.GetStaticPathConfig()

	// 加载HTML模板
	r.LoadHTMLGlob(staticConfig.TemplateGlob)

	// 设定静态文件路径
	r.Static("/static", staticConfig.StaticDir)
	r.Static("/uploads", staticConfig.UploadsDir)
}

// Services 服务集合
type Services struct {
	CourseGRPCService  *grpcClient.CourseGRPCClientService
	ContentGRPCService *grpcClient.ContentGRPCClientService
	UserGRPCService    *grpcClient.UserGRPCClientService
	UserService        service.UserServiceInterface
}

// initializeServices 初始化所有服务
func initializeServices(db *gorm.DB, rdb *redis.Client) *Services {
	// 获取服务地址配置
	addresses := configs.GetServiceAddresses()

	log.Printf("🔗 初始化gRPC客户端连接...")

	// 初始化 gRPC 客户端服务
	courseGRPCService, err := grpcClient.NewCourseGRPCClientService(addresses.CourseService)
	if err != nil {
		log.Fatalf("❌ 初始化课程gRPC客户端失败: %v", err)
	}

	contentGRPCService, err := grpcClient.NewContentGRPCClientService(addresses.ContentService)
	if err != nil {
		log.Fatalf("❌ 初始化内容gRPC客户端失败: %v", err)
	}

	userGRPCService, err := grpcClient.NewUserGRPCClientService()
	if err != nil {
		log.Fatalf("❌ 初始化用户gRPC客户端失败: %v", err)
	}

	log.Printf("✅ 所有gRPC客户端初始化成功")

	// 初始化仓储层和业务服务层
	userRepo := repository.NewUserRepository(db, rdb)
	userService := service.NewUserService(userRepo)

	return &Services{
		CourseGRPCService:  courseGRPCService,
		ContentGRPCService: contentGRPCService,
		UserGRPCService:    userGRPCService,
		UserService:        userService,
	}
}

// initializeHandlers 初始化所有处理器
func initializeHandlers(services *Services) *RouteHandlers {
	return &RouteHandlers{
		UserHandler:    userHandler.NewUserHandler(services.UserGRPCService, services.UserService),
		CourseHandler:  courseHandler.NewCourseHandler(services.CourseGRPCService),
		ContentHandler: contentHandler.NewContentHandler(services.ContentGRPCService),
	}
}

// setupAllRoutes 设置所有路由
func setupAllRoutes(r *gin.Engine, handlers *RouteHandlers, courseService *grpcClient.CourseGRPCClientService) {
	// 设置基础路由（健康检查、调试API、Swagger等）
	setupBasicRoutes(r, handlers)

	// 设置页面路由
	setupPageRoutes(r, handlers)

	// 设置API路由
	setupAPIRoutes(r, handlers)

	// 设置首页路由（使用专门的首页处理器）
	setupHomepageRoute(r, courseService)
}

// setupHomepageRoute 设置首页路由
func setupHomepageRoute(r *gin.Engine, courseService *grpcClient.CourseGRPCClientService) {
	homepageHandler := NewHomepageHandler(courseService)
	r.GET("/", homepageHandler.HandleHomepage)
}

// setupPageRoutes 设置页面路由
func setupPageRoutes(r *gin.Engine, handlers *RouteHandlers) {
	// 课程相关页面路由
	r.GET("/course/:id", handlers.CourseHandler.CourseDetailPage)
	r.GET("/courses", handlers.CourseHandler.CoursesListPage)

	// 认证页面路由
	// 静态页面路由 (无需认证)
	r.GET("/login", handlers.UserHandler.LoginPage)
	r.GET("/register", handlers.UserHandler.RegisterPage)

	// 需要可选认证的页面
	dashboardRoutes := r.Group("/")
	dashboardRoutes.Use(middleware.OptionalAuthMiddleware())
	{
		dashboardRoutes.GET("/dashboard", handlers.UserHandler.DashboardPage)
		dashboardRoutes.GET("/creator/dashboard", handlers.UserHandler.CreatorDashboardPage)
	}
}

// setupAPIRoutes 设置API路由
func setupAPIRoutes(r *gin.Engine, handlers *RouteHandlers) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 用户相关路由 (无需认证)
		v1.POST("/register", handlers.UserHandler.Register)
		v1.POST("/login", handlers.UserHandler.Login)
		v1.POST("/validate-token", handlers.UserHandler.ValidateToken)
		v1.POST("/analytics", handlers.UserHandler.Analytics)

		// 可选认证的路由 (支持演示模式)
		optional := v1.Group("/")
		optional.Use(middleware.OptionalAuthMiddleware())
		{
			// 课程相关 (支持演示模式)
			optional.GET("/courses", handlers.CourseHandler.GetCourses)
			optional.GET("/courses/:id", handlers.CourseHandler.GetCourse)
			optional.POST("/courses", handlers.CourseHandler.CreateCourse)
			optional.PUT("/courses/:id", handlers.CourseHandler.UpdateCourse)
			optional.POST("/courses/:id/publish", handlers.CourseHandler.PublishCourse)
			optional.GET("/courses/search", handlers.CourseHandler.SearchCourses)

			// 创作者相关 (支持演示模式)
			optional.GET("/creator/stats", handlers.UserHandler.GetCreatorStats)

			// 内容相关 - 文件列表支持演示模式
			optional.GET("/content/files", handlers.ContentHandler.GetFiles)
		}

		// 需要认证的路由 (必须登录)
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware())
		{
			// 用户相关
			auth.GET("/me", handlers.UserHandler.GetMe)
			auth.PUT("/user/profile", handlers.UserHandler.UpdateProfile)
			auth.PUT("/user/password", handlers.UserHandler.ChangePassword)

			// 内容相关 - 需要登录
			auth.POST("/content/upload", handlers.ContentHandler.UploadFile)
			auth.DELETE("/content/files/:id", handlers.ContentHandler.DeleteFile)
		}
	}
}

// RouteHandlers 路由处理器集合
type RouteHandlers struct {
	UserHandler    *userHandler.UserHandler
	CourseHandler  *courseHandler.CourseHandler
	ContentHandler *contentHandler.ContentHandler
}

// setupBasicRoutes 设置基础路由
func setupBasicRoutes(r *gin.Engine, handlers *RouteHandlers) {
	// Swagger API文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 图片测试页面
	r.GET("/test-image", func(c *gin.Context) {
		c.File("./test-image.html")
	})

	// favicon路由
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.File("./web/static/favicon.ico")
	})
	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
			"status":  "healthy",
		})
	})

	// 调试API
	r.GET("/debug/courses", func(c *gin.Context) {
		log.Printf("🔄 调试API: 返回完整课程数据")
		c.JSON(200, gin.H{
			"success":      true,
			"course_count": 8,
			"courses":      "mock_data",
		})
	})
}
