# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Overview

이 저장소는 Go 언어 학습을 위한 튜토리얼 및 예제 코드 모음집입니다. 각 디렉터리는 독립적인 예제로 구성되어 있으며, 실무에서 자주 사용되는 패턴과 라이브러리 활용법을 다룹니다.

## Common Commands

### Testing
```bash
# Run all tests in the repository
go test ./...

# Run tests in a specific directory
go test ./golang/testing/...

# Run a single test file
go test ./golang/testing/table_test.go

# Run a specific test function
go test -run TestFunctionName ./path/to/package

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Clean up dependencies
go mod tidy

# Verify dependencies
go mod verify
```

### Building
```bash
# Build specific example
cd <directory>
go build

# Run without building
go run main.go

# Build with ldflags (see golang/build-ldflags)
go build -ldflags "-X main.version=1.0.0"
```

### Database Operations
```bash
# PostgreSQL setup (database/postgresql)
psql -U postgres -f database/postgresql/psql.sql

# Start Docker containers for databases
docker-compose up -d  # in database/redis or other directories with docker-compose.yml
```

## Architecture Patterns

### Clean Architecture Example
`project-layout/go-clean-arch/` demonstrates a complete clean architecture implementation:

- **Layer Structure**: domain → repository → usecase → http
- **Dependency Injection**: Uses `go.uber.org/fx` for DI container
- **Domain Layer**: Core business entities and interfaces (domain/)
- **Repository Layer**: Data access implementations (*/repository/mysql/)
- **Use Case Layer**: Business logic (*/usecase/)
- **Delivery Layer**: HTTP handlers (*/http/)
- **Database Setup**: Configuration via Viper, connection pooling
- **Middleware**: CORS and custom middleware in http/middleware/

Key files:
- `main.go`: Application bootstrap with fx.New()
- `domain/*.go`: Entity definitions and repository interfaces
- `common/config/`: Viper-based configuration management
- `common/database/`: Database connection setup

### Testing Patterns

#### Unit Testing with Mocks
`go-unit-test/mockery/` shows mockery-based testing:
- Generate mocks: `mockery --name=InterfaceName --output=mocks/`
- Mock interfaces are in `mocks/` directories
- Tests use `github.com/stretchr/testify/mock` for assertions

#### Integration Testing
`go-unit-test/testcontainers/` uses testcontainers for integration tests:
- Spins up real databases (Redis, MongoDB) in Docker
- Tests run against actual database instances
- Clean setup/teardown in test functions

#### HTTP Mocking
`go-unit-test/httpmock/` demonstrates HTTP request mocking:
- Uses `github.com/jarcoal/httpmock` for stubbing HTTP responses
- Good for testing external API integrations

### Concurrency Patterns

**Mutex vs Distributed Locks** (`golang/concurrency/waitgroup/`):
- `counter_mutex.go`: Local mutex for single-process synchronization
- `counter_redislock.go`: Redis-based distributed lock
- `counter_redsync.go`: Redsync for distributed lock across multiple instances
- `counter_mongolock.go`: MongoDB-based distributed lock
- Choose based on deployment: single-process → mutex, distributed → Redis/Mongo locks

**Context Usage** (`golang/context/`):
- Timeout handling with `context.WithTimeout()`
- Cancellation propagation across goroutines
- API request cancellation patterns

### Design Patterns

- **Builder Pattern** (`golang/design-pattern/builder/`): Fluent interface for complex object construction
- **Decorator Pattern** (`golang/design-pattern/decorator/`): Dynamic behavior addition
- **Template Pattern** (`golang/design-pattern/template/`): Algorithm skeleton with customizable steps
- **Functional Options** (`golang/design-pattern/func_opts/`): Flexible configuration pattern

## Key Dependencies and Usage

### Web Frameworks
- **Echo v4**: Primary web framework (see keycloak/backend, project-layout/go-clean-arch)
  - Middleware setup in `NewEcho()` functions
  - Route grouping and versioning patterns

### Authentication
- **JWT**: JWT token handling in `jwt/` directory
  - `jwt/pem/`: PEM-based key validation
  - `jwt/jwk/`: JWK Set validation with `github.com/MicahParks/jwkset`
- **Keycloak**: Full OAuth 2.0 implementation in `keycloak/`
  - Backend: JWT token verification via JWKS endpoint
  - Frontend: Authorization Code Flow with REST API (no keycloak-js library)
  - Setup instructions in keycloak/README.md

### Database Libraries
- **GORM**: ORM for MySQL (gorm.io/gorm, gorm.io/driver/mysql)
- **MongoDB Driver**: go.mongodb.org/mongo-driver
- **Redis**: github.com/go-redis/redis/v8
- **SQL Mocking**: gopkg.in/DATA-DOG/go-sqlmock.v1

### Testing Libraries
- **Testify**: github.com/stretchr/testify (assertions, mocks, suites)
- **Testcontainers**: github.com/testcontainers/testcontainers-go
- **HTTPMock**: github.com/jarcoal/httpmock
- **Mockery**: Code generation for mocks

### Utilities
- **Viper**: Configuration management (github.com/spf13/viper)
- **Logrus/Zap**: Logging (github.com/sirupsen/logrus, go.uber.org/zap)
- **Cron**: Scheduling (github.com/robfig/cron/v3, github.com/hibiken/asynq)
- **Lo**: Functional utilities (github.com/samber/lo)

## Directory Structure Conventions

### Test Organization
- Test files: `*_test.go` in the same directory as source files
- Integration tests: May use separate `_test` package suffix
- Mocks: Generated in `mocks/` subdirectories

### Database Examples
Each database directory (`database/{mysql,postgresql,redis,mongo}`) contains:
- Connection setup examples
- Query patterns
- Testing strategies
- Docker setup (docker-compose.yml where applicable)

### Standalone Examples
Most directories under `golang/` are self-contained:
- Single topic focus (e.g., generics, context, reflect)
- Runnable test files demonstrating concepts
- README.md with Korean explanations where needed

## Development Workflow

### Working on Keycloak Example
```bash
# 1. Start Keycloak
docker run -d -p 8080:8080 \
  -e KEYCLOAK_ADMIN=admin \
  -e KEYCLOAK_ADMIN_PASSWORD=admin \
  --name keycloak-tutorial \
  quay.io/keycloak/keycloak:latest start-dev

# 2. Configure Keycloak (see keycloak/README.md)

# 3. Start Backend
cd keycloak/backend
go mod tidy
go run main.go  # Runs on :8081

# 4. Start Frontend
cd keycloak/frontend
npm install
npm start  # Runs on :3000
```

### Working on Clean Architecture Example
```bash
cd project-layout/go-clean-arch

# Setup database (MySQL required)
# Update config file with your MySQL credentials

# Run application
go run main.go

# The application uses fx for dependency injection
# Server starts on address specified in config
```

## Testing Considerations

### Database Tests
- Tests using testcontainers require Docker to be running
- Some tests may require specific database setup (see individual README.md files)
- Clean up containers after tests: `docker ps -a | grep testcontainers`

### Mock Generation
When interfaces change, regenerate mocks:
```bash
cd go-unit-test/mockery
mockery --name=Doer --dir=do_user/doer --output=do_user/mocks/doer
```

### Integration vs Unit Tests
- Unit tests: Fast, use mocks, no external dependencies
- Integration tests: Slower, use testcontainers or real services
- Tag integration tests if needed: `//go:build integration`

## Code Style Notes

### Import Aliases
The codebase uses underscore prefixes for domain-specific imports:
```go
_articleHttp "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/article/http"
_articleRepo "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch/article/repository/mysql"
```
This prevents import conflicts and improves readability in DI-heavy code.

### Configuration Management
Viper is the standard for configuration:
- Config files location varies by example
- Environment variable support with `viper.GetString()`, `viper.GetInt()`
- See `project-layout/go-clean-arch/common/config/` for patterns

### Error Handling
- Custom error types in `golang/errors/custom/`
- Error wrapping and unwrapping examples
- Structured error codes for API responses

## Language Version

- **Go Version**: 1.21.3 minimum (specified in go.mod)
- **Toolchain**: go1.22.3
- Some examples may use features from Go 1.18+ (generics, workspaces)
