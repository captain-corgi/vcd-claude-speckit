# Project Overview

## Purpose
This project is an experimental vibe coding repository using GitHub Speckit and Claude AI (with GLM 4.5 model). The main focus is building an employee management web server with GraphQL API using Go.

## Current Status
- **Active Branch**: `001-create-new-golang`
- **Development Phase**: Planning and design phase
- **Implementation Status**: Contract tests and schema defined, ready for implementation

## Tech Stack
- **Language**: Go 1.21+
- **HTTP Framework**: Gin
- **API**: GraphQL with gqlgen
- **Database**: PostgreSQL with GORM
- **Configuration**: Viper
- **Testing**: Table-driven testing with testcontainers
- **Authentication**: JWT with role-based access control

## Architecture
- Clean Architecture with Domain Driven Design
- Test Driven Development (TDD) approach
- Contract-first testing with RED-GREEN-REFACTOR cycle
- Real dependencies (actual PostgreSQL, not mocks)

## Key Entities
- **Employee**: Core entity with personal and employment information
- **User**: Authentication and authorization
- **AuditLog**: Compliance and security logging

## Development Workflow
1. Write failing test (RED)
2. Make test pass (GREEN)
3. Refactor code
4. Run all tests
5. Commit changes
6. Repeat

## Current Implementation State
- ✅ Feature specification completed
- ✅ Data model designed
- ✅ GraphQL schema defined
- ✅ Contract tests created (failing - RED phase)
- ⏳ Ready for implementation phase