package configs

import (
	"log"
	"os"
)

// ServiceConfig 服务配置结构
type ServiceConfig struct {
	Name        string
	Address     string
	DefaultAddr string
	EnvKey      string
}

// ServiceAddresses 所有服务地址配置
type ServiceAddresses struct {
	CourseService  string
	ContentService string
	UserService    string
}

// GetServiceAddresses 获取所有服务地址配置
func GetServiceAddresses() *ServiceAddresses {
	configs := []ServiceConfig{
		{
			Name:        "课程服务",
			EnvKey:      "COURSE_SERVICE_ADDR",
			DefaultAddr: "localhost:50052",
		},
		{
			Name:        "内容服务",
			EnvKey:      "CONTENT_SERVICE_ADDR",
			DefaultAddr: "localhost:50053",
		},
		{
			Name:        "用户服务",
			EnvKey:      "USER_SERVICE_ADDR",
			DefaultAddr: "localhost:50051",
		},
	}

	addresses := &ServiceAddresses{}

	// 配置课程服务地址
	addresses.CourseService = getServiceAddress(configs[0])

	// 配置内容服务地址
	addresses.ContentService = getServiceAddress(configs[1])

	// 配置用户服务地址
	addresses.UserService = getServiceAddress(configs[2])

	// 输出服务配置信息
	logServiceConfiguration(addresses)

	return addresses
}

// getServiceAddress 获取单个服务地址
func getServiceAddress(config ServiceConfig) string {
	addr := os.Getenv(config.EnvKey)
	if addr == "" {
		addr = config.DefaultAddr
		log.Printf("⚙️ %s: 使用默认地址 %s", config.Name, addr)
	} else {
		log.Printf("⚙️ %s: 使用环境变量地址 %s", config.Name, addr)
	}
	return addr
}

// logServiceConfiguration 输出服务配置日志
func logServiceConfiguration(addresses *ServiceAddresses) {
	log.Printf("🔗 gRPC服务地址配置:")
	log.Printf("   📚 课程服务: %s", addresses.CourseService)
	log.Printf("   📁 内容服务: %s", addresses.ContentService)
	log.Printf("   👥 用户服务: %s", addresses.UserService)
}

// StaticPathConfig 静态文件路径配置
type StaticPathConfig struct {
	TemplateGlob string
	StaticDir    string
	UploadsDir   string
}

// GetStaticPathConfig 获取静态文件路径配置
func GetStaticPathConfig() *StaticPathConfig {
	return &StaticPathConfig{
		TemplateGlob: "web/templates/*",
		StaticDir:    "./web/static",
		UploadsDir:   "./uploads",
	}
}
