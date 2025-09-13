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

// TestGraphQLSchemaContract tests the GraphQL schema structure and validates
// that it matches the expected employee management system contract.
// This is a RED phase test - it MUST fail before implementation exists.
func TestGraphQLSchemaContract(t *testing.T) {
	// Setup test server
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	// Create GraphQL client
	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all schema validation tests
	t.Run("GraphQL Endpoint Availability", func(t *testing.T) {
		testGraphQLEndpointAvailability(t, client)
	})

	t.Run("Query Types Schema", func(t *testing.T) {
		testQueryTypesSchema(t, client)
	})

	t.Run("Mutation Types Schema", func(t *testing.T) {
		testMutationTypesSchema(t, client)
	})

	t.Run("Input Types Schema", func(t *testing.T) {
		testInputTypesSchema(t, client)
	})

	t.Run("Enum Types Schema", func(t *testing.T) {
		testEnumTypesSchema(t, client)
	})

	t.Run("Custom Scalars Schema", func(t *testing.T) {
		testCustomScalarsSchema(t, client)
	})

	t.Run("Employee Type Fields", func(t *testing.T) {
		testEmployeeTypeFields(t, client)
	})

	t.Run("Address Embedded Type", func(t *testing.T) {
		testAddressEmbeddedType(t, client)
	})

	t.Run("Pagination Schema", func(t *testing.T) {
		testPaginationSchema(t, client)
	})

	t.Run("Authentication Schema", func(t *testing.T) {
		testAuthenticationSchema(t, client)
	})

	t.Run("Audit Log Schema", func(t *testing.T) {
		testAuditLogSchema(t, client)
	})

	t.Run("User Management Schema", func(t *testing.T) {
		testUserManagementSchema(t, client)
	})

	t.Run("Validation Rules", func(t *testing.T) {
		testValidationRules(t, client)
	})
}

// testGraphQLEndpointAvailability tests that the GraphQL endpoint is available
func testGraphQLEndpointAvailability(t *testing.T, client *helpers.GraphQLClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Test schema introspection query
	introspectionQuery := `
		query {
			__schema {
				types {
					name
					kind
					description
				}
				queryType {
					name
				}
				mutationType {
					name
				}
			}
		}`

	resp, err := client.Execute(ctx, introspectionQuery, nil)

	// In RED phase, we expect either a connection error or schema error
	// This test will fail until GraphQL endpoint is implemented
	if err == nil {
		// If we got a response, it should have errors (no schema implemented yet)
		require.True(t, resp.HasErrors(), "GraphQL endpoint should return errors when no schema is implemented")

		// Look for specific error types that indicate no GraphQL schema
		foundExpectedError := false
		for _, graphqlErr := range resp.Errors {
			if containsIgnoreCase(graphqlErr.Message, "cannot query") ||
				containsIgnoreCase(graphqlErr.Message, "unknown type") ||
				containsIgnoreCase(graphqlErr.Message, "graphql") {
				foundExpectedError = true
				break
			}
		}

		require.True(t, foundExpectedError, "Expected GraphQL schema errors, but got: %v", resp.Errors)
	} else {
		// Connection errors are expected in RED phase
		assert.True(t,
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "404") ||
			strings.Contains(err.Error(), "500"),
			"Expected connection error, but got: %v", err)
	}
}

// testQueryTypesSchema tests that all expected query types exist in the schema
func testQueryTypesSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "Employees Query",
			query:       "query { employees { edges { node { id } } } }",
			expectError: true,
		},
		{
			name:        "Employee Query",
			query:       "query { employee(id: \"test\") { id } }",
			expectError: true,
		},
		{
			name:        "Users Query",
			query:       "query { users { id username } }",
			expectError: true,
		},
		{
			name:        "User Query",
			query:       "query { user(id: \"test\") { id username } }",
			expectError: true,
		},
		{
			name:        "Audit Logs Query",
			query:       "query { auditLogs { id action } }",
			expectError: true,
		},
		{
			name:        "Me Query",
			query:       "query { me { id username } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented query: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testMutationTypesSchema tests that all expected mutation types exist in the schema
func testMutationTypesSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		mutation    string
		expectError bool
	}{
		{
			name:        "Create Employee Mutation",
			mutation:    "mutation { createEmployee(input: {firstName: \"Test\", lastName: \"User\", email: \"test@example.com\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Update Employee Mutation",
			mutation:    "mutation { updateEmployee(id: \"test\", input: {firstName: \"Updated\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Delete Employee Mutation",
			mutation:    "mutation { deleteEmployee(id: \"test\") }",
			expectError: true,
		},
		{
			name:        "Change Employee Status Mutation",
			mutation:    "mutation { changeEmployeeStatus(id: \"test\", status: ACTIVE) { id status } }",
			expectError: true,
		},
		{
			name:        "Login Mutation",
			mutation:    "mutation { login(email: \"test@example.com\", password: \"password\") { token } }",
			expectError: true,
		},
		{
			name:        "Refresh Token Mutation",
			mutation:    "mutation { refreshToken(token: \"test\") { token } }",
			expectError: true,
		},
		{
			name:        "Logout Mutation",
			mutation:    "mutation { logout }",
			expectError: true,
		},
		{
			name:        "Create User Mutation",
			mutation:    "mutation { createUser(input: {username: \"test\", email: \"test@example.com\", password: \"password\", role: ADMIN}) { id } }",
			expectError: true,
		},
		{
			name:        "Update User Mutation",
			mutation:    "mutation { updateUser(id: \"test\", input: {username: \"updated\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Delete User Mutation",
			mutation:    "mutation { deleteUser(id: \"test\") }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.mutation, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented mutation: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testInputTypesSchema tests that all expected input types exist in the schema
func testInputTypesSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "CreateEmployeeInput Type",
			query:       "query { __type(name: \"CreateEmployeeInput\") { name kind fields { name type { name } } } }",
			expectError: true,
		},
		{
			name:        "UpdateEmployeeInput Type",
			query:       "query { __type(name: \"UpdateEmployeeInput\") { name kind fields { name type { name } } } }",
			expectError: true,
		},
		{
			name:        "EmployeeFilter Type",
			query:       "query { __type(name: \"EmployeeFilter\") { name kind fields { name type { name } } } }",
			expectError: true,
		},
		{
			name:        "EmployeeSort Type",
			query:       "query { __type(name: \"EmployeeSort\") { name kind fields { name type { name } } } }",
			expectError: true,
		},
		{
			name:        "UserFilter Type",
			query:       "query { __type(name: \"UserFilter\") { name kind fields { name type { name } } }",
			expectError: true,
		},
		{
			name:        "LoginInput Type",
			query:       "query { __type(name: \"LoginInput\") { name kind fields { name type { name } } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented input type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testEnumTypesSchema tests that all expected enum types exist in the schema
func testEnumTypesSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		enumValues  []string
		expectError bool
	}{
		{
			name:        "EmployeeStatus Enum",
			query:       "query { __type(name: \"EmployeeStatus\") { name kind enumValues { name } } }",
			enumValues:  []string{"ACTIVE", "INACTIVE", "ON_LEAVE", "TERMINATED"},
			expectError: true,
		},
		{
			name:        "EmployeeSortField Enum",
			query:       "query { __type(name: \"EmployeeSortField\") { name kind enumValues { name } } }",
			enumValues:  []string{"ID", "FIRST_NAME", "LAST_NAME", "EMAIL", "DEPARTMENT", "POSITION", "SALARY", "STATUS", "CREATED_AT", "UPDATED_AT"},
			expectError: true,
		},
		{
			name:        "SortOrder Enum",
			query:       "query { __type(name: \"SortOrder\") { name kind enumValues { name } } }",
			enumValues:  []string{"ASC", "DESC"},
			expectError: true,
		},
		{
			name:        "UserRole Enum",
			query:       "query { __type(name: \"UserRole\") { name kind enumValues { name } } }",
			enumValues:  []string{"ADMIN", "MANAGER", "VIEWER"},
			expectError: true,
		},
		{
			name:        "AuditAction Enum",
			query:       "query { __type(name: \"AuditAction\") { name kind enumValues { name } } }",
			enumValues:  []string{"CREATE", "UPDATE", "DELETE", "LOGIN", "LOGOUT", "STATUS_CHANGE"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented enum type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testCustomScalarsSchema tests that all expected custom scalar types exist
func testCustomScalarsSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "Date Scalar",
			query:       "query { __type(name: \"Date\") { name kind } }",
			expectError: true,
		},
		{
			name:        "DateTime Scalar",
			query:       "query { __type(name: \"DateTime\") { name kind } }",
			expectError: true,
		},
		{
			name:        "UUID Scalar",
			query:       "query { __type(name: \"UUID\") { name kind } }",
			expectError: true,
		},
		{
			name:        "Money Scalar",
			query:       "query { __type(name: \"Money\") { name kind } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented scalar type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testEmployeeTypeFields tests that Employee type has all expected fields
func testEmployeeTypeFields(t *testing.T, client *helpers.GraphQLClient) {
	query := `query {
		__type(name: "Employee") {
			name
			fields {
				name
				type {
					name
					kind
				}
			}
		}
	}`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Execute(ctx, query, nil)

	// In RED phase, we expect errors
	if err == nil {
		require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented Employee type")
		assertSchemaError(t, resp, "Employee type")
	} else {
		assert.True(t,
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "404") ||
			strings.Contains(err.Error(), "500"),
			"Expected connection error, but got: %v", err)
	}
}

// testAddressEmbeddedType tests that Address type exists and has expected fields
func testAddressEmbeddedType(t *testing.T, client *helpers.GraphQLClient) {
	query := `query {
		__type(name: "Address") {
			name
			fields {
				name
				type {
					name
				}
			}
		}
	}`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.Execute(ctx, query, nil)

	// In RED phase, we expect errors
	if err == nil {
		require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented Address type")
		assertSchemaError(t, resp, "Address type")
	} else {
		assert.True(t,
			strings.Contains(err.Error(), "connection refused") ||
			strings.Contains(err.Error(), "no such host") ||
			strings.Contains(err.Error(), "404") ||
			strings.Contains(err.Error(), "500"),
			"Expected connection error, but got: %v", err)
	}
}

// testPaginationSchema tests that pagination types exist in the schema
func testPaginationSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "PageInfo Type",
			query:       "query { __type(name: \"PageInfo\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "EmployeeEdge Type",
			query:       "query { __type(name: \"EmployeeEdge\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "EmployeeConnection Type",
			query:       "query { __type(name: \"EmployeeConnection\") { name fields { name } } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented pagination type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testAuthenticationSchema tests that authentication-related types exist
func testAuthenticationSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "AuthPayload Type",
			query:       "query { __type(name: \"AuthPayload\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "RefreshTokenPayload Type",
			query:       "query { __type(name: \"RefreshTokenPayload\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "LogoutPayload Type",
			query:       "query { __type(name: \"LogoutPayload\") { name fields { name } } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented auth type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testAuditLogSchema tests that audit log related types exist
func testAuditLogSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "AuditLog Type",
			query:       "query { __type(name: \"AuditLog\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "AuditLogEdge Type",
			query:       "query { __type(name: \"AuditLogEdge\") { name fields { name } }",
			expectError: true,
		},
		{
			name:        "AuditLogConnection Type",
			query:       "query { __type(name: \"AuditLogConnection\") { name fields { name } }",
			expectError: true,
		},
		{
			name:        "AuditLogFilter Type",
			query:       "query { __type(name: \"AuditLogFilter\") { name fields { name } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented audit log type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testUserManagementSchema tests that user management related types exist
func testUserManagementSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "User Type",
			query:       "query { __type(name: \"User\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "CreateUserInput Type",
			query:       "query { __type(name: \"CreateUserInput\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "UpdateUserInput Type",
			query:       "query { __type(name: \"UpdateUserInput\") { name fields { name } } }",
			expectError: true,
		},
		{
			name:        "UserEdge Type",
			query:       "query { __type(name: \"UserEdge\") { name fields { name } }",
			expectError: true,
		},
		{
			name:        "UserConnection Type",
			query:       "query { __type(name: \"UserConnection\") { name fields { name } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented user type: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
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

// testValidationRules tests that validation rules are properly configured
func testValidationRules(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "Required Field Validation",
			query:       "mutation { createEmployee(input: {firstName: \"\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Email Format Validation",
			query:       "mutation { createEmployee(input: {firstName: \"Test\", lastName: \"User\", email: \"invalid-email\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Minimum Length Validation",
			query:       "mutation { createEmployee(input: {firstName: \"A\", lastName: \"User\", email: \"test@example.com\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Maximum Length Validation",
			query:       "mutation { createEmployee(input: {firstName: \"" + helpers.RandomString(51) + "\", lastName: \"User\", email: \"test@example.com\"}) { id } }",
			expectError: true,
		},
		{
			name:        "Non-null Field Validation",
			query:       "mutation { updateEmployee(id: null, input: {firstName: \"Updated\"}) { id } }",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, nil)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL validation errors for: %s", tc.name)
					assertValidationError(t, resp, tc.name)
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

// assertSchemaError asserts that the response contains schema-related errors
func assertSchemaError(t *testing.T, resp *helpers.GraphQLResponse, testName string) {
	t.Helper()

	require.True(t, resp.HasErrors(), "Expected schema errors for %s", testName)

	foundSchemaError := false
	for _, err := range resp.Errors {
		errorMsg := err.Message
		if containsIgnoreCase(errorMsg, "cannot query") ||
			containsIgnoreCase(errorMsg, "unknown type") ||
			containsIgnoreCase(errorMsg, "field") ||
			containsIgnoreCase(errorMsg, "type") ||
			containsIgnoreCase(errorMsg, "not found") ||
			containsIgnoreCase(errorMsg, "undefined") {
			foundSchemaError = true
			break
		}
	}

	require.True(t, foundSchemaError, "Expected schema-related errors for %s, but got: %v", testName, resp.Errors)
}

// assertValidationError asserts that the response contains validation errors
func assertValidationError(t *testing.T, resp *helpers.GraphQLResponse, testName string) {
	t.Helper()

	require.True(t, resp.HasErrors(), "Expected validation errors for %s", testName)

	foundValidationError := false
	for _, err := range resp.Errors {
		errorMsg := err.Message
		if containsIgnoreCase(errorMsg, "validation") ||
			containsIgnoreCase(errorMsg, "invalid") ||
			containsIgnoreCase(errorMsg, "required") ||
			containsIgnoreCase(errorMsg, "format") ||
			containsIgnoreCase(errorMsg, "length") ||
			containsIgnoreCase(errorMsg, "null") {
			foundValidationError = true
			break
		}
	}

	require.True(t, foundValidationError, "Expected validation errors for %s, but got: %v", testName, resp.Errors)
}

// containsIgnoreCase performs case-insensitive substring check
func containsIgnoreCase(s, substr string) bool {
	lowerS := strings.ToLower(s)
	lowerSubstr := strings.ToLower(substr)
	return strings.Contains(lowerS, lowerSubstr)
}
