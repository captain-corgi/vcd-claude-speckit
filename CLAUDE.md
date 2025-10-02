# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Building an employee management web server with GraphQL API using Go, following Clean Architecture and Domain Driven Design patterns with Test Driven Development.

## Architecture

- **Language**: Go 1.24.3
- **HTTP Framework**: Gin
- **API**: GraphQL (implementation in progress)
- **Database**: PostgreSQL with GORM
- **Testing**: Table-driven testing with testcontainers
- **Authentication**: JWT with role-based access control
- **Infrastructure**: Docker Compose with PostgreSQL and Redis
- **Dependencies**: Managed with Go modules

## Recent Changes (Last 3)

1. **Feature 001**: Employee Management System - Initial implementation planning
   - Created GraphQL schema for Employee, User, and AuditLog entities
   - Defined database schema with proper constraints and triggers
   - Established contract tests for all CRUD operations
   - Set up authentication and authorization patterns

2. **Architecture**: Clean Architecture with Domain Driven Design
   - Separated domain logic from infrastructure concerns
   - Implemented service layer pattern for business operations
   - Created proper package structure for maintainability

3. **Testing**: Contract-first approach with TDD
   - Contract tests define expected API behavior
   - Integration tests use real PostgreSQL via testcontainers
   - Table-driven unit tests following Go best practices

## Code Patterns to Follow

### Go Conventions

```go
// Use proper error handling
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}

// Table-driven tests
tests := []struct {
    name    string
    input   string
    want    string
    wantErr bool
}{
    {"valid input", "test", "result", false},
    {"invalid input", "", "", true},
}

for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        got, err := SomeFunction(tt.input)
        if (err != nil) != tt.wantErr {
            t.Errorf("SomeFunction() error = %v, wantErr %v", err, tt.wantErr)
            return
        }
        if got != tt.want {
            t.Errorf("SomeFunction() = %v, want %v", got, tt.want)
        }
    })
}
```

### GraphQL Schema

- Use code-first approach with gqlgen
- Implement proper input validation
- Include pagination for list queries
- Add audit logging for all mutations

### Database Patterns

```go
// Use GORM with proper model definitions
type Employee struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    FirstName string    `gorm:"size:50;not null"`
    // ... other fields
    CreatedAt time.Time
    UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// Use context for database operations
func (r *employeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
    return r.db.WithContext(ctx).Create(employee).Error
}
```

### Authentication

- JWT middleware for protected routes
- Role-based access control (ADMIN, MANAGER, VIEWER)
- Token refresh mechanism
- Proper error handling for authentication failures

### Testing Principles

- RED-GREEN-REFACTOR cycle strictly enforced
- Tests MUST fail before implementation
- Use real dependencies (actual PostgreSQL, not mocks)
- Contract tests first, then integration, then unit tests

## Architecture Overview

### Clean Architecture Implementation

The project follows Domain-Driven Design with Clean Architecture principles:

```
internal/
├── domain/            # Core business logic (entities, services, events)
│   ├── employee.go    # Employee entity with validation
│   ├── user.go        # User entity and roles
│   ├── *.go           # Other domain entities
│   └── *_test.go      # Domain unit tests
├── service/           # Application services orchestrate domain logic
│   ├── employee_service.go    # Employee business operations
│   ├── user_service.go        # User management
│   └── auth_service.go        # Authentication
├── repository/        # Data access layer (GORM models + CRUD)
│   ├── *_model.go     # GORM database models
│   └── *_repository.go # Repository implementations
├── api/              # GraphQL API layer (not yet implemented)
├── middleware/       # HTTP middleware (auth, logging, CORS)
└── config/           # Application configuration

pkg/
└── database/         # Database connection and setup

tests/
├── helpers/          # Test utilities and fixtures
├── contract/         # API contract tests (RED phase tests)
├── integration/      # Full stack tests with real DB
└── unit/            # Repository and service unit tests

cmd/
├── server/          # Main application entry point (not yet implemented)
└── migrate/         # Database migration tool
```

### Key Architecture Patterns

- **Domain-First**: Business logic lives in `internal/domain/` with no external dependencies
- **Repository Pattern**: Data access abstracted through repository interfaces
- **Service Layer**: Application services orchestrate domain operations
- **Event-Driven**: Domain events for audit logging and decoupling
- **Contract Testing**: RED-GREEN-REFACTOR TDD with contract tests defining API behavior
- **Testcontainers**: Real PostgreSQL for integration testing, no mocks

## Dependencies Management

- Use Go modules for dependency management
- Keep dependencies minimal and well-maintained
- Prefer standard library where possible
- Document non-standard dependencies

## Environment and Infrastructure

### Database Setup
- **PostgreSQL**: Primary data store (Docker Compose on port 5432)
- **Redis**: Caching and sessions (Docker Compose on port 6379)
- **Connection Pooling**: GORM with configurable pool settings
- **Migrations**: Version-controlled schema changes

### Configuration
- Environment-based configuration with `.env` file support
- Use `.env.example` as template for required variables
- Key configuration areas: Database, Redis, JWT, CORS, Rate limiting
- Docker Compose provides development database services

### Development Workflow
1. **Environment Setup**: `make dev-setup` (starts PostgreSQL + Redis + runs migrations)
2. **Test-Driven Development**: Write failing contract tests first
3. **Implementation**: Make tests pass following domain-first approach
4. **Quality Gates**: `make dev-build` ensures code quality before commits

## Current Implementation Status

### ✅ Implemented
- Domain entities with validation (Employee, User, Address, etc.)
- Repository layer with GORM models and basic CRUD operations
- Service layer with business logic orchestration
- Event system for audit logging
- Contract test suite defining expected API behavior
- Database infrastructure with Docker Compose
- Migration tooling

### 🚧 In Progress / Not Yet Implemented
- GraphQL API layer and resolvers
- HTTP server with Gin framework
- Authentication middleware with JWT
- Main application entry point (`cmd/server/`)
- HTTP middleware (CORS, logging, rate limiting)
- GraphQL schema generation with gqlgen

## Key Implementation Patterns

### Domain Entity Pattern
```go
// Entity with business logic and validation
func NewEmployee(firstName, lastName, email string) (*Employee, error) {
    employee := &Employee{...}
    return employee, employee.Validate()
}

func (e *Employee) Validate() error {
    // Business validation rules
}
```

### Repository Pattern
```go
// Interface in domain layer
type EmployeeRepository interface {
    Create(ctx context.Context, employee *Employee) error
    GetByID(ctx context.Context, id uuid.UUID) (*Employee, error)
}

// GORM implementation in repository layer
type EmployeeRepository struct {
    db *gorm.DB
}
```

### Service Layer Pattern
```go
// Orchestrates domain operations with audit logging
func (s *EmployeeService) CreateEmployee(ctx context.Context, ...) (*Employee, error) {
    // Business validation
    // Repository operations
    // Event publishing
    // Audit logging
}
```

### Test-Driven Development
- Contract tests define expected API behavior before implementation
- Integration tests use real PostgreSQL via testcontainers
- Unit tests validate domain logic in isolation
- All tests follow table-driven pattern

## Common Commands

### Development Setup
```bash
# Set up development environment
make dev-setup        # Start Docker services + run migrations

# Build and run
make build            # Build to bin/server
make run              # Run server directly
make tidy             # Tidy and vendor dependencies
```

### Testing
```bash
# Run all tests
make test             # Go test -v ./...
make test-all         # Run with coverage report (generates coverage.html)
make test-contract    # Contract tests only (./tests/contract/...)
make test-integration # Integration tests only (./tests/integration/...)
make test-unit        # Unit tests only (./tests/unit/...)
make dev-test         # Quick contract test check
```

### Database Operations
```bash
make migrate-up       # Run migrations
make migrate-down     # Rollback migrations
make migrate-create name=migration_name  # Create new migration
```

### Code Quality
```bash
make lint             # Run golangci-lint
make fmt              # Format code
make clean            # Clean build artifacts
make dev-build        # Format + lint + test + build
```

### Docker Infrastructure
```bash
make docker-up        # Start PostgreSQL and Redis
make docker-down      # Stop services
```

### Go Commands (Alternative to Makefile)
```bash
# Build
go build -o bin/server cmd/server/main.go

# Test
go test -v ./...
go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html

# Database
go run cmd/migrate/main.go up
go run cmd/migrate/main.go down
```

<!-- MANUAL ADDITIONS SECTION -->
<!-- Add any project-specific context below this line -->
<!-- PRESERVE MARKERS FOR UPDATES -->

<!-- END MANUAL ADDITIONS -->