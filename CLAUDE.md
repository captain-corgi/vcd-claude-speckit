# Claude Code Context for Employee Management System

## Project Overview
Building an employee management web server with GraphQL API using Go, following Clean Architecture and Domain Driven Design patterns with Test Driven Development.

## Architecture
- **Language**: Go 1.21+
- **HTTP Framework**: Gin
- **API**: GraphQL with gqlgen
- **Database**: PostgreSQL with GORM
- **Configuration**: Viper
- **Testing**: Table-driven testing with testcontainers
- **Authentication**: JWT with role-based access control

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

## Package Structure
```
backend/
├── cmd/
│   ├── server/         # Main server application
│   └── migrate/        # Database migrations
├── internal/
│   ├── domain/         # Business logic and entities
│   │   ├── employee.go
│   │   └── user.go
│   ├── service/        # Application services
│   ├── repository/     # Data access layer
│   ├── api/           # GraphQL resolvers and handlers
│   └── middleware/    # Authentication, logging
├── pkg/
│   ├── database/      # Database setup
│   ├── logger/        # Logging configuration
│   └── config/        # Application config
└── tests/
    ├── contract/      # Contract tests
    ├── integration/   # Integration tests
    └── unit/          # Unit tests
```

## Dependencies Management
- Use Go modules for dependency management
- Keep dependencies minimal and well-maintained
- Prefer standard library where possible
- Document non-standard dependencies

## Configuration
- Use Viper for configuration management
- Support environment variables and config files
- Provide sensible defaults
- Include configuration validation

## Logging
- Use Zap for structured logging
- Include request correlation IDs
- Log at appropriate levels (debug, info, warn, error)
- Include relevant context in log messages

## Error Handling
- Use proper error wrapping with fmt.Errorf
- Provide clear error messages
- Include error codes for client consumption
- Log errors with appropriate context

## Security
- Always use parameterized queries to prevent SQL injection
- Validate all input data
- Use HTTPS in production
- Implement proper CORS policies
- Rate limiting for API endpoints

## Performance
- Use database connection pooling
- Implement proper indexing strategy
- Cache frequently accessed data when appropriate
- Monitor and optimize query performance

## Development Workflow
1. Write failing test (RED)
2. Make test pass (GREEN)
3. Refactor code
4. Run all tests
5. Commit changes with descriptive message
6. Repeat

## Testing Commands
```bash
# Run all tests
go test ./...

# Run specific test type
go test ./tests/contract/...
go test ./tests/integration/...
go test ./tests/unit/...

# Run with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## Build Commands
```bash
# Build the application
go build -o bin/server cmd/server/main.go

# Run with hot reload (during development)
go run cmd/server/main.go

# Run database migrations
go run cmd/migrate/main.go
```

<!-- MANUAL ADDITIONS SECTION -->
<!-- Add any project-specific context below this line -->
<!-- PRESERVE MARKERS FOR UPDATES -->

<!-- END MANUAL ADDITIONS -->