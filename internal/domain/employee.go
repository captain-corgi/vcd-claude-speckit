package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
)

// Employee represents an employee in the system
type Employee struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Phone     string
	Department string
	Position  string
	HireDate  time.Time
	Salary    float64
	Status    EmployeeStatus
	ManagerID *uuid.UUID
	Address   *Address
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewEmployee creates a new Employee with validation
func NewEmployee(firstName, lastName, email, phone, department, position string, hireDate time.Time, salary float64, managerID *uuid.UUID, address *Address) (*Employee, error) {
	employee := &Employee{
		ID:        uuid.New(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		Department: department,
		Position:  position,
		HireDate:  hireDate,
		Salary:    salary,
		Status:    EmployeeStatusActive,
		ManagerID: managerID,
		Address:   address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := employee.Validate(); err != nil {
		return nil, err
	}

	return employee, nil
}

// Validate validates all employee fields
func (e *Employee) Validate() error {
	if e == nil {
		return errors.New("employee cannot be nil")
	}

	// Validate ID
	if e.ID == uuid.Nil {
		return errors.New("employee ID cannot be empty")
	}

	// Validate first name
	if err := validateName(e.FirstName, "first name"); err != nil {
		return err
	}

	// Validate last name
	if err := validateName(e.LastName, "last name"); err != nil {
		return err
	}

	// Validate email
	if err := validateEmail(e.Email); err != nil {
		return err
	}

	// Validate phone (optional)
	if e.Phone != "" {
		if err := validatePhone(e.Phone); err != nil {
			return err
		}
	}

	// Validate department
	if err := validateDepartment(e.Department); err != nil {
		return err
	}

	// Validate position
	if err := validatePosition(e.Position); err != nil {
		return err
	}

	// Validate hire date
	if err := validateHireDate(e.HireDate); err != nil {
		return err
	}

	// Validate salary
	if err := validateSalary(e.Salary); err != nil {
		return err
	}

	// Validate status
	if !e.Status.IsValid() {
		return fmt.Errorf("invalid employee status: %s", e.Status)
	}

	// Validate address
	if e.Address != nil {
		if err := e.Address.Validate(); err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}
	}

	return nil
}

// validateName validates a name field
func validateName(name, fieldName string) error {
	if name == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if len(name) < 2 {
		return fmt.Errorf("%s must be at least 2 characters long", fieldName)
	}
	if len(name) > 50 {
		return fmt.Errorf("%s cannot exceed 50 characters", fieldName)
	}

	// Check for valid characters (letters, spaces, hyphens, apostrophes)
	for _, r := range name {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' && r != '\'' {
			return fmt.Errorf("%s contains invalid characters", fieldName)
		}
	}

	return nil
}

// validateEmail validates email format
func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > 100 {
		return errors.New("email cannot exceed 100 characters")
	}

	// Basic email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("email format is invalid")
	}

	return nil
}

// validatePhone validates phone format
func validatePhone(phone string) error {
	if len(phone) > 20 {
		return errors.New("phone cannot exceed 20 characters")
	}

	// Allow digits, spaces, hyphens, plus, and parentheses
	phoneRegex := regexp.MustCompile(`^[\d\s\-\+\(\)]+$`)
	if !phoneRegex.MatchString(phone) {
		return errors.New("phone format is invalid")
	}

	return nil
}

// validateDepartment validates department
func validateDepartment(department string) error {
	if department == "" {
		return errors.New("department is required")
	}
	if len(department) < 2 {
		return errors.New("department must be at least 2 characters long")
	}
	if len(department) > 50 {
		return errors.New("department cannot exceed 50 characters")
	}

	// Check for valid characters (letters, spaces, and ampersands)
	for _, r := range department {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '&' {
			return errors.New("department contains invalid characters")
		}
	}

	return nil
}

// validatePosition validates position
func validatePosition(position string) error {
	if position == "" {
		return errors.New("position is required")
	}
	if len(position) < 2 {
		return errors.New("position must be at least 2 characters long")
	}
	if len(position) > 50 {
		return errors.New("position cannot exceed 50 characters")
	}

	// Check for valid characters (letters, spaces, hyphens, slashes)
	for _, r := range position {
		if !unicode.IsLetter(r) && !unicode.IsSpace(r) && r != '-' && r != '/' {
			return errors.New("position contains invalid characters")
		}
	}

	return nil
}

// validateHireDate validates hire date
func validateHireDate(hireDate time.Time) error {
	if hireDate.IsZero() {
		return errors.New("hire date is required")
	}

	// Check if hire date is in the future
	if hireDate.After(time.Now()) {
		return errors.New("hire date cannot be in the future")
	}

	// Check if hire date is too far in the past (more than 50 years)
	minDate := time.Now().AddDate(-50, 0, 0)
	if hireDate.Before(minDate) {
		return errors.New("hire date cannot be more than 50 years in the past")
	}

	return nil
}

// validateSalary validates salary
func validateSalary(salary float64) error {
	if salary < 0 {
		return errors.New("salary cannot be negative")
	}
	if salary == 0 {
		return errors.New("salary is required")
	}
	if salary > 1000000 {
		return errors.New("salary cannot exceed $1,000,000")
	}

	return nil
}

// FullName returns the employee's full name
func (e *Employee) FullName() string {
	return fmt.Sprintf("%s %s", e.FirstName, e.LastName)
}

// IsActive checks if the employee is active
func (e *Employee) IsActive() bool {
	return e.Status == EmployeeStatusActive
}

// IsTerminated checks if the employee is terminated
func (e *Employee) IsTerminated() bool {
	return e.Status == EmployeeStatusTerminated
}

// IsOnLeave checks if the employee is on leave
func (e *Employee) IsOnLeave() bool {
	return e.Status == EmployeeStatusOnLeave
}

// CanBeManagedBy checks if the employee can be managed by the given manager
func (e *Employee) CanBeManagedBy(managerID uuid.UUID) bool {
	if e.ManagerID == nil {
		return false
	}
	return *e.ManagerID == managerID
}

// HasManager checks if the employee has a manager
func (e *Employee) HasManager() bool {
	return e.ManagerID != nil
}

// ChangeStatus changes the employee's status with validation
func (e *Employee) ChangeStatus(newStatus EmployeeStatus) error {
	if err := e.Status.ValidateTransition(newStatus); err != nil {
		return err
	}

	e.Status = newStatus
	e.UpdatedAt = time.Now()
	return nil
}

// UpdateSalary updates the employee's salary with validation
func (e *Employee) UpdateSalary(newSalary float64) error {
	if err := validateSalary(newSalary); err != nil {
		return err
	}

	if e.IsTerminated() {
		return errors.New("cannot update salary for terminated employees")
	}

	e.Salary = newSalary
	e.UpdatedAt = time.Now()
	return nil
}

// UpdateContactInfo updates contact information with validation
func (e *Employee) UpdateContactInfo(email, phone string) error {
	if err := validateEmail(email); err != nil {
		return err
	}

	if phone != "" {
		if err := validatePhone(phone); err != nil {
			return err
		}
	}

	e.Email = email
	e.Phone = phone
	e.UpdatedAt = time.Now()
	return nil
}

// UpdatePosition updates position and department with validation
func (e *Employee) UpdatePosition(position, department string) error {
	if err := validatePosition(position); err != nil {
		return err
	}
	if err := validateDepartment(department); err != nil {
		return err
	}

	if e.IsTerminated() {
		return errors.New("cannot update position for terminated employees")
	}

	e.Position = position
	e.Department = department
	e.UpdatedAt = time.Now()
	return nil
}

// UpdateAddress updates the employee's address
func (e *Employee) UpdateAddress(address *Address) error {
	if address != nil {
		if err := address.Validate(); err != nil {
			return fmt.Errorf("invalid address: %w", err)
		}
	}

	e.Address = address
	e.UpdatedAt = time.Now()
	return nil
}

// SetManager sets the employee's manager
func (e *Employee) SetManager(managerID *uuid.UUID) error {
	if managerID != nil && *managerID == e.ID {
		return errors.New("employee cannot be their own manager")
	}

	e.ManagerID = managerID
	e.UpdatedAt = time.Now()
	return nil
}

// YearsOfService calculates years of service
func (e *Employee) YearsOfService() float64 {
	now := time.Now()
	duration := now.Sub(e.HireDate)
	return duration.Hours() / (24 * 365.25)
}

// TenureString returns a human-readable tenure string
func (e *Employee) TenureString() string {
	years := e.YearsOfService()
	if years < 1 {
		months := int(years * 12)
		return fmt.Sprintf("%d month%s", months, pluralS(months))
	}
	return fmt.Sprintf("%.1f year%s", years, pluralS(int(years)))
}

func pluralS(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

// Clone returns a copy of the employee
func (e *Employee) Clone() *Employee {
	if e == nil {
		return nil
	}

	var managerID *uuid.UUID
	if e.ManagerID != nil {
		id := *e.ManagerID
		managerID = &id
	}

	return &Employee{
		ID:        e.ID,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Email:     e.Email,
		Phone:     e.Phone,
		Department: e.Department,
		Position:  e.Position,
		HireDate:  e.HireDate,
		Salary:    e.Salary,
		Status:    e.Status,
		ManagerID: managerID,
		Address:   e.Address.Clone(),
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

var (
	ErrInvalidEmployee = errors.New("invalid employee")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidPhone = errors.New("invalid phone")
	ErrInvalidSalary = errors.New("invalid salary")
	ErrInvalidHireDate = errors.New("invalid hire date")
	ErrEmployeeTerminated = errors.New("employee is terminated")
	ErrCannotSelfManage = errors.New("employee cannot be their own manager")
)