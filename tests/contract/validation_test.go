package contract

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"employee-management-system/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// TestValidationContract tests all validation rules through GraphQL API.
// This is a RED phase test - it MUST fail before implementation exists.
// Tests validate API contract and expected validation behavior before actual implementation.
func TestValidationContract(t *testing.T) {
	// Setup test server and GraphQL client
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all validation tests
	t.Run("Employee Field Validation", func(t *testing.T) {
		testEmployeeFieldValidation(t, client)
	})

	t.Run("Employee Email Uniqueness Validation", func(t *testing.T) {
		testEmployeeEmailUniqueness(t, client)
	})

	t.Run("Employee Manager Validation", func(t *testing.T) {
		testEmployeeManagerValidation(t, client)
	})

	t.Run("Employee Status Transition Validation", func(t *testing.T) {
		testEmployeeStatusTransitionValidation(t, client)
	})

	t.Run("Employee Salary Validation", func(t *testing.T) {
		testEmployeeSalaryValidation(t, client)
	})

	t.Run("Employee Hire Date Validation", func(t *testing.T) {
		testEmployeeHireDateValidation(t, client)
	})

	t.Run("Employee Address Validation", func(t *testing.T) {
		testEmployeeAddressValidation(t, client)
	})

	t.Run("User Field Validation", func(t *testing.T) {
		testUserFieldValidation(t, client)
	})

	t.Run("User Email Uniqueness Validation", func(t *testing.T) {
		testUserEmailUniqueness(t, client)
	})

	t.Run("User Role Validation", func(t *testing.T) {
		testUserRoleValidation(t, client)
	})

	t.Run("Authentication Input Validation", func(t *testing.T) {
		testAuthenticationInputValidation(t, client)
	})
}

// testEmployeeFieldValidation validates individual employee field constraints
func testEmployeeFieldValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Required Fields Validation", func(t *testing.T) {
		// Test missing required fields
		query := `
		mutation {
			createEmployee(input: {
				firstName: ""
				lastName: "Doe"
				email: "test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
			}) {
				id
				firstName
				lastName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "required"), "required")
	})

	t.Run("String Length Validation", func(t *testing.T) {
		// Test firstName too short
		query := `
		mutation {
			createEmployee(input: {
				firstName: "A"
				lastName: "Doe"
				email: "test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "minimum length"), "minimum length")
	})

	t.Run("Email Format Validation", func(t *testing.T) {
		// Test invalid email format
		query := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "invalid-email"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "valid email"), "valid email")
	})

	t.Run("Numeric Range Validation", func(t *testing.T) {
		// Test negative salary
		query := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: -1000
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "positive"), "positive")
	})
}

// testEmployeeEmailUniqueness validates email uniqueness constraint
func testEmployeeEmailUniqueness(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Duplicate Email Prevention", func(t *testing.T) {
		// Create first employee
		createQuery := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "duplicate@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
			}) {
				id
				email
			}
		}`

		resp1, _ := client.Execute(context.Background(), createQuery, nil)
		resp1.AssertNoError(t)

		// Try to create second employee with same email
		resp2, _ := client.Execute(context.Background(), createQuery, nil)
		resp2.AssertHasError(t)
		assert.Contains(t, resp2.AssertErrorContains(t, "unique"), "unique")
		assert.Contains(t, resp2.AssertErrorContains(t, "email"), "email")
	})
}

// testEmployeeManagerValidation validates manager relationship constraints
func testEmployeeManagerValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Non-existent Manager", func(t *testing.T) {
		query := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "manager-test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
				managerId: "00000000-0000-0000-0000-000000000000"
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "manager"), "manager")
		assert.Contains(t, resp.AssertErrorContains(t, "exists"), "exists")
	})
}

// testEmployeeStatusTransitionValidation validates status transition rules
func testEmployeeStatusTransitionValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Invalid Status Transitions", func(t *testing.T) {
		// Create an active employee
		createQuery := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "status-test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
			}) {
				id
				status
			}
		}`

		createResp, _ := client.Execute(context.Background(), createQuery, nil)
		createResp.AssertNoError(t)
		employeeID, _ := createResp.GetField("createEmployee.id")

		t.Run("From ACTIVE to ACTIVE (invalid)", func(t *testing.T) {
			query := fmt.Sprintf(`
			mutation {
				changeEmployeeStatus(id: "%s", status: ACTIVE) {
					id
					status
				}
			}`, employeeID)

			resp, _ := client.Execute(context.Background(), query, nil)
			resp.AssertHasError(t)
		})
	})
}

// testEmployeeSalaryValidation validates salary constraints
func testEmployeeSalaryValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Salary Range Validation", func(t *testing.T) {
		// Test zero salary
		query := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "salary-zero@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 0
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "positive"), "positive")

		// Test excessive salary
		query = `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "salary-max@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 2000000
			}) {
				id
				firstName
			}
		}`

		resp, _ = client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "maximum"), "maximum")
	})
}

// testEmployeeHireDateValidation validates date constraints
func testEmployeeHireDateValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Future Hire Date", func(t *testing.T) {
		// Test future hire date
		futureDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")
		query := fmt.Sprintf(`
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "future-date@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "%s"
				salary: 50000
			}) {
				id
				firstName
			}
		}`, futureDate)

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "future"), "future")
	})

	t.Run("Invalid Date Format", func(t *testing.T) {
		query := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "invalid-date@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "invalid-date"
				salary: 50000
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "date"), "date")
	})
}

// testEmployeeAddressValidation validates address field constraints
func testEmployeeAddressValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Address Length Validation", func(t *testing.T) {
		// Test overly long address
		longStreet := strings.Repeat("A", 201)
		query := fmt.Sprintf(`
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "address-test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
				address: {
					street: "%s"
					city: "City"
					state: "State"
				}
			}) {
				id
				firstName
			}
		}`, longStreet)

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "maximum length"), "maximum length")
	})

	t.Run("Country Code Validation", func(t *testing.T) {
		// Test invalid country code
		query := `
		mutation {
			createEmployee(input: {
				firstName: "John"
				lastName: "Doe"
				email: "country-test@example.com"
				department: "Engineering"
				position: "Developer"
				hireDate: "2023-01-01"
				salary: 50000
				address: {
					country: "USA" // Should be 2-letter code
					city: "City"
					state: "State"
				}
			}) {
				id
				firstName
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "2-letter"), "2-letter")
	})
}

// testUserFieldValidation validates user field constraints
func testUserFieldValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Username Validation", func(t *testing.T) {
		// Test username too short
		query := `
		mutation {
			createUser(input: {
				username: "ab"
				email: "user-test@example.com"
				password: "password123"
			}) {
				id
				username
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "minimum length"), "minimum length")

		// Test username with invalid characters
		query = `
		mutation {
			createUser(input: {
				username: "user@name"
				email: "user-test2@example.com"
				password: "password123"
			}) {
				id
				username
			}
		}`

		resp, _ = client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "alphanumeric"), "alphanumeric")
	})

	t.Run("Password Validation", func(t *testing.T) {
		// Test weak password
		query := `
		mutation {
			createUser(input: {
				username: "testuser"
				email: "user-test3@example.com"
				password: "123" // Too short
			}) {
				id
				username
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "minimum length"), "minimum length")

		// Test password complexity
		query = `
		mutation {
			createUser(input: {
				username: "testuser2"
				email: "user-test4@example.com"
				password: "password" // No numbers
			}) {
				id
				username
			}
		}`

		resp, _ = client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "complexity"), "complexity")
	})
}

// testUserEmailUniqueness validates user email uniqueness
func testUserEmailUniqueness(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Duplicate User Email", func(t *testing.T) {
		// Create first user
		createQuery := `
		mutation {
			createUser(input: {
				username: "user1"
				email: "duplicate-user@example.com"
				password: "password123"
				role: ADMIN
			}) {
				id
				email
			}
		}`

		resp1, _ := client.Execute(context.Background(), createQuery, nil)
		resp1.AssertNoError(t)

		// Try to create second user with same email
		resp2, _ := client.Execute(context.Background(), createQuery, nil)
		resp2.AssertHasError(t)
		assert.Contains(t, resp2.AssertErrorContains(t, "unique"), "unique")
		assert.Contains(t, resp2.AssertErrorContains(t, "email"), "email")
	})
}

// testUserRoleValidation validates role constraints
func testUserRoleValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Invalid Role", func(t *testing.T) {
		query := `
		mutation {
			createUser(input: {
				username: "invalidrole"
				email: "role-test@example.com"
				password: "password123"
				role: SUPER_ADMIN // Invalid role
			}) {
				id
				username
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "valid role"), "valid role")
	})
}

// testAuthenticationInputValidation validates authentication inputs
func testAuthenticationInputValidation(t *testing.T, client *helpers.GraphQLClient) {
	t.Run("Empty Login Credentials", func(t *testing.T) {
		query := `
		mutation {
			login(username: "", password: "") {
				token
				user {
					id
					username
				}
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "required"), "required")
	})

	t.Run("Weak Password", func(t *testing.T) {
		query := `
		mutation {
			login(username: "testuser", password: "123") {
				token
				user {
					id
					username
				}
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "weak password"), "weak password")
	})

	t.Run("Invalid Refresh Token", func(t *testing.T) {
		query := `
		mutation {
			refreshToken(token: "invalid-token-format") {
				token
				expiresAt
			}
		}`

		resp, _ := client.Execute(context.Background(), query, nil)
		resp.AssertHasError(t)
		assert.Contains(t, resp.AssertErrorContains(t, "invalid"), "invalid")
	})
}
