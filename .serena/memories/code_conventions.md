# Code Conventions and Style Guide

## Go Conventions

### Error Handling
```go
// Use proper error handling with wrapping
if err != nil {
    return nil, fmt.Errorf("operation failed: %w", err)
}

// Always check errors
result, err := someFunction()
if err != nil {
    // Handle error appropriately
}
```

### Testing Patterns
```go
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

### GraphQL Schema
- Use code-first approach with gqlgen
- Implement proper input validation
- Include pagination for list queries
- Add audit logging for all mutations

### Package Structure
```
backend/
├── cmd/
│   ├── server/         # Main server application
│   └── migrate/        # Database migrations
├── internal/
│   ├── domain/         # Business logic and entities
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

### Naming Conventions
- **Packages**: lowercase, single words when possible
- **Functions**: PascalCase for exported, camelCase for private
- **Variables**: camelCase
- **Constants**: PascalCase or SCREAMING_SNAKE_CASE
- **Files**: lowercase_with_underscores.go
- **Interfaces**: single method interfaces named with -er suffix (Reader, Writer)

### Authentication Patterns
- JWT middleware for protected routes
- Role-based access control (ADMIN, MANAGER, VIEWER)
- Token refresh mechanism
- Proper error handling for authentication failures

### TDD Principles
- Tests MUST fail before implementation (RED phase)
- Use real dependencies (actual PostgreSQL, not mocks)
- Follow contract → integration → unit test order
- Each test should be independent and isolated

### Configuration Management
- Use Viper for configuration management
- Support environment variables and config files
- Provide sensible defaults
- Include configuration validation

### Logging
- Use Zap for structured logging
- Include request correlation IDs
- Log at appropriate levels (debug, info, warn, error)
- Include relevant context in log messages

### Security
- Always use parameterized queries to prevent SQL injection
- Validate all input data
- Use HTTPS in production
- Implement proper CORS policies
- Rate limiting for API endpoints

### Performance
- Use database connection pooling
- Implement proper indexing strategy
- Cache frequently accessed data when appropriate
- Monitor and optimize query performance