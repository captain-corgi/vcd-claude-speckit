package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID
	Username  string
	Email     string
	Password  string // Hashed password
	Role      UserRole
	IsActive  bool
	LastLogin *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User with validation
func NewUser(username, email, password string, role UserRole) (*User, error) {
	user := &User{
		ID:        uuid.New(),
		Username:  username,
		Email:     email,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Validate all fields first
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// Hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = hashedPassword

	return user, nil
}

// Validate validates all user fields
func (u *User) Validate() error {
	if u == nil {
		return errors.New("user cannot be nil")
	}

	// Validate ID
	if u.ID == uuid.Nil {
		return errors.New("user ID cannot be empty")
	}

	// Validate username
	if err := validateUsername(u.Username); err != nil {
		return err
	}

	// Validate email
	if err := validateUserEmail(u.Email); err != nil {
		return err
	}

	// Validate role
	if !u.Role.IsValid() {
		return fmt.Errorf("invalid user role: %s", u.Role)
	}

	return nil
}

// validateUsername validates username format
func validateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters long")
	}
	if len(username) > 50 {
		return errors.New("username cannot exceed 50 characters")
	}

	// Username should contain only letters, numbers, and underscores
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, and underscores")
	}

	// Username should not start or end with underscore
	if username[0] == '_' || username[len(username)-1] == '_' {
		return errors.New("username cannot start or end with underscore")
	}

	return nil
}

// validateUserEmail validates email format for users
func validateUserEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if len(email) > 100 {
		return errors.New("email cannot exceed 100 characters")
	}

	// More comprehensive email regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return errors.New("email format is invalid")
	}

	return nil
}

// Authenticate verifies the provided password against the stored hash
func (u *User) Authenticate(password string) bool {
	if !u.IsActive {
		return false
	}
	return CheckPassword(password, u.Password)
}

// UpdateLastLogin updates the last login timestamp
func (u *User) UpdateLastLogin() {
	now := time.Now()
	u.LastLogin = &now
	u.UpdatedAt = now
}

// Activate activates the user account
func (u *User) Activate() error {
	if u.IsActive {
		return errors.New("user is already active")
	}

	u.IsActive = true
	u.UpdatedAt = time.Now()
	return nil
}

// Deactivate deactivates the user account
func (u *User) Deactivate() error {
	if !u.IsActive {
		return errors.New("user is already inactive")
	}

	u.IsActive = false
	u.UpdatedAt = time.Now()
	return nil
}

// ChangePassword changes the user's password
func (u *User) ChangePassword(currentPassword, newPassword string) error {
	// Verify current password
	if !u.Authenticate(currentPassword) {
		return errors.New("current password is incorrect")
	}

	// Validate new password
	if err := validatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("new password is too weak: %w", err)
	}

	// Hash new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// ResetPassword resets the user's password (admin function)
func (u *User) ResetPassword(newPassword string) error {
	// Validate new password
	if err := validatePasswordStrength(newPassword); err != nil {
		return fmt.Errorf("new password is too weak: %w", err)
	}

	// Hash new password
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash new password: %w", err)
	}

	u.Password = hashedPassword
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateRole updates the user's role
func (u *User) UpdateRole(newRole UserRole) error {
	if !newRole.IsValid() {
		return fmt.Errorf("invalid user role: %s", newRole)
	}

	u.Role = newRole
	u.UpdatedAt = time.Now()
	return nil
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(newEmail string) error {
	if err := validateUserEmail(newEmail); err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}

	u.Email = newEmail
	u.UpdatedAt = time.Now()
	return nil
}

// HasPermission checks if the user has a specific permission
func (u *User) HasPermission(permission string) bool {
	return u.Role.HasPermission(permission)
}

// HasAnyPermission checks if the user has any of the specified permissions
func (u *User) HasAnyPermission(permissions ...string) bool {
	for _, permission := range permissions {
		if u.HasPermission(permission) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if the user has all of the specified permissions
func (u *User) HasAllPermissions(permissions ...string) bool {
	for _, permission := range permissions {
		if !u.HasPermission(permission) {
			return false
		}
	}
	return true
}

// IsAdmin checks if the user is an admin
func (u *User) IsAdmin() bool {
	return u.Role == UserRoleAdmin
}

// IsManager checks if the user is a manager
func (u *User) IsManager() bool {
	return u.Role == UserRoleManager
}

// IsViewer checks if the user is a viewer
func (u *User) IsViewer() bool {
	return u.Role == UserRoleViewer
}

// CanAccessSalary checks if the user can access salary information
func (u *User) CanAccessSalary() bool {
	return u.Role.CanAccessSalary()
}

// CanManageUsers checks if the user can manage users
func (u *User) CanManageUsers() bool {
	return u.Role.CanManageUsers()
}

// CanDeleteEmployees checks if the user can delete employees
func (u *User) CanDeleteEmployees() bool {
	return u.Role.CanDeleteEmployees()
}

// CanViewAuditLogs checks if the user can view audit logs
func (u *User) CanViewAuditLogs() bool {
	return u.Role.CanViewAuditLogs()
}

// IsOnline checks if the user was recently active (within 30 minutes)
func (u *User) IsOnline() bool {
	if u.LastLogin == nil {
		return false
	}
	return time.Since(*u.LastLogin) < 30*time.Minute
}

// LastSeenString returns a human-readable last seen string
func (u *User) LastSeenString() string {
	if u.LastLogin == nil {
		return "Never"
	}

	duration := time.Since(*u.LastLogin)
	if duration < time.Minute {
		return "Just now"
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minute%s ago", int(duration.Minutes()), userPluralS(int(duration.Minutes())))
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hour%s ago", int(duration.Hours()), userPluralS(int(duration.Hours())))
	} else if duration < 7*24*time.Hour {
		return fmt.Sprintf("%d day%s ago", int(duration.Hours()/24), userPluralS(int(duration.Hours()/24)))
	} else {
		return u.LastLogin.Format("Jan 2, 2006")
	}
}

// AccountAge returns the account age in days
func (u *User) AccountAge() int {
	return int(time.Since(u.CreatedAt).Hours() / 24)
}

// ValidatePasswordResetToken validates a password reset token
func (u *User) ValidatePasswordResetToken(token string, expiry time.Duration) bool {
	// This is a simplified version. In a real implementation,
	// you would validate a JWT or similar token
	return token != "" && u.IsActive
}

// Clone returns a copy of the user (without password for security)
func (u *User) Clone() *User {
	if u == nil {
		return nil
	}

	var lastLogin *time.Time
	if u.LastLogin != nil {
		t := *u.LastLogin
		lastLogin = &t
	}

	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  "", // Don't clone password for security
		Role:      u.Role,
		IsActive:  u.IsActive,
		LastLogin: lastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// CloneWithPassword returns a copy of the user including password (for internal use)
func (u *User) CloneWithPassword() *User {
	if u == nil {
		return nil
	}

	var lastLogin *time.Time
	if u.LastLogin != nil {
		t := *u.LastLogin
		lastLogin = &t
	}

	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Password:  u.Password,
		Role:      u.Role,
		IsActive:  u.IsActive,
		LastLogin: lastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Password utility functions

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	if err := validatePasswordStrength(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to generate password hash: %w", err)
	}
	return string(hash), nil
}

// CheckPassword checks if a password matches the hash
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// validatePasswordStrength validates password strength requirements
func validatePasswordStrength(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password cannot exceed 128 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsNumber(r):
			hasNumber = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func userPluralS(n int) string {
	if n != 1 {
		return "s"
	}
	return ""
}

var (
	ErrInvalidUser = errors.New("invalid user")
	ErrInvalidUsername = errors.New("invalid username")
	ErrUserNotActive = errors.New("user is not active")
	ErrInvalidPassword = errors.New("invalid password")
	ErrPasswordMismatch = errors.New("password does not match")
	ErrWeakPassword = errors.New("password is too weak")
	ErrUserAlreadyActive = errors.New("user is already active")
	ErrUserAlreadyInactive = errors.New("user is already inactive")
)