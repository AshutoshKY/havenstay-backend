---
name: "Airbnb Backend Architecture"
description: "Context and architectural patterns for the domain-driven GoLang Airbnb backend."
---

# Airbnb Clone Backend Context

This project is a GoLang backend for a hotel and room booking service, structured strictly around **Domain-Driven Design (DDD)** concepts and utilizing a **Server, Core, Repo architecture**.

## Technology Stack
- **Language**: Go 1.21+
- **Database**: MySQL 8.0 (managed via GORM)
- **APIs**: Provides BOTH **REST** (via `gin-gonic`) and **gRPC** interfaces.
- **DTO Layer**: Powered entirely by Protobufs (`*.proto`). Proto compilers generate `*.pb.go` structures which act as Data Transfer Objects passing across network boundaries.
- **Containerization**: `docker-compose.yml` provides the local MySQL instance.

## Directory Structure
The architecture is divided into `cmd`, `config`, `internal/platform`, and `internal/domain`.

- `api/proto/v1/`: Contains `*.proto` files defining the gRPC schemas and services.
- `cmd/server/main.go`: Entry point. Injects dependencies, configures platform tools, and boots both HTTP and gRPC listeners.
- `internal/platform/`: Generic infrastructure logic like DB connections (`database/mysql.go`), config setup, and API routers (`server/grpc.go`, `server/http.go`).
- `internal/domain/`: Broken into bounded contexts: **Host**, **Property**, **Guest**, **Booking**, and **Wishlist**.

## Domain Architecture Pattern
Each folder inside `internal/domain/{bounded_context}` explicitly isolates the following responsibilities:

1. **`model.go`**: Defines the GORM structural maps representing raw SQL tables.
2. **`entity.go`**: Defines the pure, framework-agnostic business Entity used strictly within the `Core` layer.
3. **`dto.go`**: Defines explicit mappers:
   - `Proto <-> Entity`
   - `Model <-> Entity`
4. **`repo.go`**: Holds the internal abstract interface (e.g., `Repository`) and the specific infrastructure implementation (`MySQLRepository`). The Repo pattern acts *strictly* as a translation layer converting abstract `Entity` data to/from SQL `Model` data.
5. **`core.go`**: Holds the `Service` interface where all primary business logic resides. The core service depends on the `Repository` interface (Dependency Inversion). It accepts and returns gRPC Proto structures, utilizing DTOs to convert to Entities for logic processing.
6. **`server.go`**: Contains both `GRPCServer` and `HTTPServer` structs, acting as adapters. They validate external HTTP/gRPC requests and pass arguments directly into the `Service`.

## Extending the Codebase
When adding new endpoints:
1. Add the definition to the relevant `api/proto/v1/*.proto` file.
2. Run `make proto` from the root directory to re-generate the `*.pb.go` files.
3. Add the logic to the respective `core.go` interface and implementation using `Entity` conversions.
4. Implement the controller logic in the specific `server.go` file (handling the new pb/HTTP request).
