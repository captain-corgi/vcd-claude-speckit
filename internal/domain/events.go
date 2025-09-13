package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DomainEvent represents a domain event
type DomainEvent interface {
	GetID() uuid.UUID
	GetAggregateID() uuid.UUID
	GetType() string
	GetTimestamp() time.Time
	GetVersion() int
	GetData() map[string]interface{}
}

// BaseDomainEvent provides common functionality for domain events
type BaseDomainEvent struct {
	ID          uuid.UUID
	AggregateID uuid.UUID
	Type        string
	Timestamp   time.Time
	Version     int
	Data        map[string]interface{}
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(aggregateID uuid.UUID, eventType string, data map[string]interface{}) BaseDomainEvent {
	return BaseDomainEvent{
		ID:          uuid.New(),
		AggregateID: aggregateID,
		Type:        eventType,
		Timestamp:   time.Now(),
		Version:     1,
		Data:        data,
	}
}

// GetID returns the event ID
func (e BaseDomainEvent) GetID() uuid.UUID {
	return e.ID
}

// GetAggregateID returns the aggregate ID
func (e BaseDomainEvent) GetAggregateID() uuid.UUID {
	return e.AggregateID
}

// GetType returns the event type
func (e BaseDomainEvent) GetType() string {
	return e.Type
}

// GetTimestamp returns the event timestamp
func (e BaseDomainEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

// GetVersion returns the event version
func (e BaseDomainEvent) GetVersion() int {
	return e.Version
}

// GetData returns the event data
func (e BaseDomainEvent) GetData() map[string]interface{} {
	return e.Data
}

// EmployeeCreatedEvent is fired when an employee is created
type EmployeeCreatedEvent struct {
	BaseDomainEvent
	Employee *Employee
}

// NewEmployeeCreatedEvent creates a new EmployeeCreatedEvent
func NewEmployeeCreatedEvent(employee *Employee) *EmployeeCreatedEvent {
	data := map[string]interface{}{
		"employeeId":   employee.ID,
		"firstName":    employee.FirstName,
		"lastName":     employee.LastName,
		"email":        employee.Email,
		"department":   employee.Department,
		"position":     employee.Position,
		"salary":       employee.Salary,
		"status":       string(employee.Status),
		"hasManager":   employee.ManagerID != nil,
		"hasAddress":   employee.Address != nil,
	}

	return &EmployeeCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent(employee.ID, "employee.created", data),
		Employee:       employee,
	}
}

// EmployeeUpdatedEvent is fired when an employee is updated
type EmployeeUpdatedEvent struct {
	BaseDomainEvent
	Employee    *Employee
	ChangedFields []string
}

// NewEmployeeUpdatedEvent creates a new EmployeeUpdatedEvent
func NewEmployeeUpdatedEvent(employee *Employee, changedFields []string) *EmployeeUpdatedEvent {
	data := map[string]interface{}{
		"employeeId":    employee.ID,
		"changedFields": changedFields,
		"firstName":     employee.FirstName,
		"lastName":      employee.LastName,
		"email":         employee.Email,
		"department":    employee.Department,
		"position":      employee.Position,
		"salary":        employee.Salary,
		"status":        string(employee.Status),
	}

	return &EmployeeUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent(employee.ID, "employee.updated", data),
		Employee:        employee,
		ChangedFields:   changedFields,
	}
}

// EmployeeDeletedEvent is fired when an employee is deleted
type EmployeeDeletedEvent struct {
	BaseDomainEvent
	EmployeeID uuid.UUID
	OldData   map[string]interface{}
}

// NewEmployeeDeletedEvent creates a new EmployeeDeletedEvent
func NewEmployeeDeletedEvent(employeeID uuid.UUID, oldData map[string]interface{}) *EmployeeDeletedEvent {
	data := map[string]interface{}{
		"employeeId": employeeID,
		"oldData":    oldData,
	}

	return &EmployeeDeletedEvent{
		BaseDomainEvent: NewBaseDomainEvent(employeeID, "employee.deleted", data),
		EmployeeID:     employeeID,
		OldData:       oldData,
	}
}

// EmployeeStatusChangedEvent is fired when an employee's status changes
type EmployeeStatusChangedEvent struct {
	BaseDomainEvent
	EmployeeID    uuid.UUID
	OldStatus     EmployeeStatus
	NewStatus     EmployeeStatus
	ChangedBy     string // User ID who made the change
}

// NewEmployeeStatusChangedEvent creates a new EmployeeStatusChangedEvent
func NewEmployeeStatusChangedEvent(employeeID uuid.UUID, oldStatus, newStatus EmployeeStatus, changedBy string) *EmployeeStatusChangedEvent {
	data := map[string]interface{}{
		"employeeId": employeeID,
		"oldStatus":  string(oldStatus),
		"newStatus":  string(newStatus),
		"changedBy":  changedBy,
	}

	return &EmployeeStatusChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent(employeeID, "employee.status_changed", data),
		EmployeeID:     employeeID,
		OldStatus:      oldStatus,
		NewStatus:      newStatus,
		ChangedBy:      changedBy,
	}
}

// EmployeeSalaryChangedEvent is fired when an employee's salary changes
type EmployeeSalaryChangedEvent struct {
	BaseDomainEvent
	EmployeeID  uuid.UUID
	OldSalary   float64
	NewSalary   float64
	ChangeType  string // "increase", "decrease", "same"
	ChangedBy   string
}

// NewEmployeeSalaryChangedEvent creates a new EmployeeSalaryChangedEvent
func NewEmployeeSalaryChangedEvent(employeeID uuid.UUID, oldSalary, newSalary float64, changedBy string) *EmployeeSalaryChangedEvent {
	changeType := "same"
	if newSalary > oldSalary {
		changeType = "increase"
	} else if newSalary < oldSalary {
		changeType = "decrease"
	}

	data := map[string]interface{}{
		"employeeId": employeeID,
		"oldSalary":  oldSalary,
		"newSalary":  newSalary,
		"changeType": changeType,
		"changeAmount": newSalary - oldSalary,
		"changePercent": ((newSalary - oldSalary) / oldSalary) * 100,
		"changedBy":  changedBy,
	}

	return &EmployeeSalaryChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent(employeeID, "employee.salary_changed", data),
		EmployeeID:      employeeID,
		OldSalary:       oldSalary,
		NewSalary:       newSalary,
		ChangeType:      changeType,
		ChangedBy:       changedBy,
	}
}

// UserLoggedInEvent is fired when a user logs in
type UserLoggedInEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID
	Username  string
	IPAddress string
	UserAgent string
}

// NewUserLoggedInEvent creates a new UserLoggedInEvent
func NewUserLoggedInEvent(userID uuid.UUID, username, ipAddress, userAgent string) *UserLoggedInEvent {
	data := map[string]interface{}{
		"userId":    userID,
		"username":  username,
		"ipAddress": ipAddress,
		"userAgent": userAgent,
	}

	return &UserLoggedInEvent{
		BaseDomainEvent: NewBaseDomainEvent(userID, "user.logged_in", data),
		UserID:          userID,
		Username:        username,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
	}
}

// UserPasswordChangedEvent is fired when a user's password is changed
type UserPasswordChangedEvent struct {
	BaseDomainEvent
	UserID    uuid.UUID
	Username  string
	ChangedBy string // "self" or user ID
	Method    string // "reset" or "change"
}

// NewUserPasswordChangedEvent creates a new UserPasswordChangedEvent
func NewUserPasswordChangedEvent(userID uuid.UUID, username, changedBy, method string) *UserPasswordChangedEvent {
	data := map[string]interface{}{
		"userId":    userID,
		"username":  username,
		"changedBy": changedBy,
		"method":    method,
	}

	return &UserPasswordChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent(userID, "user.password_changed", data),
		UserID:          userID,
		Username:        username,
		ChangedBy:       changedBy,
		Method:          method,
	}
}

// AuditLogCreatedEvent is fired when an audit log entry is created
type AuditLogCreatedEvent struct {
	BaseDomainEvent
	AuditLog *AuditLog
}

// NewAuditLogCreatedEvent creates a new AuditLogCreatedEvent
func NewAuditLogCreatedEvent(auditLog *AuditLog) *AuditLogCreatedEvent {
	data := map[string]interface{}{
		"auditLogId":  auditLog.ID,
		"employeeId":  auditLog.EmployeeID,
		"operation":   auditLog.Operation,
		"userId":      auditLog.UserID,
		"ipAddress":   auditLog.IPAddress,
		"isCreation":  auditLog.IsCreation(),
		"isDeletion":  auditLog.IsDeletion(),
		"isUpdate":    auditLog.IsUpdate(),
		"changedFields": auditLog.GetChangedFields(),
	}

	return &AuditLogCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent(auditLog.ID, "audit_log.created", data),
		AuditLog:        auditLog,
	}
}

// EventHandler is an interface for handling domain events
type EventHandler interface {
	Handle(event DomainEvent) error
	CanHandle(eventType string) bool
}

// EventDispatcher manages event handlers and dispatches events
type EventDispatcher struct {
	handlers map[string][]EventHandler
}

// NewEventDispatcher creates a new event dispatcher
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[string][]EventHandler),
	}
}

// RegisterHandler registers an event handler for a specific event type
func (d *EventDispatcher) RegisterHandler(eventType string, handler EventHandler) {
	d.handlers[eventType] = append(d.handlers[eventType], handler)
}

// Dispatch dispatches an event to all registered handlers
func (d *EventDispatcher) Dispatch(event DomainEvent) error {
	handlers, exists := d.handlers[event.GetType()]
	if !exists {
		return nil // No handlers registered for this event type
	}

	var errs []error
	for _, handler := range handlers {
		if handler.CanHandle(event.GetType()) {
			if err := handler.Handle(event); err != nil {
				errs = append(errs, fmt.Errorf("handler failed for event %s: %w", event.GetType(), err))
			}
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("event dispatch had %d errors: %v", len(errs), errs)
	}

	return nil
}

// EventStore is an interface for storing and retrieving domain events
type EventStore interface {
	Save(event DomainEvent) error
	GetByAggregateID(aggregateID uuid.UUID) ([]DomainEvent, error)
	GetByType(eventType string) ([]DomainEvent, error)
	GetByTimeRange(start, end time.Time) ([]DomainEvent, error)
}

var (
	ErrInvalidEvent = errors.New("invalid domain event")
	ErrEventHandlingFailed = errors.New("event handling failed")
	ErrNoHandlersRegistered = errors.New("no handlers registered for event type")
)