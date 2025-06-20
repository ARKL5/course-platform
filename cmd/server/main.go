package main

import (
	"log"

	_ "course-platform/docs"
	"course-platform/internal/configs"
	contentModel "course-platform/internal/domain/content/model"
	courseModel "course-platform/internal/domain/course/model"
	userModel "course-platform/internal/domain/user/model"
	"course-platform/internal/infrastructure/db"
	router "course-platform/internal/transport/http"
)

// @title 课程学习网站 API
// @version 1.0
// @description 这是课程学习网站的 API 文件。
// @host localhost:8083
// @BasePath /api/v1
func main() {
	// 加载配置文件
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 初始化数据库连接
	database, err := db.InitMySQL(config.MySQL)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("✅ 成功连接到 MySQL 数据库")

	// 初始化 Redis 连接
	redisClient, err := db.InitRedis(config.Redis)
	if err != nil {
		log.Printf("⚠️ Redis 连接失败，将使用无缓存模式: %v", err)
		redisClient = nil
	} else {
		log.Println("✅ 成功连接到 Redis")
	}

	// 自动迁移数据库结构
	if err := database.AutoMigrate(
		&userModel.User{},
		&courseModel.Course{},
		&contentModel.FileInfo{},
		&contentModel.File{},
	); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	log.Println("✅ 数据库迁移完成")

	// 设置路由和依赖注入
	r := router.SetupRouter(database, redisClient)

	// 使用配置文件中的端口启动服务器
	log.Printf("🚀 服务器启动在端口: %s", config.Server.Port)
	r.Run(config.Server.Port)
}
