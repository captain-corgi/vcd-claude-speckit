package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
)

// Employee represents the GORM model for the Employee entity
type Employee struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	FirstName string    `gorm:"size:50;not null"`
	LastName  string    `gorm:"size:50;not null"`
	Email     string    `gorm:"size:100;not null;uniqueIndex"`
	Phone     string    `gorm:"size:20"`
	Department string   `gorm:"size:50;not null"`
	Position  string    `gorm:"size:50;not null"`
	HireDate  time.Time `gorm:"not null"`
	Salary    float64   `gorm:"not null"`
	Status    string    `gorm:"size:20;not null;index"`
	ManagerID *uuid.UUID `gorm:"type:uuid;index"`

	// Embedded Address
	Street     string `gorm:"size:200"`
	City       string `gorm:"size:100"`
	State      string `gorm:"size:50"`
	PostalCode string `gorm:"size:20"`
	Country    string `gorm:"size:100"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for the Employee model
func (Employee) TableName() string {
	return "employees"
}

// ToDomain converts the GORM model to the domain entity
func (e *Employee) ToDomain() (*domain.Employee, error) {
	var address *domain.Address
	if e.Street != "" || e.City != "" || e.State != "" || e.PostalCode != "" || e.Country != "" {
		addr, err := domain.NewAddress(e.Street, e.City, e.State, e.PostalCode, e.Country)
		if err != nil {
			return nil, err
		}
		address = addr
	}

	// Parse status
	_, err := domain.ParseEmployeeStatus(e.Status)
	if err != nil {
		return nil, err
	}

	return domain.NewEmployee(
		e.FirstName,
		e.LastName,
		e.Email,
		e.Phone,
		e.Department,
		e.Position,
		e.HireDate,
		e.Salary,
		e.ManagerID,
		address,
	)
}

// FromDomain converts the domain entity to the GORM model
func (e *Employee) FromDomain(employee *domain.Employee) {
	e.ID = employee.ID
	e.FirstName = employee.FirstName
	e.LastName = employee.LastName
	e.Email = employee.Email
	e.Phone = employee.Phone
	e.Department = employee.Department
	e.Position = employee.Position
	e.HireDate = employee.HireDate
	e.Salary = employee.Salary
	e.Status = string(employee.Status)
	e.ManagerID = employee.ManagerID

	if employee.Address != nil {
		e.Street = employee.Address.Street
		e.City = employee.Address.City
		e.State = employee.Address.State
		e.PostalCode = employee.Address.PostalCode
		e.Country = employee.Address.Country
	}

	e.UpdatedAt = employee.UpdatedAt
}

// BeforeCreate is called before creating a record
func (e *Employee) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}
	e.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is called before updating a record
func (e *Employee) BeforeUpdate(tx *gorm.DB) error {
	e.UpdatedAt = time.Now()
	return nil
}