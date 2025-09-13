# Suggested Commands for Development

## Testing Commands
```bash
# Run all tests
go test ./...

# Run specific test types
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

## Development Workflow Commands
```bash
# Create new feature (uses Speckit framework)
./.specify/scripts/bash/create-new-feature.sh "feature description"

# Setup project plan
./.specify/scripts/bash/setup-plan.sh

# Check task prerequisites
./.specify/scripts/bash/check-task-prerequisites.sh

# Update agent context
./.specify/scripts/bash/update-agent-context.sh claude
```

## Git Commands
```bash
# Check current status
git status

# View changes
git diff

# View recent commits
git log --oneline -10

# Create commit with proper message
git add .
git commit -m "feat: descriptive message

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"
```

## Database Commands
```bash
# Start PostgreSQL with Docker (for local development)
docker run --name vcd-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=vcd -p 5432:5432 -d postgres:15

# Connect to database
psql -h localhost -p 5432 -U postgres -d vcd

# Run database migrations (when implemented)
go run cmd/migrate/main.go up
go run cmd/migrate/main.go down
```

## GraphQL Commands
```bash
# Generate GraphQL code (when gqlgen is set up)
go generate ./...

# Start GraphQL playground (when server is running)
# Open http://localhost:8080/graphql in browser
```

## System Commands (Darwin/macOS)
```bash
# Find files by pattern
find . -name "*.go" -type f

# Search in files
grep -r "pattern" . --include="*.go"

# Check file structure
tree . -I '__pycache__|node_modules|.git'

# Monitor file changes
fswatch -o . | xargs -n1 -I{} go test ./...
```

## Project Status Commands
```bash
# Check current feature branch
git branch --show-current

# View feature structure
ls -la specs/$(git branch --show-current)/

# Check implementation progress
./.specify/scripts/bash/get-feature-paths.sh
```

## Utility Commands
```bash
# Format Go code
go fmt ./...

# Lint Go code (when golangci-lint is installed)
golangci-lint run

# Tidy dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify
```