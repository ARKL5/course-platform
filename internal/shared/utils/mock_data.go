package utils

import "github.com/gin-gonic/gin"

// GetMockCourses 获取模拟课程数据
func GetMockCourses() []gin.H {
	return []gin.H{
		{
			"Title":        "Go语言进阶开发",
			"TeacherName":  "张三",
			"CoverImage":   "https://images.unsplash.com/photo-1555066931-4365d14bab8c?w=400&h=225&fit=crop",
			"Rating":       4.9,
			"StudentCount": 1250,
			"ID":           1,
			"Price":        299.00,
			"Category":     "编程开发",
		},
		{
			"Title":        "微服务架构设计",
			"TeacherName":  "李四",
			"CoverImage":   "https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=400&h=225&fit=crop",
			"Rating":       4.8,
			"StudentCount": 850,
			"ID":           2,
			"Price":        399.00,
			"Category":     "架构设计",
		},
		{
			"Title":        "Docker容器技术",
			"TeacherName":  "王五",
			"CoverImage":   "https://images.unsplash.com/photo-1605810230434-7631ac76ec81?w=400&h=225&fit=crop",
			"Rating":       4.7,
			"StudentCount": 1580,
			"ID":           3,
			"Price":        249.00,
			"Category":     "运维部署",
		},
		{
			"Title":        "Kubernetes实战",
			"TeacherName":  "赵六",
			"CoverImage":   "https://images.unsplash.com/photo-1667372393086-9d4001d51cf1?w=400&h=225&fit=crop",
			"Rating":       4.9,
			"StudentCount": 720,
			"ID":           4,
			"Price":        459.00,
			"Category":     "容器编排",
		},
		{
			"Title":        "React高级开发",
			"TeacherName":  "李明",
			"CoverImage":   "https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=400&h=225&fit=crop",
			"Rating":       4.8,
			"StudentCount": 950,
			"ID":           5,
			"Price":        329.00,
			"Category":     "前端开发",
		},
		{
			"Title":        "Vue.js 3.0全栈开发",
			"TeacherName":  "陈小红",
			"CoverImage":   "https://images.unsplash.com/photo-1593720213428-28a5b9e94613?w=400&h=225&fit=crop",
			"Rating":       4.7,
			"StudentCount": 1100,
			"ID":           6,
			"Price":        359.00,
			"Category":     "前端开发",
		},
		{
			"Title":        "Python机器学习入门",
			"TeacherName":  "刘博士",
			"CoverImage":   "https://images.unsplash.com/photo-1555949963-aa79dcee981c?w=400&h=225&fit=crop",
			"Rating":       4.9,
			"StudentCount": 680,
			"ID":           7,
			"Price":        399.00,
			"Category":     "人工智能",
		},
		{
			"Title":        "数据库设计与优化",
			"TeacherName":  "周工",
			"CoverImage":   "https://images.unsplash.com/photo-1544383835-bda2bc66a55d?w=400&h=225&fit=crop",
			"Rating":       4.6,
			"StudentCount": 890,
			"ID":           8,
			"Price":        279.00,
			"Category":     "数据库",
		},
	}
}
