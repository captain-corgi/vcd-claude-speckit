package domain

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// EmployeeRepository defines the interface for employee data access
type EmployeeRepository interface {
	// CRUD operations
	Create(ctx context.Context, employee *Employee) error
	GetByID(ctx context.Context, id uuid.UUID) (*Employee, error)
	Update(ctx context.Context, employee *Employee) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	FindByEmail(ctx context.Context, email string) (*Employee, error)
	FindByManagerID(ctx context.Context, managerID uuid.UUID) ([]*Employee, error)
	FindByDepartment(ctx context.Context, department string) ([]*Employee, error)
	FindByStatus(ctx context.Context, status EmployeeStatus) ([]*Employee, error)

	// Pagination and filtering
	List(ctx context.Context, filter EmployeeFilter, sort EmployeeSort, pagination Pagination) (*EmployeeResult, error)
	Count(ctx context.Context, filter EmployeeFilter) (int64, error)

	// Exists checks
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)

	// Bulk operations
	CreateMany(ctx context.Context, employees []*Employee) error
	UpdateMany(ctx context.Context, employees []*Employee) error
	DeleteMany(ctx context.Context, ids []uuid.UUID) error
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	// CRUD operations
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Query operations
	FindByUsername(ctx context.Context, username string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByRole(ctx context.Context, role UserRole) ([]*User, error)
	FindByActiveStatus(ctx context.Context, isActive bool) ([]*User, error)

	// Pagination and filtering
	List(ctx context.Context, filter UserFilter, sort UserSort, pagination Pagination) (*UserResult, error)
	Count(ctx context.Context, filter UserFilter) (int64, error)

	// Exists checks
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	ExistsByID(ctx context.Context, id uuid.UUID) (bool, error)

	// Authentication related
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
	GetInactiveUsers(ctx context.Context, inactiveSince time.Time) ([]*User, error)

	// Bulk operations
	CreateMany(ctx context.Context, users []*User) error
}

// AuditLogRepository defines the interface for audit log data access
type AuditLogRepository interface {
	// CRUD operations
	Create(ctx context.Context, auditLog *AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*AuditLog, error)

	// Query operations
	FindByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]*AuditLog, error)
	FindByOperation(ctx context.Context, operation string) ([]*AuditLog, error)
	FindByUserID(ctx context.Context, userID string) ([]*AuditLog, error)
	FindByTimeRange(ctx context.Context, start, end time.Time) ([]*AuditLog, error)

	// Complex queries
	Find(ctx context.Context, filter AuditLogFilter, sort AuditLogSort, pagination Pagination) (*AuditLogResult, error)
	Count(ctx context.Context, filter AuditLogFilter) (int64, error)

	// Analytics
	GetOperationsSummary(ctx context.Context, start, end time.Time) (map[string]int64, error)
	GetUserActivity(ctx context.Context, userID string, start, end time.Time) ([]*AuditLog, error)

	// Bulk operations
	CreateMany(ctx context.Context, auditLogs []*AuditLog) error

	// Cleanup
	DeleteOldLogs(ctx context.Context, olderThan time.Time) error
}

// EventStoreRepository defines the interface for event store data access
type EventStoreRepository interface {
	// Event storage
	SaveEvent(ctx context.Context, event DomainEvent) error
	GetEventsByAggregateID(ctx context.Context, aggregateID uuid.UUID) ([]DomainEvent, error)
	GetEventsByType(ctx context.Context, eventType string) ([]DomainEvent, error)
	GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]DomainEvent, error)

	// Event versioning
	GetEventVersion(ctx context.Context, aggregateID uuid.UUID) (int64, error)
	UpdateEventVersion(ctx context.Context, aggregateID uuid.UUID, version int64) error

	// Event replay
	GetAllEvents(ctx context.Context) ([]DomainEvent, error)
	GetEventsAfterSequence(ctx context.Context, sequence int64) ([]DomainEvent, error)
}

// Filter and sort types

// EmployeeFilter defines filtering criteria for employee queries
type EmployeeFilter struct {
	Department  *string
	Status      *EmployeeStatus
	ManagerID   *uuid.UUID
	Search      *string
	MinSalary   *float64
	MaxSalary   *float64
	HireDateFrom *time.Time
	HireDateTo   *time.Time
}

// EmployeeSort defines sorting criteria for employee queries
type EmployeeSort struct {
	Field     EmployeeSortField
	Direction SortDirection
}

// EmployeeSortField defines sortable fields for employees
type EmployeeSortField string

const (
	SortByEmployeeID      EmployeeSortField = "id"
	SortByEmployeeName    EmployeeSortField = "name"
	SortByEmployeeEmail   EmployeeSortField = "email"
	SortByEmployeeDept    EmployeeSortField = "department"
	SortByEmployeePosition EmployeeSortField = "position"
	SortByEmployeeHireDate EmployeeSortField = "hire_date"
	SortByEmployeeSalary  EmployeeSortField = "salary"
	SortByEmployeeStatus  EmployeeSortField = "status"
	SortByEmployeeCreated EmployeeSortField = "created_at"
)

// UserFilter defines filtering criteria for user queries
type UserFilter struct {
	Role      *UserRole
	IsActive  *bool
	Search    *string
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	LastLoginFrom *time.Time
	LastLoginTo   *time.Time
}

// UserSort defines sorting criteria for user queries
type UserSort struct {
	Field     UserSortField
	Direction SortDirection
}

// UserSortField defines sortable fields for users
type UserSortField string

const (
	SortByUserID       UserSortField = "id"
	SortByUsername     UserSortField = "username"
	SortByUserEmail    UserSortField = "email"
	SortByUserRole     UserSortField = "role"
	SortByUserActive   UserSortField = "is_active"
	SortByUserCreated  UserSortField = "created_at"
	SortByUserLastLogin UserSortField = "last_login"
)

// AuditLogFilter defines filtering criteria for audit log queries
type AuditLogFilter struct {
	EmployeeID  *uuid.UUID
	Operation   *string
	UserID      *string
	IPAddress   *string
	FromTime    *time.Time
	ToTime      *time.Time
	Operations  []string
}

// AuditLogSort defines sorting criteria for audit log queries
type AuditLogSort struct {
	Field     AuditLogSortField
	Direction SortDirection
}

// AuditLogSortField defines sortable fields for audit logs
type AuditLogSortField string

const (
	SortByAuditLogID        AuditLogSortField = "id"
	SortByAuditLogTimestamp AuditLogSortField = "timestamp"
	SortByAuditLogOperation AuditLogSortField = "operation"
	SortByAuditLogUserID    AuditLogSortField = "user_id"
	SortByAuditLogEmployeeID AuditLogSortField = "employee_id"
)

// SortDirection defines the sort direction
type SortDirection string

const (
	SortDirectionASC  SortDirection = "ASC"
	SortDirectionDESC SortDirection = "DESC"
)

// Pagination defines pagination parameters
type Pagination struct {
	Page     int
	PageSize int
	Cursor   *string
}

// EmployeeResult represents the result of an employee query
type EmployeeResult struct {
	Employees []*Employee
	Total     int64
	Page      int
	PageSize  int
	HasNext   bool
	HasPrev   bool
	Cursors   *PageCursors
}

// UserResult represents the result of a user query
type UserResult struct {
	Users    []*User
	Total    int64
	Page     int
	PageSize int
	HasNext  bool
	HasPrev  bool
	Cursors  *PageCursors
}

// AuditLogResult represents the result of an audit log query
type AuditLogResult struct {
	AuditLogs []*AuditLog
	Total     int64
	Page      int
	PageSize  int
	HasNext   bool
	HasPrev   bool
	Cursors   *PageCursors
}

// PageCursors contains pagination cursors
type PageCursors struct {
	StartCursor *string
	EndCursor   *string
}

// NewPagination creates a new pagination instance
func NewPagination(page, pageSize int) Pagination {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100 // Max page size
	}

	return Pagination{
		Page:     page,
		PageSize: pageSize,
	}
}

// NewCursorPagination creates a new cursor-based pagination instance
func NewCursorPagination(pageSize int, cursor string) Pagination {
	return Pagination{
		PageSize: pageSize,
		Cursor:   &cursor,
	}
}

// GetOffset calculates the offset for page-based pagination
func (p Pagination) GetOffset() int {
	if p.Cursor != nil {
		return 0 // Cursor-based pagination doesn't use offset
	}
	return (p.Page - 1) * p.PageSize
}

// IsCursorBased checks if this is cursor-based pagination
func (p Pagination) IsCursorBased() bool {
	return p.Cursor != nil
}

var (
	ErrRepository = errors.New("repository error")
	ErrNotFound = errors.New("record not found")
	ErrDuplicate = errors.New("duplicate record")
	ErrInvalidFilter = errors.New("invalid filter")
	ErrInvalidSort = errors.New("invalid sort")
	ErrInvalidPagination = errors.New("invalid pagination")
)