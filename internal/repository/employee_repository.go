package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
)

// EmployeeRepository implements the domain repository interface using GORM
type EmployeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository creates a new EmployeeRepository
func NewEmployeeRepository(db *gorm.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// Create creates a new employee record
func (r *EmployeeRepository) Create(ctx context.Context, employee *domain.Employee) error {
	model := &Employee{}
	model.FromDomain(employee)

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return fmt.Errorf("failed to create employee: %w", err)
	}

	employee.ID = model.ID
	employee.CreatedAt = model.CreatedAt
	employee.UpdatedAt = model.UpdatedAt

	return nil
}

// GetByID retrieves an employee by ID
func (r *EmployeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Employee, error) {
	var model Employee
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get employee by ID: %w", err)
	}

	return model.ToDomain()
}

// Update updates an existing employee record
func (r *EmployeeRepository) Update(ctx context.Context, employee *domain.Employee) error {
	model := &Employee{}
	model.FromDomain(employee)

	err := r.db.WithContext(ctx).Where("id = ?", employee.ID).Save(model).Error
	if err != nil {
		return fmt.Errorf("failed to update employee: %w", err)
	}

	employee.UpdatedAt = model.UpdatedAt

	return nil
}

// Delete deletes an employee record
func (r *EmployeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&Employee{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// FindByEmail retrieves an employee by email
func (r *EmployeeRepository) FindByEmail(ctx context.Context, email string) (*domain.Employee, error) {
	var model Employee
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find employee by email: %w", err)
	}

	return model.ToDomain()
}

// FindByManagerID retrieves all employees managed by the specified manager
func (r *EmployeeRepository) FindByManagerID(ctx context.Context, managerID uuid.UUID) ([]*domain.Employee, error) {
	var models []Employee
	err := r.db.WithContext(ctx).Where("manager_id = ?", managerID).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find employees by manager ID: %w", err)
	}

	result := make([]*domain.Employee, len(models))
	for i, model := range models {
		employee, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = employee
	}

	return result, nil
}

// FindByDepartment retrieves all employees in the specified department
func (r *EmployeeRepository) FindByDepartment(ctx context.Context, department string) ([]*domain.Employee, error) {
	var models []Employee
	err := r.db.WithContext(ctx).Where("department = ?", department).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find employees by department: %w", err)
	}

	result := make([]*domain.Employee, len(models))
	for i, model := range models {
		employee, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = employee
	}

	return result, nil
}

// FindByStatus retrieves all employees with the specified status
func (r *EmployeeRepository) FindByStatus(ctx context.Context, status domain.EmployeeStatus) ([]*domain.Employee, error) {
	var models []Employee
	err := r.db.WithContext(ctx).Where("status = ?", string(status)).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find employees by status: %w", err)
	}

	result := make([]*domain.Employee, len(models))
	for i, model := range models {
		employee, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = employee
	}

	return result, nil
}

// List retrieves employees with filtering, sorting, and pagination
func (r *EmployeeRepository) List(ctx context.Context, filter domain.EmployeeFilter, sort domain.EmployeeSort, pagination domain.Pagination) (*domain.EmployeeResult, error) {
	query := r.db.WithContext(ctx).Model(&Employee{})

	// Apply filters
	if filter.Department != nil {
		query = query.Where("department = ?", *filter.Department)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.ManagerID != nil {
		query = query.Where("manager_id = ?", *filter.ManagerID)
	}
	if filter.Search != nil {
		searchTerm := *filter.Search
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR department LIKE ? OR position LIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}
	if filter.MinSalary != nil {
		query = query.Where("salary >= ?", *filter.MinSalary)
	}
	if filter.MaxSalary != nil {
		query = query.Where("salary <= ?", *filter.MaxSalary)
	}
	if filter.HireDateFrom != nil {
		query = query.Where("hire_date >= ?", *filter.HireDateFrom)
	}
	if filter.HireDateTo != nil {
		query = query.Where("hire_date <= ?", *filter.HireDateTo)
	}

	// Apply sorting
	sortField := map[string]string{
		string(domain.SortByEmployeeID):      "id",
		string(domain.SortByEmployeeName):    "first_name, last_name",
		string(domain.SortByEmployeeEmail):   "email",
		string(domain.SortByEmployeeDept):    "department",
		string(domain.SortByEmployeePosition): "position",
		string(domain.SortByEmployeeHireDate): "hire_date",
		string(domain.SortByEmployeeSalary):  "salary",
		string(domain.SortByEmployeeStatus):  "status",
		string(domain.SortByEmployeeCreated):  "created_at",
	}

	orderBy := sortField[string(sort.Field)]
	if orderBy == "" {
		orderBy = "created_at"
	}

	sortDirection := domain.SortDirectionASC
	if sort.Direction == domain.SortDirectionDESC {
		sortDirection = domain.SortDirectionDESC
	}

	query = query.Order(orderBy + " " + string(sortDirection))

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count employees: %w", err)
	}

	// Apply pagination
	var models []Employee
	offset := pagination.GetOffset()
	query = query.Offset(offset).Limit(pagination.PageSize)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list employees: %w", err)
	}

	// Convert to domain entities
	employees := make([]*domain.Employee, len(models))
	for i, model := range models {
		employee, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		employees[i] = employee
	}

	// Calculate pagination metadata
	hasNext := int64(offset+pagination.PageSize) < total
	hasPrev := offset > 0

	result := &domain.EmployeeResult{
		Employees: employees,
		Total:     total,
		Page:      pagination.Page,
		PageSize:  pagination.PageSize,
		HasNext:   hasNext,
		HasPrev:   hasPrev,
	}

	return result, nil
}

// Count counts employees with the specified filter
func (r *EmployeeRepository) Count(ctx context.Context, filter domain.EmployeeFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&Employee{})

	// Apply filters (same as in List method)
	if filter.Department != nil {
		query = query.Where("department = ?", *filter.Department)
	}
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.ManagerID != nil {
		query = query.Where("manager_id = ?", *filter.ManagerID)
	}
	if filter.Search != nil {
		searchTerm := *filter.Search
		query = query.Where("first_name LIKE ? OR last_name LIKE ? OR email LIKE ? OR department LIKE ? OR position LIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}
	if filter.MinSalary != nil {
		query = query.Where("salary >= ?", *filter.MinSalary)
	}
	if filter.MaxSalary != nil {
		query = query.Where("salary <= ?", *filter.MaxSalary)
	}
	if filter.HireDateFrom != nil {
		query = query.Where("hire_date >= ?", *filter.HireDateFrom)
	}
	if filter.HireDateTo != nil {
		query = query.Where("hire_date <= ?", *filter.HireDateTo)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count employees: %w", err)
	}

	return count, nil
}

// ExistsByEmail checks if an employee with the given email exists
func (r *EmployeeRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Employee{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if employee exists by email: %w", err)
	}
	return count > 0, nil
}

// ExistsByID checks if an employee with the given ID exists
func (r *EmployeeRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Employee{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if employee exists by ID: %w", err)
	}
	return count > 0, nil
}

// CreateMany creates multiple employee records
func (r *EmployeeRepository) CreateMany(ctx context.Context, employees []*domain.Employee) error {
	models := make([]Employee, len(employees))
	for i, employee := range employees {
		models[i] = Employee{}
		models[i].FromDomain(employee)
	}

	err := r.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return fmt.Errorf("failed to create multiple employees: %w", err)
	}

	// Update domain entities with generated IDs and timestamps
	for i, model := range models {
		employees[i].ID = model.ID
		employees[i].CreatedAt = model.CreatedAt
		employees[i].UpdatedAt = model.UpdatedAt
	}

	return nil
}

// UpdateMany updates multiple employee records
func (r *EmployeeRepository) UpdateMany(ctx context.Context, employees []*domain.Employee) error {
	for _, employee := range employees {
		if err := r.Update(ctx, employee); err != nil {
			return fmt.Errorf("failed to update employee %s: %w", employee.ID, err)
		}
	}
	return nil
}

// DeleteMany deletes multiple employee records
func (r *EmployeeRepository) DeleteMany(ctx context.Context, ids []uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&Employee{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete multiple employees: %w", result.Error)
	}

	if result.RowsAffected != int64(len(ids)) {
		return domain.ErrNotFound
	}

	return nil
}

// FindEmployeesWithoutManager retrieves all employees without a manager
func (r *EmployeeRepository) FindEmployeesWithoutManager(ctx context.Context) ([]*domain.Employee, error) {
	var models []Employee
	err := r.db.WithContext(ctx).Where("manager_id IS NULL").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find employees without manager: %w", err)
	}

	result := make([]*domain.Employee, len(models))
	for i, model := range models {
		employee, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = employee
	}

	return result, nil
}

// GetTopPaidEmployees retrieves the top N highest-paid employees
func (r *EmployeeRepository) GetTopPaidEmployees(ctx context.Context, limit int) ([]*domain.Employee, error) {
	var models []Employee
	err := r.db.WithContext(ctx).Order("salary DESC").Limit(limit).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get top paid employees: %w", err)
	}

	result := make([]*domain.Employee, len(models))
	for i, model := range models {
		employee, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = employee
	}

	return result, nil
}