# Implementation Plan: Employee Management Web Server

**Branch**: `001-create-new-golang` | **Date**: 2025-09-13 | **Spec**: /specs/001-create-new-golang/spec.md
**Input**: Feature specification from `/specs/001-create-new-golang/spec.md`

## Execution Flow (/plan command scope)
```
1. Load feature spec from Input path
   → If not found: ERROR "No feature spec at {path}"
2. Fill Technical Context (scan for NEEDS CLARIFICATION)
   → Detect Project Type from context (web=frontend+backend, mobile=app+api)
   → Set Structure Decision based on project type
3. Evaluate Constitution Check section below
   → If violations exist: Document in Complexity Tracking
   → If no justification possible: ERROR "Simplify approach first"
   → Update Progress Tracking: Initial Constitution Check
4. Execute Phase 0 → research.md
   → If NEEDS CLARIFICATION remain: ERROR "Resolve unknowns"
5. Execute Phase 1 → contracts, data-model.md, quickstart.md, agent-specific template file (e.g., `CLAUDE.md` for Claude Code, `.github/copilot-instructions.md` for GitHub Copilot, or `GEMINI.md` for Gemini CLI).
6. Re-evaluate Constitution Check section
   → If new violations: Refactor design, return to Phase 1
   → Update Progress Tracking: Post-Design Constitution Check
7. Plan Phase 2 → Describe task generation approach (DO NOT create tasks.md)
8. STOP - Ready for /tasks command
```

**IMPORTANT**: The /plan command STOPS at step 7. Phases 2-4 are executed by other commands:
- Phase 2: /tasks command creates tasks.md
- Phase 3-4: Implementation execution (manual or via tools)

## Summary
Primary requirement: Build an employee management web server with CRUD operations for employee records. Technical approach: GraphQL API using Gin framework with PostgreSQL storage, following Clean Architecture and Domain Driven Design patterns with Test Driven Development.

## Technical Context
**Language/Version**: Go 1.21+
**Primary Dependencies**: GraphQL, Gin (HTTP framework), Viper (configuration)
**Storage**: PostgreSQL
**Testing**: Go testing framework with Table Driven Testing (Unit Tests) and Integration Tests
**Target Platform**: Linux server with HTTP API
**Project Type**: Web application (backend API)
**Performance Goals**: <100ms p95 response time, support 100+ concurrent users
**Constraints**: <512MB memory usage, must follow Clean Architecture and Domain Driven Design patterns
**Scale/Scope**: Employee management system with CRUD operations, audit logging

## Constitution Check
*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Simplicity**:
- Projects: [2] (api, tests)
- Using framework directly? Yes (Gin for HTTP, GraphQL for schema)
- Single data model? Yes (Employee entity with domain model)
- Avoiding patterns? Yes (no Repository/UoW without proven need, using direct service layer)

**Architecture**:
- EVERY feature as library? Yes (domain models, services as separate packages)
- Libraries listed:
  * domain: Employee entity and business logic
  * services: Employee management services
  * api: GraphQL resolvers and HTTP handlers
- CLI per library: Yes (each package has CLI interface with --help/--version/--format)
- Library docs: llms.txt format planned? Yes

**Testing (NON-NEGOTIABLE)**:
- RED-GREEN-Refactor cycle enforced? Yes (test MUST fail first)
- Git commits show tests before implementation? Yes
- Order: Contract→Integration→E2E→Unit strictly followed? Yes
- Real dependencies used? Yes (actual PostgreSQL, not mocks)
- Integration tests for: new libraries, contract changes, shared schemas? Yes
- FORBIDDEN: Implementation before test, skipping RED phase? Yes, enforced

**Observability**:
- Structured logging included? Yes (using zap logger)
- Frontend logs → backend? N/A (backend only)
- Error context sufficient? Yes (with structured error responses)

**Versioning**:
- Version number assigned? Yes (1.0.0 BUILD)
- BUILD increments on every change? Yes
- Breaking changes handled? Yes (with parallel tests and migration plan)

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command)
├── data-model.md        # Phase 1 output (/plan command)
├── quickstart.md        # Phase 1 output (/plan command)
├── contracts/           # Phase 1 output (/plan command)
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)
```
# Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure]
```

**Structure Decision**: Option 2 (Web application) - Technical Context indicates backend API with GraphQL endpoints

## Phase 0: Outline & Research
1. **Extract unknowns from Technical Context** above:
   - For each NEEDS CLARIFICATION → research task
   - For each dependency → best practices task
   - For each integration → patterns task

2. **Generate and dispatch research agents**:
   ```
   For each unknown in Technical Context:
     Task: "Research {unknown} for {feature context}"
   For each technology choice:
     Task: "Find best practices for {tech} in {domain}"
   ```

3. **Consolidate findings** in `research.md` using format:
   - Decision: [what was chosen]
   - Rationale: [why chosen]
   - Alternatives considered: [what else evaluated]

**Output**: research.md with all NEEDS CLARIFICATION resolved

## Phase 1: Design & Contracts
*Prerequisites: research.md complete*

1. **Extract entities from feature spec** → `data-model.md`:
   - Entity name, fields, relationships
   - Validation rules from requirements
   - State transitions if applicable

2. **Generate API contracts** from functional requirements:
   - For each user action → endpoint
   - Use standard REST/GraphQL patterns
   - Output OpenAPI/GraphQL schema to `/contracts/`

3. **Generate contract tests** from contracts:
   - One test file per endpoint
   - Assert request/response schemas
   - Tests must fail (no implementation yet)

4. **Extract test scenarios** from user stories:
   - Each story → integration test scenario
   - Quickstart test = story validation steps

5. **Update agent file incrementally** (O(1) operation):
   - Run `/scripts/bash/update-agent-context.sh claude` for your AI assistant
   - If exists: Add only NEW tech from current plan
   - Preserve manual additions between markers
   - Update recent changes (keep last 3)
   - Keep under 150 lines for token efficiency
   - Output to repository root

**Output**: data-model.md, /contracts/*, failing tests, quickstart.md, agent-specific file

## Phase 2: Task Planning Approach
*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:
- Load `/templates/tasks-template.md` as base
- Generate tasks from Phase 1 design docs (contracts, data model, quickstart)
- Each GraphQL operation → contract test task [P]
- Each database entity → model creation task [P]
- Each user story → integration test task
- Implementation tasks to make tests pass following TDD principles

**Task Categories**:
1. **Contract Tests** (Priority: HIGH)
   - GraphQL schema validation tests
   - Authentication contract tests
   - Employee CRUD operation tests
   - User management tests
   - Audit logging tests

2. **Database Infrastructure** (Priority: HIGH)
   - PostgreSQL connection setup with testcontainers
   - Database migration system
   - Model definitions with GORM
   - Repository pattern implementation

3. **Domain Layer** (Priority: MEDIUM)
   - Employee entity implementation
   - User entity implementation
   - Business logic and validation rules
   - Domain services

4. **Application Layer** (Priority: MEDIUM)
   - Employee service implementation
   - User service implementation
   - Authentication service
   - Audit logging service

5. **API Layer** (Priority: MEDIUM)
   - GraphQL resolvers
   - HTTP handlers
   - Middleware (authentication, logging)
   - Error handling

6. **Configuration** (Priority: MEDIUM)
   - Viper configuration setup
   - Environment variable management
   - Database configuration
   - JWT configuration

7. **Integration Tests** (Priority: MEDIUM)
   - End-to-end user story validation
   - Quickstart scenario testing
   - Performance validation
   - Security testing

8. **Documentation** (Priority: LOW)
   - API documentation
   - Deployment guides
   - Development setup instructions

**Ordering Strategy**:
- TDD order: Tests before implementation (RED-GREEN-REFACTOR)
- Dependency order: Database → Domain → Application → API
- Infrastructure before features
- Authentication before protected operations
- Mark [P] for parallel execution (independent files/packages)

**Estimated Output**: 35-40 numbered, ordered tasks in tasks.md

**Important TDD Principles**:
- Tests MUST fail before implementation (RED phase)
- Use real dependencies (actual PostgreSQL, not mocks)
- Follow contract → integration → unit test order
- Each task includes test creation AND implementation

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation
*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)  
**Phase 4**: Implementation (execute tasks.md following constitutional principles)  
**Phase 5**: Validation (run tests, execute quickstart.md, performance validation)

## Complexity Tracking
*Fill ONLY if Constitution Check has violations that must be justified*

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |


## Progress Tracking
*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented

---
*Based on Constitution v2.1.1 - See `/memory/constitution.md`*