# Project Structure

## Root Directory Structure
```
vcd-claude-speckit/
├── README.md                    # Project overview
├── CLAUDE.md                   # Claude Code specific context
├── .mcp.json                   # MCP server configuration
├── .serena/                    # Serena project configuration
├── .specify/                   # Speckit framework configuration
│   ├── memory/                 # Framework memory files
│   ├── scripts/                # Development scripts
│   │   └── bash/              # Bash utility scripts
│   └── templates/             # Project templates
└── specs/                      # Feature specifications
    └── 001-create-new-golang/   # Current feature branch
        ├── spec.md             # Feature specification
        ├── plan.md             # Implementation plan
        ├── data-model.md       # Data model design
        ├── quickstart.md       # Quick start guide
        ├── tasks.md            # Implementation tasks
        ├── research.md         # Research findings
        └── contracts/          # API contracts
            ├── schema.graphql  # GraphQL schema
            └── employee-contract-test.go  # Contract tests
```

## Expected Implementation Structure
```
backend/                          # Main application directory
├── cmd/                         # Command line applications
│   ├── server/                 # Main server
│   │   └── main.go             # Server entry point
│   └── migrate/                # Database migrations
│       └── main.go             # Migration tool
├── internal/                   # Internal application code
│   ├── domain/                 # Domain layer (business logic)
│   │   ├── employee.go         # Employee entity and business rules
│   │   ├── user.go             # User entity and authentication
│   │   └── audit.go            # Audit logging domain
│   ├── service/                # Application services
│   │   ├── employee_service.go # Employee management service
│   │   ├── user_service.go     # User management service
│   │   └── auth_service.go     # Authentication service
│   ├── repository/             # Data access layer
│   │   ├── employee_repository.go
│   │   ├── user_repository.go
│   │   └── audit_repository.go
│   ├── api/                    # API layer
│   │   ├── graphql/            # GraphQL resolvers
│   │   │   ├── resolver.go     # Root resolver
│   │   │   ├── employee_resolver.go
│   │   │   └── user_resolver.go
│   │   └── middleware/         # HTTP middleware
│   │       ├── auth.go         # Authentication middleware
│   │       ├── logging.go      # Logging middleware
│   │       └── cors.go         # CORS middleware
│   └── infrastructure/         # Infrastructure concerns
│       ├── database/           # Database setup and migrations
│       ├── config/             # Configuration management
│       └── logger/             # Logging setup
├── pkg/                        # Public packages
│   ├── database/               # Database utilities
│   ├── validator/              # Validation helpers
│   └── auth/                   # Authentication utilities
├── tests/                      # Test suites
│   ├── contract/               # Contract tests (API validation)
│   │   ├── employee_test.go
│   │   └── user_test.go
│   ├── integration/            # Integration tests
│   │   ├── employee_service_test.go
│   │   └── end_to_end_test.go
│   └── unit/                   # Unit tests
│       ├── employee_test.go
│       └── user_test.go
├── migrations/                 # Database migration files
│   ├── 001_initial_schema.sql
│   └── 002_add_audit_logs.sql
├── configs/                    # Configuration files
│   ├── config.yaml             # Default configuration
│   └── config.dev.yaml         # Development override
├── graphql/                    # GraphQL generated code
│   ├── generated.go            # gqlgen generated code
│   └── models_gen.go           # Generated models
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
└── tools.go                    # Build tools
```

## Development Workflow Structure
```
.specify/                        # Speckit framework
├── scripts/bash/               # Development scripts
│   ├── create-new-feature.sh    # Create new feature branch
│   ├── setup-plan.sh           # Setup implementation plan
│   ├── check-task-prerequisites.sh  # Validate task readiness
│   ├── update-agent-context.sh    # Update AI assistant context
│   ├── get-feature-paths.sh    # Get current feature paths
│   └── common.sh              # Common utilities
└── templates/                  # Project templates
    ├── spec-template.md       # Feature spec template
    ├── plan-template.md       # Implementation plan template
    ├── tasks-template.md      # Tasks template
    └── agent-file-template.md # AI assistant template
```

## Feature Branch Structure
Each feature branch follows this structure:
```
specs/[###-feature-name]/
├── spec.md             # Feature specification
├── plan.md             # Implementation plan
├── data-model.md       # Data model design
├── quickstart.md       # Quick start guide
├── tasks.md            # Implementation tasks
├── research.md         # Research findings
└── contracts/          # API contracts
    ├── schema.graphql  # GraphQL schema
    └── *.go            # Contract test files
```