package main

import (
	"log"
	"net"
	"os"

	"course-platform/internal/configs"
	"course-platform/internal/domain/content/model"
	"course-platform/internal/domain/content/repository"
	"course-platform/internal/domain/content/service"
	"course-platform/internal/infrastructure/db"
	"course-platform/internal/shared/pb/contentpb"
	"course-platform/internal/transport/grpc"

	grpcServer "google.golang.org/grpc"
)

func main() {
	log.Println("🚀 启动内容微服务...")

	// 加载配置
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	// 连接数据库
	database, err := db.InitMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}
	log.Println("✅ 成功连接到 MySQL 数据库")

	// 数据库迁移
	if err := database.AutoMigrate(&model.FileInfo{}); err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 连接Redis
	rdb, err := db.InitRedis(cfg.Redis)
	if err != nil {
		log.Printf("⚠️ Redis 初始化失败，将在无缓存模式下运行: %v", err)
		rdb = nil
	} else {
		log.Println("✅ 成功连接到 Redis")
	}

	// 创建文件上传目录
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("❌ 创建上传目录失败: %v", err)
	}

	// 初始化仓库层
	contentRepo := repository.NewContentRepository(database, rdb)

	// 初始化服务层
	baseURL := "http://localhost:8083/uploads" // API Gateway作为文件访问代理
	contentService := service.NewContentService(contentRepo, uploadDir, baseURL)

	// 初始化gRPC处理器
	contentHandler := grpc.NewContentHandler(contentService)

	// 创建gRPC服务器
	grpcSrv := grpcServer.NewServer()
	contentpb.RegisterContentServiceServer(grpcSrv, contentHandler)

	// 启动gRPC服务器
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("❌ 监听端口失败: %v", err)
	}

	log.Println("🎯 内容微服务正在监听端口 :50053")
	log.Println("📁 文件上传目录:", uploadDir)
	log.Println("🌐 文件访问URL前缀:", baseURL)

	if err := grpcSrv.Serve(lis); err != nil {
		log.Fatalf("❌ gRPC服务启动失败: %v", err)
	}
}
