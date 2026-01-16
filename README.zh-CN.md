# Go gRPC REST 演示项目

一个 Go 后端模板，演示用户和产品服务的 CRUD API，支持三种访问方式：

1. **REST API** - 基于 Gin 框架，带 Swagger 文档
2. **gRPC API** - 原生 gRPC，使用 Protocol Buffers
3. **CLI 客户端** - 基于 Cobra 的命令行工具

## 特性

- **用户服务**：增删改查、列表查询，支持过滤和排序
- **产品服务**：增删改查、多条件搜索
- **双协议支持**：REST (HTTP/JSON) 和 gRPC
- **Swagger 文档**：自动生成 API 文档
- **优雅关闭**：正确处理系统信号
- **线程安全**：并发安全的内存存储

## 快速开始

### 前置要求

- Go 1.25+
- （可选）`protoc` 用于重新生成 protobuf 代码
- （可选）`jq` 用于格式化 JSON 输出

### 运行服务器

```bash
make run-server
```

服务端点：

- REST API：<http://localhost:8080/api/v1/>
- gRPC：localhost:9090
- Swagger UI：<http://localhost:8080/swagger/index.html>

## API 端点

### REST API (`/api/v1`)

| 方法   | 端点               | 描述                               |
|--------|--------------------|------------------------------------|
| GET    | `/health`          | 健康检查                           |
| POST   | `/users`           | 创建用户                           |
| GET    | `/users`           | 用户列表（支持分页、过滤、排序）   |
| GET    | `/users/:id`       | 获取用户                           |
| PUT    | `/users/:id`       | 更新用户                           |
| DELETE | `/users/:id`       | 删除用户                           |
| POST   | `/products`        | 创建产品                           |
| GET    | `/products/:id`    | 获取产品                           |
| GET    | `/products/search` | 搜索产品（关键词、类别、价格范围） |

### gRPC 服务 (端口 9090)

| 服务           | 方法                                                   |
|----------------|--------------------------------------------------------|
| UserService    | CreateUser, GetUser, UpdateUser, DeleteUser, ListUsers |
| ProductService | CreateProduct, GetProduct, SearchProducts              |

### CLI 命令

```bash
go run cmd/client/main.go user create <用户名> <邮箱> <全名>
go run cmd/client/main.go user get <id>
go run cmd/client/main.go user list [--filter] [--sort-by]
go run cmd/client/main.go product create <名称> <描述> <价格> <数量> <类别>
go run cmd/client/main.go product get <id>
go run cmd/client/main.go product search [--query] [--category] [--min-price] [--max-price]
```

## Make 命令

```bash
make help           # 显示所有可用命令
make deps           # 安装依赖
make build          # 构建服务器和客户端
make run-server     # 运行服务器
make run-dev        # 热重载运行（需要 air）
make test-unit      # 运行单元测试
make test-coverage  # 生成覆盖率报告
make lint           # 运行 golangci-lint
make swagger        # 生成 Swagger 文档
make proto          # 生成 protobuf 代码
make endpoints      # 显示所有 API 端点
```

## 项目结构

```text
go-grpc-rest-demo/
├── api/
│   ├── grpc/server/    # gRPC 服务实现
│   └── rest/           # REST 处理器和路由
├── cmd/
│   ├── client/         # CLI 客户端 (Cobra)
│   └── server/         # 服务器入口
├── internal/
│   ├── client/         # 客户端实现 (gRPC/REST)
│   ├── errors/         # 标准化错误处理
│   ├── model/          # 数据模型和 DTO
│   ├── response/       # API 响应辅助
│   └── service/        # 业务逻辑层
├── proto/              # Protocol Buffer 定义
├── docs/               # Swagger 文档
├── Makefile            # 构建命令
└── go.mod              # Go 模块
```

## 许可证

此项目仅供演示目的。
