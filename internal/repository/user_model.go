package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
)

// User represents the GORM model for the User entity
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Username  string    `gorm:"size:50;not null;uniqueIndex"`
	Email     string    `gorm:"size:100;not null;uniqueIndex"`
	Password  string    `gorm:"size:255;not null"`
	Role      string    `gorm:"size:20;not null;index"`
	IsActive  bool      `gorm:"not null;default:true;index"`
	LastLogin *time.Time `gorm:"type:timestamp"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TableName returns the table name for the User model
func (User) TableName() string {
	return "users"
}

// ToDomain converts the GORM model to the domain entity
func (u *User) ToDomain() (*domain.User, error) {
	// Parse role
	role, err := domain.ParseUserRole(u.Role)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password, // This is the hashed password
		Role:      role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.LastLogin != nil {
		user.LastLogin = u.LastLogin
	}

	return user, nil
}

// FromDomain converts the domain entity to the GORM model
func (u *User) FromDomain(user *domain.User) {
	u.ID = user.ID
	u.Username = user.Username
	u.Email = user.Email
	u.Password = user.Password
	u.Role = string(user.Role)
	u.IsActive = user.IsActive
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt

	if user.LastLogin != nil {
		u.LastLogin = user.LastLogin
	}
}

// BeforeCreate is called before creating a record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is called before updating a record
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// FindByUsernameOrEmail is a helper method to find a user by username or email
func FindByUsernameOrEmail(tx *gorm.DB, identifier string) (*User, error) {
	var user User
	err := tx.Where("username = ? OR email = ?", identifier, identifier).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}