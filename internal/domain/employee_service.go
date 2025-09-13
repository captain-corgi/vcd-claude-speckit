package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EmployeeService handles business logic for employee operations
type EmployeeService struct {
	employeeRepo EmployeeRepository
	auditRepo    AuditLogRepository
	eventStore   EventStoreRepository
	eventDispatcher *EventDispatcher
}

// NewEmployeeService creates a new employee service
func NewEmployeeService(
	employeeRepo EmployeeRepository,
	auditRepo AuditLogRepository,
	eventStore EventStoreRepository,
	eventDispatcher *EventDispatcher,
) *EmployeeService {
	return &EmployeeService{
		employeeRepo:    employeeRepo,
		auditRepo:       auditRepo,
		eventStore:      eventStore,
		eventDispatcher: eventDispatcher,
	}
}

// CreateEmployee creates a new employee with validation and audit logging
func (s *EmployeeService) CreateEmployee(
	ctx context.Context,
	firstName, lastName, email, phone, department, position string,
	hireDate time.Time,
	salary float64,
	managerID *uuid.UUID,
	address *Address,
	userID string,
	ipAddress, userAgent string,
) (*Employee, error) {
	// Validate business rules
	if err := s.validateEmployeeCreation(ctx, email, managerID); err != nil {
		return nil, err
	}

	// Create employee
	employee, err := NewEmployee(firstName, lastName, email, phone, department, position, hireDate, salary, managerID, address)
	if err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}

	// Save to repository
	if err := s.employeeRepo.Create(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to save employee: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForCreation(employee, userID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		// Log error but don't fail the operation
		// In production, you might want to handle this more robustly
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewEmployeeCreatedEvent(employee)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	// Create audit log event
	auditEvent := NewAuditLogCreatedEvent(auditLog)
	if err := s.eventStore.SaveEvent(ctx, auditEvent); err != nil {
		fmt.Printf("Warning: failed to save audit event: %v\n", err)
	}

	return employee, nil
}

// UpdateEmployee updates an existing employee with validation and audit logging
func (s *EmployeeService) UpdateEmployee(
	ctx context.Context,
	id uuid.UUID,
	updates map[string]interface{},
	userID string,
	ipAddress, userAgent string,
) (*Employee, error) {
	// Get existing employee
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	// Create snapshot of old values for audit logging
	oldValues := s.createEmployeeSnapshot(employee)

	// Apply updates
	changedFields, err := s.applyEmployeeUpdates(employee, updates)
	if err != nil {
		return nil, err
	}

	// Validate business rules
	if err := s.validateEmployeeUpdate(ctx, employee); err != nil {
		return nil, err
	}

	// Update employee
	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	// Create audit log
	newValues := s.createEmployeeSnapshot(employee)
	auditLog := s.createAuditLogForUpdate(employee, oldValues, newValues, userID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewEmployeeUpdatedEvent(employee, changedFields)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return employee, nil
}

// DeleteEmployee deletes an employee with validation and audit logging
func (s *EmployeeService) DeleteEmployee(
	ctx context.Context,
	id uuid.UUID,
	userID string,
	ipAddress, userAgent string,
) error {
	// Get existing employee
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get employee: %w", err)
	}

	if employee == nil {
		return ErrEmployeeNotFound
	}

	// Validate business rules
	if err := s.validateEmployeeDeletion(ctx, employee); err != nil {
		return err
	}

	// Create snapshot for audit logging
	oldValues := s.createEmployeeSnapshot(employee)

	// Delete employee
	if err := s.employeeRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete employee: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForDeletion(employee, oldValues, userID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewEmployeeDeletedEvent(id, oldValues)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return nil
}

// ChangeEmployeeStatus changes an employee's status with validation and audit logging
func (s *EmployeeService) ChangeEmployeeStatus(
	ctx context.Context,
	id uuid.UUID,
	newStatus EmployeeStatus,
	userID string,
	ipAddress, userAgent string,
) (*Employee, error) {
	// Get existing employee
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	// Validate status transition
	if err := employee.ChangeStatus(newStatus); err != nil {
		return nil, err
	}

	// Update employee
	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForStatusChange(employee, newStatus, userID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewEmployeeStatusChangedEvent(id, employee.Status, newStatus, userID)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return employee, nil
}

// UpdateEmployeeSalary updates an employee's salary with validation and audit logging
func (s *EmployeeService) UpdateEmployeeSalary(
	ctx context.Context,
	id uuid.UUID,
	newSalary float64,
	userID string,
	ipAddress, userAgent string,
) (*Employee, error) {
	// Get existing employee
	employee, err := s.employeeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	if employee == nil {
		return nil, ErrEmployeeNotFound
	}

	oldSalary := employee.Salary

	// Update salary
	if err := employee.UpdateSalary(newSalary); err != nil {
		return nil, err
	}

	// Update employee
	if err := s.employeeRepo.Update(ctx, employee); err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForSalaryChange(employee, oldSalary, newSalary, userID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewEmployeeSalaryChangedEvent(id, oldSalary, newSalary, userID)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return employee, nil
}

// GetEmployeeByID retrieves an employee by ID
func (s *EmployeeService) GetEmployeeByID(ctx context.Context, id uuid.UUID) (*Employee, error) {
	return s.employeeRepo.GetByID(ctx, id)
}

// ListEmployees returns a paginated list of employees
func (s *EmployeeService) ListEmployees(
	ctx context.Context,
	filter EmployeeFilter,
	sort EmployeeSort,
	pagination Pagination,
) (*EmployeeResult, error) {
	return s.employeeRepo.List(ctx, filter, sort, pagination)
}

// SearchEmployees searches for employees by various criteria
func (s *EmployeeService) SearchEmployees(
	ctx context.Context,
	searchTerm string,
	department *string,
	status *EmployeeStatus,
	pagination Pagination,
) (*EmployeeResult, error) {
	filter := EmployeeFilter{
		Search:     &searchTerm,
		Department: department,
		Status:     status,
	}

	sort := EmployeeSort{
		Field:     SortByEmployeeName,
		Direction: SortDirectionASC,
	}

	return s.employeeRepo.List(ctx, filter, sort, pagination)
}

// Validation methods

func (s *EmployeeService) validateEmployeeCreation(ctx context.Context, email string, managerID *uuid.UUID) error {
	// Check for duplicate email
	if exists, err := s.employeeRepo.ExistsByEmail(ctx, email); err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	} else if exists {
		return ErrEmailAlreadyExists
	}

	// Validate manager exists if provided
	if managerID != nil {
		if exists, err := s.employeeRepo.ExistsByID(ctx, *managerID); err != nil {
			return fmt.Errorf("failed to check manager existence: %w", err)
		} else if !exists {
			return ErrManagerNotFound
		}
	}

	return nil
}

func (s *EmployeeService) validateEmployeeUpdate(ctx context.Context, employee *Employee) error {
	// Check for email uniqueness if email was changed
	if exists, err := s.employeeRepo.ExistsByEmail(ctx, employee.Email); err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	} else if exists {
		// Get existing employee to check if it's the same employee
		existing, err := s.employeeRepo.GetByID(ctx, employee.ID)
		if err != nil {
			return fmt.Errorf("failed to get existing employee: %w", err)
		}
		if existing == nil || existing.Email != employee.Email {
			return ErrEmailAlreadyExists
		}
	}

	// Validate manager exists if manager was changed
	if employee.ManagerID != nil {
		if exists, err := s.employeeRepo.ExistsByID(ctx, *employee.ManagerID); err != nil {
			return fmt.Errorf("failed to check manager existence: %w", err)
		} else if !exists {
			return ErrManagerNotFound
		}
	}

	return nil
}

func (s *EmployeeService) validateEmployeeDeletion(ctx context.Context, employee *Employee) error {
	// Check if employee has direct reports
	reports, err := s.employeeRepo.FindByManagerID(ctx, employee.ID)
	if err != nil {
		return fmt.Errorf("failed to check for direct reports: %w", err)
	}
	if len(reports) > 0 {
		return ErrEmployeeHasDirectReports
	}

	return nil
}

// Helper methods

func (s *EmployeeService) createEmployeeSnapshot(employee *Employee) map[string]interface{} {
	snapshot := map[string]interface{}{
		"id":         employee.ID,
		"firstName":  employee.FirstName,
		"lastName":   employee.LastName,
		"email":      employee.Email,
		"phone":      employee.Phone,
		"department": employee.Department,
		"position":   employee.Position,
		"hireDate":   employee.HireDate,
		"salary":     employee.Salary,
		"status":     string(employee.Status),
		"updatedAt":  employee.UpdatedAt,
	}

	if employee.ManagerID != nil {
		snapshot["managerId"] = *employee.ManagerID
	}

	if employee.Address != nil {
		snapshot["address"] = map[string]interface{}{
			"street":     employee.Address.Street,
			"city":       employee.Address.City,
			"state":      employee.Address.State,
			"postalCode": employee.Address.PostalCode,
			"country":    employee.Address.Country,
		}
	}

	return snapshot
}

func (s *EmployeeService) applyEmployeeUpdates(employee *Employee, updates map[string]interface{}) ([]string, error) {
	var changedFields []string

	for field, value := range updates {
		switch field {
		case "firstName":
			if val, ok := value.(string); ok {
				employee.FirstName = val
				changedFields = append(changedFields, "firstName")
			}
		case "lastName":
			if val, ok := value.(string); ok {
				employee.LastName = val
				changedFields = append(changedFields, "lastName")
			}
		case "email":
			if val, ok := value.(string); ok {
				employee.Email = val
				changedFields = append(changedFields, "email")
			}
		case "phone":
			if val, ok := value.(string); ok {
				employee.Phone = val
				changedFields = append(changedFields, "phone")
			}
		case "department":
			if val, ok := value.(string); ok {
				employee.Department = val
				changedFields = append(changedFields, "department")
			}
		case "position":
			if val, ok := value.(string); ok {
				employee.Position = val
				changedFields = append(changedFields, "position")
			}
		case "salary":
			if val, ok := value.(float64); ok {
				employee.Salary = val
				changedFields = append(changedFields, "salary")
			}
		case "managerId":
			if val, ok := value.(string); ok {
				if val == "" {
					employee.ManagerID = nil
				} else {
					id, err := uuid.Parse(val)
					if err != nil {
						return nil, fmt.Errorf("invalid manager ID: %w", err)
					}
					employee.ManagerID = &id
				}
				changedFields = append(changedFields, "managerId")
			}
		case "address":
			if val, ok := value.(map[string]interface{}); ok {
				address, err := s.parseAddress(val)
				if err != nil {
					return nil, fmt.Errorf("invalid address: %w", err)
				}
				employee.Address = address
				changedFields = append(changedFields, "address")
			}
		}
	}

	employee.UpdatedAt = time.Now()
	return changedFields, nil
}

func (s *EmployeeService) parseAddress(data map[string]interface{}) (*Address, error) {
	street, _ := data["street"].(string)
	city, _ := data["city"].(string)
	state, _ := data["state"].(string)
	postalCode, _ := data["postalCode"].(string)
	country, _ := data["country"].(string)

	return NewAddress(street, city, state, postalCode, country)
}

// Audit log creation methods

func (s *EmployeeService) createAuditLogForCreation(employee *Employee, userID, ipAddress, userAgent string) *AuditLog {
	newValues := s.createEmployeeSnapshot(employee)

	auditLog, err := NewAuditLog(
		employee.ID,
		OperationCreateEmployee,
		userID,
		nil,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// If audit log creation fails, we still want to continue with the main operation
		// Log the error but return nil to not break the flow
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
		return nil
	}
	return auditLog
}

func (s *EmployeeService) createAuditLogForUpdate(employee *Employee, oldValues, newValues map[string]interface{}, userID, ipAddress, userAgent string) *AuditLog {
	auditLog, err := NewAuditLog(
		employee.ID,
		OperationUpdateEmployee,
		userID,
		oldValues,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// If audit log creation fails, we still want to continue with the main operation
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
		return nil
	}
	return auditLog
}

func (s *EmployeeService) createAuditLogForDeletion(employee *Employee, oldValues map[string]interface{}, userID, ipAddress, userAgent string) *AuditLog {
	auditLog, err := NewAuditLog(
		employee.ID,
		OperationDeleteEmployee,
		userID,
		oldValues,
		nil,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// If audit log creation fails, we still want to continue with the main operation
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
		return nil
	}
	return auditLog
}

func (s *EmployeeService) createAuditLogForStatusChange(employee *Employee, newStatus EmployeeStatus, userID, ipAddress, userAgent string) *AuditLog {
	oldValues := map[string]interface{}{
		"status": string(employee.Status),
	}
	newValues := map[string]interface{}{
		"status": string(newStatus),
	}

	auditLog, err := NewAuditLog(
		employee.ID,
		OperationChangeStatus,
		userID,
		oldValues,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// If audit log creation fails, we still want to continue with the main operation
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
		return nil
	}
	return auditLog
}

func (s *EmployeeService) createAuditLogForSalaryChange(employee *Employee, oldSalary, newSalary float64, userID, ipAddress, userAgent string) *AuditLog {
	oldValues := map[string]interface{}{
		"salary": oldSalary,
	}
	newValues := map[string]interface{}{
		"salary": newSalary,
	}

	auditLog, err := NewAuditLog(
		employee.ID,
		OperationUpdateSalary,
		userID,
		oldValues,
		newValues,
		ipAddress,
		userAgent,
	)
	if err != nil {
		// If audit log creation fails, we still want to continue with the main operation
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
		return nil
	}
	return auditLog
}

var (
	ErrEmployeeNotFound = errors.New("employee not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrManagerNotFound = errors.New("manager not found")
	ErrEmployeeHasDirectReports = errors.New("employee has direct reports and cannot be deleted")
)