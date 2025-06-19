package main

import (
	"log"
	"net"

	"course-platform/internal/configs"
	"course-platform/internal/domain/user/model"
	"course-platform/internal/domain/user/repository"
	"course-platform/internal/domain/user/service"
	"course-platform/internal/infrastructure/db"
	"course-platform/internal/shared/pb/userpb"
	"course-platform/internal/transport/grpc"

	grpcServer "google.golang.org/grpc"
)

// @title 用户微服务 gRPC API
// @version 1.0
// @description 用户管理微服务，处理用户注册、登录、资料管理等功能
func main() {
	log.Println("🚀 启动用户微服务...")

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
		&model.User{},
	)
	if err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 4. 初始化 Redis 连接
	redisClient, err := db.InitRedis(config.Redis)
	if err != nil {
		log.Printf("⚠️ Redis 初始化失敗，將在無快取模式下運行: %v", err)
		redisClient = nil
	} else {
		log.Println("✅ 成功連接到 Redis")
	}

	// 5. 初始化仓储层
	userRepo := repository.NewUserRepository(database, redisClient)

	// 6. 初始化服务层
	userService := service.NewUserService(userRepo)

	// 7. 初始化gRPC处理器
	userHandler := grpc.NewUserHandler(userService)

	// 8. 创建gRPC服务器
	grpcSrv := grpcServer.NewServer()

	// 9. 注册用户服务
	userpb.RegisterUserServiceServer(grpcSrv, userHandler)

	// 10. 创建监听器
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("❌ 创建监听器失败: %v", err)
	}

	log.Println("✅ 用户微服务启动成功，监听端口 :50051")
	log.Println("🔍 等待gRPC请求...")

	// 11. 启动gRPC服务器
	if err := grpcSrv.Serve(listener); err != nil {
		log.Fatalf("❌ gRPC服务器启动失败: %v", err)
	}
}
