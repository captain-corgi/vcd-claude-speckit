package domain

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo      UserRepository
	auditRepo     AuditLogRepository
	eventStore    EventStoreRepository
	eventDispatcher *EventDispatcher
}

// NewUserService creates a new user service
func NewUserService(
	userRepo UserRepository,
	auditRepo AuditLogRepository,
	eventStore EventStoreRepository,
	eventDispatcher *EventDispatcher,
) *UserService {
	return &UserService{
		userRepo:        userRepo,
		auditRepo:       auditRepo,
		eventStore:      eventStore,
		eventDispatcher: eventDispatcher,
	}
}

// CreateUser creates a new user with validation and audit logging
func (s *UserService) CreateUser(
	ctx context.Context,
	username, email, password string,
	role UserRole,
	userID string,
	ipAddress, userAgent string,
) (*User, error) {
	// Validate business rules
	if err := s.validateUserCreation(ctx, username, email); err != nil {
		return nil, err
	}

	// Create user
	user, err := NewUser(username, email, password, role)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Save to repository
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForUserCreation(user, userID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewUserLoggedInEvent(user.ID, username, ipAddress, userAgent)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return user, nil
}

// AuthenticateUser authenticates a user with username/password
func (s *UserService) AuthenticateUser(
	ctx context.Context,
	username, password string,
	ipAddress, userAgent string,
) (*User, error) {
	// Find user by username
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	if !user.Authenticate(password) {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	user.UpdateLastLogin()
	if err := s.userRepo.UpdateLastLogin(ctx, user.ID); err != nil {
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Create domain event
	event := NewUserLoggedInEvent(user.ID, username, ipAddress, userAgent)
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(
	ctx context.Context,
	userID uuid.UUID,
	currentPassword, newPassword string,
	ipAddress, userAgent string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Change password
	if err := user.ChangePassword(currentPassword, newPassword); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForPasswordChange(user, "self", "change", ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewUserPasswordChangedEvent(user.ID, user.Username, "self", "change")
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return nil
}

// ResetPassword resets a user's password (admin function)
func (s *UserService) ResetPassword(
	ctx context.Context,
	userID uuid.UUID,
	newPassword string,
	adminID string,
	ipAddress, userAgent string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Reset password
	if err := user.ResetPassword(newPassword); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForPasswordChange(user, adminID, "reset", ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	// Create domain event
	event := NewUserPasswordChangedEvent(user.ID, user.Username, adminID, "reset")
	if err := s.eventStore.SaveEvent(ctx, event); err != nil {
		fmt.Printf("Warning: failed to save event: %v\n", err)
	}

	// Dispatch event
	if err := s.eventDispatcher.Dispatch(event); err != nil {
		fmt.Printf("Warning: failed to dispatch event: %v\n", err)
	}

	return nil
}

// UpdateUserProfile updates user profile information
func (s *UserService) UpdateUserProfile(
	ctx context.Context,
	userID uuid.UUID,
	updates map[string]interface{},
	ipAddress, userAgent string,
) (*User, error) {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Apply updates
	changedFields, err := s.applyUserUpdates(user, updates)
	if err != nil {
		return nil, err
	}

	// Validate business rules
	if err := s.validateUserUpdate(ctx, user); err != nil {
		return nil, err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForUserUpdate(user, changedFields, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	return user, nil
}

// DeactivateUser deactivates a user account
func (s *UserService) DeactivateUser(
	ctx context.Context,
	userID uuid.UUID,
	adminID string,
	ipAddress, userAgent string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Deactivate user
	if err := user.Deactivate(); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForUserDeactivation(user, adminID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	return nil
}

// ActivateUser activates a user account
func (s *UserService) ActivateUser(
	ctx context.Context,
	userID uuid.UUID,
	adminID string,
	ipAddress, userAgent string,
) error {
	// Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return ErrUserNotFound
	}

	// Activate user
	if err := user.Activate(); err != nil {
		return err
	}

	// Update user
	if err := s.userRepo.Update(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Create audit log
	auditLog := s.createAuditLogForUserActivation(user, adminID, ipAddress, userAgent)
	if err := s.auditRepo.Create(ctx, auditLog); err != nil {
		fmt.Printf("Warning: failed to create audit log: %v\n", err)
	}

	return nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return s.userRepo.FindByUsername(ctx, username)
}

// ListUsers returns a paginated list of users
func (s *UserService) ListUsers(
	ctx context.Context,
	filter UserFilter,
	sort UserSort,
	pagination Pagination,
) (*UserResult, error) {
	return s.userRepo.List(ctx, filter, sort, pagination)
}

// SearchUsers searches for users by various criteria
func (s *UserService) SearchUsers(
	ctx context.Context,
	searchTerm string,
	role *UserRole,
	isActive *bool,
	pagination Pagination,
) (*UserResult, error) {
	filter := UserFilter{
		Search:   &searchTerm,
		Role:     role,
		IsActive: isActive,
	}

	sort := UserSort{
		Field:     SortByUsername,
		Direction: SortDirectionASC,
	}

	return s.userRepo.List(ctx, filter, sort, pagination)
}

// Validation methods

func (s *UserService) validateUserCreation(ctx context.Context, username, email string) error {
	// Check for duplicate username
	if exists, err := s.userRepo.ExistsByUsername(ctx, username); err != nil {
		return fmt.Errorf("failed to check username uniqueness: %w", err)
	} else if exists {
		return ErrUsernameAlreadyExists
	}

	// Check for duplicate email
	if exists, err := s.userRepo.ExistsByEmail(ctx, email); err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	} else if exists {
		return ErrEmailAlreadyExists
	}

	return nil
}

func (s *UserService) validateUserUpdate(ctx context.Context, user *User) error {
	// Check for username uniqueness if username was changed
	if exists, err := s.userRepo.ExistsByUsername(ctx, user.Username); err != nil {
		return fmt.Errorf("failed to check username uniqueness: %w", err)
	} else if exists {
		// Get existing user to check if it's the same user
		existing, err := s.userRepo.GetByID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to get existing user: %w", err)
		}
		if existing == nil || existing.Username != user.Username {
			return ErrUsernameAlreadyExists
		}
	}

	// Check for email uniqueness if email was changed
	if exists, err := s.userRepo.ExistsByEmail(ctx, user.Email); err != nil {
		return fmt.Errorf("failed to check email uniqueness: %w", err)
	} else if exists {
		// Get existing user to check if it's the same user
		existing, err := s.userRepo.GetByID(ctx, user.ID)
		if err != nil {
			return fmt.Errorf("failed to get existing user: %w", err)
		}
		if existing == nil || existing.Email != user.Email {
			return ErrEmailAlreadyExists
		}
	}

	return nil
}

// Helper methods

func (s *UserService) applyUserUpdates(user *User, updates map[string]interface{}) ([]string, error) {
	var changedFields []string

	for field, value := range updates {
		switch field {
		case "username":
			if val, ok := value.(string); ok {
				user.Username = val
				changedFields = append(changedFields, "username")
			}
		case "email":
			if val, ok := value.(string); ok {
				user.Email = val
				changedFields = append(changedFields, "email")
			}
		case "role":
			if val, ok := value.(string); ok {
				role, err := ParseUserRole(val)
				if err != nil {
					return nil, fmt.Errorf("invalid role: %w", err)
				}
				user.Role = role
				changedFields = append(changedFields, "role")
			}
		}
	}

	user.UpdatedAt = time.Now()
	return changedFields, nil
}

// Audit log creation methods

func (s *UserService) createAuditLogForUserCreation(user *User, userID, ipAddress, userAgent string) *AuditLog {
	newValues := map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"role":      string(user.Role),
		"isActive":  user.IsActive,
		"createdAt": user.CreatedAt,
	}

	auditLog, err := NewAuditLog(
		uuid.New(), // Using a new UUID for system operations
		OperationCreateUser,
		userID,
		nil,
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

func (s *UserService) createAuditLogForPasswordChange(user *User, changedBy, method, ipAddress, userAgent string) *AuditLog {
	newValues := map[string]interface{}{
		"passwordChanged": true,
		"changedBy":      changedBy,
		"method":         method,
	}

	auditLog, err := NewAuditLog(
		user.ID,
		OperationPasswordChange,
		user.ID.String(),
		nil,
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

func (s *UserService) createAuditLogForUserUpdate(user *User, changedFields []string, ipAddress, userAgent string) *AuditLog {
	newValues := map[string]interface{}{
		"username":  user.Username,
		"email":     user.Email,
		"role":      string(user.Role),
		"isActive":  user.IsActive,
		"updatedAt": user.UpdatedAt,
	}

	auditLog, err := NewAuditLog(
		user.ID,
		OperationUpdateUser,
		user.ID.String(),
		nil,
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

func (s *UserService) createAuditLogForUserDeactivation(user *User, adminID, ipAddress, userAgent string) *AuditLog {
	oldValues := map[string]interface{}{
		"isActive": true,
	}
	newValues := map[string]interface{}{
		"isActive": false,
	}

	auditLog, err := NewAuditLog(
		user.ID,
		OperationUpdateUser,
		adminID,
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

func (s *UserService) createAuditLogForUserActivation(user *User, adminID, ipAddress, userAgent string) *AuditLog {
	oldValues := map[string]interface{}{
		"isActive": false,
	}
	newValues := map[string]interface{}{
		"isActive": true,
	}

	auditLog, err := NewAuditLog(
		user.ID,
		OperationUpdateUser,
		adminID,
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

// GetInactiveUsers returns users who have been inactive for a certain period
func (s *UserService) GetInactiveUsers(ctx context.Context, inactiveSince time.Time) ([]*User, error) {
	return s.userRepo.GetInactiveUsers(ctx, inactiveSince)
}

// GetUserActivity returns user activity from audit logs
func (s *UserService) GetUserActivity(ctx context.Context, userID string, start, end time.Time) ([]*AuditLog, error) {
	return s.auditRepo.GetUserActivity(ctx, userID, start, end)
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)