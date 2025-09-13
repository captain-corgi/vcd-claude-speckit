package helpers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{}            `json:"data"`
	Errors []GraphQLError          `json:"errors,omitempty"`
}

// GraphQLError represents a GraphQL error
type GraphQLError struct {
	Message   string                 `json:"message"`
	Locations []GraphQLLocation      `json:"locations,omitempty"`
	Path      []string               `json:"path,omitempty"`
	Extensions map[string]interface{} `json:"extensions,omitempty"`
}

// GraphQLLocation represents a location in a GraphQL query
type GraphQLLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// GraphQLClient represents a GraphQL test client
type GraphQLClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewGraphQLClient creates a new GraphQL test client
func NewGraphQLClient(baseURL string) *GraphQLClient {
	return &GraphQLClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// WithAuth sets the authentication token for the client
func (c *GraphQLClient) WithAuth(token string) *GraphQLClient {
	c.AuthToken = token
	return c
}

// Execute executes a GraphQL query
func (c *GraphQLClient) Execute(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	reqBody := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal GraphQL request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/graphql", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	var result GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode GraphQL response: %w", err)
	}

	return &result, nil
}

// ExecuteMutation executes a GraphQL mutation
func (c *GraphQLClient) ExecuteMutation(ctx context.Context, mutation string, variables map[string]interface{}) (*GraphQLResponse, error) {
	return c.Execute(ctx, mutation, variables)
}

// ExecuteQuery executes a GraphQL query
func (c *GraphQLClient) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	return c.Execute(ctx, query, variables)
}

// HasErrors checks if the GraphQL response contains errors
func (r *GraphQLResponse) HasErrors() bool {
	return len(r.Errors) > 0
}

// FirstError returns the first error from the GraphQL response
func (r *GraphQLResponse) FirstError() *GraphQLError {
	if r.HasErrors() {
		return &r.Errors[0]
	}
	return nil
}

// GetField extracts a field from the GraphQL response data
func (r *GraphQLResponse) GetField(path string) (interface{}, error) {
	data, ok := r.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("response data is not an object")
	}

	// Simple path traversal - could be enhanced for nested paths
	if val, exists := data[path]; exists {
		return val, nil
	}

	return nil, fmt.Errorf("field '%s' not found in response", path)
}

// AssertNoError asserts that the GraphQL response has no errors
func (r *GraphQLResponse) AssertNoError(t *testing.T) *GraphQLResponse {
	t.Helper()
	if r.HasErrors() {
		t.Errorf("GraphQL response has errors: %v", r.Errors)
	}
	return r
}

// AssertHasError asserts that the GraphQL response has errors
func (r *GraphQLResponse) AssertHasError(t *testing.T) *GraphQLResponse {
	t.Helper()
	if !r.HasErrors() {
		t.Error("Expected GraphQL response to have errors, but none found")
	}
	return r
}

// AssertErrorContains asserts that the GraphQL response contains an error with the specified message
func (r *GraphQLResponse) AssertErrorContains(t *testing.T, message string) *GraphQLResponse {
	t.Helper()
	r.AssertHasError(t)

	found := false
	for _, err := range r.Errors {
		if err.Message == message {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Expected error containing '%s', but got: %v", message, r.Errors)
	}

	return r
}

// CreateGraphQLTestClient creates a GraphQL test client for testing
func CreateGraphQLTestClient(t *testing.T, baseURL string) *GraphQLClient {
	t.Helper()
	return NewGraphQLClient(baseURL)
}

// RunGraphQLTest is a helper function for GraphQL tests
func RunGraphQLTest(t *testing.T, baseURL string, testFunc func(*testing.T, *GraphQLClient)) {
	t.Helper()

	client := CreateGraphQLTestClient(t, baseURL)
	testFunc(t, client)
}

// CreateTestGraphQLRequest creates a GraphQL request structure for testing
func CreateTestGraphQLRequest(query string, variables map[string]interface{}) *GraphQLRequest {
	return &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}
}

// Common GraphQL queries and mutations for testing
const (
	// Employee queries
	QueryEmployees = `
		query Employees($first: Int, $after: String, $filter: EmployeeFilter, $sortBy: EmployeeSort) {
			employees(first: $first, after: $after, filter: $filter, sortBy: $sortBy) {
				edges {
					node {
						id
						firstName
						lastName
						email
						department
						position
						status
						createdAt
					}
				}
				pageInfo {
					hasNextPage
					hasPreviousPage
					startCursor
					endCursor
				}
			}
		}`

	QueryEmployee = `
		query Employee($id: ID!) {
			employee(id: $id) {
				id
				firstName
				lastName
				email
				department
				position
				salary
				status
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

	// Employee mutations
	MutationCreateEmployee = `
		mutation CreateEmployee($input: CreateEmployeeInput!) {
			createEmployee(input: $input) {
				id
				firstName
				lastName
				email
				department
				position
				status
				createdAt
			}
		}`

	MutationUpdateEmployee = `
		mutation UpdateEmployee($id: ID!, $input: UpdateEmployeeInput!) {
			updateEmployee(id: $id, input: $input) {
				id
				firstName
				lastName
				email
				updatedAt
			}
		}`

	MutationDeleteEmployee = `
		mutation DeleteEmployee($id: ID!) {
			deleteEmployee(id: $id)
		}`
)