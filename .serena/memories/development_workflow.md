# Development Workflow

## Project Development Process

### 1. Feature Creation
```bash
# Create new feature with proper branch structure
./.specify/scripts/bash/create-new-feature.sh "feature description"
# Creates: feature branch, directory structure, and template files
```

### 2. Feature Planning
- **Phase 0**: Research and resolve unknowns
- **Phase 1**: Design contracts and data models
- **Phase 2**: Create implementation tasks
- **Phase 3-4**: Implementation and validation

### 3. Test-Driven Development (TDD) Cycle

#### RED Phase (MUST fail first)
1. Write failing test
2. Verify test fails with expected error
3. Commit failing test
4. **NO implementation allowed**

#### GREEN Phase
1. Write minimal code to make test pass
2. Run test - should now pass
3. No refactoring allowed yet

#### REFACTOR Phase
1. Improve code quality
2. Ensure tests still pass
3. Run all tests
4. Commit improvements

### 4. Implementation Order
1. **Contract Tests** - Define API behavior
2. **Database Infrastructure** - Set up data layer
3. **Domain Layer** - Business logic
4. **Application Layer** - Services
5. **API Layer** - GraphQL resolvers
6. **Integration Tests** - End-to-end validation

### 5. Quality Gates

#### Before Commit
- [ ] All tests pass
- [ ] Code follows conventions
- [ ] No TODO comments without issues
- [ ] Proper error handling
- [ ] Security checks passed

#### Code Review
- [ ] TDD cycle followed correctly
- [ ] Tests are comprehensive
- [ ] Code is maintainable
- [ ] Performance considered
- [ ] Security implications reviewed

### 6. Branch Management
```bash
# Feature branches follow pattern: ###-feature-name
# Examples: 001-employee-management, 002-user-authentication

# Always work on feature branches
git checkout -b 001-new-feature

# Merge to main only after full validation
git checkout main
git merge --no-ff 001-new-feature
git branch -d 001-new-feature
```

### 7. Database Workflow

#### Migrations
1. Create migration file in `migrations/`
2. Write up() and down() functions
3. Test migration with testcontainers
4. Apply to development database
5. Include in deployment process

#### Data Seeding
1. Create seed data in `tests/testdata/`
2. Use for integration tests
3. Document seed structure
4. Keep seed data minimal

### 8. Testing Strategy

#### Test Types and Order
1. **Contract Tests** - API validation (must fail first)
2. **Integration Tests** - Service layer validation
3. **Unit Tests** - Individual function testing

#### Test Dependencies
- Use real PostgreSQL (testcontainers)
- No mocks for external dependencies
- Clean database state for each test
- Parallel test execution when possible

### 9. Configuration Management

#### Environment-Specific Configs
- `configs/config.yaml` - Base configuration
- `configs/config.dev.yaml` - Development overrides
- `configs/config.prod.yaml` - Production settings
- Environment variables take precedence

#### Configuration Validation
- Validate all required fields on startup
- Provide clear error messages for missing config
- Use sensible defaults where appropriate

### 10. Documentation Requirements

#### Code Documentation
- Package-level documentation
- Function documentation for public APIs
- Inline comments for complex logic
- Example usage in docstrings

#### Process Documentation
- Update `CLAUDE.md` with new technologies
- Document breaking changes
- Keep README.md current
- Update project memory files

### 11. Continuous Integration

#### Pre-Commit Hooks
- Run tests
- Check code formatting
- Validate configuration
- Security scans

#### Build Process
1. Run all tests
2. Build application
3. Run integration tests
4. Generate artifacts
5. Deploy to staging

### 12. Deployment Workflow

#### Staging
- Deploy to staging environment
- Run end-to-end tests
- Performance testing
- Security validation

#### Production
- Database migrations
- Zero-downtime deployment
- Health checks
- Monitoring setup

### 13. Monitoring and Observability

#### Logging
- Structured logging with correlation IDs
- Log levels: DEBUG, INFO, WARN, ERROR
- Contextual information in logs
- Log aggregation setup

#### Metrics
- API response times
- Database query performance
- Error rates
- System resource usage

### 14. Security Practices

#### Development
- Input validation on all endpoints
- Parameterized queries only
- Secure password handling
- Regular security scanning

#### Deployment
- Environment secrets management
- Network security policies
- Regular security audits
- Incident response plan