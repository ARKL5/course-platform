package main

import (
	"fmt"
	"log"
	"time"

	"course-platform/internal/configs"
	"course-platform/internal/domain/course/model"
	courseModel "course-platform/internal/domain/course/model"
	userModel "course-platform/internal/domain/user/model"
	"course-platform/internal/infrastructure/db"
	"course-platform/internal/shared/utils"
)

func main() {
	fmt.Println("🚀 开始初始化课程平台示例数据...")

	// 加载配置
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("❌ 配置加载失败: %v", err)
	}

	// 初始化数据库连接
	database, err := db.InitMySQL(config.MySQL)
	if err != nil {
		log.Fatalf("❌ 数据库连接失败: %v", err)
	}

	fmt.Println("✅ 数据库连接成功")

	// 清理现有数据（可选）
	fmt.Println("🧹 清理现有数据...")
	database.Unscoped().Where("1 = 1").Delete(&courseModel.Course{})
	database.Unscoped().Where("1 = 1").Delete(&userModel.User{})

	// 插入示例用户数据（讲师）
	fmt.Println("👥 插入示例讲师数据...")

	// 创建密码哈希
	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		log.Fatalf("❌ 密码哈希失败: %v", err)
	}

	users := []userModel.User{
		{
			Username:     "zhangsan",
			Email:        "zhangsan@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "张三",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "lisi",
			Email:        "lisi@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "李四",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "wangwu",
			Email:        "wangwu@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "王五",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "zhaoliu",
			Email:        "zhaoliu@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "赵六",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "liming",
			Email:        "liming@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "李明",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "chenxiaohong",
			Email:        "chenxiaohong@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "陈小红",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "liuboshi",
			Email:        "liuboshi@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "刘博士",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Username:     "zhougong",
			Email:        "zhougong@example.com",
			PasswordHash: hashedPassword,
			Nickname:     "周工",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// 批量插入用户数据
	result := database.Create(&users)
	if result.Error != nil {
		log.Fatalf("❌ 插入用户数据失败: %v", result.Error)
	}

	fmt.Printf("✅ 成功插入 %d 位讲师\n", result.RowsAffected)

	// 插入示例课程数据
	courses := []courseModel.Course{
		{
			Title:         "Go语言进阶开发",
			Description:   "从零开始学习Go语言，掌握现代化后端开发技能。包含微服务、并发编程、数据库操作等核心内容。",
			InstructorID:  1,
			TeacherName:   "张三",
			Category:      "编程开发",
			CategoryID:    1,
			Price:         299.00,
			CoverImage:    "https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1516116216624-53e697fedbea?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  1250,
			Rating:        4.9,
			ViewCount:     3200,
		},
		{
			Title:         "微服务架构设计",
			Description:   "深入学习微服务架构设计原理，掌握分布式系统开发。涵盖服务拆分、API网关、服务治理等关键技术。",
			InstructorID:  2,
			TeacherName:   "李四",
			Category:      "架构设计",
			CategoryID:    2,
			Price:         399.00,
			CoverImage:    "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  850,
			Rating:        4.8,
			ViewCount:     2100,
		},
		{
			Title:         "Docker容器技术",
			Description:   "全面掌握Docker容器技术，从基础概念到生产环境部署。学习镜像构建、容器编排、集群管理等实用技能。",
			InstructorID:  3,
			TeacherName:   "王五",
			Category:      "运维部署",
			CategoryID:    3,
			Price:         199.00,
			CoverImage:    "https://images.unsplash.com/photo-1605745341112-85968b19335b?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1605745341112-85968b19335b?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  1580,
			Rating:        4.7,
			ViewCount:     4200,
		},
		{
			Title:         "Kubernetes实战",
			Description:   "Kubernetes集群管理与应用部署实战教程。涵盖Pod管理、服务发现、滚动更新、监控告警等企业级应用场景。",
			InstructorID:  4,
			TeacherName:   "赵六",
			Category:      "容器编排",
			CategoryID:    4,
			Price:         499.00,
			CoverImage:    "https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1667372393119-3d4c48d07fc9?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  720,
			Rating:        4.9,
			ViewCount:     1800,
		},
		{
			Title:         "React高级开发技巧",
			Description:   "React进阶开发课程，学习Hooks、状态管理、性能优化等高级特性。结合实际项目，提升前端开发能力。",
			InstructorID:  5,
			TeacherName:   "李明",
			Category:      "前端开发",
			CategoryID:    5,
			Price:         349.00,
			CoverImage:    "https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  950,
			Rating:        4.8,
			ViewCount:     2800,
		},
		{
			Title:         "Vue.js 3.0全栈开发",
			Description:   "最新Vue.js 3.0框架全栈开发教程。涵盖Composition API、TypeScript集成、服务端渲染等现代化开发技术。",
			InstructorID:  6,
			TeacherName:   "陈小红",
			Category:      "前端开发",
			CategoryID:    5,
			Price:         329.00,
			CoverImage:    "https://images.unsplash.com/photo-1627398242454-45a1465c2479?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1627398242454-45a1465c2479?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  1100,
			Rating:        4.7,
			ViewCount:     3100,
		},
		{
			Title:         "Python机器学习入门",
			Description:   "Python机器学习基础教程，学习数据分析、机器学习算法、深度学习框架等AI核心技术。",
			InstructorID:  7,
			TeacherName:   "刘博士",
			Category:      "人工智能",
			CategoryID:    6,
			Price:         599.00,
			CoverImage:    "https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  680,
			Rating:        4.9,
			ViewCount:     1900,
		},
		{
			Title:         "数据库设计与优化",
			Description:   "MySQL数据库设计与性能优化实战课程。学习索引优化、查询优化、分库分表等高级数据库技术。",
			InstructorID:  8,
			TeacherName:   "周工",
			Category:      "数据库",
			CategoryID:    7,
			Price:         279.00,
			CoverImage:    "https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=400&h=225&fit=crop",
			CoverImageURL: "https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=400&h=225&fit=crop",
			Status:        "published",
			StudentCount:  890,
			Rating:        4.6,
			ViewCount:     2400,
		},
	}

	// 批量插入课程数据
	fmt.Println("📚 插入示例课程数据...")
	result = database.Create(&courses)
	if result.Error != nil {
		log.Fatalf("❌ 插入课程数据失败: %v", result.Error)
	}

	fmt.Printf("✅ 成功插入 %d 门课程\n", result.RowsAffected)

	// 验证数据插入
	var count int64
	database.Model(&model.Course{}).Count(&count)
	fmt.Printf("📊 数据库中共有 %d 门课程\n", count)

	// 插入用户学习进度数据（模拟Continue Watching）
	fmt.Println("📈 创建学习进度示例数据...")

	// 这里可以添加用户学习进度表的示例数据
	// 由于当前主要关注课程显示，我们先专注于课程数据

	fmt.Println("🎉 数据初始化完成！")
	fmt.Println("💡 现在可以启动服务器查看带有真实数据的课程平台了")
}
