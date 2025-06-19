package main

import (
	"log"
	"net"

	"course-platform/internal/configs"
	"course-platform/internal/domain/course/model"
	"course-platform/internal/domain/course/repository"
	"course-platform/internal/domain/course/service"
	userRepository "course-platform/internal/domain/user/repository"
	"course-platform/internal/infrastructure/db"
	"course-platform/internal/shared/pb/coursepb"
	"course-platform/internal/transport/grpc"

	grpcServer "google.golang.org/grpc"
)

// @title 课程微服务 gRPC API
// @version 1.0
// @description 课程管理微服务，处理课程列表、详情、搜索等功能
func main() {
	log.Println("🚀 启动课程微服务...")

	// 1. 载入配置文件
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ 载入配置文件失败: %v", err)
	}

	// 2. 初始化数据库连接
	database, err := db.InitMySQL(config.MySQL)
	if err != nil {
		log.Fatalf("❌ MySQL 数据库连接失败: %v", err)
	}
	log.Println("✅ 成功连接到 MySQL 数据库")

	// 3. 数据库自动迁移
	err = database.AutoMigrate(
		&model.Course{},
	)
	if err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 4. 初始化 Redis 连接
	redisClient, err := db.InitRedis(config.Redis)
	if err != nil {
		log.Printf("⚠️ Redis 初始化失败，将在无缓存模式下运行: %v", err)
		redisClient = nil
	} else {
		log.Println("✅ 成功连接到 Redis")
	}

	// 5. 初始化仓储层
	courseRepo := repository.NewCourseRepository(database, redisClient)
	userRepo := userRepository.NewUserRepository(database, redisClient)

	// 6. 初始化服务层
	courseService := service.NewCourseService(courseRepo, userRepo)

	// 7. 初始化gRPC处理器
	courseHandler := grpc.NewCourseHandler(courseService)

	// 8. 创建gRPC服务器
	grpcSrv := grpcServer.NewServer()

	// 9. 注册课程服务
	coursepb.RegisterCourseServiceServer(grpcSrv, courseHandler)

	// 10. 创建监听器
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ 创建监听器失败: %v", err)
	}

	log.Println("✅ 课程微服务启动成功，监听端口 :50052")
	log.Println("🔍 等待gRPC请求...")

	// 11. 启动gRPC服务器
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatalf("❌ gRPC服务器启动失败: %v", err)
	}
}
