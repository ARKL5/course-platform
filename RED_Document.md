# Course Platform - RED 文档
*Rapid Engineering Documentation*

## 📋 目录

1. [项目概述](#1-项目概述)
2. [架构设计](#2-架构设计)
3. [技术栈](#3-技术栈)
4. [模块与功能详解](#4-模块与功能详解)
5. [数据模型与管理](#5-数据模型与管理)
6. [API 设计与使用](#6-api-设计与使用)
7. [关键设计决策](#7-关键设计决策)
8. [安装配置与部署](#8-安装配置与部署)
9. [测试策略](#9-测试策略)
10. [未来工作与路线图](#10-未来工作与路线图)

---

## 1. 项目概述

### 1.1 项目简介
Course Platform 是一个基于微服务架构的在线课程平台，采用 Go 语言开发，支持课程创建、内容管理、用户认证等核心功能。

### 1.2 主要目标
- **教育内容管理**：为教育机构和个人讲师提供完整的课程内容管理解决方案
- **可扩展性**：采用微服务架构，支持水平扩展和模块化开发
- **用户体验**：提供直观的前端界面和完善的API接口

### 1.3 核心功能
- 用户注册、登录和资料管理
- 课程创建、编辑和发布
- 文件上传和内容管理
- 课程浏览和搜索
- 创作者工作台

### 1.4 目标用户
- **讲师/创作者**：创建和管理在线课程
- **学员**：浏览和学习课程内容
- **管理员**：平台运营和内容审核

### 1.5 应用场景
- 在线教育平台
- 企业内训系统
- 知识付费平台
- 技能培训机构

---

## 2. 架构设计

### 2.1 高层架构

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Web Frontend  │    │   Mobile App    │    │  Third Party    │
│   (HTML/CSS/JS) │    │                 │    │    Clients      │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌─────────────▼───────────┐
                    │     API Gateway        │
                    │   (Port: 8083)         │
                    │   HTTP REST + gRPC     │
                    └─────────────┬───────────┘
                                  │
              ┌───────────────────┼───────────────────┐
              │                   │                   │
    ┌─────────▼─────────┐ ┌───────▼───────┐ ┌─────────▼─────────┐
    │  User Service     │ │ Course Service│ │ Content Service   │
    │  (Port: 50051)    │ │ (Port: 50052) │ │ (Port: 50053)     │
    │  gRPC Server      │ │ gRPC Server   │ │ gRPC Server       │
    └─────────┬─────────┘ └───────┬───────┘ └─────────┬─────────┘
              │                   │                   │
              └───────────────────┼───────────────────┘
                                  │
                    ┌─────────────▼───────────┐
                    │    Data Layer          │
                    │  MySQL + Redis         │
                    └─────────────────────────┘
```

### 2.2 微服务组件

#### 2.2.1 API Gateway (端口: 8083)
- **职责**：HTTP请求路由、认证鉴权、静态文件服务
- **技术**：Gin Web Framework
- **功能**：
  - HTTP到gRPC的协议转换
  - JWT token验证
  - 静态资源服务
  - 请求日志和监控

#### 2.2.2 User Service (端口: 50051)
- **职责**：用户管理、认证授权
- **功能**：
  - 用户注册和登录
  - 用户资料管理
  - JWT token生成和验证
  - 密码加密和验证

#### 2.2.3 Course Service (端口: 50052)
- **职责**：课程管理和业务逻辑
- **功能**：
  - 课程CRUD操作
  - 课程分类管理
  - 课程搜索和过滤
  - 课程状态管理

#### 2.2.4 Content Service (端口: 50053)
- **职责**：文件和内容管理
- **功能**：
  - 文件上传和存储
  - 文件元数据管理
  - 文件访问控制
  - 多媒体内容处理

### 2.3 架构模式选择

#### 微服务架构优势
1. **模块化**：每个服务独立开发和部署
2. **可扩展性**：可根据负载独立扩展单个服务
3. **技术多样性**：不同服务可选择最适合的技术栈
4. **容错性**：单个服务故障不影响整个系统

#### gRPC通信选择
1. **性能**：二进制协议，比HTTP/JSON更高效
2. **类型安全**：强类型定义，减少运行时错误
3. **多语言支持**：便于未来扩展不同语言的服务
4. **流式处理**：支持双向流和服务端推送

---

## 3. 技术栈

### 3.1 后端技术
- **语言**：Go 1.24.0
- **Web框架**：Gin v1.10.1
- **RPC框架**：gRPC v1.73.0
- **ORM**：GORM v1.25.12
- **认证**：JWT v5.2.2
- **配置管理**：Viper v1.20.1

### 3.2 数据存储
- **主数据库**：MySQL 8.0
- **缓存**：Redis v8.11.5
- **文件存储**：本地文件系统

### 3.3 前端技术
- **模板引擎**：Go HTML Template
- **样式**：CSS3 + 响应式设计
- **脚本**：原生JavaScript ES6+
- **图标**：Font Awesome

### 3.4 开发工具
- **API文档**：Swagger/OpenAPI
- **代码生成**：Protocol Buffers
- **依赖管理**：Go Modules
- **容器化**：Docker + Docker Compose

---

## 4. 模块与功能详解

### 4.1 用户模块 (User Module)

#### 4.1.1 功能概述
- 用户注册、登录、登出
- 用户资料管理
- 头像上传
- 密码修改

#### 4.1.2 核心组件
```go
// 用户服务接口
type UserServiceInterface interface {
    CreateUser(email, password, nickname string) (*model.User, error)
    GetUserByEmail(email string) (*model.User, error)
    ValidatePassword(user *model.User, password string) bool
    GenerateToken(user *model.User) (string, error)
    UpdateProfile(userID uint, updates map[string]interface{}) error
}
```

#### 4.1.3 关键特性
- **密码安全**：使用bcrypt加密
- **JWT认证**：无状态token机制
- **缓存优化**：Redis缓存用户信息
- **兼容性设计**：支持email/username双重标识

### 4.2 课程模块 (Course Module)

#### 4.2.1 功能概述
- 课程创建和编辑
- 课程分类管理
- 课程搜索和过滤
- 课程状态管理（草稿/发布）

#### 4.2.2 核心组件
```go
// 课程服务接口
type CourseServiceInterface interface {
    CreateCourse(title, description string, instructorID, categoryID uint, price float32, coverImage string) (*model.Course, error)
    GetCourseByID(id uint) (*model.Course, error)
    GetCoursesList(page, pageSize, categoryID uint) ([]*model.Course, uint, error)
    UpdateCourse(id uint, title, description string, categoryID uint, price float32, coverImage string) (*model.Course, error)
    PublishCourse(id uint) error
}
```

#### 4.2.3 业务规则
- 课程状态：draft（草稿）→ published（已发布）
- 封面图片自动处理和兼容性同步
- 分页查询支持
- 分类过滤功能

### 4.3 内容模块 (Content Module)

#### 4.3.1 功能概述
- 文件上传（图片、文档、视频）
- 文件元数据管理
- 文件访问控制
- 存储空间管理

#### 4.3.2 文件类型支持
- **图片**：JPG, PNG, GIF（课程封面、用户头像）
- **文档**：PDF, DOC, DOCX, PPT, PPTX
- **音视频**：MP4, AVI, MP3（课程内容）

#### 4.3.3 存储策略
```
uploads/
├── avatars/           # 用户头像
├── course_{id}/       # 课程文件
│   ├── cover.jpg      # 课程封面
│   ├── video/         # 视频文件
│   └── documents/     # 文档文件
```

---

## 5. 数据模型与管理

### 5.1 数据库设计

#### 5.1.1 用户表 (users)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    nickname VARCHAR(100),
    avatar_url VARCHAR(500),
    phone VARCHAR(20),
    bio VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);
```

#### 5.1.2 课程表 (courses)
```sql
CREATE TABLE courses (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    instructor_id BIGINT NOT NULL,
    category_id BIGINT,
    price DECIMAL(10,2) DEFAULT 0,
    cover_image VARCHAR(500),
    status VARCHAR(20) DEFAULT 'draft',
    student_count INT DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0,
    view_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    INDEX idx_instructor_id (instructor_id),
    INDEX idx_category_id (category_id),
    INDEX idx_status (status)
);
```

#### 5.1.3 文件表 (course_files)
```sql
CREATE TABLE course_files (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    file_name VARCHAR(255) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_size BIGINT NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    course_id BIGINT NOT NULL,
    uploader_id BIGINT NOT NULL,
    upload_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_course_id (course_id),
    INDEX idx_uploader_id (uploader_id)
);
```

### 5.2 数据关系

```
User (1) ──────── (N) Course
  │                    │
  │                    │
  └─── (N) Files ──────┘
```

### 5.3 CRUD操作流程

#### 5.3.1 课程创建流程
1. **验证用户权限**：检查JWT token
2. **数据验证**：验证必填字段和格式
3. **创建记录**：插入课程基础信息
4. **处理封面**：如有封面文件，调用内容服务
5. **返回结果**：返回完整课程信息

#### 5.3.2 文件上传流程
1. **接收文件**：通过multipart/form-data
2. **文件验证**：检查类型、大小限制
3. **生成路径**：按规则生成存储路径
4. **保存文件**：写入文件系统
5. **记录元数据**：保存到数据库
6. **返回URL**：返回访问地址

---

## 6. API 设计与使用

### 6.1 API 架构

#### 6.1.1 RESTful 设计原则
- **资源导向**：URL表示资源，HTTP方法表示操作
- **无状态**：每个请求包含完整信息
- **统一接口**：标准HTTP状态码和响应格式
- **分层系统**：API Gateway作为统一入口

#### 6.1.2 响应格式标准
```json
{
  "code": 200,
  "message": "操作成功",
  "data": {
    // 具体数据
  }
}
```

### 6.2 核心API端点

#### 6.2.1 用户认证API
```http
POST /api/v1/register
POST /api/v1/login
POST /api/v1/validate-token
GET  /api/v1/me
PUT  /api/v1/user/profile
PUT  /api/v1/user/password
```

#### 6.2.2 课程管理API
```http
GET    /api/v1/courses              # 获取课程列表
GET    /api/v1/courses/{id}         # 获取课程详情
POST   /api/v1/courses              # 创建课程
PUT    /api/v1/courses/{id}         # 更新课程
POST   /api/v1/courses/{id}/publish # 发布课程
GET    /api/v1/courses/search       # 搜索课程
```

#### 6.2.3 内容管理API
```http
POST   /api/v1/content/upload       # 上传文件
GET    /api/v1/content/files        # 获取文件列表
DELETE /api/v1/content/files/{id}   # 删除文件
```

### 6.3 认证授权机制

#### 6.3.1 JWT Token 结构
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "user_id": 123,
    "email": "user@example.com",
    "exp": 1640995200
  }
}
```

#### 6.3.2 权限控制
- **公开接口**：课程列表、课程详情
- **可选认证**：支持演示模式的创建操作
- **必需认证**：文件上传、用户资料修改

### 6.4 API 使用示例

#### 6.4.1 创建课程
```bash
curl -X POST http://localhost:8083/api/v1/courses \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Go语言入门",
    "description": "从零开始学习Go语言",
    "category_id": 1,
    "price": 99.00,
    "cover_image": ""
  }'
```

#### 6.4.2 上传文件
```bash
curl -X POST http://localhost:8083/api/v1/content/upload \
  -H "Authorization: Bearer <token>" \
  -F "file=@course-cover.jpg" \
  -F "course_id=17" \
  -F "file_type=image"
```

---

## 7. 关键设计决策

### 7.1 微服务拆分策略

#### 7.1.1 按业务领域拆分
- **用户域**：用户管理、认证授权
- **课程域**：课程业务逻辑
- **内容域**：文件和媒体管理

#### 7.1.2 权衡考虑
**优势**：
- 独立开发和部署
- 技术栈灵活性
- 水平扩展能力

**挑战**：
- 分布式事务复杂性
- 服务间通信开销
- 运维复杂度增加

### 7.2 技术选型决策

#### 7.2.1 Go语言选择
**原因**：
- 高并发性能优秀
- 微服务生态成熟
- 部署简单（单一可执行文件）
- 强类型安全

#### 7.2.2 gRPC vs HTTP
**选择gRPC的原因**：
- 性能优于HTTP/JSON
- 强类型接口定义
- 支持流式处理
- 多语言客户端生成

#### 7.2.3 MySQL vs NoSQL
**选择MySQL的原因**：
- 事务ACID保证
- 复杂查询支持
- 运维成熟度高
- 教育场景关系型需求

### 7.3 安全设计决策

#### 7.3.1 JWT vs Session
**选择JWT**：
- 无状态设计
- 微服务友好
- 水平扩展容易

#### 7.3.2 文件存储安全
- 文件类型验证
- 大小限制控制
- 路径遍历防护
- 访问权限控制

---

## 8. 安装配置与部署

### 8.1 环境要求

#### 8.1.1 系统要求
- **操作系统**：Windows 10/11, macOS 10.15+, Linux
- **Go版本**：Go 1.24.0 或更高
- **数据库**：MySQL 8.0+
- **缓存**：Redis 6.0+

#### 8.1.2 开发工具
- **编译器**：protoc (Protocol Buffers)
- **代码编辑器**：VS Code, GoLand
- **API测试**：Postman, curl

### 8.2 安装步骤

#### 8.2.1 克隆项目
```bash
git clone <repository-url>
cd course-platform
```

#### 8.2.2 安装依赖
```bash
go mod download
```

#### 8.2.3 配置数据库
```sql
CREATE DATABASE course_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

#### 8.2.4 配置文件
编辑 `configs/config.yaml`：
```yaml
server:
  port: ":8083"
mysql:
  user: "root"
  password: "your_password"
  host: "127.0.0.1"
  port: "3306"
  db_name: "course_db"
redis:
  addr: "127.0.0.1:6379"
  password: ""
  db: 0
```

### 8.3 构建与运行

#### 8.3.1 编译服务
```bash
# 编译所有服务
go build -o bin/user-service cmd/user-service/main.go
go build -o bin/course-service cmd/course-service/main.go
go build -o bin/content-service cmd/content-service/main.go
go build -o bin/server cmd/server/main.go
```

#### 8.3.2 启动服务
```bash
# Windows
.\start-all.bat

# Linux/macOS
chmod +x start-all.sh
./start-all.sh
```

### 8.4 Docker 部署

#### 8.4.1 Docker Compose
```yaml
version: '3.8'
services:
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: course_db
    ports:
      - "3306:3306"
  
  redis:
    image: redis:6-alpine
    ports:
      - "6379:6379"
  
  user-service:
    build:
      context: .
      dockerfile: docker/Dockerfile.user-service
    ports:
      - "50051:50051"
    depends_on:
      - mysql
      - redis
  
  course-service:
    build:
      context: .
      dockerfile: docker/Dockerfile.course-service
    ports:
      - "50052:50052"
    depends_on:
      - mysql
      - redis
  
  content-service:
    build:
      context: .
      dockerfile: docker/Dockerfile.content-service
    ports:
      - "50053:50053"
    depends_on:
      - mysql
      - redis
  
  api-gateway:
    build:
      context: .
      dockerfile: docker/Dockerfile.server
    ports:
      - "8083:8083"
    depends_on:
      - user-service
      - course-service
      - content-service
```

#### 8.4.2 部署命令
```bash
docker-compose up -d
```

### 8.5 验证部署

#### 8.5.1 健康检查
```bash
curl http://localhost:8083/ping
```

#### 8.5.2 API测试
```bash
# 运行API测试脚本
powershell -ExecutionPolicy Bypass -File api-test.ps1
```

---

## 9. 测试策略

### 9.1 测试架构

#### 9.1.1 测试金字塔
```
    ┌─────────────────┐
    │   E2E Tests     │  ← 少量端到端测试
    ├─────────────────┤
    │Integration Tests│  ← 中等数量集成测试
    ├─────────────────┤
    │   Unit Tests    │  ← 大量单元测试
    └─────────────────┘
```

#### 9.1.2 测试范围
- **单元测试**：业务逻辑、工具函数
- **集成测试**：数据库操作、gRPC通信
- **端到端测试**：完整用户流程

### 9.2 测试工具

#### 9.2.1 Go测试框架
```go
// 使用标准库testing
func TestCreateUser(t *testing.T) {
    // 测试用户创建逻辑
}

// 使用testify库
func TestCreateUserWithTestify(t *testing.T) {
    assert := assert.New(t)
    user, err := userService.CreateUser("test@example.com", "password", "Test User")
    assert.NoError(err)
    assert.NotNil(user)
    assert.Equal("test@example.com", user.Email)
}
```

#### 9.2.2 数据库测试
```go
func setupTestDB() *gorm.DB {
    db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    db.AutoMigrate(&model.User{}, &model.Course{})
    return db
}
```

### 9.3 API测试

#### 9.3.1 自动化测试脚本
项目包含 `api-test.ps1` 脚本，测试所有API端点：
```powershell
# 测试用户注册
$registerResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/register" -Method POST -Body $registerData -ContentType "application/json"

# 测试课程创建
$courseResponse = Invoke-RestMethod -Uri "$baseUrl/api/v1/courses" -Method POST -Headers $headers -Body $courseData -ContentType "application/json"
```

#### 9.3.2 性能测试
```bash
# 使用ab进行压力测试
ab -n 1000 -c 10 http://localhost:8083/api/v1/courses

# 使用wrk进行负载测试
wrk -t12 -c400 -d30s http://localhost:8083/api/v1/courses
```

### 9.4 运行测试

#### 9.4.1 单元测试
```bash
# 运行所有测试
go test ./...

# 运行特定包测试
go test ./internal/domain/user/service

# 生成覆盖率报告
go test -cover ./...
```

#### 9.4.2 集成测试
```bash
# 启动测试数据库
docker run -d --name test-mysql -e MYSQL_ROOT_PASSWORD=test -p 3307:3306 mysql:8.0

# 运行集成测试
go test -tags=integration ./...
```

---

## 10. 未来工作与路线图

### 10.1 短期目标 (1-3个月)

#### 10.1.1 功能完善
- [ ] 课程评论和评分系统
- [ ] 用户学习进度跟踪
- [ ] 课程分类管理后台
- [ ] 邮件通知系统

#### 10.1.2 性能优化
- [ ] 数据库查询优化
- [ ] Redis缓存策略完善
- [ ] 图片压缩和CDN集成
- [ ] API响应时间监控

#### 10.1.3 安全增强
- [ ] API限流和防护
- [ ] 文件上传安全扫描
- [ ] 敏感信息加密存储
- [ ] 审计日志系统

### 10.2 中期目标 (3-6个月)

#### 10.2.1 新功能开发
- [ ] 视频播放和进度记录
- [ ] 在线考试和测验系统
- [ ] 课程购买和支付集成
- [ ] 社区讨论功能

#### 10.2.2 架构升级
- [ ] 服务网格(Service Mesh)集成
- [ ] 事件驱动架构
- [ ] 分布式事务处理
- [ ] 微服务监控和追踪

#### 10.2.3 运维改进
- [ ] CI/CD流水线
- [ ] 自动化部署
- [ ] 容器编排(Kubernetes)
- [ ] 日志聚合和分析

### 10.3 长期愿景 (6-12个月)

#### 10.3.1 平台扩展
- [ ] 移动端APP开发
- [ ] 多语言国际化
- [ ] 多租户支持
- [ ] 白标解决方案

#### 10.3.2 智能化功能
- [ ] AI推荐系统
- [ ] 智能内容审核
- [ ] 学习路径规划
- [ ] 个性化学习建议

#### 10.3.3 生态建设
- [ ] 开放API平台
- [ ] 第三方插件系统
- [ ] 开发者社区
- [ ] 合作伙伴集成

### 10.4 技术债务

#### 10.4.1 代码重构
- [ ] 统一错误处理机制
- [ ] 配置管理优化
- [ ] 日志格式标准化
- [ ] 代码注释完善

#### 10.4.2 测试覆盖
- [ ] 单元测试覆盖率提升至80%+
- [ ] 集成测试自动化
- [ ] 性能测试基准建立
- [ ] 安全测试流程

---

## 附录

### A. 术语表

| 术语 | 定义 |
|------|------|
| gRPC | Google Remote Procedure Call，高性能RPC框架 |
| JWT | JSON Web Token，用于身份验证的令牌格式 |
| GORM | Go语言的ORM库 |
| Gin | Go语言的Web框架 |
| Redis | 内存数据结构存储，用作缓存 |
| protobuf | Protocol Buffers，Google的数据序列化格式 |

### B. 重要依赖库

| 库名 | 版本 | 用途 |
|------|------|------|
| gin-gonic/gin | v1.10.1 | Web框架 |
| google.golang.org/grpc | v1.73.0 | RPC框架 |
| gorm.io/gorm | v1.25.12 | ORM |
| golang-jwt/jwt/v5 | v5.2.2 | JWT处理 |
| go-redis/redis/v8 | v8.11.5 | Redis客户端 |
| spf13/viper | v1.20.1 | 配置管理 |

### C. 配置参考

#### 生产环境配置示例
```yaml
server:
  port: ":8083"
  mode: "release"
mysql:
  user: "course_user"
  password: "${MYSQL_PASSWORD}"
  host: "mysql.example.com"
  port: "3306"
  db_name: "course_production"
  max_idle_conns: 10
  max_open_conns: 100
redis:
  addr: "redis.example.com:6379"
  password: "${REDIS_PASSWORD}"
  db: 0
  pool_size: 50
```

### D. 监控指标

#### 关键性能指标(KPI)
- API响应时间: < 200ms (P95)
- 系统可用性: > 99.9%
- 数据库连接池使用率: < 80%
- 内存使用率: < 70%
- CPU使用率: < 60%

---

*文档版本: v1.0*  
*最后更新: 2025年1月3日*  
*维护者: Course Platform Team* 