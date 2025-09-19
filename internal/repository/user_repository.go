package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
)

// UserRepository implements the domain repository interface using GORM
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user record
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	model := &User{}
	model.FromDomain(user)

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt

	return nil
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var model User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return model.ToDomain()
}

// Update updates an existing user record
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	model := &User{}
	model.FromDomain(user)

	err := r.db.WithContext(ctx).Where("id = ?", user.ID).Save(model).Error
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	user.UpdatedAt = model.UpdatedAt

	return nil
}

// Delete deletes a user record
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&User{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// FindByUsername retrieves a user by username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var model User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}

	return model.ToDomain()
}

// FindByEmail retrieves a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var model User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return model.ToDomain()
}

// FindByRole retrieves all users with the specified role
func (r *UserRepository) FindByRole(ctx context.Context, role domain.UserRole) ([]*domain.User, error) {
	var models []User
	err := r.db.WithContext(ctx).Where("role = ?", string(role)).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users by role: %w", err)
	}

	result := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = user
	}

	return result, nil
}

// FindByActiveStatus retrieves users by their active status
func (r *UserRepository) FindByActiveStatus(ctx context.Context, isActive bool) ([]*domain.User, error) {
	var models []User
	err := r.db.WithContext(ctx).Where("is_active = ?", isActive).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users by active status: %w", err)
	}

	result := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = user
	}

	return result, nil
}

// List retrieves users with filtering, sorting, and pagination
func (r *UserRepository) List(ctx context.Context, filter domain.UserFilter, sort domain.UserSort, pagination domain.Pagination) (*domain.UserResult, error) {
	query := r.db.WithContext(ctx).Model(&User{})

	// Apply filters
	if filter.Role != nil {
		query = query.Where("role = ?", string(*filter.Role))
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != nil {
		searchTerm := *filter.Search
		query = query.Where("username LIKE ? OR email LIKE ? OR role LIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}
	if filter.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *filter.CreatedFrom)
	}
	if filter.CreatedTo != nil {
		query = query.Where("created_at <= ?", *filter.CreatedTo)
	}
	if filter.LastLoginFrom != nil {
		query = query.Where("last_login >= ?", *filter.LastLoginFrom)
	}
	if filter.LastLoginTo != nil {
		query = query.Where("last_login <= ?", *filter.LastLoginTo)
	}

	// Apply sorting
	sortField := map[string]string{
		string(domain.SortByUserID):        "id",
		string(domain.SortByUsername):      "username",
		string(domain.SortByUserEmail):     "email",
		string(domain.SortByUserRole):     "role",
		string(domain.SortByUserActive):   "is_active",
		string(domain.SortByUserCreated):  "created_at",
		string(domain.SortByUserLastLogin): "last_login",
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
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	// Apply pagination
	var models []User
	offset := pagination.GetOffset()
	query = query.Offset(offset).Limit(pagination.PageSize)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	// Convert to domain entities
	users := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		users[i] = user
	}

	// Calculate pagination metadata
	hasNext := int64(offset+pagination.PageSize) < total
	hasPrev := offset > 0

	result := &domain.UserResult{
		Users:    users,
		Total:    total,
		Page:     pagination.Page,
		PageSize: pagination.PageSize,
		HasNext:  hasNext,
		HasPrev:  hasPrev,
	}

	return result, nil
}

// Count counts users with the specified filter
func (r *UserRepository) Count(ctx context.Context, filter domain.UserFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&User{})

	// Apply filters (same as in List method)
	if filter.Role != nil {
		query = query.Where("role = ?", string(*filter.Role))
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	if filter.Search != nil {
		searchTerm := *filter.Search
		query = query.Where("username LIKE ? OR email LIKE ? OR role LIKE ?",
			"%"+searchTerm+"%", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}
	if filter.CreatedFrom != nil {
		query = query.Where("created_at >= ?", *filter.CreatedFrom)
	}
	if filter.CreatedTo != nil {
		query = query.Where("created_at <= ?", *filter.CreatedTo)
	}
	if filter.LastLoginFrom != nil {
		query = query.Where("last_login >= ?", *filter.LastLoginFrom)
	}
	if filter.LastLoginTo != nil {
		query = query.Where("last_login <= ?", *filter.LastLoginTo)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// ExistsByUsername checks if a user with the given username exists
func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists by username: %w", err)
	}
	return count > 0, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists by email: %w", err)
	}
	return count > 0, nil
}

// ExistsByID checks if a user with the given ID exists
func (r *UserRepository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists by ID: %w", err)
	}
	return count > 0, nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("last_login", now).Error
}

// GetInactiveUsers retrieves users who have been inactive since the specified time
func (r *UserRepository) GetInactiveUsers(ctx context.Context, inactiveSince time.Time) ([]*domain.User, error) {
	var models []User
	err := r.db.WithContext(ctx).
		Where("last_login < ? OR (last_login IS NULL AND created_at < ?)", inactiveSince, inactiveSince).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get inactive users: %w", err)
	}

	result := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = user
	}

	return result, nil
}

// FindUsersWithoutLastLogin retrieves users who have never logged in
func (r *UserRepository) FindUsersWithoutLastLogin(ctx context.Context) ([]*domain.User, error) {
	var models []User
	err := r.db.WithContext(ctx).Where("last_login IS NULL").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find users without last login: %w", err)
	}

	result := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = user
	}

	return result, nil
}

// GetMostActiveUsers retrieves users with the most recent activity
func (r *UserRepository) GetMostActiveUsers(ctx context.Context, limit int) ([]*domain.User, error) {
	var models []User
	err := r.db.WithContext(ctx).
		Where("last_login IS NOT NULL").
		Order("last_login DESC").
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get most active users: %w", err)
	}

	result := make([]*domain.User, len(models))
	for i, model := range models {
		user, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = user
	}

	return result, nil
}

// CreateMany creates multiple user records
func (r *UserRepository) CreateMany(ctx context.Context, users []*domain.User) error {
	models := make([]User, len(users))
	for i, user := range users {
		models[i] = User{}
		models[i].FromDomain(user)
	}

	err := r.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return fmt.Errorf("failed to create multiple users: %w", err)
	}

	// Update domain entities with generated IDs and timestamps
	for i, model := range models {
		users[i].ID = model.ID
		users[i].CreatedAt = model.CreatedAt
		users[i].UpdatedAt = model.UpdatedAt
	}

	return nil
}

// CountUsersByRole counts users by their role
func (r *UserRepository) CountUsersByRole(ctx context.Context) (map[string]int64, error) {
	var results []struct {
	Role string
	Count int64
	}

	err := r.db.WithContext(ctx).
		Model(&User{}).
		Select("role, count(*) as count").
		Group("role").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to count users by role: %w", err)
	}

	counts := make(map[string]int64)
	for _, result := range results {
		counts[result.Role] = result.Count
	}

	return counts, nil
}