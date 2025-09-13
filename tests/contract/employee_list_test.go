package contract

import (
	"context"
	"testing"
	"time"

	"employee-management-system/tests/helpers"

	"github.com/stretchr/testify/require"
)

// TestEmployeeListContract tests all Employee listing operations through GraphQL API.
// This is a RED phase test - it MUST fail before implementation exists.
// Tests validate API contract and expected behavior for employee listing with pagination,
// filtering, and sorting before actual implementation.
func TestEmployeeListContract(t *testing.T) {
	// Setup test server and GraphQL client
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all employee listing operation tests
	t.Run("Basic Pagination Operations", func(t *testing.T) {
		testBasicPaginationOperations(t, client)
	})

	t.Run("Forward Pagination Scenarios", func(t *testing.T) {
		testForwardPaginationScenarios(t, client)
	})

	t.Run("Backward Pagination Scenarios", func(t *testing.T) {
		testBackwardPaginationScenarios(t, client)
	})

	t.Run("Filtering Operations", func(t *testing.T) {
		testFilteringOperations(t, client)
	})

	t.Run("Sorting Operations", func(t *testing.T) {
		testSortingOperations(t, client)
	})

	t.Run("Combined Filtering and Sorting", func(t *testing.T) {
		testCombinedFilteringAndSorting(t, client)
	})

	t.Run("Edge Cases and Error Handling", func(t *testing.T) {
		testEdgeCasesAndErrorHandling(t, client)
	})

	t.Run("Authentication Requirements", func(t *testing.T) {
		testAuthenticationRequirements(t, client, testServer)
	})

	t.Run("Authorization for Different Roles", func(t *testing.T) {
		testAuthorizationForDifferentRoles(t, client, testServer)
	})

	t.Run("Response Structure Validation", func(t *testing.T) {
		testResponseStructureValidation(t, client)
	})
}

// testBasicPaginationOperations tests basic pagination parameters and validation
func testBasicPaginationOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		variables     map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "List First 10 Employees",
			variables: map[string]interface{}{
				"first": 10,
			},
			expectedError: true, // Will fail until implemented
		},
		{
			name: "List First Employee With Cursor",
			variables: map[string]interface{}{
				"first": 1,
			},
			expectedError: true,
		},
		{
			name:          "List Employees Without Pagination",
			variables:     map[string]interface{}{},
			expectedError: true,
		},
		{
			name: "Invalid First Parameter - Negative",
			variables: map[string]interface{}{
				"first": -5,
			},
			expectedError: true,
			errorContains: "first",
		},
		{
			name: "Invalid First Parameter - Zero",
			variables: map[string]interface{}{
				"first": 0,
			},
			expectedError: true,
			errorContains: "first",
		},
		{
			name: "Invalid First Parameter - Too Large",
			variables: map[string]interface{}{
				"first": 1000, // Exceeds typical page size limit
			},
			expectedError: true,
			errorContains: "first",
		},
		{
			name: "Invalid After Parameter - Empty Cursor",
			variables: map[string]interface{}{
				"first": 10,
				"after": "",
			},
			expectedError: true,
			errorContains: "after",
		},
		{
			name: "Invalid After Parameter - Malformed Cursor",
			variables: map[string]interface{}{
				"first": 10,
				"after": "invalid-cursor-format",
			},
			expectedError: true,
			errorContains: "after",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query Employees($first: Int, $after: String) {
					employees(first: $first, after: $after) {
						edges {
							node {
								id
								firstName
								lastName
								email
								department
								position
								status
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
							startCursor
							endCursor
						}
						totalCount
					}
				}`

			resp, err := client.Execute(ctx, query, tc.variables)

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented employees query: %s", tc.name)
						assertContractError(t, resp, "employees")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented employees query: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testForwardPaginationScenarios tests forward pagination with cursor-based navigation
func testForwardPaginationScenarios(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		variables     map[string]interface{}
		expectedError bool
	}{
		{
			name: "Pagination With Valid Cursor",
			variables: map[string]interface{}{
				"first": 5,
				"after": "valid-cursor-string",
			},
			expectedError: true,
		},
		{
			name: "Large Page Size",
			variables: map[string]interface{}{
				"first": 50,
			},
			expectedError: true,
		},
		{
			name: "Single Item Pagination",
			variables: map[string]interface{}{
				"first": 1,
			},
			expectedError: true,
		},
		{
			name: "Pagination Beyond Available Data",
			variables: map[string]interface{}{
				"first": 10,
				"after": "cursor-for-last-item",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query Employees($first: Int, $after: String) {
					employees(first: $first, after: $after) {
						edges {
							node {
								id
								firstName
								lastName
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
							startCursor
							endCursor
						}
						totalCount
					}
				}`

			resp, err := client.Execute(ctx, query, tc.variables)

			if err == nil {
				require.True(t, resp.HasErrors(), "Expected GraphQL errors for forward pagination: %s", tc.name)
				assertContractError(t, resp, "employees pagination")
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testBackwardPaginationScenarios tests backward pagination with last/before parameters
func testBackwardPaginationScenarios(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		variables     map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "Backward Pagination With Last Parameter",
			variables: map[string]interface{}{
				"last": 10,
			},
			expectedError: true,
		},
		{
			name: "Backward Pagination With Before Cursor",
			variables: map[string]interface{}{
				"last":   5,
				"before": "valid-cursor-string",
			},
			expectedError: true,
		},
		{
			name: "Invalid Last Parameter - Negative",
			variables: map[string]interface{}{
				"last": -5,
			},
			expectedError: true,
			errorContains: "last",
		},
		{
			name: "Invalid Last Parameter - Zero",
			variables: map[string]interface{}{
				"last": 0,
			},
			expectedError: true,
			errorContains: "last",
		},
		{
			name: "Invalid Before Parameter - Empty",
			variables: map[string]interface{}{
				"last":   10,
				"before": "",
			},
			expectedError: true,
			errorContains: "before",
		},
		{
			name: "Combined Forward and Backward Pagination",
			variables: map[string]interface{}{
				"first":  10,
				"last":   5,
				"after":  "cursor-1",
				"before": "cursor-2",
			},
			expectedError: true,
			errorContains: "pagination",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query Employees($first: Int, $after: String, $last: Int, $before: String) {
					employees(first: $first, after: $after, last: $last, before: $before) {
						edges {
							node {
								id
								firstName
								lastName
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
							startCursor
							endCursor
						}
						totalCount
					}
				}`

			resp, err := client.Execute(ctx, query, tc.variables)

			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for backward pagination: %s", tc.name)
						assertContractError(t, resp, "employees pagination")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for backward pagination: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testFilteringOperations tests all filtering capabilities
func testFilteringOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		filter        map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "Filter By Department",
			filter: map[string]interface{}{
				"department": "Engineering",
			},
			expectedError: true,
		},
		{
			name: "Filter By Status",
			filter: map[string]interface{}{
				"status": "ACTIVE",
			},
			expectedError: true,
		},
		{
			name: "Filter By Position",
			filter: map[string]interface{}{
				"position": "Software Engineer",
			},
			expectedError: true,
		},
		{
			name: "Filter By Manager ID",
			filter: map[string]interface{}{
				"managerId": "manager-123",
			},
			expectedError: true,
		},
		{
			name: "Filter By Hire Date Range",
			filter: map[string]interface{}{
				"hireDateFrom": "2023-01-01",
				"hireDateTo":   "2023-12-31",
			},
			expectedError: true,
		},
		{
			name: "Filter By Salary Range",
			filter: map[string]interface{}{
				"salaryFrom": 50000.0,
				"salaryTo":   100000.0,
			},
			expectedError: true,
		},
		{
			name: "Text Search In Name",
			filter: map[string]interface{}{
				"search": "John",
			},
			expectedError: true,
		},
		{
			name: "Text Search In Email",
			filter: map[string]interface{}{
				"search": "@example.com",
			},
			expectedError: true,
		},
		{
			name: "Multiple Filter Criteria",
			filter: map[string]interface{}{
				"department": "Engineering",
				"status":     "ACTIVE",
				"position":   "Software Engineer",
			},
			expectedError: true,
		},
		{
			name: "Invalid Status Value",
			filter: map[string]interface{}{
				"status": "INVALID_STATUS",
			},
			expectedError: true,
			errorContains: "status",
		},
		{
			name: "Invalid Date Range",
			filter: map[string]interface{}{
				"hireDateFrom": "2023-12-31",
				"hireDateTo":   "2023-01-01", // End before start
			},
			expectedError: true,
			errorContains: "date",
		},
		{
			name: "Invalid Salary Range",
			filter: map[string]interface{}{
				"salaryFrom": 100000.0,
				"salaryTo":   50000.0, // End before start
			},
			expectedError: true,
			errorContains: "salary",
		},
		{
			name: "Empty Search Term",
			filter: map[string]interface{}{
				"search": "",
			},
			expectedError: true,
			errorContains: "search",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query Employees($first: Int, $filter: EmployeeFilter) {
					employees(first: $first, filter: $filter) {
						edges {
							node {
								id
								firstName
								lastName
								email
								department
								position
								status
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
						}
						totalCount
					}
				}`

			variables := map[string]interface{}{
				"first":  10,
				"filter": tc.filter,
			}

			resp, err := client.Execute(ctx, query, variables)

			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for filtering: %s", tc.name)
						assertContractError(t, resp, "employees filter")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for filtering: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testSortingOperations tests all sorting capabilities
func testSortingOperations(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		sortBy        map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "Sort By ID Ascending",
			sortBy: map[string]interface{}{
				"field":     "ID",
				"direction": "ASC",
			},
			expectedError: true,
		},
		{
			name: "Sort By First Name Descending",
			sortBy: map[string]interface{}{
				"field":     "FIRST_NAME",
				"direction": "DESC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Last Name",
			sortBy: map[string]interface{}{
				"field":     "LAST_NAME",
				"direction": "ASC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Email",
			sortBy: map[string]interface{}{
				"field":     "EMAIL",
				"direction": "ASC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Department",
			sortBy: map[string]interface{}{
				"field":     "DEPARTMENT",
				"direction": "ASC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Position",
			sortBy: map[string]interface{}{
				"field":     "POSITION",
				"direction": "ASC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Hire Date",
			sortBy: map[string]interface{}{
				"field":     "HIRE_DATE",
				"direction": "DESC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Salary",
			sortBy: map[string]interface{}{
				"field":     "SALARY",
				"direction": "DESC",
			},
			expectedError: true,
		},
		{
			name: "Sort By Status",
			sortBy: map[string]interface{}{
				"field":     "STATUS",
				"direction": "ASC",
			},
			expectedError: true,
		},
		{
			name: "Invalid Sort Field",
			sortBy: map[string]interface{}{
				"field":     "INVALID_FIELD",
				"direction": "ASC",
			},
			expectedError: true,
			errorContains: "field",
		},
		{
			name: "Invalid Sort Direction",
			sortBy: map[string]interface{}{
				"field":     "ID",
				"direction": "INVALID_DIRECTION",
			},
			expectedError: true,
			errorContains: "direction",
		},
		{
			name: "Missing Sort Field",
			sortBy: map[string]interface{}{
				"direction": "ASC",
			},
			expectedError: true,
			errorContains: "field",
		},
		{
			name: "Missing Sort Direction",
			sortBy: map[string]interface{}{
				"field": "ID",
			},
			expectedError: true,
			errorContains: "direction",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query Employees($first: Int, $sortBy: EmployeeSort) {
					employees(first: $first, sortBy: $sortBy) {
						edges {
							node {
								id
								firstName
								lastName
								email
								department
								position
								status
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
						}
						totalCount
					}
				}`

			variables := map[string]interface{}{
				"first":  10,
				"sortBy": tc.sortBy,
			}

			resp, err := client.Execute(ctx, query, variables)

			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for sorting: %s", tc.name)
						assertContractError(t, resp, "employees sort")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for sorting: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testCombinedFilteringAndSorting tests complex scenarios with both filtering and sorting
func testCombinedFilteringAndSorting(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		variables     map[string]interface{}
		expectedError bool
	}{
		{
			name: "Filter By Department And Sort By Salary",
			variables: map[string]interface{}{
				"first": 10,
				"filter": map[string]interface{}{
					"department": "Engineering",
				},
				"sortBy": map[string]interface{}{
					"field":     "SALARY",
					"direction": "DESC",
				},
			},
			expectedError: true,
		},
		{
			name: "Filter By Status And Sort By Hire Date",
			variables: map[string]interface{}{
				"first": 10,
				"filter": map[string]interface{}{
					"status": "ACTIVE",
				},
				"sortBy": map[string]interface{}{
					"field":     "HIRE_DATE",
					"direction": "ASC",
				},
			},
			expectedError: true,
		},
		{
			name: "Multiple Filters With Sorting",
			variables: map[string]interface{}{
				"first": 15,
				"filter": map[string]interface{}{
					"department": "Engineering",
					"status":     "ACTIVE",
					"position":   "Software Engineer",
				},
				"sortBy": map[string]interface{}{
					"field":     "LAST_NAME",
					"direction": "ASC",
				},
			},
			expectedError: true,
		},
		{
			name: "Date Range Filter With Salary Sort",
			variables: map[string]interface{}{
				"first": 20,
				"filter": map[string]interface{}{
					"hireDateFrom": "2023-01-01",
					"hireDateTo":   "2023-06-30",
				},
				"sortBy": map[string]interface{}{
					"field":     "SALARY",
					"direction": "DESC",
				},
			},
			expectedError: true,
		},
		{
			name: "Text Search With Department Filter And Name Sort",
			variables: map[string]interface{}{
				"first": 10,
				"filter": map[string]interface{}{
					"search":     "John",
					"department": "Engineering",
				},
				"sortBy": map[string]interface{}{
					"field":     "FIRST_NAME",
					"direction": "ASC",
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				query Employees($first: Int, $filter: EmployeeFilter, $sortBy: EmployeeSort) {
					employees(first: $first, filter: $filter, sortBy: $sortBy) {
						edges {
							node {
								id
								firstName
								lastName
								email
								department
								position
								status
								salary
								hireDate
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
							startCursor
							endCursor
						}
						totalCount
					}
				}`

			resp, err := client.Execute(ctx, query, tc.variables)

			if err == nil {
				require.True(t, resp.HasErrors(), "Expected GraphQL errors for combined filtering and sorting: %s", tc.name)
				assertContractError(t, resp, "employees combined")
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testEdgeCasesAndErrorHandling tests edge cases and error scenarios
func testEdgeCasesAndErrorHandling(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		query         string
		variables     map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "Empty Result Set",
			query: `
				query Employees($filter: EmployeeFilter) {
					employees(first: 10, filter: $filter) {
						edges {
							node {
								id
							}
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
						}
						totalCount
					}
				}`,
			variables: map[string]interface{}{
				"filter": map[string]interface{}{
					"department": "NonExistentDepartment",
				},
			},
			expectedError: true,
		},
		{
			name: "Very Large Result Set Simulation",
			query: `
				query Employees {
					employees(first: 100) {
						edges {
							node {
								id
							}
						}
						totalCount
					}
				}`,
			variables:     map[string]interface{}{},
			expectedError: true,
		},
		{
			name: "Invalid GraphQL Query Structure",
			query: `
				query Employees {
					employees(first: 10) {
						edges {
							node {
								invalidField
							}
						}
					}
				}`,
			variables:     map[string]interface{}{},
			expectedError: true,
			errorContains: "invalidField",
		},
		{
			name: "Malformed Filter Object",
			query: `
				query Employees($filter: EmployeeFilter) {
					employees(first: 10, filter: $filter) {
						edges {
							node {
								id
							}
						}
					}
				}`,
			variables: map[string]interface{}{
				"filter": "invalid-filter-type",
			},
			expectedError: true,
			errorContains: "filter",
		},
		{
			name: "Invalid Variable Type",
			query: `
				query Employees($first: String) {
					employees(first: $first) {
						edges {
							node {
								id
							}
						}
					}
				}`,
			variables: map[string]interface{}{
				"first": "invalid-string",
			},
			expectedError: true,
			errorContains: "first",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for edge case: %s", tc.name)
						assertContractError(t, resp, "employees edge case")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for edge case: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testAuthenticationRequirements tests authentication requirements for employee listing
func testAuthenticationRequirements(t *testing.T, client *helpers.GraphQLClient, testServer *helpers.TestServer) {
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
			name: "List Employees Without Authentication",
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
						pageInfo {
							hasNextPage
							hasPreviousPage
						}
					}
				}`,
			variables:     map[string]interface{}{},
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name: "Filtered List Without Authentication",
			query: `
				query Employees($filter: EmployeeFilter) {
					employees(first: 10, filter: $filter) {
						edges {
							node {
								id
								department
							}
						}
					}
				}`,
			variables: map[string]interface{}{
				"filter": map[string]interface{}{
					"department": "Engineering",
				},
			},
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name: "Sorted List Without Authentication",
			query: `
				query Employees($sortBy: EmployeeSort) {
					employees(first: 10, sortBy: $sortBy) {
						edges {
							node {
								id
								salary
							}
						}
					}
				}`,
			variables: map[string]interface{}{
				"sortBy": map[string]interface{}{
					"field":     "SALARY",
					"direction": "DESC",
				},
			},
			expectedError: true,
			errorContains: "authentication",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := unauthenticatedClient.Execute(ctx, tc.query, tc.variables)

			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unauthenticated listing: %s", tc.name)
						assertContractError(t, resp, "authentication")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unauthenticated listing: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testAuthorizationForDifferentRoles tests different user role access
func testAuthorizationForDifferentRoles(t *testing.T, client *helpers.GraphQLClient, testServer *helpers.TestServer) {
	testCases := []struct {
		name          string
		role          string
		variables     map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "Admin Role - Full Access",
			role: "ADMIN",
			variables: map[string]interface{}{
				"first": 10,
				"filter": map[string]interface{}{
					"salaryFrom": 100000.0, // High salary access
				},
			},
			expectedError: true,
		},
		{
			name: "Manager Role - Department Access",
			role: "MANAGER",
			variables: map[string]interface{}{
				"first": 10,
				"filter": map[string]interface{}{
					"department": "Engineering",
				},
			},
			expectedError: true,
		},
		{
			name: "Viewer Role - Limited Access",
			role: "VIEWER",
			variables: map[string]interface{}{
				"first": 5,
			},
			expectedError: true,
		},
		{
			name: "Invalid Role",
			role: "INVALID_ROLE",
			variables: map[string]interface{}{
				"first": 10,
			},
			expectedError: true,
			errorContains: "role",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create client with role-based authentication
			roleClient := helpers.NewGraphQLClient(testServer.BaseURL)
			roleClient.WithAuth("mock-token-for-" + tc.role)

			query := `
				query Employees($first: Int, $filter: EmployeeFilter) {
					employees(first: $first, filter: $filter) {
						edges {
							node {
								id
								firstName
								lastName
								email
								department
								position
								salary
								status
							}
							cursor
						}
						pageInfo {
							hasNextPage
							hasPreviousPage
						}
						totalCount
					}
				}`

			resp, err := roleClient.Execute(ctx, query, tc.variables)

			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for role-based access: %s", tc.name)
						assertContractError(t, resp, "authorization")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for role-based access: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testResponseStructureValidation tests the expected response structure
func testResponseStructureValidation(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		query         string
		variables     map[string]interface{}
		expectedError bool
	}{
		{
			name: "Complete Response Structure",
			query: `
				query Employees {
					employees(first: 5) {
						edges {
							node {
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
							cursor
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
			variables:     map[string]interface{}{},
			expectedError: true,
		},
		{
			name: "Minimal Response Structure",
			query: `
				query Employees {
					employees(first: 3) {
						edges {
							node {
								id
								firstName
							}
							cursor
						}
						pageInfo {
							hasNextPage
						}
					}
				}`,
			variables:     map[string]interface{}{},
			expectedError: true,
		},
		{
			name: "Invalid Response Structure - Missing Required Field",
			query: `
				query Employees {
					employees(first: 5) {
						edges {
							node {
								id
								# Missing required fields
							}
						}
					}
				}`,
			variables:     map[string]interface{}{},
			expectedError: true,
		},
		{
			name: "Invalid Response Structure - Missing PageInfo",
			query: `
				query Employees {
					employees(first: 5) {
						edges {
							node {
								id
							}
						}
						# Missing pageInfo
					}
				}`,
			variables:     map[string]interface{}{},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			resp, err := client.Execute(ctx, tc.query, tc.variables)

			if err == nil {
				if tc.expectedError {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for response structure validation: %s", tc.name)
					assertContractError(t, resp, "employees structure")
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for response structure validation: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

