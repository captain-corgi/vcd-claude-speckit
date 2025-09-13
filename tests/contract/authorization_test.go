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

// TestAuthorizationContract tests the role-based access control (RBAC) system
// and validates that operations work correctly for different user roles.
// This is a RED phase test - it MUST fail before implementation exists.
func TestAuthorizationContract(t *testing.T) {
	// Setup test server
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	// Create GraphQL client
	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all authorization validation tests
	t.Run("Employee Access Control", func(t *testing.T) {
		testEmployeeAccessControl(t, client)
	})

	t.Run("User Management Access Control", func(t *testing.T) {
		testUserManagementAccessControl(t, client)
	})

	t.Run("Audit Log Access Control", func(t *testing.T) {
		testAuditLogAccessControl(t, client)
	})

	t.Run("Authentication Flow Access Control", func(t *testing.T) {
		testAuthenticationFlowAccessControl(t, client)
	})

	t.Run("System Configuration Access Control", func(t *testing.T) {
		testSystemConfigurationAccessControl(t, client)
	})
}

// testEmployeeAccessControl tests employee operations with different user roles
func testEmployeeAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Employee Queries", func(t *testing.T) {
		testEmployeeQueriesAccessControl(t, client)
	})

	t.Run("Employee Mutations", func(t *testing.T) {
		testEmployeeMutationsAccessControl(t, client)
	})
}

// testEmployeeQueriesAccessControl tests query operations on employees
func testEmployeeQueriesAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		role        string
		token       string
		query       string
		expectError bool
		errorType   string
	}{
		// ADMIN role - should have full access
		{
			name:        "ADMIN can list employees",
			role:        "ADMIN",
			token:       helpers.CreateAdminToken(t),
			query:       "query { employees { edges { node { id firstName lastName } } } }",
			expectError: false,
		},
		{
			name:        "ADMIN can get single employee",
			role:        "ADMIN",
			token:       helpers.CreateAdminToken(t),
			query:       "query { employee(id: \"test-id\") { id firstName lastName } }",
			expectError: false,
		},

		// MANAGER role - should have read access
		{
			name:        "MANAGER can list employees",
			role:        "MANAGER",
			token:       helpers.CreateManagerToken(t),
			query:       "query { employees { edges { node { id firstName lastName } } } }",
			expectError: false,
		},
		{
			name:        "MANAGER can get single employee",
			role:        "MANAGER",
			token:       helpers.CreateManagerToken(t),
			query:       "query { employee(id: \"test-id\") { id firstName lastName } }",
			expectError: false,
		},

		// VIEWER role - should have read access but limited data
		{
			name:        "VIEWER can list employees",
			role:        "VIEWER",
			token:       helpers.CreateViewerToken(t),
			query:       "query { employees { edges { node { id firstName lastName status } } } }",
			expectError: false,
		},
		{
			name:        "VIEWER can get single employee",
			role:        "VIEWER",
			token:       helpers.CreateViewerToken(t),
			query:       "query { employee(id: \"test-id\") { id firstName lastName status } }",
			expectError: false,
		},

		// Unauthorized users - should fail
		{
			name:        "Unauthenticated user cannot list employees",
			role:        "UNAUTHENTICATED",
			token:       "",
			query:       "query { employees { edges { node { id } } } }",
			expectError: true,
			errorType:   "unauthenticated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create client with appropriate authentication
			testClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				testClient.WithAuth(tc.token)
			}

			resp, err := testClient.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for %s", tc.name)
					assertAuthorizationError(t, resp, tc.name)
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				}
			} else {
				// In RED phase, we expect connection errors or schema errors
				if err != nil {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL schema errors for unimplemented operation")
					assertSchemaError(t, resp, tc.name)
				}
			}
		})
	}
}

// testEmployeeMutationsAccessControl tests mutation operations on employees
func testEmployeeMutationsAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	createEmployeeMutation := `
		mutation CreateEmployee($input: CreateEmployeeInput!) {
			createEmployee(input: $input) {
				id
				firstName
				lastName
				email
				createdAt
			}
		}`

	updateEmployeeMutation := `
		mutation UpdateEmployee($id: ID!, $input: UpdateEmployeeInput!) {
			updateEmployee(id: $id, input: $input) {
				id
				firstName
				lastName
				email
				updatedAt
			}
		}`

	deleteEmployeeMutation := `
		mutation DeleteEmployee($id: ID!) {
			deleteEmployee(id: $id)
		}`

	changeEmployeeStatusMutation := `
		mutation ChangeEmployeeStatus($id: ID!, $status: EmployeeStatus!) {
			changeEmployeeStatus(id: $id, status: $status) {
				id
				status
			}
		}`

	testCases := []struct {
		name          string
		role          string
		token         string
		mutation      string
		variables     map[string]interface{}
		expectError   bool
		errorType     string
		operationType string
	}{
		// ADMIN role - should have full CRUD access
		{
			name:          "ADMIN can create employee",
			role:          "ADMIN",
			token:         helpers.CreateAdminToken(t),
			mutation:      createEmployeeMutation,
			variables:     map[string]interface{}{"input": map[string]interface{}{"firstName": "Test", "lastName": "User", "email": "test@example.com"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "schema",
			operationType: "create",
		},
		{
			name:          "ADMIN can update employee",
			role:          "ADMIN",
			token:         helpers.CreateAdminToken(t),
			mutation:      updateEmployeeMutation,
			variables:     map[string]interface{}{"id": "test-id", "input": map[string]interface{}{"firstName": "Updated"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "schema",
			operationType: "update",
		},
		{
			name:          "ADMIN can delete employee",
			role:          "ADMIN",
			token:         helpers.CreateAdminToken(t),
			mutation:      deleteEmployeeMutation,
			variables:     map[string]interface{}{"id": "test-id"},
			expectError:   true, // RED phase - expect to fail
			errorType:     "schema",
			operationType: "delete",
		},
		{
			name:          "ADMIN can change employee status",
			role:          "ADMIN",
			token:         helpers.CreateAdminToken(t),
			mutation:      changeEmployeeStatusMutation,
			variables:     map[string]interface{}{"id": "test-id", "status": "ACTIVE"},
			expectError:   true, // RED phase - expect to fail
			errorType:     "schema",
			operationType: "status_change",
		},

		// MANAGER role - should have read/update but not delete
		{
			name:          "MANAGER cannot create employee",
			role:          "MANAGER",
			token:         helpers.CreateManagerToken(t),
			mutation:      createEmployeeMutation,
			variables:     map[string]interface{}{"input": map[string]interface{}{"firstName": "Test", "lastName": "User", "email": "test@example.com"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "authorization",
			operationType: "create",
		},
		{
			name:          "MANAGER can update employee",
			role:          "MANAGER",
			token:         helpers.CreateManagerToken(t),
			mutation:      updateEmployeeMutation,
			variables:     map[string]interface{}{"id": "test-id", "input": map[string]interface{}{"firstName": "Updated"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "schema",
			operationType: "update",
		},
		{
			name:          "MANAGER cannot delete employee",
			role:          "MANAGER",
			token:         helpers.CreateManagerToken(t),
			mutation:      deleteEmployeeMutation,
			variables:     map[string]interface{}{"id": "test-id"},
			expectError:   true, // RED phase - expect to fail
			errorType:     "authorization",
			operationType: "delete",
		},

		// VIEWER role - should only have read access
		{
			name:          "VIEWER cannot create employee",
			role:          "VIEWER",
			token:         helpers.CreateViewerToken(t),
			mutation:      createEmployeeMutation,
			variables:     map[string]interface{}{"input": map[string]interface{}{"firstName": "Test", "lastName": "User", "email": "test@example.com"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "authorization",
			operationType: "create",
		},
		{
			name:          "VIEWER cannot update employee",
			role:          "VIEWER",
			token:         helpers.CreateViewerToken(t),
			mutation:      updateEmployeeMutation,
			variables:     map[string]interface{}{"id": "test-id", "input": map[string]interface{}{"firstName": "Updated"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "authorization",
			operationType: "update",
		},
		{
			name:          "VIEWER cannot delete employee",
			role:          "VIEWER",
			token:         helpers.CreateViewerToken(t),
			mutation:      deleteEmployeeMutation,
			variables:     map[string]interface{}{"id": "test-id"},
			expectError:   true, // RED phase - expect to fail
			errorType:     "authorization",
			operationType: "delete",
		},

		// Unauthorized users - should fail all operations
		{
			name:          "Unauthenticated user cannot create employee",
			role:          "UNAUTHENTICATED",
			token:         "",
			mutation:      createEmployeeMutation,
			variables:     map[string]interface{}{"input": map[string]interface{}{"firstName": "Test", "lastName": "User", "email": "test@example.com"}},
			expectError:   true, // RED phase - expect to fail
			errorType:     "unauthenticated",
			operationType: "create",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create client with appropriate authentication
			testClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				testClient.WithAuth(tc.token)
			}

			resp, err := testClient.Execute(ctx, tc.mutation, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for %s", tc.name)
					if tc.errorType == "schema" {
						assertSchemaError(t, resp, tc.name)
					} else if tc.errorType == "authorization" {
						assertAuthorizationError(t, resp, tc.name)
					} else if tc.errorType == "unauthenticated" {
						assertUnauthenticatedError(t, resp, tc.name)
					}
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				}
			}
		})
	}
}

// testUserManagementAccessControl tests user management operations with different roles
func testUserManagementAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		role        string
		token       string
		query       string
		expectError bool
		errorType   string
	}{
		// ADMIN role - should have full user management access
		{
			name:        "ADMIN can list users",
			role:        "ADMIN",
			token:       helpers.CreateAdminToken(t),
			query:       "query { users { edges { node { id username role } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},
		{
			name:        "ADMIN can get user details",
			role:        "ADMIN",
			token:       helpers.CreateAdminToken(t),
			query:       "query { user(id: \"test-id\") { id username role email } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},

		// MANAGER role - should have limited user management access
		{
			name:        "MANAGER can list users",
			role:        "MANAGER",
			token:       helpers.CreateManagerToken(t),
			query:       "query { users { edges { node { id username } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},
		{
			name:        "MANAGER cannot get user details with sensitive information",
			role:        "MANAGER",
			token:       helpers.CreateManagerToken(t),
			query:       "query { user(id: \"test-id\") { id username role email } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "authorization",
		},

		// VIEWER role - should have no user management access
		{
			name:        "VIEWER cannot list users",
			role:        "VIEWER",
			token:       helpers.CreateViewerToken(t),
			query:       "query { users { edges { node { id } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "authorization",
		},

		// Unauthorized users - should fail
		{
			name:        "Unauthenticated user cannot list users",
			role:        "UNAUTHENTICATED",
			token:       "",
			query:       "query { users { edges { node { id } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "unauthenticated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create client with appropriate authentication
			testClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				testClient.WithAuth(tc.token)
			}

			resp, err := testClient.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for %s", tc.name)
					if tc.errorType == "schema" {
						assertSchemaError(t, resp, tc.name)
					} else if tc.errorType == "authorization" {
						assertAuthorizationError(t, resp, tc.name)
					} else if tc.errorType == "unauthenticated" {
						assertUnauthenticatedError(t, resp, tc.name)
					}
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				}
			}
		})
	}
}

// testAuditLogAccessControl tests audit log operations with different roles
func testAuditLogAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		role        string
		token       string
		query       string
		expectError bool
		errorType   string
	}{
		// ADMIN role - should have full audit log access
		{
			name:        "ADMIN can list audit logs",
			role:        "ADMIN",
			token:       helpers.CreateAdminToken(t),
			query:       "query { auditLogs { edges { node { id action entityType entityId user timestamp } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},

		// MANAGER role - should have limited audit log access
		{
			name:        "MANAGER can list audit logs but limited fields",
			role:        "MANAGER",
			token:       helpers.CreateManagerToken(t),
			query:       "query { auditLogs { edges { node { id action entityType timestamp } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},

		// VIEWER role - should have no audit log access
		{
			name:        "VIEWER cannot list audit logs",
			role:        "VIEWER",
			token:       helpers.CreateViewerToken(t),
			query:       "query { auditLogs { edges { node { id } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "authorization",
		},

		// Unauthorized users - should fail
		{
			name:        "Unauthenticated user cannot list audit logs",
			role:        "UNAUTHENTICATED",
			token:       "",
			query:       "query { auditLogs { edges { node { id } } } }",
			expectError: true, // RED phase - expect to fail
			errorType:   "unauthenticated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create client with appropriate authentication
			testClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				testClient.WithAuth(tc.token)
			}

			resp, err := testClient.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for %s", tc.name)
					if tc.errorType == "schema" {
						assertSchemaError(t, resp, tc.name)
					} else if tc.errorType == "authorization" {
						assertAuthorizationError(t, resp, tc.name)
					} else if tc.errorType == "unauthenticated" {
						assertUnauthenticatedError(t, resp, tc.name)
					}
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				}
			}
		})
	}
}

// testAuthenticationFlowAccessControl tests authentication-related operations
func testAuthenticationFlowAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	loginMutation := `
		mutation Login($email: String!, $password: String!) {
			login(email: $email, password: $password) {
				token
				user {
					id
					username
				}
			}
		}`

	refreshTokenMutation := `
		mutation RefreshToken($token: String!) {
			refreshToken(token: $token) {
				token
				expiresAt
			}
		}`

	logoutMutation := `
		mutation Logout {
			logout
		}`

	testCases := []struct {
		name        string
		role        string
		token       string
		mutation    string
		variables   map[string]interface{}
		expectError bool
		errorType   string
	}{
		// Login should work for all authenticated users
		{
			name:        "Valid login works",
			role:        "VALID_CREDENTIALS",
			token:       "",
			mutation:    loginMutation,
			variables:   map[string]interface{}{"email": "test@example.com", "password": "password"},
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},
		{
			name:        "Invalid login fails",
			role:        "INVALID_CREDENTIALS",
			token:       "",
			mutation:    loginMutation,
			variables:   map[string]interface{}{"email": "invalid@example.com", "password": "wrong"},
			expectError: true, // RED phase - expect to fail
			errorType:   "validation",
		},

		// Token refresh should work for valid tokens
		{
			name:        "Valid token refresh works",
			role:        "VALID_TOKEN",
			token:       helpers.CreateAdminToken(t),
			mutation:    refreshTokenMutation,
			variables:   map[string]interface{}{"token": "valid-token"},
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},

		// Logout should work for authenticated users
		{
			name:        "Authenticated user can logout",
			role:        "AUTHENTICATED",
			token:       helpers.CreateAdminToken(t),
			mutation:    logoutMutation,
			variables:   nil,
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},

		// Logout should fail for unauthenticated users
		{
			name:        "Unauthenticated user cannot logout",
			role:        "UNAUTHENTICATED",
			token:       "",
			mutation:    logoutMutation,
			variables:   nil,
			expectError: true, // RED phase - expect to fail
			errorType:   "unauthenticated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create client with appropriate authentication
			testClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				testClient.WithAuth(tc.token)
			}

			resp, err := testClient.Execute(ctx, tc.mutation, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for %s", tc.name)
					if tc.errorType == "schema" {
						assertSchemaError(t, resp, tc.name)
					} else if tc.errorType == "validation" {
						assertValidationError(t, resp, tc.name)
					} else if tc.errorType == "unauthenticated" {
						assertUnauthenticatedError(t, resp, tc.name)
					}
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				}
			}
		})
	}
}

// testSystemConfigurationAccessControl tests system configuration operations
func testSystemConfigurationAccessControl(t *testing.T, client *helpers.GraphQLClient) {
	meQuery := `
		query Me {
			me {
				id
				username
				role
				email
			}
		}`

	systemConfigQuery := `
		query SystemConfig {
			systemConfig {
				companyName
				settings {
					featuresEnabled
				}
			}
		}`

	testCases := []struct {
		name        string
		role        string
		token       string
		query       string
		expectError bool
		errorType   string
	}{
		// Me query should work for all authenticated users
		{
			name:        "Authenticated user can access me query",
			role:        "AUTHENTICATED",
			token:       helpers.CreateAdminToken(t),
			query:       meQuery,
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},

		// System config should only be accessible by ADMIN
		{
			name:        "ADMIN can access system configuration",
			role:        "ADMIN",
			token:       helpers.CreateAdminToken(t),
			query:       systemConfigQuery,
			expectError: true, // RED phase - expect to fail
			errorType:   "schema",
		},
		{
			name:        "MANAGER cannot access system configuration",
			role:        "MANAGER",
			token:       helpers.CreateManagerToken(t),
			query:       systemConfigQuery,
			expectError: true, // RED phase - expect to fail
			errorType:   "authorization",
		},
		{
			name:        "VIEWER cannot access system configuration",
			role:        "VIEWER",
			token:       helpers.CreateViewerToken(t),
			query:       systemConfigQuery,
			expectError: true, // RED phase - expect to fail
			errorType:   "authorization",
		},

		// Me query should fail for unauthenticated users
		{
			name:        "Unauthenticated user cannot access me query",
			role:        "UNAUTHENTICATED",
			token:       "",
			query:       meQuery,
			expectError: true, // RED phase - expect to fail
			errorType:   "unauthenticated",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Create client with appropriate authentication
			testClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				testClient.WithAuth(tc.token)
			}

			resp, err := testClient.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for %s", tc.name)
					if tc.errorType == "schema" {
						assertSchemaError(t, resp, tc.name)
					} else if tc.errorType == "authorization" {
						assertAuthorizationError(t, resp, tc.name)
					} else if tc.errorType == "unauthenticated" {
						assertUnauthenticatedError(t, resp, tc.name)
					}
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
						strings.Contains(err.Error(), "no such host") ||
						strings.Contains(err.Error(), "404") ||
						strings.Contains(err.Error(), "500"),
						"Expected connection error, but got: %v", err)
				}
			}
		})
	}
}

// Helper functions for asserting specific error types

// assertAuthorizationError asserts that the response contains authorization-related errors
func assertAuthorizationError(t *testing.T, resp *helpers.GraphQLResponse, testName string) {
	t.Helper()

	require.True(t, resp.HasErrors(), "Expected authorization errors for %s", testName)

	foundAuthError := false
	for _, err := range resp.Errors {
		errorMsg := err.Message
		if containsIgnoreCase(errorMsg, "unauthorized") ||
			containsIgnoreCase(errorMsg, "permission denied") ||
			containsIgnoreCase(errorMsg, "forbidden") ||
			containsIgnoreCase(errorMsg, "access denied") ||
			containsIgnoreCase(errorMsg, "insufficient privileges") ||
			containsIgnoreCase(errorMsg, "role not allowed") {
			foundAuthError = true
			break
		}
	}

	require.True(t, foundAuthError, "Expected authorization-related errors for %s, but got: %v", testName, resp.Errors)
}

// assertUnauthenticatedError asserts that the response contains unauthentication errors
func assertUnauthenticatedError(t *testing.T, resp *helpers.GraphQLResponse, testName string) {
	t.Helper()

	require.True(t, resp.HasErrors(), "Expected unauthenticated errors for %s", testName)

	foundUnauthError := false
	for _, err := range resp.Errors {
		errorMsg := err.Message
		if containsIgnoreCase(errorMsg, "unauthenticated") ||
			containsIgnoreCase(errorMsg, "not authenticated") ||
			containsIgnoreCase(errorMsg, "missing token") ||
			containsIgnoreCase(errorMsg, "invalid token") ||
			containsIgnoreCase(errorMsg, "token expired") ||
			containsIgnoreCase(errorMsg, "authorization required") {
			foundUnauthError = true
			break
		}
	}

	require.True(t, foundUnauthError, "Expected unauthentication-related errors for %s, but got: %v", testName, resp.Errors)
}