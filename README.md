# Go gRPC REST Demo

A Go backend template demonstrating User and Product services with CRUD APIs, accessible via:

1. **REST API** - Gin framework with Swagger documentation
2. **gRPC API** - Native gRPC with Protocol Buffers
3. **CLI Client** - Cobra-based command line interface

## Features

- **UserService**: Create, Read, Update, Delete, List users with filtering and sorting
- **ProductService**: Create, Read, Search products with multi-condition filtering
- **Dual Protocol**: REST (HTTP/JSON) and gRPC support
- **Swagger Documentation**: Auto-generated API docs
- **Graceful Shutdown**: Proper signal handling
- **Thread-Safe**: Concurrent-safe in-memory storage

## Quick Start

### Prerequisites

- Go 1.25+
- (Optional) `protoc` for regenerating protobuf code
- (Optional) `jq` for formatted JSON output

### Run the Server

```bash
make run-server
```

Server endpoints:

- REST API: <http://localhost:8080/api/v1/>
- gRPC: localhost:9090
- Swagger UI: <http://localhost:8080/swagger/index.html>

## API Endpoints

### REST API (`/api/v1`)

| Method | Endpoint           | Description                                    |
|--------|--------------------|------------------------------------------------|
| GET    | `/health`          | Health check                                   |
| POST   | `/users`           | Create user                                    |
| GET    | `/users`           | List users (with pagination, filter, sort)     |
| GET    | `/users/:id`       | Get user by ID                                 |
| PUT    | `/users/:id`       | Update user                                    |
| DELETE | `/users/:id`       | Delete user                                    |
| POST   | `/products`        | Create product                                 |
| GET    | `/products/:id`    | Get product by ID                              |
| GET    | `/products/search` | Search products (query, category, price range) |

### gRPC Services (port 9090)

| Service        | Methods                                                |
|----------------|--------------------------------------------------------|
| UserService    | CreateUser, GetUser, UpdateUser, DeleteUser, ListUsers |
| ProductService | CreateProduct, GetProduct, SearchProducts              |

### CLI Commands

```bash
go run cmd/client/main.go user create <username> <email> <full_name>
go run cmd/client/main.go user get <id>
go run cmd/client/main.go user list [--filter] [--sort-by]
go run cmd/client/main.go product create <name> <desc> <price> <qty> <category>
go run cmd/client/main.go product get <id>
go run cmd/client/main.go product search [--query] [--category] [--min-price] [--max-price]
```

## Make Commands

```bash
make help           # Show all available commands
make deps           # Install dependencies
make build          # Build server and client binaries
make run-server     # Run the server
make run-dev        # Run with hot reload (requires air)
make test-unit      # Run unit tests
make test-coverage  # Generate coverage report
make lint           # Run golangci-lint
make swagger        # Generate Swagger docs
make proto          # Generate protobuf code
make endpoints      # Show all API endpoints
```

## Project Structure

```text
go-grpc-rest-demo/
├── api/
│   ├── grpc/server/    # gRPC service implementations
│   └── rest/           # REST handlers and router
├── cmd/
│   ├── client/         # CLI client (Cobra)
│   └── server/         # Server entry point
├── internal/
│   ├── client/         # Client implementations (gRPC/REST)
│   ├── errors/         # Standardized error handling
│   ├── model/          # Data models and DTOs
│   ├── response/       # API response helpers
│   └── service/        # Business logic layer
├── proto/              # Protocol Buffer definitions
├── docs/               # Swagger documentation
├── Makefile            # Build commands
└── go.mod              # Go module
```

## License

This project is for demonstration purposes.
