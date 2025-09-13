# Tasks: Employee Management Web Server

**Input**: Design documents from `/specs/001-create-new-golang/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md
**Tech Stack**: Go 1.21+, GraphQL, Gin, Viper, PostgreSQL, Clean Architecture, DDD, TDD

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → Extract: tech stack, libraries, structure (backend API)
2. Load design documents:
   → data-model.md: Extract Employee, User, AuditLog entities
   → contracts/: GraphQL schema + contract tests
   → research.md: PostgreSQL, gqlgen, JWT decisions
   → quickstart.md: Test scenarios and validation requirements
3. Generate tasks by category:
   → Setup: Go modules, PostgreSQL, Docker
   → Tests: Contract tests → Integration tests → Unit tests
   → Core: Domain → Services → API → Middleware
   → Integration: Database → Auth → Logging
   → Polish: Performance, docs, optimization
4. Apply TDD rules:
   → Tests MUST fail before implementation
   → Real dependencies (actual PostgreSQL)
   → Contract → Integration → Unit order
5. Generate 38 numbered tasks with dependency graph
6. Create parallel execution examples
7. Validate: All contracts tested, all entities modeled
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **TDD**: Tests written first, MUST fail, then implementation
- Paths follow plan.md structure: `cmd/`, `internal/`, `pkg/`, `src/`, `tests/`

## Phase 3.1: Infrastructure Setup
**Goal: Basic project structure and dependencies**

- [x] **T001** Initialize Go module with required dependencies
  `go mod init employee-management-system && go get github.com/gin-gonic/gin github.com/99designs/gqlgen github.com/spf13/viper github.com/golang-jwt/jwt/v5`

- [x] **T002** [P] Create project structure per plan.md
  `cmd/server/`, `cmd/migrate/`, `internal/domain/`, `internal/service/`, `internal/repository/`, `internal/api/`, `internal/middleware/`, `pkg/database/`, `pkg/logger/`, `pkg/config/`, `src/models/`, `src/services/`, `src/cli/`, `src/lib/`, `tests/`

- [x] **T003** [P] Configure Go development tools
  `.golangci.yml`, `Makefile` with test/lint/build targets, `.gitignore`

## Phase 3.2: Database Infrastructure ⚠️ MUST COMPLETE BEFORE 3.3
**Goal: PostgreSQL setup with migrations and connection management**

- [x] **T004** Create Docker Compose configuration for PostgreSQL
  `docker-compose.yml` with PostgreSQL service, volumes, environment variables

- [x] **T005** [P] Set up GORM with PostgreSQL connection pooling
  `pkg/database/database.go` with connection management, health checks

- [x] **T006** [P] Create database migration system
  `cmd/migrate/main.go` with golang-migrate, SQL files for employees/users/audit_logs tables

- [x] **T007** [P] Configure testcontainers for integration tests
  `tests/helpers/database.go` with test PostgreSQL setup/teardown

## Phase 3.3: Contract Tests ⚠️ MUST COMPLETE BEFORE 3.4 (RED PHASE)
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**

- [x] **T008** [P] GraphQL schema contract test in `tests/contract/test_graphql_schema.go`
  Validate schema.graphql matches gqlgen generated schema

- [x] **T009** [P] Employee CRUD contract tests in `tests/contract/test_employee_crud.go`
  createEmployee, updateEmployee, deleteEmployee, getEmployee mutations

- [x] **T010** [P] Employee listing contract test in `tests/contract/test_employee_list.go`
  employees query with pagination, filtering, sorting

- [x] **T011** [P] Authentication contract tests in `tests/contract/test_auth.go`
  login, refreshToken, logout mutations with valid/invalid credentials

- [x] **T012** [P] Authorization contract tests in `tests/contract/test_authorization.go`
  Role-based access control for protected operations

- [x] **T013** [P] Validation contract tests in `tests/contract/test_validation.go`
  Required fields, email uniqueness, invalid data handling

- [x] **T014** [P] Audit logging contract test in `tests/contract/test_audit.go`
  Verify all operations create audit log entries

## Phase 3.4: Domain Layer Implementation (GREEN PHASE)
**Goal: Business entities and logic (only after tests are failing)**

- [ ] **T015** [P] Employee domain entity in `internal/domain/employee.go`
  Business logic, validation rules, state transitions per data-model.md

- [ ] **T016** [P] User domain entity in `internal/domain/user.go`
  Authentication logic, role validation, password hashing

- [ ] **T017** [P] AuditLog domain entity in `internal/domain/audit.go`
  Immutable audit entries, operation tracking

- [ ] **T018** [P] Domain services and interfaces in `internal/domain/service.go`
  EmployeeService, UserService, AuthService interfaces

## Phase 3.5: Data Access Layer
**Goal: Repository pattern with GORM models**

- [ ] **T019** [P] Employee GORM model in `internal/repository/employee_model.go`
  Database mapping, constraints, indexes per data-model.md

- [ ] **T020** [P] User GORM model in `internal/repository/user_model.go`
  Authentication fields, role constraints

- [ ] **T021** [P] Employee repository implementation in `internal/repository/employee_repository.go`
  CRUD operations, pagination, filtering with GORM

- [ ] **T022** [P] User repository implementation in `internal/repository/user_repository.go`
  Authentication queries, user management

- [ ] **T023** [P] Audit log repository in `internal/repository/audit_repository.go`
  Immutable audit entry creation, querying

## Phase 3.6: Application Services
**Goal: Business logic and orchestration**

- [ ] **T024** Employee service implementation in `internal/service/employee_service.go`
  Business rules, validation, audit logging per domain rules

- [ ] **T025** User service implementation in `internal/service/user_service.go`
  Authentication, user management, role handling

- [ ] **T026** JWT authentication service in `internal/service/auth_service.go`
  Token generation, validation, refresh logic per research.md decisions

## Phase 3.7: GraphQL API Layer
**Goal: Resolvers and HTTP handlers**

- [ ] **T027** [P] Generate gqlgen resolvers from schema.graphql
  Run `gqlgen generate` to create resolver interfaces

- [ ] **T028** [P] Employee resolvers in `internal/api/resolver_employee.go`
  Implement all Employee-related GraphQL operations

- [ ] **T029** [P] User resolvers in `internal/api/resolver_user.go`
  Authentication and user management resolvers

- [ ] **T030** [P] Audit log resolvers in `internal/api/resolver_audit.go`
  Audit log querying with filtering/pagination

- [ ] **T031** Health check resolver in `internal/api/resolver_health.go`

## Phase 3.8: HTTP Server and Middleware
**Goal: Gin server with security and logging**

- [ ] **T032** [P] Gin server setup in `cmd/server/main.go`
  HTTP server, GraphQL endpoint, playground, graceful shutdown

- [ ] **T033** [P] JWT authentication middleware in `internal/middleware/auth.go`
  Token validation, user context injection per research.md

- [ ] **T034** [P] Logging middleware in `internal/middleware/logging.go`
  Structured logging with zap per research.md, request correlation

- [ ] **T035** [P] Error handling middleware in `internal/middleware/error.go`
  GraphQL error formatting, HTTP status codes

## Phase 3.9: Configuration and Environment
**Goal: Viper configuration for all environments**

- [ ] **T036** [P] Configuration system in `pkg/config/config.go`
  Environment variables, config files, validation per research.md

## Phase 3.10: Integration Tests (E2E Validation)
**Goal: Validate complete user stories from quickstart.md**

- [ ] **T037** [P] End-to-end employee management test in `tests/integration/test_employee_lifecycle.go`
  Create → Update → List → Delete workflow from quickstart.md

- [ ] **T038** [P] Authentication flow test in `tests/integration/test_auth_flow.go`
  Login → Protected operation → Token refresh → Logout

## Phase 3.11: Polish and Optimization
**Goal: Performance, documentation, and code quality**

- [ ] **T039** [P] Unit tests for domain logic in `tests/unit/test_domain.go`
  Table-driven tests for business rules and validation

- [ ] **T040** [P] Performance and load testing in `tests/performance/`
  Verify <100ms p95 response time, 100+ concurrent users

- [ ] **T041** [P] Update CLAUDE.md with project-specific patterns
  Add employee management examples, GraphQL patterns, testing setup

## Dependencies
```
Infrastructure (T001-T007) → Contract Tests (T008-T014) → Domain (T015-T018) →
Repository (T019-T023) → Services (T024-T026) → API (T027-T031) →
Server (T032-T035) → Config (T036) → Integration (T037-T038) → Polish (T039-T041)
```

## Critical Dependencies (Sequential Execution Required)
- **T008-T014** (Contract Tests) MUST fail before **T015-T041** (Implementation)
- **T015-T018** (Domain) blocks **T019-T023** (Repository)
- **T019-T023** (Repository) blocks **T024-T026** (Services)
- **T024-T026** (Services) blocks **T027-T031** (API)
- **T027-T031** (API) blocks **T032-T036** (Infrastructure)

## Parallel Execution Groups

### Group 1: Infrastructure Setup (Can run together)
```
Task: "Initialize Go module with required dependencies"
Task: "Create project structure per plan.md"
Task: "Configure Go development tools"
```

### Group 2: Database Infrastructure (Can run together)
```
Task: "Create Docker Compose configuration for PostgreSQL"
Task: "Set up GORM with PostgreSQL connection pooling"
Task: "Create database migration system"
Task: "Configure testcontainers for integration tests"
```

### Group 3: Contract Tests (Can run together - MUST ALL FAIL)
```
Task: "GraphQL schema contract test"
Task: "Employee CRUD contract tests"
Task: "Employee listing contract test"
Task: "Authentication contract tests"
Task: "Authorization contract tests"
Task: "Validation contract tests"
Task: "Audit logging contract test"
```

### Group 4: Domain Models (Can run together)
```
Task: "Employee domain entity"
Task: "User domain entity"
Task: "AuditLog domain entity"
Task: "Domain services and interfaces"
```

### Group 5: Repository Layer (Can run together)
```
Task: "Employee GORM model"
Task: "User GORM model"
Task: "Employee repository implementation"
Task: "User repository implementation"
Task: "Audit log repository"
```

## TDD Workflow Example
```bash
# Phase 1: Write failing tests (RED)
go test ./tests/contract/...  # MUST ALL FAIL
git add . && git commit -m "feat: write failing contract tests"

# Phase 2: Implement to make tests pass (GREEN)
# Implement domain layer
# Implement repository layer
# Implement service layer
# Implement API layer
go test ./tests/contract/...  # NOW SHOULD PASS
git add . && git commit -m "feat: implement core functionality"

# Phase 3: Refactor (REFACTOR)
# Optimize, improve code quality, add unit tests
git add . && git commit -m "refactor: optimize and add unit tests"
```

## Quickstart Validation
After completing all tasks, the quickstart.md scenarios should work:
1. ✅ PostgreSQL starts with migrations
2. ✅ GraphQL server starts on port 8080
3. ✅ Admin user creation works
4. ✅ JWT authentication functions
5. ✅ Employee CRUD operations work
6. ✅ Data validation prevents invalid operations
7. ✅ Audit logging captures all changes
8. ✅ Integration tests pass

## Notes
- **[P]** tasks = different files, safe for parallel execution
- **TDD Non-negotiable**: Tests MUST fail before implementation
- **Real dependencies**: Use actual PostgreSQL via testcontainers, NOT mocks
- **Constitution compliance**: Follow all clean architecture, DDD, testing principles
- **Verify quickstart.md**: All scenarios should work after implementation
- **Commit after each task**: Small, focused commits with descriptive messages

## Validation Checklist
- [x] All GraphQL operations have contract tests
- [x] All domain entities have models
- [x] Tests come before implementation (TDD)
- [x] Parallel tasks are truly independent
- [x] Each task specifies exact file path
- [x] No conflicting file modifications in [P] tasks
- [x] Follows Clean Architecture and DDD principles
- [x] Uses real dependencies (PostgreSQL, not mocks)
- [x] Validates against quickstart.md scenarios

## 📊 Progress Summary (Updated: 2025-01-13)

### ✅ **COMPLETED: 20/41 Tasks (49%)**

**Phase 3.1: Infrastructure Setup (3/3) ✅ COMPLETE**
- ✅ T001: Go module with all dependencies
- ✅ T002: Project structure per plan.md
- ✅ T003: Go development tools (.golangci.yml, Makefile, .gitignore)

**Phase 3.2: Database Infrastructure (4/4) ✅ COMPLETE**
- ✅ T004: Docker Compose for PostgreSQL
- ✅ T005: GORM with PostgreSQL connection pooling
- ✅ T006: Database migration system
- ✅ T007: Testcontainers for integration tests

**Phase 3.3: Contract Tests (7/7) - ✅ COMPLETE**
- ✅ T008: GraphQL schema contract test
- ✅ T009: Employee CRUD contract tests
- ✅ T010: Employee listing contract test
- ✅ T011: Authentication contract tests (recently fixed)
- ✅ T012: Authorization contract tests
- ✅ T013: Validation contract tests
- ✅ T014: Audit logging contract test

**Phase 3.4-3.11: Implementation Layers (0/27) - ⚠️ BLOCKED BY DATABASE & T013**

### 🚨 **NEXT STEPS**
1. **Verify RED state** - Run contract tests to ensure they fail
2. **Begin Phase 3.4** - Domain layer implementation

### 🎯 **READY FOR GREEN PHASE**
All contract tests are complete and ready to validate the RED phase before implementation begins.