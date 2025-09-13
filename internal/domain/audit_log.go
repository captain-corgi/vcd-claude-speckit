package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// AuditLog represents an audit log entry for tracking operations
type AuditLog struct {
	ID         uuid.UUID
	EmployeeID uuid.UUID
	Operation  string
	UserID     string // Can be user ID or system identifier
	Timestamp  time.Time
	OldValues  map[string]interface{}
	NewValues  map[string]interface{}
	IPAddress  string
	UserAgent  string
}

// NewAuditLog creates a new AuditLog entry
func NewAuditLog(employeeID uuid.UUID, operation, userID string, oldValues, newValues map[string]interface{}, ipAddress, userAgent string) (*AuditLog, error) {
	auditLog := &AuditLog{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		Operation:  operation,
		UserID:     userID,
		Timestamp:  time.Now(),
		OldValues:  oldValues,
		NewValues:  newValues,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}

	if err := auditLog.Validate(); err != nil {
		return nil, err
	}

	return auditLog, nil
}

// Validate validates the audit log entry
func (a *AuditLog) Validate() error {
	if a == nil {
		return errors.New("audit log cannot be nil")
	}

	// Validate ID
	if a.ID == uuid.Nil {
		return errors.New("audit log ID cannot be empty")
	}

	// Validate EmployeeID
	if a.EmployeeID == uuid.Nil {
		return errors.New("employee ID cannot be empty")
	}

	// Validate Operation
	if err := validateOperation(a.Operation); err != nil {
		return err
	}

	// Validate UserID
	if a.UserID == "" {
		return errors.New("user ID cannot be empty")
	}

	// Validate Timestamp
	if a.Timestamp.IsZero() {
		return errors.New("timestamp cannot be empty")
	}

	// Validate IPAddress (required)
	if a.IPAddress == "" {
		return errors.New("IP address cannot be empty")
	}

	// Validate IPAddress format
	if err := validateIPAddress(a.IPAddress); err != nil {
		return fmt.Errorf("invalid IP address: %w", err)
	}

	// Validate values are not both empty
	if len(a.OldValues) == 0 && len(a.NewValues) == 0 {
		return errors.New("at least one of old values or new values must be provided")
	}

	return nil
}

// validateOperation validates the operation string
func validateOperation(operation string) error {
	if operation == "" {
		return errors.New("operation cannot be empty")
	}
	if len(operation) > 50 {
		return errors.New("operation cannot exceed 50 characters")
	}

	// Operation should contain only letters, numbers, underscores, and colons
	if !isValidOperationFormat(operation) {
		return errors.New("operation contains invalid characters")
	}

	return nil
}

// isValidOperationFormat checks if the operation format is valid
func isValidOperationFormat(operation string) bool {
	for _, r := range operation {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' || r == ':') {
			return false
		}
	}
	return true
}

// validateIPAddress validates the IP address format
func validateIPAddress(ip string) error {
	if ip == "" {
		return errors.New("IP address cannot be empty")
	}
	if len(ip) > 45 { // IPv6 address can be up to 45 characters
		return errors.New("IP address is too long")
	}

	// This is a basic validation. In a real implementation,
	// you might want to use a proper IP address validation library
	// For now, we'll accept common formats
	if len(ip) < 7 { // Minimum valid IP length (e.g., "1.1.1.1")
		return errors.New("IP address is too short")
	}

	return nil
}

// GetChangeSummary returns a summary of the changes
func (a *AuditLog) GetChangeSummary() string {
	if len(a.OldValues) == 0 {
		return fmt.Sprintf("Created: %s", a.getFieldsSummary(a.NewValues))
	}
	if len(a.NewValues) == 0 {
		return fmt.Sprintf("Deleted: %s", a.getFieldsSummary(a.OldValues))
	}
	return fmt.Sprintf("Updated: %s", a.getChangedFieldsSummary())
}

// getFieldsSummary returns a summary of the fields in the values map
func (a *AuditLog) getFieldsSummary(values map[string]interface{}) string {
	if len(values) == 0 {
		return "no fields"
	}

	fields := make([]string, 0, len(values))
	for field := range values {
		fields = append(fields, field)
	}

	if len(fields) <= 3 {
		return fmt.Sprintf("%s", joinFields(fields))
	}
	return fmt.Sprintf("%s and %d more", joinFields(fields[:3]), len(fields)-3)
}

// getChangedFieldsSummary returns a summary of changed fields
func (a *AuditLog) getChangedFieldsSummary() string {
	changedFields := make([]string, 0)

	for field := range a.OldValues {
		if newValue, exists := a.NewValues[field]; exists {
			oldValue := a.OldValues[field]
			if !valuesEqual(oldValue, newValue) {
				changedFields = append(changedFields, field)
			}
		}
	}

	// Check for new fields
	for field := range a.NewValues {
		if _, exists := a.OldValues[field]; !exists {
			changedFields = append(changedFields, field)
		}
	}

	if len(changedFields) == 0 {
		return "no changes"
	}

	if len(changedFields) <= 3 {
		return joinFields(changedFields)
	}
	return fmt.Sprintf("%s and %d more", joinFields(changedFields[:3]), len(changedFields)-3)
}

// valuesEqual checks if two values are equal
func valuesEqual(a, b interface{}) bool {
	// Handle basic types
	if a == b {
		return true
	}

	// Handle string representations
	aStr, aOK := a.(string)
	bStr, bOK := b.(string)
	if aOK && bOK {
		return aStr == bStr
	}

	// Handle numeric types
	if a != nil && b != nil {
		return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
	}

	return false
}

// joinFields joins field names with commas and "and"
func joinFields(fields []string) string {
	if len(fields) == 0 {
		return ""
	}
	if len(fields) == 1 {
		return fields[0]
	}
	if len(fields) == 2 {
		return fmt.Sprintf("%s and %s", fields[0], fields[1])
	}

	result := ""
	for i, field := range fields {
		if i > 0 {
			result += ", "
		}
		if i == len(fields)-1 {
			result += "and "
		}
		result += field
	}
	return result
}

// IsCreation checks if this audit log represents a creation operation
func (a *AuditLog) IsCreation() bool {
	return len(a.OldValues) == 0 && len(a.NewValues) > 0
}

// IsDeletion checks if this audit log represents a deletion operation
func (a *AuditLog) IsDeletion() bool {
	return len(a.NewValues) == 0 && len(a.OldValues) > 0
}

// IsUpdate checks if this audit log represents an update operation
func (a *AuditLog) IsUpdate() bool {
	return len(a.OldValues) > 0 && len(a.NewValues) > 0
}

// GetFieldChange returns the old and new values for a specific field
func (a *AuditLog) GetFieldChange(field string) (oldValue, newValue interface{}, changed bool) {
	oldValue, oldExists := a.OldValues[field]
	newValue, newExists := a.NewValues[field]

	if !oldExists && !newExists {
		return nil, nil, false
	}

	if !oldExists {
		return nil, newValue, true
	}

	if !newExists {
		return oldValue, nil, true
	}

	return oldValue, newValue, !valuesEqual(oldValue, newValue)
}

// GetChangedFields returns a list of fields that were changed
func (a *AuditLog) GetChangedFields() []string {
	changedFields := make([]string, 0)

	for field := range a.OldValues {
		if newValue, exists := a.NewValues[field]; exists {
			oldValue := a.OldValues[field]
			if !valuesEqual(oldValue, newValue) {
				changedFields = append(changedFields, field)
			}
		}
	}

	// Check for new fields
	for field := range a.NewValues {
		if _, exists := a.OldValues[field]; !exists {
			changedFields = append(changedFields, field)
		}
	}

	return changedFields
}

// ToJSON converts the audit log to JSON
func (a *AuditLog) ToJSON() (string, error) {
	data := map[string]interface{}{
		"id":         a.ID,
		"employeeId": a.EmployeeID,
		"operation":  a.Operation,
		"userId":     a.UserID,
		"timestamp":  a.Timestamp,
		"oldValues":  a.OldValues,
		"newValues":  a.NewValues,
		"ipAddress":  a.IPAddress,
		"userAgent":  a.UserAgent,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal audit log to JSON: %w", err)
	}

	return string(jsonData), nil
}

// Clone returns a copy of the audit log
func (a *AuditLog) Clone() *AuditLog {
	if a == nil {
		return nil
	}

	oldValues := make(map[string]interface{})
	for k, v := range a.OldValues {
		oldValues[k] = v
	}

	newValues := make(map[string]interface{})
	for k, v := range a.NewValues {
		newValues[k] = v
	}

	return &AuditLog{
		ID:         a.ID,
		EmployeeID: a.EmployeeID,
		Operation:  a.Operation,
		UserID:     a.UserID,
		Timestamp:  a.Timestamp,
		OldValues:  oldValues,
		NewValues:  newValues,
		IPAddress:  a.IPAddress,
		UserAgent:  a.UserAgent,
	}
}

// Standard operations for audit logging
const (
	OperationCreateEmployee      = "employee:create"
	OperationUpdateEmployee      = "employee:update"
	OperationDeleteEmployee      = "employee:delete"
	OperationChangeStatus        = "employee:change_status"
	OperationUpdateSalary        = "employee:update_salary"
	OperationUpdatePosition      = "employee:update_position"
	OperationUpdateAddress       = "employee:update_address"
	OperationSetManager          = "employee:set_manager"
	OperationCreateUser          = "user:create"
	OperationUpdateUser          = "user:update"
	OperationDeleteUser          = "user:delete"
	OperationUserLogin           = "user:login"
	OperationUserLogout          = "user:logout"
	OperationPasswordChange      = "user:password_change"
	OperationPasswordReset       = "user:password_reset"
	OperationSystemAction        = "system:action"
)

var (
	ErrInvalidAuditLog = errors.New("invalid audit log")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrInvalidIPAddress = errors.New("invalid IP address")
)