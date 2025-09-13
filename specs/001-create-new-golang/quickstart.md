# Quickstart Guide: Employee Management System

This guide demonstrates how to set up and use the Employee Management System. Follow these steps to verify the system works correctly.

## Prerequisites

- Go 1.21+ installed
- Docker and Docker Compose (for PostgreSQL)
- curl or GraphQL client for API testing

## Step 1: Clone and Setup

```bash
# Clone the repository
git clone <repository-url>
cd employee-management-system

# Checkout the feature branch
git checkout 001-create-new-golang

# Install Go dependencies
go mod download
go mod tidy
```

## Step 2: Start PostgreSQL

```bash
# Start PostgreSQL using Docker Compose
docker-compose up -d postgres

# Wait for PostgreSQL to be ready
sleep 10

# Run database migrations
go run cmd/migrate/main.go
```

## Step 3: Start the Server

```bash
# Start the GraphQL server
go run cmd/server/main.go
```

The server should start on `http://localhost:8080` with the GraphQL playground available at `http://localhost:8080/playground`.

## Step 4: Create an Admin User

```bash
# Create admin user (password: admin123)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { createUser(input: {username: \"admin\", email: \"admin@example.com\", password: \"admin123\", role: ADMIN}) { id username email role } }"
  }'
```

## Step 5: Authentication

### Login and Get Token

```bash
# Login as admin
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "mutation { login(username: \"admin\", password: \"admin123\") { token refreshToken user { id username role } } }"
  }'
```

Save the token from the response for subsequent requests.

## Step 6: Employee Operations

### Create an Employee

```bash
# Create employee (replace YOUR_TOKEN_HERE with actual token)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "mutation { createEmployee(input: {firstName: \"John\", lastName: \"Doe\", email: \"john.doe@example.com\", department: \"Engineering\", position: \"Software Engineer\", hireDate: \"2023-01-15\", salary: 75000}) { id firstName lastName email department position hireDate salary status } }"
  }'
```

### Get All Employees

```bash
# List all employees
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "{ employees(first: 10) { edges { node { id firstName lastName email department position } } pageInfo { hasNextPage endCursor } totalCount } }"
  }'
```

### Update an Employee

```bash
# Update employee (replace EMPLOYEE_ID with actual ID)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "mutation { updateEmployee(id: \"EMPLOYEE_ID\", input: {salary: 80000, position: \"Senior Software Engineer\"}) { id salary position updatedAt } }"
  }'
```

### Change Employee Status

```bash
# Put employee on leave
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "mutation { changeEmployeeStatus(id: \"EMPLOYEE_ID\", status: ON_LEAVE) { id status updatedAt } }"
  }'
```

### View Audit Logs

```bash
# View audit logs for an employee
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "{ auditLogs(employeeId: \"EMPLOYEE_ID\", first: 10) { edges { node { operation userId timestamp oldValues newValues ipAddress } } totalCount } }"
  }'
```

## Step 7: Test Validation

### Test Duplicate Email

```bash
# Try to create employee with same email (should fail)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "mutation { createEmployee(input: {firstName: \"Jane\", lastName: \"Smith\", email: \"john.doe@example.com\", department: \"Engineering\", position: \"Software Engineer\", hireDate: \"2023-01-15\", salary: 75000}) { id email } }"
  }'
```

Should return an error: "email already exists"

### Test Required Fields

```bash
# Try to create employee without required fields (should fail)
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "mutation { createEmployee(input: {firstName: \"Jane\"}) { id } }"
  }'
```

Should return validation errors for missing required fields.

## Step 8: Pagination Test

```bash
# Test pagination with more employees
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "query": "{ employees(first: 5) { edges { node { id firstName } } pageInfo { hasNextPage hasPreviousPage startCursor endCursor } totalCount } }"
  }'
```

## Step 9: Health Check

```bash
# Check system health
curl -X POST http://localhost:8080/graphql \
  -H "Content-Type: application/json" \
  -d '{
    "query": "{ health { status version timestamp } }"
  }'
```

## Step 10: Cleanup

```bash
# Stop PostgreSQL
docker-compose down

# Remove database volume (optional)
docker volume rm employee-management-system_postgres_data
```

## Expected Results

After following these steps, you should have:

1. ✅ Running PostgreSQL database with proper schema
2. ✅ GraphQL server responding on port 8080
3. ✅ Admin user created successfully
4. ✅ Authentication working with JWT tokens
5. ✅ Employee CRUD operations working
6. ✅ Data validation preventing duplicates and invalid data
7. ✅ Audit logging capturing all operations
8. ✅ Pagination working for large datasets
9. ✅ Health check endpoint responding

## Troubleshooting

### Database Connection Issues
```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Verify PostgreSQL is running
docker-compose ps

# Restart PostgreSQL
docker-compose restart postgres
```

### Server Issues
```bash
# Check server logs (should be visible in terminal)
# Verify no port conflicts
lsof -i :8080
```

### Authentication Issues
```bash
# Verify token is valid and not expired
# Check server logs for authentication errors
```

## Next Steps

Once you've verified the system works:

1. Run the full test suite: `go test ./...`
2. Review the generated contracts in `specs/001-create-new-golang/contracts/`
3. Implement additional features as needed
4. Deploy to production environment