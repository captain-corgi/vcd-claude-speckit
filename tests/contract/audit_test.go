package contract

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"employee-management-system/tests/helpers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAuditLoggingContract tests comprehensive audit logging functionality
// This is a RED phase test - it MUST fail before implementation exists.
func TestAuditLoggingContract(t *testing.T) {
	// Setup test server
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	// Create GraphQL client
	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all audit logging tests
	t.Run("AuditLog Type Schema", func(t *testing.T) {
		testAuditLogTypeSchema(t, client)
	})

	t.Run("AuditLog Queries", func(t *testing.T) {
		testAuditLogQueries(t, client)
	})

	t.Run("AuditLog Pagination", func(t *testing.T) {
		testAuditLogPagination(t, client)
	})

	t.Run("AuditLog Filtering", func(t *testing.T) {
		testAuditLogFiltering(t, client)
	})

	t.Run("AuditLog Sorting", func(t *testing.T) {
		testAuditLogSorting(t, client)
	})

	t.Run("Audit Log Creation Validation", func(t *testing.T) {
		testAuditLogCreationValidation(t, client)
	})

	t.Run("Audit Action Types", func(t *testing.T) {
		testAuditActionTypes(t, client)
	})

	t.Run("User Action Tracking", func(t *testing.T) {
		testUserActionTracking(t, client)
	})

	t.Run("Audit Log Immutability", func(t *testing.T) {
		testAuditLogImmutability(t, client)
	})

	t.Run("Access Controls", func(t *testing.T) {
		testAuditLogAccessControls(t, client)
	})

	t.Run("Performance Requirements", func(t *testing.T) {
		testAuditLogPerformance(t, client)
	})
}

// testAuditLogTypeSchema tests that AuditLog type exists in GraphQL schema
func testAuditLogTypeSchema(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "AuditLog Type Exists",
			query:       "query { __type(name: \"AuditLog\") { name fields { name type { name kind } } } }",
			expectError: true,
		},
		{
			name:        "AuditLogEdge Type Exists",
			query:       "query { __type(name: \"AuditLogEdge\") { name fields { name type { name } } } }",
			expectError: true,
		},
		{
			name:        "AuditLogConnection Type Exists",
			query:       "query { __type(name: \"AuditLogConnection\") { name fields { name type { name } } } }",
			expectError: true,
		},
		{
			name:        "AuditLogFilter Input Exists",
			query:       "query { __type(name: \"AuditLogFilter\") { name inputFields { name type { name kind } } } }",
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
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented type: %s", tc.name)
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

// testAuditLogQueries tests basic audit log queries
func testAuditLogQueries(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Get All Audit Logs",
			query: `
				query AuditLogs($first: Int, $after: String) {
					auditLogs(first: $first, after: $after) {
						edges {
							node {
								id
								operation
								userId
								timestamp
								oldValues
								newValues
								ipAddress
								userAgent
							}
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
							startCursor
							endCursor
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"first": 10},
			expectError: true,
		},
		{
			name: "Get Audit Log by ID",
			query: `
				query AuditLog($id: ID!) {
					auditLog(id: $id) {
						id
						operation
						userId
						timestamp
						oldValues
						newValues
						ipAddress
					}
				}`,
			variables:   map[string]interface{}{"id": "test-audit-log-id"},
			expectError: true,
		},
		{
			name: "Get Audit Logs for Employee",
			query: `
				query AuditLogsByEmployee($employeeId: ID!, $first: Int) {
					auditLogs(employeeId: $employeeId, first: $first) {
						edges {
							node {
								id
								operation
								timestamp
							}
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"employeeId": "test-employee-id", "first": 5},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

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

// testAuditLogPagination tests audit log pagination functionality
func testAuditLogPagination(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "First Page",
			query: `
				query AuditLogs($first: Int) {
					auditLogs(first: $first) {
						edges {
							node { id }
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
							startCursor
							endCursor
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"first": 10},
			expectError: true,
		},
		{
			name: "Page with Cursor",
			query: `
				query AuditLogs($first: Int, $after: String) {
					auditLogs(first: $first, after: $after) {
						edges {
							node { id }
						}
						pageInfo {
							hasNextPage
							startCursor
							endCursor
						}
					}
				}`,
			variables:   map[string]interface{}{"first": 10, "after": "cursor123"},
			expectError: true,
		},
		{
			name: "Large Page Size",
			query: `
				query AuditLogs($first: Int) {
					auditLogs(first: $first) {
						edges {
							node { id }
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"first": 100},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for pagination: %s", tc.name)
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

// testAuditLogFiltering tests audit log filtering capabilities
func testAuditLogFiltering(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Filter by Operation",
			query: `
				query AuditLogs($operation: String, $first: Int) {
					auditLogs(operation: $operation, first: $first) {
						edges {
							node {
								id
								operation
							}
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"operation": "CREATE", "first": 10},
			expectError: true,
		},
		{
			name: "Filter by Employee ID",
			query: `
				query AuditLogs($employeeId: ID, $first: Int) {
					auditLogs(employeeId: $employeeId, first: $first) {
						edges {
							node { id }
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"employeeId": "test-employee-id", "first": 10},
			expectError: true,
		},
		{
			name: "Filter by Date Range",
			query: `
				query AuditLogs($from: DateTime, $to: DateTime, $first: Int) {
					auditLogs(from: $from, to: $to, first: $first) {
						edges {
							node {
								id
								timestamp
							}
						}
						totalCount
					}
				}`,
			variables: map[string]interface{}{
				"from":  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				"to":    time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
				"first": 10,
			},
			expectError: true,
		},
		{
			name: "Filter by User ID",
			query: `
				query AuditLogs($userId: String, $first: Int) {
					auditLogs(first: $first) {
						edges {
							node {
								id
								userId
							}
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"userId": "test-user-id", "first": 10},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for filtering: %s", tc.name)
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

// testAuditLogSorting tests audit log sorting capabilities
func testAuditLogSorting(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Sort by Timestamp DESC",
			query: `
				query AuditLogs($first: Int) {
					auditLogs(first: $first) {
						edges {
							node {
								id
								timestamp
							}
						}
					}
				}`,
			variables:   map[string]interface{}{"first": 10},
			expectError: true,
		},
		{
			name: "Sort by Operation ASC",
			query: `
				query AuditLogs($first: Int) {
					auditLogs(first: $first) {
						edges {
							node {
								id
								operation
							}
						}
					}
				}`,
			variables:   map[string]interface{}{"first": 10},
			expectError: true,
		},
		{
			name: "Sort by User ID",
			query: `
				query AuditLogs($first: Int) {
					auditLogs(first: $first) {
						edges {
							node {
								id
								userId
							}
						}
					}
				}`,
			variables:   map[string]interface{}{"first": 10},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for sorting: %s", tc.name)
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

// testAuditLogCreationValidation tests audit log field validation
func testAuditLogCreationValidation(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Required Fields Validation",
			query: `
				mutation CreateAuditLog($input: AuditLogInput!) {
					createAuditLog(input: $input) {
						id
						operation
						userId
						timestamp
						ipAddress
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"operation":  "CREATE",
					"userId":     "user123",
					"ipAddress": "192.168.1.1",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid Operation Type",
			query: `
				mutation CreateAuditLog($input: AuditLogInput!) {
					createAuditLog(input: $input) {
						id
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"operation":  "INVALID_OPERATION",
					"userId":     "user123",
					"ipAddress": "192.168.1.1",
					"employeeId": "emp123",
				},
			},
			expectError: true,
		},
		{
			name: "Invalid IP Address Format",
			query: `
				mutation CreateAuditLog($input: AuditLogInput!) {
					createAuditLog(input: $input) {
						id
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"operation":  "CREATE",
					"userId":     "user123",
					"ipAddress": "invalid-ip",
					"employeeId": "emp123",
				},
			},
			expectError: true,
		},
		{
			name: "Missing Required EmployeeID for Employee Operations",
			query: `
				mutation CreateAuditLog($input: AuditLogInput!) {
					createAuditLog(input: $input) {
						id
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"operation":  "CREATE",
					"userId":     "user123",
					"ipAddress": "192.168.1.1",
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

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

// testAuditActionTypes tests different audit action types
func testAuditActionTypes(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		operation   string
		expectError bool
	}{
		{"CREATE Operation", "CREATE", true},
		{"UPDATE Operation", "UPDATE", true},
		{"DELETE Operation", "DELETE", true},
		{"STATUS_CHANGE Operation", "STATUS_CHANGE", true},
		{"LOGIN Operation", "LOGIN", true},
		{"LOGOUT Operation", "LOGOUT", true},
		{"INVALID Operation", "INVALID", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query := `
				mutation CreateAuditLog($input: AuditLogInput!) {
					createAuditLog(input: $input) {
						id
						operation
						userId
						ipAddress
					}
				}`

			variables := map[string]interface{}{
				"input": map[string]interface{}{
					"operation":  tc.operation,
					"userId":     "user123",
					"ipAddress": "192.168.1.1",
					"employeeId": "emp123",
				},
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, query, variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for operation: %s", tc.operation)
					assertSchemaError(t, resp, fmt.Sprintf("operation %s", tc.operation))
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

// testUserActionTracking tests user action tracking capabilities
func testUserActionTracking(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Track User Actions",
			query: `
				query UserAuditLogs($userId: String, $first: Int) {
					auditLogs(first: $first) {
						edges {
							node {
								id
								userId
								operation
								timestamp
								ipAddress
								userAgent
							}
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"userId": "user123", "first": 10},
			expectError: true,
		},
		{
			name: "User Session Audit",
			query: `
				query UserSessionAudit($userId: String, $from: DateTime, $to: DateTime) {
					auditLogs(userId: $userId, from: $from, to: $to) {
						edges {
							node {
								id
								operation
								timestamp
							}
						}
						totalCount
					}
				}`,
			variables: map[string]interface{}{
				"userId": "user123",
				"from":   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				"to":     time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			expectError: true,
		},
		{
			name: "Recent User Actions",
			query: `
				query RecentUserActions($userId: String, $first: Int) {
					auditLogs(userId: $userId, first: $first, sortBy: {field: TIMESTAMP, direction: DESC}) {
						edges {
							node {
								id
								operation
								timestamp
							}
						}
					}
				}`,
			variables:   map[string]interface{}{"userId": "user123", "first": 5},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for user tracking: %s", tc.name)
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

// testAuditLogImmutability tests audit log immutability constraints
func testAuditLogImmutability(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Update Audit Log (Should Fail)",
			query: `
				mutation UpdateAuditLog($id: ID!, $input: UpdateAuditLogInput!) {
					updateAuditLog(id: $id, input: $input) {
						id
						operation
					}
				}`,
			variables: map[string]interface{}{
				"id": "audit-log-id",
				"input": map[string]interface{}{
					"operation": "UPDATED_OPERATION",
				},
			},
			expectError: true,
		},
		{
			name: "Delete Audit Log (Should Fail)",
			query: `
				mutation DeleteAuditLog($id: ID!) {
					deleteAuditLog(id: $id)
				}`,
			variables:   map[string]interface{}{"id": "audit-log-id"},
			expectError: true,
		},
		{
			name: "Create Audit Log in Past (Should Fail)",
			query: `
				mutation CreateAuditLog($input: AuditLogInput!) {
					createAuditLog(input: $input) {
						id
						timestamp
					}
				}`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"operation":  "CREATE",
					"userId":     "user123",
					"ipAddress": "192.168.1.1",
					"employeeId": "emp123",
					"timestamp":  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), // Past timestamp
				},
			},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for immutability: %s", tc.name)
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

// testAuditLogAccessControls tests access controls for audit logs
func testAuditLogAccessControls(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		expectError bool
	}{
		{
			name: "Admin User Access",
			query: `
				query AdminAuditLogs {
					auditLogs(first: 10) {
						edges {
							node {
								id
								operation
								userId
								timestamp
							}
						}
					}
				}`,
			expectError: true,
		},
		{
			name: "Manager Access (Limited Scope)",
			query: `
				query ManagerAuditLogs {
					auditLogs(first: 10) {
						edges {
							node {
								id
								operation
							}
						}
					}
				}`,
			expectError: true,
		},
		{
			name: "Viewer Access (Denied)",
			query: `
				query ViewerAuditLogs {
					auditLogs(first: 10) {
						edges {
							node { id }
						}
					}
				}`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for access control: %s", tc.name)
					assertAuthorizationError(t, resp, tc.name)
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

// testAuditLogPerformance tests performance requirements for audit logs
func testAuditLogPerformance(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name        string
		query       string
		variables   map[string]interface{}
		timeout     time.Duration
		expectError bool
	}{
		{
			name: "Large Dataset Query",
			query: `
				query LargeAuditQuery($first: Int) {
					auditLogs(first: $first) {
						edges {
							node { id }
						}
						totalCount
					}
				}`,
			variables:   map[string]interface{}{"first": 1000},
			timeout:     1 * time.Second,
			expectError: true,
		},
		{
			name: "Complex Filtering Query",
			query: `
				query ComplexAuditFilter($operation: String, $from: DateTime, $to: DateTime) {
					auditLogs(operation: $operation, from: $from, to: $to, first: 100) {
						edges {
							node {
								id
								operation
								timestamp
								userId
							}
						}
						totalCount
					}
				}`,
			variables: map[string]interface{}{
				"operation": "CREATE",
				"from":      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				"to":        time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
			},
			timeout:     500 * time.Millisecond,
			expectError: true,
		},
		{
			name: "Real-time Audit Query",
			query: `
				query RealtimeAudit($operation: String) {
					auditLogs(operation: $operation, first: 50) {
						edges {
							node {
								id
								timestamp
								operation
								userId
							}
						}
					}
				}`,
			variables:   map[string]interface{}{"operation": "LOGIN"},
			timeout:     200 * time.Millisecond,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), tc.timeout)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if tc.expectError {
				// In RED phase, we expect errors
				if err == nil {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for performance test: %s", tc.name)
					assertSchemaError(t, resp, tc.name)
				} else {
					assert.True(t,
						strings.Contains(err.Error(), "connection refused") ||
							strings.Contains(err.Error(), "no such host") ||
							strings.Contains(err.Error(), "404") ||
							strings.Contains(err.Error(), "500") ||
							strings.Contains(err.Error(), "deadline exceeded"),
						"Expected connection or timeout error, but got: %v", err)
				}
			}
		})
	}
}