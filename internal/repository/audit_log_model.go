package repository

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"employee-management-system/internal/domain"
)

// AuditLog represents the GORM model for the AuditLog entity
type AuditLog struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	EmployeeID uuid.UUID `gorm:"type:uuid;not null;index"`
	Operation  string    `gorm:"size:50;not null;index"`
	UserID     string    `gorm:"size:255;not null"`
	Timestamp  time.Time `gorm:"not null;index"`

	// JSON fields for storing old and new values as JSON
	OldValues  string `gorm:"type:text;not null"`
	NewValues  string `gorm:"type:text;not null"`

	IPAddress  string `gorm:"size:45;not null;index"`
	UserAgent  string `gorm:"size:500"`

	CreatedAt time.Time
}

// TableName returns the table name for the AuditLog model
func (AuditLog) TableName() string {
	return "audit_logs"
}

// ToDomain converts the GORM model to the domain entity
func (a *AuditLog) ToDomain() (*domain.AuditLog, error) {
	// Parse JSON fields
	oldValues, err := jsonToMap(a.OldValues)
	if err != nil {
		return nil, err
	}

	newValues, err := jsonToMap(a.NewValues)
	if err != nil {
		return nil, err
	}

	return domain.NewAuditLog(
		a.EmployeeID,
		a.Operation,
		a.UserID,
		oldValues,
		newValues,
		a.IPAddress,
		a.UserAgent,
	)
}

// FromDomain converts the domain entity to the GORM model
func (a *AuditLog) FromDomain(log *domain.AuditLog) {
	a.ID = log.ID
	a.EmployeeID = log.EmployeeID
	a.Operation = log.Operation
	a.UserID = log.UserID
	a.Timestamp = log.Timestamp
	a.IPAddress = log.IPAddress
	a.UserAgent = log.UserAgent
	a.CreatedAt = log.Timestamp

	// Convert maps to JSON strings
	var err error
	a.OldValues, err = mapToJSON(log.OldValues)
	if err != nil {
		// Log error but don't fail the conversion
		a.OldValues = "{}"
	}

	a.NewValues, err = mapToJSON(log.NewValues)
	if err != nil {
		// Log error but don't fail the conversion
		a.NewValues = "{}"
	}
}

// BeforeCreate is called before creating a record
func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	if a.Timestamp.IsZero() {
		a.Timestamp = time.Now()
	}
	a.CreatedAt = time.Now()
	return nil
}

// Helper function to convert JSON string to map
func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	var result map[string]interface{}
	if jsonStr == "" {
		return make(map[string]interface{}), nil
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Helper function to convert map to JSON string
func mapToJSON(data map[string]interface{}) (string, error) {
	if data == nil {
		return "{}", nil
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// FindByOperationType is a helper method to find audit logs by operation type
func FindByOperationType(tx *gorm.DB, operation string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := tx.Where("operation = ?", operation).Limit(limit).Order("timestamp DESC").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// FindByEmployeeIDRange is a helper method to find audit logs by employee ID and date range
func FindByEmployeeIDRange(tx *gorm.DB, employeeID uuid.UUID, startTime, endTime time.Time) ([]AuditLog, error) {
	var logs []AuditLog
	err := tx.Where("employee_id = ? AND timestamp BETWEEN ? AND ?", employeeID, startTime, endTime).
		Order("timestamp DESC").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}

// FindByUserID is a helper method to find audit logs by user ID
func FindByUserID(tx *gorm.DB, userID string, limit int) ([]AuditLog, error) {
	var logs []AuditLog
	err := tx.Where("user_id = ?", userID).Limit(limit).Order("timestamp DESC").Find(&logs).Error
	if err != nil {
		return nil, err
	}
	return logs, nil
}