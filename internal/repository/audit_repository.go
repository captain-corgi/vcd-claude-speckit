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

// AuditLogRepository implements the domain repository interface using GORM
type AuditLogRepository struct {
	db *gorm.DB
}

// NewAuditLogRepository creates a new AuditLogRepository
func NewAuditLogRepository(db *gorm.DB) *AuditLogRepository {
	return &AuditLogRepository{db: db}
}

// Create creates a new audit log record
func (r *AuditLogRepository) Create(ctx context.Context, auditLog *domain.AuditLog) error {
	model := &AuditLog{}
	model.FromDomain(auditLog)

	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	auditLog.ID = model.ID

	return nil
}

// GetByID retrieves an audit log by ID
func (r *AuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	var model AuditLog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get audit log by ID: %w", err)
	}

	return model.ToDomain()
}

// FindByEmployeeID retrieves all audit logs for a specific employee
func (r *AuditLogRepository) FindByEmployeeID(ctx context.Context, employeeID uuid.UUID) ([]*domain.AuditLog, error) {
	var models []AuditLog
	err := r.db.WithContext(ctx).Where("employee_id = ?", employeeID).Order("timestamp DESC").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find audit logs by employee ID: %w", err)
	}

	result := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = log
	}

	return result, nil
}

// FindByOperation retrieves all audit logs for a specific operation
func (r *AuditLogRepository) FindByOperation(ctx context.Context, operation string) ([]*domain.AuditLog, error) {
	var models []AuditLog
	err := r.db.WithContext(ctx).Where("operation = ?", operation).Order("timestamp DESC").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find audit logs by operation: %w", err)
	}

	result := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = log
	}

	return result, nil
}

// FindByUserID retrieves all audit logs for a specific user
func (r *AuditLogRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.AuditLog, error) {
	var models []AuditLog
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("timestamp DESC").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find audit logs by user ID: %w", err)
	}

	result := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = log
	}

	return result, nil
}

// FindByTimeRange retrieves audit logs within a specific time range
func (r *AuditLogRepository) FindByTimeRange(ctx context.Context, start, end time.Time) ([]*domain.AuditLog, error) {
	var models []AuditLog
	err := r.db.WithContext(ctx).Where("timestamp BETWEEN ? AND ?", start, end).Order("timestamp DESC").Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find audit logs by time range: %w", err)
	}

	result := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = log
	}

	return result, nil
}

// Find retrieves audit logs with filtering, sorting, and pagination
func (r *AuditLogRepository) Find(ctx context.Context, filter domain.AuditLogFilter, sort domain.AuditLogSort, pagination domain.Pagination) (*domain.AuditLogResult, error) {
	query := r.db.WithContext(ctx).Model(&AuditLog{})

	// Apply filters
	if filter.EmployeeID != nil {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}
	if filter.Operation != nil {
		query = query.Where("operation = ?", *filter.Operation)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.IPAddress != nil {
		query = query.Where("ip_address = ?", *filter.IPAddress)
	}
	if filter.FromTime != nil {
		query = query.Where("timestamp >= ?", *filter.FromTime)
	}
	if filter.ToTime != nil {
		query = query.Where("timestamp <= ?", *filter.ToTime)
	}
	if len(filter.Operations) > 0 {
		query = query.Where("operation IN ?", filter.Operations)
	}

	// Apply sorting
	sortField := map[string]string{
		string(domain.SortByAuditLogID):        "id",
		string(domain.SortByAuditLogTimestamp): "timestamp",
		string(domain.SortByAuditLogOperation): "operation",
		string(domain.SortByAuditLogUserID):    "user_id",
		string(domain.SortByAuditLogEmployeeID): "employee_id",
	}

	orderBy := sortField[string(sort.Field)]
	if orderBy == "" {
		orderBy = "timestamp"
	}

	sortDirection := domain.SortDirectionASC
	if sort.Direction == domain.SortDirectionDESC {
		sortDirection = domain.SortDirectionDESC
	}

	query = query.Order(orderBy + " " + string(sortDirection))

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count audit logs: %w", err)
	}

	// Apply pagination
	var models []AuditLog
	offset := pagination.GetOffset()
	query = query.Offset(offset).Limit(pagination.PageSize)

	if err := query.Find(&models).Error; err != nil {
		return nil, fmt.Errorf("failed to find audit logs: %w", err)
	}

	// Convert to domain entities
	auditLogs := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		auditLogs[i] = log
	}

	// Calculate pagination metadata
	hasNext := int64(offset+pagination.PageSize) < total
	hasPrev := offset > 0

	result := &domain.AuditLogResult{
		AuditLogs: auditLogs,
		Total:     total,
		Page:      pagination.Page,
		PageSize:  pagination.PageSize,
		HasNext:   hasNext,
		HasPrev:   hasPrev,
	}

	return result, nil
}

// Count counts audit logs with the specified filter
func (r *AuditLogRepository) Count(ctx context.Context, filter domain.AuditLogFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&AuditLog{})

	// Apply filters (same as in Find method)
	if filter.EmployeeID != nil {
		query = query.Where("employee_id = ?", *filter.EmployeeID)
	}
	if filter.Operation != nil {
		query = query.Where("operation = ?", *filter.Operation)
	}
	if filter.UserID != nil {
		query = query.Where("user_id = ?", *filter.UserID)
	}
	if filter.IPAddress != nil {
		query = query.Where("ip_address = ?", *filter.IPAddress)
	}
	if filter.FromTime != nil {
		query = query.Where("timestamp >= ?", *filter.FromTime)
	}
	if filter.ToTime != nil {
		query = query.Where("timestamp <= ?", *filter.ToTime)
	}
	if len(filter.Operations) > 0 {
		query = query.Where("operation IN ?", filter.Operations)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return count, nil
}

// GetOperationsSummary returns a summary of operations within a time range
func (r *AuditLogRepository) GetOperationsSummary(ctx context.Context, start, end time.Time) (map[string]int64, error) {
	var results []struct {
	Operation string
	Count     int64
	}

	err := r.db.WithContext(ctx).
		Model(&AuditLog{}).
		Select("operation, count(*) as count").
		Where("timestamp BETWEEN ? AND ?", start, end).
		Group("operation").
		Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get operations summary: %w", err)
	}

	summary := make(map[string]int64)
	for _, result := range results {
		summary[result.Operation] = result.Count
	}

	return summary, nil
}

// GetUserActivity retrieves audit logs for a specific user within a time range
func (r *AuditLogRepository) GetUserActivity(ctx context.Context, userID string, start, end time.Time) ([]*domain.AuditLog, error) {
	var models []AuditLog
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND timestamp BETWEEN ? AND ?", userID, start, end).
		Order("timestamp DESC").
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get user activity: %w", err)
	}

	result := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = log
	}

	return result, nil
}

// CreateMany creates multiple audit log records
func (r *AuditLogRepository) CreateMany(ctx context.Context, auditLogs []*domain.AuditLog) error {
	models := make([]AuditLog, len(auditLogs))
	for i, log := range auditLogs {
		models[i] = AuditLog{}
		models[i].FromDomain(log)
	}

	err := r.db.WithContext(ctx).Create(&models).Error
	if err != nil {
		return fmt.Errorf("failed to create multiple audit logs: %w", err)
	}

	// Update domain entities with generated IDs
	for i, model := range models {
		auditLogs[i].ID = model.ID
	}

	return nil
}

// DeleteOldLogs deletes audit logs older than the specified time
func (r *AuditLogRepository) DeleteOldLogs(ctx context.Context, olderThan time.Time) error {
	result := r.db.WithContext(ctx).
		Where("timestamp < ?", olderThan).
		Delete(&AuditLog{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete old audit logs: %w", result.Error)
	}

	return nil
}

// GetRecentAuditLogs retrieves the most recent audit logs
func (r *AuditLogRepository) GetRecentAuditLogs(ctx context.Context, limit int) ([]*domain.AuditLog, error) {
	var models []AuditLog
	err := r.db.WithContext(ctx).
		Order("timestamp DESC").
		Limit(limit).
		Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get recent audit logs: %w", err)
	}

	result := make([]*domain.AuditLog, len(models))
	for i, model := range models {
		log, err := model.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert model to domain: %w", err)
		}
		result[i] = log
	}

	return result, nil
}

// GetAuditLogStats returns statistics about audit logs
func (r *AuditLogRepository) GetAuditLogStats(ctx context.Context) (*AuditLogStats, error) {
	var total int64
	var todayCount int64
	var thisWeekCount int64
	var thisMonthCount int64

	// Total count
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to get total audit log count: %w", err)
	}

	// Today's count
	today := time.Now().Truncate(24 * time.Hour)
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Where("timestamp >= ?", today).Count(&todayCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get today's audit log count: %w", err)
	}

	// This week's count
	weekStart := today.AddDate(0, 0, -int(today.Weekday()))
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Where("timestamp >= ?", weekStart).Count(&thisWeekCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get this week's audit log count: %w", err)
	}

	// This month's count
	monthStart := today.AddDate(0, -1, 0)
	if err := r.db.WithContext(ctx).Model(&AuditLog{}).
		Where("timestamp >= ?", monthStart).Count(&thisMonthCount).Error; err != nil {
		return nil, fmt.Errorf("failed to get this month's audit log count: %w", err)
	}

	return &AuditLogStats{
		Total:       total,
		TodayCount:  todayCount,
		ThisWeekCount: thisWeekCount,
		ThisMonthCount: thisMonthCount,
	}, nil
}

// AuditLogStats contains statistics about audit logs
type AuditLogStats struct {
	Total       int64
	TodayCount  int64
	ThisWeekCount int64
	ThisMonthCount int64
}