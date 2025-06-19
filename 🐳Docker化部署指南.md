# 🐳 课程平台Docker化部署指南

## 📊 改善效果对比

### ❌ 原有问题
- gRPC微服务通信不稳定 (500错误)
- 手动启动多个服务，容易出错
- 端口冲突和服务发现问题
- 环境配置不一致

### ✅ Docker化解决方案
- **统一容器网络**: 微服务通过内部网络稳定通信
- **一键启动**: `./docker-start.ps1 start` 启动所有服务
- **服务发现**: 通过服务名而非IP地址通信
- **环境一致性**: 开发、测试、生产环境完全一致

## 🚀 快速开始

### 1. 前置要求
- Docker Desktop (已安装 ✅)
- PowerShell 5.0+ (Windows自带)

### 2. 一键启动
```powershell
# 启动所有服务
./docker-start.ps1 start

# 查看服务状态
./docker-start.ps1 status

# 查看日志
./docker-start.ps1 logs

# 停止所有服务
./docker-start.ps1 stop
```

### 3. 访问服务
- 🌐 **前端页面**: http://localhost:8083
- 📚 **API文档**: http://localhost:8083/swagger/index.html
- 🔍 **健康检查**: http://localhost:8083/ping

## 🏗️ 架构说明

### 容器化服务列表
```
course-platform-gateway  → API网关 (端口 8083)
├── course-platform-user    → 用户服务 (gRPC 50051)
├── course-platform-course  → 课程服务 (gRPC 50052)
└── course-platform-content → 内容服务 (gRPC 50053)
```

### 网络通信改善
```
原来: API网关 → localhost:50051 (不稳定)
现在: API网关 → user-service:50051 (Docker内部网络)
```

## 🔧 配置文件说明

### Docker Compose配置
```yaml
# docker-compose.yml
services:
  api-gateway:
    environment:
      - USER_SERVICE_ADDR=user-service:50051
      - COURSE_SERVICE_ADDR=course-service:50052
      - CONTENT_SERVICE_ADDR=content-service:50053
```

### 环境变量支持
gRPC客户端现在支持环境变量配置:
```go
serviceAddr := os.Getenv("USER_SERVICE_ADDR")
if serviceAddr == "" {
    serviceAddr = "localhost:50051" // 向后兼容
}
```

## 🛠️ 常用命令

### Docker Compose命令
```powershell
# 构建镜像
docker-compose build

# 启动服务 (后台)
docker-compose up -d

# 查看容器状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 停止并删除容器
docker-compose down
```

### 故障排除
```powershell
# 重新构建并启动
./docker-start.ps1 build

# 查看特定服务日志
docker-compose logs user-service

# 进入容器调试
docker exec -it course-platform-user sh
```

## 📈 性能优化

### 1. 健康检查
每个容器都配置了健康检查:
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD pgrep user-service || exit 1
```

### 2. 资源限制
可以在docker-compose.yml中添加:
```yaml
deploy:
  resources:
    limits:
      cpus: '0.5'
      memory: 512M
```

### 3. 数据持久化
重要数据通过volumes持久化:
```yaml
volumes:
  - ./uploads:/root/uploads
  - ./configs:/root/configs
  - ./logs:/root/logs
```

## 🔄 部署流程

### 开发环境
1. 修改代码
2. `./docker-start.ps1 build` (重新构建)
3. 测试验证

### 生产部署
1. 推送代码到Git仓库
2. 在服务器上执行相同的Docker命令
3. 通过环境变量配置不同环境

## 🎯 核心优势

### 1. **解决通信问题**
- Docker内部网络确保gRPC服务间稳定通信
- 消除了原有的500错误和连接问题

### 2. **简化部署流程**
- 从复杂的多步骤启动脚本 → 一键启动
- 自动处理服务依赖关系

### 3. **提高可维护性**
- 统一的容器管理
- 标准化的日志和监控
- 易于扩展和升级

### 4. **向后兼容**
- 保持所有API接口不变
- 支持传统部署方式
- 渐进式迁移

## 🚧 下一步优化建议

1. **添加数据库容器**: 将MySQL和Redis也容器化
2. **负载均衡**: 添加Nginx容器做负载均衡
3. **服务监控**: 集成Prometheus和Grafana
4. **CI/CD集成**: 自动化构建和部署
5. **集群部署**: 支持Kubernetes部署

## 📝 总结

Docker化改善方案成功解决了：
- ✅ gRPC微服务通信不稳定问题
- ✅ 复杂的手动启动流程
- ✅ 环境配置不一致问题
- ✅ 服务发现和端口冲突

同时保持了：
- ✅ 现有微服务架构不变
- ✅ 所有API接口向后兼容
- ✅ 渐进式部署能力

这是一个**风险可控、收益明显**的改善方案！ 