# Data Model: Employee Management System

## Entity: Employee

### Overview
The Employee entity represents an individual working for the company. It contains both personal and employment information necessary for employee management operations.

### Fields

#### Core Fields
| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| ID | UUID | Yes | Unique identifier for the employee | Auto-generated, UUID v4 |
| FirstName | String | Yes | Employee's first name | 2-50 characters, letters only |
| LastName | String | Yes | Employee's last name | 2-50 characters, letters only |
| Email | String | Yes | Employee's work email | Valid email format, unique |
| Phone | String | No | Employee's phone number | E.164 format or empty |
| Department | String | Yes | Employee's department | 2-50 characters |
| Position | String | Yes | Employee's job title | 2-100 characters |
| HireDate | Date | Yes | Date of employment start | Must be past or current date |
| Salary | Decimal | Yes | Employee's annual salary | > 0, max 1,000,000 |
| Status | Enum | Yes | Employment status | ACTIVE, TERMINATED, ON_LEAVE |
| ManagerID | UUID | No | ID of employee's manager | Must reference existing employee |
| CreatedAt | DateTime | Yes | Record creation timestamp | Auto-generated |
| UpdatedAt | DateTime | Yes | Record last update timestamp | Auto-generated |

#### Address Fields (Embedded)
| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| Street | String | No | Street address | Max 200 characters |
| City | String | No | City | Max 100 characters |
| State | String | No | State/Province | Max 50 characters |
| PostalCode | String | No | Postal/ZIP code | Max 20 characters |
| Country | String | No | Country | ISO 3166-1 alpha-2 |

### Relationships
- **Manager**: Self-referential relationship (many-to-one)
- **Direct Reports**: Self-referential relationship (one-to-many)
- **Audit Logs**: One-to-many relationship with AuditLog entity

### State Transitions
```
ACTIVE ──┐
         │
         └─── TERMINATED (final)
         │
         └─── ON_LEAVE ──┐
                         │
                         └─── ACTIVE
```

**State Transition Rules**:
- Can only transition from ACTIVE to TERMINATED or ON_LEAVE
- Can only transition from ON_LEAVE to ACTIVE
- TERMINATED is a final state
- All transitions require appropriate authorization

## Entity: AuditLog

### Overview
The AuditLog entity records all employee management operations for compliance and security purposes.

### Fields
| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| ID | UUID | Yes | Unique identifier for the audit entry | Auto-generated, UUID v4 |
| EmployeeID | UUID | Yes | ID of affected employee | Must reference existing employee |
| Operation | String | Yes | Operation performed | CREATE, UPDATE, DELETE, STATUS_CHANGE |
| UserID | String | Yes | User who performed the action | From JWT token |
| Timestamp | DateTime | Yes | When the operation occurred | Auto-generated |
| OldValues | JSON | No | Previous values before change | Valid JSON |
| NewValues | JSON | No | New values after change | Valid JSON |
| IPAddress | String | Yes | Client IP address | Valid IP format |
| UserAgent | String | No | Client user agent | Max 500 characters |

### Relationships
- **Employee**: Many-to-one relationship with Employee entity

## Entity: User (Authentication)

### Overview
The User entity handles authentication and authorization for the system.

### Fields
| Field | Type | Required | Description | Validation Rules |
|-------|------|----------|-------------|------------------|
| ID | UUID | Yes | Unique identifier | Auto-generated, UUID v4 |
| Username | String | Yes | Login username | 3-50 characters, alphanumeric + _ |
| Email | String | Yes | User email | Valid email format, unique |
| PasswordHash | String | Yes | Hashed password | Bcrypt hash |
| Role | String | Yes | User role | ADMIN, MANAGER, VIEWER |
| IsActive | Boolean | Yes | Account status | Default: true |
| LastLogin | DateTime | No | Last successful login | Auto-updated |
| CreatedAt | DateTime | Yes | Account creation timestamp | Auto-generated |
| UpdatedAt | DateTime | Yes | Account last update timestamp | Auto-generated |

## Domain Rules

### Employee Validation Rules
1. **Email Uniqueness**: Email must be unique across all employees
2. **Manager Validation**: Manager must be an existing employee with ACTIVE status
3. **Hierarchy Validation**: Cannot create circular reporting relationships
4. **Termination Rules**: Terminated employees cannot have direct reports
5. **Salary Validation**: Manager salary must be >= direct reports' salaries

### Data Integrity Rules
1. **Cascading Deletes**: When an employee is deleted, all associated audit logs are preserved
2. **Soft Deletes**: Employees are marked as TERMINATED rather than physically deleted
3. **Immutable History**: Audit logs cannot be modified once created
4. **Timestamp Consistency**: All timestamps use UTC timezone

### Business Rules
1. **Department Consistency**: All departments must have at least one manager
2. **Hiring Date Validation**: Hire date cannot be in the future
3. **Status Transitions**: Follow defined state transition diagram
4. **Access Control**: Users can only access employees within their organizational scope

## Database Schema

### Tables

#### employees
```sql
CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    department VARCHAR(50) NOT NULL,
    position VARCHAR(100) NOT NULL,
    hire_date DATE NOT NULL,
    salary DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'TERMINATED', 'ON_LEAVE')),
    manager_id UUID REFERENCES employees(id),
    street VARCHAR(200),
    city VARCHAR(100),
    state VARCHAR(50),
    postal_code VARCHAR(20),
    country VARCHAR(2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_employees_department ON employees(department);
CREATE INDEX idx_employees_manager_id ON employees(manager_id);
CREATE INDEX idx_employees_status ON employees(status);
```

#### audit_logs
```sql
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    operation VARCHAR(50) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    old_values JSONB,
    new_values JSONB,
    ip_address INET NOT NULL,
    user_agent VARCHAR(500)
);

-- Indexes for performance
CREATE INDEX idx_audit_logs_employee_id ON audit_logs(employee_id);
CREATE INDEX idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX idx_audit_logs_operation ON audit_logs(operation);
```

#### users
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'VIEWER' CHECK (role IN ('ADMIN', 'MANAGER', 'VIEWER')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    last_login TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
```

### Constraints and Triggers
```sql
-- Prevent circular manager relationships
CREATE OR REPLACE FUNCTION prevent_circular_manager()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.manager_id IS NOT NULL THEN
        WITH RECURSIVE manager_chain AS (
            SELECT id, manager_id FROM employees WHERE id = NEW.manager_id
            UNION
            SELECT e.id, e.manager_id FROM employees e
            INNER JOIN manager_chain mc ON e.id = mc.manager_id
        )
        SELECT EXISTS(SELECT 1 FROM manager_chain WHERE id = NEW.id) INTO circular;

        IF circular THEN
            RAISE EXCEPTION 'Circular manager relationship detected';
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_prevent_circular_manager
BEFORE INSERT OR UPDATE ON employees
FOR EACH ROW EXECUTE FUNCTION prevent_circular_manager();

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_employees_updated_at
BEFORE UPDATE ON employees
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER trg_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

## GraphQL Schema Types

```graphql
scalar Date
scalar DateTime
scalar Upload

type Employee {
    id: ID!
    firstName: String!
    lastName: String!
    email: String!
    phone: String
    department: String!
    position: String!
    hireDate: Date!
    salary: Float!
    status: EmployeeStatus!
    manager: Employee
    directReports: [Employee!]!
    address: Address
    createdAt: DateTime!
    updatedAt: DateTime!
}

type Address {
    street: String
    city: String
    state: String
    postalCode: String
    country: String
}

enum EmployeeStatus {
    ACTIVE
    TERMINATED
    ON_LEAVE
}

type AuditLog {
    id: ID!
    employee: Employee!
    operation: String!
    userId: String!
    timestamp: DateTime!
    oldValues: JSON
    newValues: JSON
    ipAddress: String!
    userAgent: String
}

type User {
    id: ID!
    username: String!
    email: String!
    role: UserRole!
    isActive: Boolean!
    lastLogin: DateTime
    createdAt: DateTime!
    updatedAt: DateTime!
}

enum UserRole {
    ADMIN
    MANAGER
    VIEWER
}

type Query {
    # Employee queries
    employee(id: ID!): Employee
    employees(
        first: Int
        after: String
        filter: EmployeeFilter
        sort: EmployeeSort
    ): EmployeeConnection!

    # Audit queries
    auditLogs(
        employeeId: ID
        operation: String
        from: DateTime
        to: DateTime
        first: Int
        after: String
    ): AuditLogConnection!

    # User queries (admin only)
    user(id: ID!): User
    users: [User!]!
}

type Mutation {
    # Employee mutations
    createEmployee(input: CreateEmployeeInput!): Employee!
    updateEmployee(id: ID!, input: UpdateEmployeeInput!): Employee!
    deleteEmployee(id: ID!): Boolean!
    changeEmployeeStatus(id: ID!, status: EmployeeStatus!): Employee!

    # User mutations (admin only)
    createUser(input: CreateUserInput!): User!
    updateUser(id: ID!, input: UpdateUserInput!): User!
    deleteUser(id: ID!): Boolean!
}

input CreateEmployeeInput {
    firstName: String!
    lastName: String!
    email: String!
    phone: String
    department: String!
    position: String!
    hireDate: Date!
    salary: Float!
    managerId: ID
    address: AddressInput
}

input UpdateEmployeeInput {
    firstName: String
    lastName: String
    email: String
    phone: String
    department: String
    position: String
    salary: Float
    managerId: ID
    address: AddressInput
}

input AddressInput {
    street: String
    city: String
    state: String
    postalCode: String
    country: String
}

input CreateUserInput {
    username: String!
    email: String!
    password: String!
    role: UserRole! = VIEWER
}

input UpdateUserInput {
    username: String
    email: String
    password: String
    role: UserRole
    isActive: Boolean
}

input EmployeeFilter {
    department: String
    status: EmployeeStatus
    managerId: ID
    search: String
}

input EmployeeSort {
    field: EmployeeSortField = NAME
    direction: SortDirection = ASC
}

enum EmployeeSortField {
    NAME
    EMAIL
    DEPARTMENT
    POSITION
    HIRE_DATE
    SALARY
    STATUS
}

enum SortDirection {
    ASC
    DESC
}

# Pagination types
type EmployeeConnection {
    edges: [EmployeeEdge!]!
    pageInfo: PageInfo!
    totalCount: Int!
}

type EmployeeEdge {
    node: Employee!
    cursor: String!
}

type AuditLogConnection {
    edges: [AuditLogEdge!]!
    pageInfo: PageInfo!
    totalCount: Int!
}

type AuditLogEdge {
    node: AuditLog!
    cursor: String!
}

type PageInfo {
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
    startCursor: String
    endCursor: String
}

scalar JSON
```

This data model provides a comprehensive foundation for the employee management system, supporting all the requirements identified in the feature specification while following clean architecture and domain-driven design principles.