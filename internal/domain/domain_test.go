package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestEmployeeStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   EmployeeStatus
		expected string
	}{
		{"Active status", EmployeeStatusActive, "ACTIVE"},
		{"Terminated status", EmployeeStatusTerminated, "TERMINATED"},
		{"On leave status", EmployeeStatusOnLeave, "ON_LEAVE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.status))
		})
	}
}

func TestEmployeeStatusValidation(t *testing.T) {
	tests := []struct {
		name    string
		status  EmployeeStatus
		isValid bool
	}{
		{"Valid active", EmployeeStatusActive, true},
		{"Valid terminated", EmployeeStatusTerminated, true},
		{"Valid on leave", EmployeeStatusOnLeave, true},
		{"Invalid empty", EmployeeStatus(""), false},
		{"Invalid unknown", EmployeeStatus("UNKNOWN"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.isValid, tt.status.IsValid())
		})
	}
}

func TestUserRole(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		expected string
	}{
		{"Admin role", UserRoleAdmin, "ADMIN"},
		{"Manager role", UserRoleManager, "MANAGER"},
		{"Viewer role", UserRoleViewer, "VIEWER"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, string(tt.role))
		})
	}
}

func TestUserRolePermissions(t *testing.T) {
	tests := []struct {
		name            string
		role            UserRole
		canAccessSalary bool
		canManageUsers  bool
	}{
		{"Admin permissions", UserRoleAdmin, true, true},
		{"Manager permissions", UserRoleManager, true, false},
		{"Viewer permissions", UserRoleViewer, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.canAccessSalary, tt.role.CanAccessSalary())
			assert.Equal(t, tt.canManageUsers, tt.role.CanManageUsers())
		})
	}
}

func TestAddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		address Address
		wantErr bool
	}{
		{
			name: "Valid address",
			address: Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
			wantErr: false,
		},
		{
			name: "Invalid empty street",
			address: Address{
				Street:     "",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "12345",
				Country:    "USA",
			},
			wantErr: true,
		},
		{
			name: "Invalid postal code",
			address: Address{
				Street:     "123 Main St",
				City:       "Anytown",
				State:      "CA",
				PostalCode: "!",
				Country:    "USA",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.address.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEmployeeCreation(t *testing.T) {
	address, err := NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	assert.NoError(t, err)

	hireDate := time.Now()
	employee, err := NewEmployee(
		"John",
		"Doe",
		"john.doe@example.com",
		"555-123-4567",
		"Engineering",
		"Software Engineer",
		hireDate,
		75000.00,
		nil,
		address,
	)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, employee.ID)
	assert.Equal(t, "John", employee.FirstName)
	assert.Equal(t, "Doe", employee.LastName)
	assert.Equal(t, "john.doe@example.com", employee.Email)
	assert.Equal(t, "Engineering", employee.Department)
	assert.Equal(t, "Software Engineer", employee.Position)
	assert.Equal(t, 75000.00, employee.Salary)
	assert.Equal(t, EmployeeStatusActive, employee.Status)
}

func TestEmployeeValidation(t *testing.T) {
	address, err := NewAddress("123 Main St", "Anytown", "CA", "12345", "USA")
	assert.NoError(t, err)

	hireDate := time.Now()

	tests := []struct {
		name        string
		firstName   string
		lastName    string
		email       string
		phone       string
		department  string
		position    string
		salary      float64
		expectError bool
	}{
		{
			name:        "Valid employee",
			firstName:   "John",
			lastName:    "Doe",
			email:       "john.doe@example.com",
			phone:       "555-123-4567",
			department:  "Engineering",
			position:    "Software Engineer",
			salary:      75000.00,
			expectError: false,
		},
		{
			name:        "Invalid first name",
			firstName:   "",
			lastName:    "Doe",
			email:       "john.doe@example.com",
			phone:       "555-123-4567",
			department:  "Engineering",
			position:    "Software Engineer",
			salary:      75000.00,
			expectError: true,
		},
		{
			name:        "Invalid email",
			firstName:   "John",
			lastName:    "Doe",
			email:       "invalid-email",
			phone:       "555-123-4567",
			department:  "Engineering",
			position:    "Software Engineer",
			salary:      75000.00,
			expectError: true,
		},
		{
			name:        "Negative salary",
			firstName:   "John",
			lastName:    "Doe",
			email:       "john.doe@example.com",
			phone:       "555-123-4567",
			department:  "Engineering",
			position:    "Software Engineer",
			salary:      -1000.00,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			employee, err := NewEmployee(
				tt.firstName,
				tt.lastName,
				tt.email,
				tt.phone,
				tt.department,
				tt.position,
				hireDate,
				tt.salary,
				nil,
				address,
			)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, employee)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, employee)
			}
		})
	}
}

func TestUserCreation(t *testing.T) {
	user, err := NewUser(
		"johndoe",
		"john.doe@example.com",
		"SecurePassword123!",
		UserRoleViewer,
	)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
	assert.Equal(t, "johndoe", user.Username)
	assert.Equal(t, "john.doe@example.com", user.Email)
	assert.Equal(t, UserRoleViewer, user.Role)
	assert.True(t, user.IsActive)
	assert.NotEmpty(t, user.Password) // Should be hashed
}

func TestUserAuthentication(t *testing.T) {
	user, err := NewUser(
		"johndoe",
		"john.doe@example.com",
		"SecurePassword123!",
		UserRoleViewer,
	)

	assert.NoError(t, err)

	// Test correct password
	assert.True(t, user.Authenticate("SecurePassword123!"))

	// Test incorrect password
	assert.False(t, user.Authenticate("wrongpassword"))

	// Test inactive user
	user.Deactivate()
	assert.False(t, user.Authenticate("SecurePassword123!"))
}

func TestUserPermissions(t *testing.T) {
	tests := []struct {
		name         string
		role         UserRole
		canManage    bool
		canDelete    bool
		canViewLogs  bool
	}{
		{"Admin permissions", UserRoleAdmin, true, true, true},
		{"Manager permissions", UserRoleManager, false, false, true},
		{"Viewer permissions", UserRoleViewer, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser("testuser", "test@example.com", "SecurePassword123!", tt.role)
			assert.NoError(t, err)

			assert.Equal(t, tt.canManage, user.CanManageUsers())
			assert.Equal(t, tt.canDelete, user.CanDeleteEmployees())
			assert.Equal(t, tt.canViewLogs, user.CanViewAuditLogs())
		})
	}
}

func TestAuditLogCreation(t *testing.T) {
	employeeID := uuid.New()
	userID := "user123"

	auditLog, err := NewAuditLog(
		employeeID,
		"CREATE_EMPLOYEE",
		userID,
		nil,
		map[string]interface{}{
			"name": "John Doe",
			"email": "john.doe@example.com",
		},
		"192.168.1.1",
		"Mozilla/5.0",
	)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, auditLog.ID)
	assert.Equal(t, employeeID, auditLog.EmployeeID)
	assert.Equal(t, "CREATE_EMPLOYEE", auditLog.Operation)
	assert.Equal(t, userID, auditLog.UserID)
	assert.Equal(t, "192.168.1.1", auditLog.IPAddress)
}

func TestAuditLogChangeDetection(t *testing.T) {
	oldValues := map[string]interface{}{
		"salary": 50000.0,
		"status": "ACTIVE",
	}

	newValues := map[string]interface{}{
		"salary": 60000.0,
		"status": "ACTIVE",
	}

	auditLog := &AuditLog{
		OldValues: oldValues,
		NewValues: newValues,
	}

	summary := auditLog.GetChangeSummary()
	assert.Contains(t, summary, "salary")
	assert.NotContains(t, summary, "status")
}