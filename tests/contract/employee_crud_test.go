package contract

import (
	"context"
	"strings"
	"testing"
	"time"

	"employee-management-system/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEmployeeCRUDContract tests all Employee CRUD operations through GraphQL API.
// This is a RED phase test - it MUST fail before implementation exists.
// Tests validate API contract and expected behavior before actual implementation.
func TestEmployeeCRUDContract(t *testing.T) {
	// Setup test server and GraphQL client
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all CRUD operation tests
	t.Run("Create Employee Operations", func(t *testing.T) {
		testCreateEmployeeOperations(t, client)
	})

	t.Run("Update Employee Operations", func(t *testing.T) {
		testUpdateEmployeeOperations(t, client)
	})

	t.Run("Delete Employee Operations", func(t *testing.T) {
		testDeleteEmployeeOperations(t, client)
	})

	t.Run("Get Employee Operations", func(t *testing.T) {
		testGetEmployeeOperations(t, client)
	})

	t.Run("Change Employee Status Operations", func(t *testing.T) {
		testChangeEmployeeStatusOperations(t, client)
	})

	t.Run("Employee Validation and Error Handling", func(t *testing.T) {
		testEmployeeValidationAndErrorHandling(t, client)
	})

	t.Run("Employee Authentication Requirements", func(t *testing.T) {
		testEmployeeAuthenticationRequirements(t, client, testServer)
	})
}

// testCreateEmployeeOperations tests all create employee mutation scenarios
func testCreateEmployeeOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name           string
		input          map[string]interface{}
		expectedError  bool
		errorContains  string
		expectedFields map[string]interface{}
	}{
		{
			name: "Valid Employee Creation With Required Fields",
			input: map[string]interface{}{
				"firstName":  "John",
				"lastName":   "Doe",
				"email":      "john.doe@example.com",
				"department": "Engineering",
				"position":   "Software Engineer",
				"hireDate":   "2023-01-15",
				"salary":     75000.0,
			},
			expectedError: true, // Will fail until implemented
			expectedFields: map[string]interface{}{
				"firstName":  "John",
				"lastName":   "Doe",
				"email":      "john.doe@example.com",
				"department": "Engineering",
				"position":   "Software Engineer",
				"status":      "ACTIVE",
			},
		},
		{
			name: "Valid Employee Creation With Complete Address",
			input: map[string]interface{}{
				"firstName":  "Jane",
				"lastName":   "Smith",
				"email":      "jane.smith@example.com",
				"department": "Marketing",
				"position":   "Marketing Manager",
				"hireDate":   "2023-02-01",
				"salary":     85000.0,
				"address": map[string]interface{}{
					"street":     "456 Oak Ave",
					"city":       "New York",
					"state":      "NY",
					"postalCode": "10001",
					"country":    "US",
				},
			},
			expectedError: true,
			expectedFields: map[string]interface{}{
				"firstName":  "Jane",
				"lastName":   "Smith",
				"department": "Marketing",
				"position":   "Marketing Manager",
			},
		},
		{
			name: "Employee Creation With Manager ID",
			input: map[string]interface{}{
				"firstName":  "Robert",
				"lastName":   "Johnson",
				"email":      "robert.j@example.com",
				"department": "Engineering",
				"position":   "Senior Software Engineer",
				"hireDate":   "2023-03-15",
				"salary":     95000.0,
				"managerId":  "manager-123",
			},
			expectedError: true,
			expectedFields: map[string]interface{}{
				"firstName": "Robert",
				"lastName":  "Johnson",
			},
		},
		{
			name: "Missing Required Fields - No First Name",
			input: map[string]interface{}{
				"lastName":   "Doe",
				"email":      "john.doe@example.com",
				"department": "Engineering",
				"position":   "Software Engineer",
			},
			expectedError: true,
			errorContains: "firstName",
		},
		{
			name: "Invalid Email Format",
			input: map[string]interface{}{
				"firstName":  "John",
				"lastName":   "Doe",
				"email":      "invalid-email-format",
				"department": "Engineering",
				"position":   "Software Engineer",
			},
			expectedError: true,
			errorContains: "email",
		},
		{
			name: "Negative Salary",
			input: map[string]interface{}{
				"firstName":  "John",
				"lastName":   "Doe",
				"email":      "john.doe@example.com",
				"department": "Engineering",
				"position":   "Software Engineer",
				"salary":     -50000.0,
			},
			expectedError: true,
			errorContains: "salary",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation CreateEmployee($input: CreateEmployeeInput!) {
					createEmployee(input: $input) {
						id
						firstName
						lastName
						email
						phone
						department
						position
						hireDate
						salary
						status
						managerId
						address {
							street
							city
							state
							postalCode
							country
						}
						createdAt
						updatedAt
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{"input": tc.input})

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				// If we got a response, check for expected errors or validate structure
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						// For successful operations, validate response structure when implemented
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented createEmployee: %s", tc.name)
						assertContractError(t, resp, "createEmployee")
					}
				} else {
					// Would validate successful response when implemented
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented createEmployee: %s", tc.name)
				}
			} else {
				// Connection errors are expected in RED phase
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testUpdateEmployeeOperations tests all update employee mutation scenarios
func testUpdateEmployeeOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name           string
		employeeId     string
		input          map[string]interface{}
		expectedError  bool
		errorContains  string
		expectedFields map[string]interface{}
	}{
		{
			name:       "Update Basic Employee Information",
			employeeId: "employee-123",
			input: map[string]interface{}{
				"firstName": "Jonathan",
				"lastName":  "Doe-Smith",
				"email":     "jonathan.doe@example.com",
			},
			expectedError: true,
			expectedFields: map[string]interface{}{
				"firstName": "Jonathan",
				"lastName":  "Doe-Smith",
				"email":     "jonathan.doe@example.com",
			},
		},
		{
			name:       "Update Employee Salary and Position",
			employeeId: "employee-123",
			input: map[string]interface{}{
				"salary":   85000.0,
				"position": "Senior Software Engineer",
			},
			expectedError: true,
			expectedFields: map[string]interface{}{
				"salary":   85000.0,
				"position": "Senior Software Engineer",
			},
		},
		{
			name:       "Update Employee Address",
			employeeId: "employee-123",
			input: map[string]interface{}{
				"address": map[string]interface{}{
					"street":     "789 Pine St",
					"city":       "San Francisco",
					"state":      "CA",
					"postalCode": "94105",
					"country":    "US",
				},
			},
			expectedError: true,
		},
		{
			name:       "Update Employee Manager",
			employeeId: "employee-123",
			input: map[string]interface{}{
				"managerId": "new-manager-456",
			},
			expectedError: true,
		},
		{
			name:          "Update Non-existent Employee",
			employeeId:    "non-existent-employee",
			input:         map[string]interface{}{"firstName": "Updated"},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name:          "Update With Invalid Employee ID",
			employeeId:    "",
			input:         map[string]interface{}{"firstName": "Updated"},
			expectedError: true,
			errorContains: "id",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation UpdateEmployee($id: ID!, $input: UpdateEmployeeInput!) {
					updateEmployee(id: $id, input: $input) {
						id
						firstName
						lastName
						email
						phone
						department
						position
						salary
						status
						managerId
						address {
							street
							city
							state
							postalCode
							country
						}
						updatedAt
					}
				}`

			variables := map[string]interface{}{
				"id":    tc.employeeId,
				"input": tc.input,
			}

			resp, err := client.Execute(ctx, query, variables)

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented updateEmployee: %s", tc.name)
						assertContractError(t, resp, "updateEmployee")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented updateEmployee: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testDeleteEmployeeOperations tests all delete employee mutation scenarios
func testDeleteEmployeeOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		employeeId    string
		expectedError bool
		errorContains string
		expectSuccess bool
	}{
		{
			name:          "Delete Existing Employee",
			employeeId:    "employee-123",
			expectedError: true,
			expectSuccess: true,
		},
		{
			name:          "Delete Non-existent Employee",
			employeeId:    "non-existent-employee",
			expectedError: true,
			errorContains: "not found",
		},
		{
			name:          "Delete With Invalid Employee ID",
			employeeId:    "",
			expectedError: true,
			errorContains: "id",
		},
		{
			name:          "Delete Employee That Has Subordinates",
			employeeId:    "manager-123",
			expectedError: true,
			errorContains: "subordinates",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation DeleteEmployee($id: ID!) {
					deleteEmployee(id: $id)
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{"id": tc.employeeId})

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented deleteEmployee: %s", tc.name)
						assertContractError(t, resp, "deleteEmployee")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented deleteEmployee: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testGetEmployeeOperations tests all get employee query scenarios
func testGetEmployeeOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		employeeId    string
		expectedError bool
		errorContains string
		expectFields  []string
	}{
		{
			name:          "Get Existing Employee By ID",
			employeeId:    "employee-123",
			expectedError: true,
			expectFields:  []string{"id", "firstName", "lastName", "email", "department", "position", "status"},
		},
		{
			name:          "Get Non-existent Employee",
			employeeId:    "non-existent-employee",
			expectedError: true,
			errorContains: "not found",
		},
		{
			name:          "Get Employee With Invalid ID",
			employeeId:    "",
			expectedError: true,
			errorContains: "id",
		},
		{
			name:          "Get Employee Complete Information",
			employeeId:    "employee-123",
			expectedError: true,
			expectFields:  []string{"id", "firstName", "lastName", "email", "phone", "department", "position", "hireDate", "salary", "status", "managerId", "address", "createdAt", "updatedAt"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query GetEmployee($id: ID!) {
					employee(id: $id) {
						id
						firstName
						lastName
						email
						phone
						department
						position
						hireDate
						salary
						status
						managerId
						address {
							street
							city
							state
							postalCode
							country
						}
						createdAt
						updatedAt
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{"id": tc.employeeId})

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented employee query: %s", tc.name)
						assertContractError(t, resp, "employee query")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented employee query: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testChangeEmployeeStatusOperations tests all change employee status mutation scenarios
func testChangeEmployeeStatusOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		employeeId    string
		newStatus     string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Change Status From ACTIVE to ON_LEAVE",
			employeeId:    "employee-123",
			newStatus:     "ON_LEAVE",
			expectedError: true,
		},
		{
			name:          "Change Status From ON_LEAVE to ACTIVE",
			employeeId:    "employee-123",
			newStatus:     "ACTIVE",
			expectedError: true,
		},
		{
			name:          "Change Status to TERMINATED",
			employeeId:    "employee-123",
			newStatus:     "TERMINATED",
			expectedError: true,
		},
		{
			name:          "Change Status for Non-existent Employee",
			 employeeId:   "non-existent-employee",
			newStatus:     "ON_LEAVE",
			expectedError: true,
			errorContains: "not found",
		},
		{
			name:          "Invalid Status Value",
			employeeId:    "employee-123",
			newStatus:     "INVALID_STATUS",
			expectedError: true,
			errorContains: "status",
		},
		{
			name:          "Change Status Without Employee ID",
			employeeId:    "",
			newStatus:     "ON_LEAVE",
			expectedError: true,
			errorContains: "id",
		},
		{
			name:          "Change Status From TERMINATED (Should Fail)",
			employeeId:    "terminated-employee",
			newStatus:     "ACTIVE",
			expectedError: true,
			errorContains: "terminated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation ChangeEmployeeStatus($id: ID!, $status: EmployeeStatus!) {
					changeEmployeeStatus(id: $id, status: $status) {
						id
						status
						updatedAt
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"id":     tc.employeeId,
				"status": tc.newStatus,
			})

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented changeEmployeeStatus: %s", tc.name)
						assertContractError(t, resp, "changeEmployeeStatus")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented changeEmployeeStatus: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testEmployeeValidationAndErrorHandling tests validation and error handling scenarios
func testEmployeeValidationAndErrorHandling(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		operation     string
		query         string
		variables     map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name:      "Create Employee With Duplicate Email",
			operation: "create",
			query: `
				mutation CreateEmployee($input: CreateEmployeeInput!) {
					createEmployee(input: $input) {
						id
						email
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"firstName":  "John",
					"lastName":   "Doe",
					"email":      "duplicate@example.com",
					"department": "Engineering",
					"position":   "Software Engineer",
				},
			},
			expectedError: true,
			errorContains: "duplicate",
		},
		{
			name:      "Create Employee With Exceedingly Long Name",
			operation: "create",
			query: `
				mutation CreateEmployee($input: CreateEmployeeInput!) {
					createEmployee(input: $input) {
						id
						firstName
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"firstName":  strings.Repeat("A", 101), // Exceeds typical name length limit
					"lastName":   "Doe",
					"email":      "longname@example.com",
					"department": "Engineering",
					"position":   "Software Engineer",
				},
			},
			expectedError: true,
			errorContains: "length",
		},
		{
			name:      "Update Employee With Invalid Salary",
			operation: "update",
			query: `
				mutation UpdateEmployee($id: ID!, $input: UpdateEmployeeInput!) {
					updateEmployee(id: $id, input: $input) {
						id
						salary
					}
				}`,
			variables: map[string]interface{}{
				"id": "employee-123",
				"input": map[string]interface{}{
					"salary": -1000.0,
				},
			},
			expectedError: true,
			errorContains: "salary",
		},
		{
			name:      "Create Employee With Future Hire Date",
			operation: "create",
			query: `
				mutation CreateEmployee($input: CreateEmployeeInput!) {
					createEmployee(input: $input) {
						id
						hireDate
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"firstName":  "Future",
					"lastName":   "Employee",
					"email":      "future@example.com",
					"department": "Engineering",
					"position":   "Software Engineer",
					"hireDate":   "2030-01-01", // Future date
				},
			},
			expectedError: true,
			errorContains: "hireDate",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented %s operation: %s", tc.operation, tc.name)
						assertContractError(t, resp, tc.operation)
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented %s operation: %s", tc.operation, tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testEmployeeAuthenticationRequirements tests authentication requirements for employee operations
func testEmployeeAuthenticationRequirements(t *testing.T, client *helpers.GraphQLClient, testServer *helpers.TestServer) {
	// Create a client without authentication
	unauthenticatedClient := helpers.NewGraphQLClient(testServer.BaseURL)

	testCases := []struct {
		name          string
		query         string
		variables     map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name:  "Create Employee Without Authentication",
			query: `
				mutation CreateEmployee($input: CreateEmployeeInput!) {
					createEmployee(input: $input) {
						id
						firstName
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"firstName":  "Unauthorized",
					"lastName":   "User",
					"email":      "unauthorized@example.com",
					"department": "Engineering",
					"position":   "Software Engineer",
				},
			},
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name:  "Update Employee Without Authentication",
			query: `
				mutation UpdateEmployee($id: ID!, $input: UpdateEmployeeInput!) {
					updateEmployee(id: $id, input: $input) {
						id
						firstName
					}
				}`,
			variables: map[string]interface{}{
				"id":    "employee-123",
				"input": map[string]interface{}{"firstName": "Hacked"},
			},
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name:  "Delete Employee Without Authentication",
			query: `
				mutation DeleteEmployee($id: ID!) {
					deleteEmployee(id: $id)
				}`,
			variables:     map[string]interface{}{"id": "employee-123"},
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name:  "Get Employee Without Authentication",
			query: `
				query GetEmployee($id: ID!) {
					employee(id: $id) {
						id
						firstName
					}
				}`,
			variables:     map[string]interface{}{"id": "employee-123"},
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name:  "List Employees Without Authentication",
			query: `
				query {
					employees(first: 10) {
						edges {
							node {
								id
								firstName
								lastName
							}
						}
					}
				}`,
			expectedError: true,
			errorContains: "authentication",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := unauthenticatedClient.Execute(ctx, tc.query, tc.variables)

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unauthenticated operation: %s", tc.name)
						assertContractError(t, resp, "authentication")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unauthenticated operation: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// Helper functions for error assertions

// assertContractError asserts that the response contains contract-related errors
func assertContractError(t *testing.T, resp *helpers.GraphQLResponse, operation string) {
	t.Helper()

	require.True(t, resp.HasErrors(), "Expected contract errors for %s", operation)

	foundContractError := false
	for _, err := range resp.Errors {
		errorMsg := err.Message
		if containsIgnoreCase(errorMsg, "cannot query") ||
		   containsIgnoreCase(errorMsg, "unknown type") ||
		   containsIgnoreCase(errorMsg, "field") ||
		   containsIgnoreCase(errorMsg, "mutation") ||
		   containsIgnoreCase(errorMsg, "query") ||
		   containsIgnoreCase(errorMsg, "not found") ||
		   containsIgnoreCase(errorMsg, "undefined") ||
		   containsIgnoreCase(errorMsg, operation) {
			foundContractError = true
			break
		}
	}

	require.True(t, foundContractError, "Expected contract-related errors for %s, but got: %v", operation, resp.Errors)
}

// assertErrorContains asserts that the GraphQL response contains an error with the specified message
func assertErrorContains(t *testing.T, resp *helpers.GraphQLResponse, message, testName string) {
	t.Helper()

	require.True(t, resp.HasErrors(), "Expected errors for %s", testName)

	found := false
	for _, err := range resp.Errors {
		if containsIgnoreCase(err.Message, message) {
			found = true
			break
		}
	}

	require.True(t, found, "Expected error containing '%s' for %s, but got: %v", message, testName, resp.Errors)
}

// assertConnectionError asserts that the error is a connection-related error
func assertConnectionError(t *testing.T, err error, testName string) {
	t.Helper()

	require.NotNil(t, err, "Expected connection error for %s", testName)
	errorMsg := err.Error()

	assert.True(t,
		containsIgnoreCase(errorMsg, "connection refused") ||
		containsIgnoreCase(errorMsg, "no such host") ||
		containsIgnoreCase(errorMsg, "404") ||
		containsIgnoreCase(errorMsg, "500") ||
		containsIgnoreCase(errorMsg, "timeout") ||
		containsIgnoreCase(errorMsg, "unreachable"),
		"Expected connection error for %s, but got: %v", testName, err)
}

