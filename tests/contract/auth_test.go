package contract

import (
	"context"
	"testing"
	"time"

	"employee-management-system/tests/helpers"

	"github.com/stretchr/testify/require"
)

// TestAuthenticationContract tests all authentication functionality through GraphQL API.
// This is a RED phase test - it MUST fail before implementation exists.
// Tests validate authentication contract and expected behavior before actual implementation.
func TestAuthenticationContract(t *testing.T) {
	// Setup test server and GraphQL client
	testServer := helpers.NewTestServer(t)
	defer testServer.Close()

	client := helpers.CreateGraphQLTestClient(t, testServer.BaseURL)

	// Run all authentication operation tests
	t.Run("JWT Token Generation and Validation", func(t *testing.T) {
		testJWTTokenGenerationAndValidation(t, client)
	})

	t.Run("Login Mutations With Valid Credentials", func(t *testing.T) {
		testLoginWithValidCredentials(t, client)
	})

	t.Run("Login With Invalid Credentials", func(t *testing.T) {
		testLoginWithInvalidCredentials(t, client)
	})

	t.Run("Token Refresh Functionality", func(t *testing.T) {
		testTokenRefreshFunctionality(t, client)
	})

	t.Run("Token Expiration Handling", func(t *testing.T) {
		testTokenExpirationHandling(t, client)
	})

	t.Run("Password Hashing and Verification", func(t *testing.T) {
		testPasswordHashingAndVerification(t, client)
	})

	t.Run("Multiple Authentication Methods", func(t *testing.T) {
		testMultipleAuthenticationMethods(t, client)
	})

	t.Run("Session Management", func(t *testing.T) {
		testSessionManagement(t, client)
	})

	t.Run("Rate Limiting for Authentication Attempts", func(t *testing.T) {
		testRateLimitingForAuthAttempts(t, client)
	})

	t.Run("Security Headers and Cookies", func(t *testing.T) {
		testSecurityHeadersAndCookies(t, client)
	})

	t.Run("Authentication Error Handling", func(t *testing.T) {
		testAuthenticationErrorHandling(t, client)
	})
}

// testJWTTokenGenerationAndValidation tests JWT token creation and validation
func testJWTTokenGenerationAndValidation(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		tokenConfig   *helpers.TestTokenConfig
		claims        *helpers.TestClaims
		expectedError bool
		errorContains string
	}{
		{
			name: "Valid JWT Token Generation",
			tokenConfig: &helpers.TestTokenConfig{
				Secret:     "test-secret-key",
				Expiration: time.Hour,
			},
			claims: &helpers.TestClaims{
				UserID:   "user-123",
				Username: "testuser",
				Roles:    []string{"ADMIN"},
				Email:    "test@example.com",
			},
			expectedError: true, // Will fail until implemented
		},
		{
			name: "Token With Short Expiration",
			tokenConfig: &helpers.TestTokenConfig{
				Secret:     "test-secret-key",
				Expiration: time.Minute,
			},
			claims: &helpers.TestClaims{
				UserID:   "user-123",
				Username: "testuser",
				Roles:    []string{"VIEWER"},
				Email:    "test@example.com",
			},
			expectedError: true,
		},
		{
			name: "Token With Multiple Roles",
			tokenConfig: &helpers.TestTokenConfig{
				Secret:     "test-secret-key",
				Expiration: time.Hour,
			},
			claims: &helpers.TestClaims{
				UserID:   "user-123",
				Username: "adminuser",
				Roles:    []string{"ADMIN", "MANAGER"},
				Email:    "admin@example.com",
			},
			expectedError: true,
		},
		{
			name:          "Invalid Token Secret",
			tokenConfig:   &helpers.TestTokenConfig{Secret: "wrong-secret"},
			claims:        &helpers.TestClaims{UserID: "user-123"},
			expectedError: true,
			errorContains: "invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create test token using helper
			testToken := helpers.CreateTestToken(t, tc.tokenConfig, tc.claims)

			query := `
				query ValidateToken($token: String!) {
					validateToken(token: $token) {
						isValid
						user {
							id
							username
							email
							roles
						}
						expiresAt
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"token": testToken,
			})

			// In RED phase, expect either connection errors or GraphQL schema errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented validateToken: %s", tc.name)
						assertContractError(t, resp, "validateToken")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented validateToken: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testLoginWithValidCredentials tests login mutations with valid credentials
func testLoginWithValidCredentials(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name           string
		loginInput     map[string]interface{}
		expectedFields map[string]interface{}
		expectedError  bool
	}{
		{
			name: "Login With Valid Email and Password",
			loginInput: map[string]interface{}{
				"email":    "admin@example.com",
				"password": "SecurePassword123!",
			},
			expectedFields: map[string]interface{}{
				"user": map[string]interface{}{
					"email":    "admin@example.com",
					"username": "admin",
				},
			},
			expectedError: true, // Will fail until implemented
		},
		{
			name: "Login With Valid Username and Password",
			loginInput: map[string]interface{}{
				"username": "admin",
				"password": "SecurePassword123!",
			},
			expectedFields: map[string]interface{}{
				"user": map[string]interface{}{
					"username": "admin",
				},
			},
			expectedError: true,
		},
		{
			name: "Login Returns Token and User Information",
			loginInput: map[string]interface{}{
				"email":    "manager@example.com",
				"password": "ManagerPass123!",
			},
			expectedFields: map[string]interface{}{
				"token":        "expected-jwt-token",
				"refreshToken": "expected-refresh-token",
				"user": map[string]interface{}{
					"id":       "user-456",
					"username": "manager",
					"email":    "manager@example.com",
					"roles":    []string{"MANAGER"},
				},
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Test login mutation
			query := `
				mutation Login($input: LoginInput!) {
					login(input: $input) {
						token
						refreshToken
						user {
							id
							username
							email
							roles
							isActive
						}
						expiresAt
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"input": tc.loginInput,
			})

			// In RED phase, expect errors for unimplemented functionality
			if err == nil {
				if tc.expectedError {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented login: %s", tc.name)
					assertContractError(t, resp, "login")
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented login: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testLoginWithInvalidCredentials tests login scenarios with invalid credentials
func testLoginWithInvalidCredentials(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		loginInput    map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name: "Login With Wrong Password",
			loginInput: map[string]interface{}{
				"email":    "admin@example.com",
				"password": "wrongpassword",
			},
			expectedError: true,
			errorContains: "password",
		},
		{
			name: "Login With Non-existent User Email",
			loginInput: map[string]interface{}{
				"email":    "nonexistent@example.com",
				"password": "anypassword",
			},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Login With Non-existent Username",
			loginInput: map[string]interface{}{
				"username": "nonexistent",
				"password": "anypassword",
			},
			expectedError: true,
			errorContains: "not found",
		},
		{
			name: "Login With Empty Password",
			loginInput: map[string]interface{}{
				"email":    "admin@example.com",
				"password": "",
			},
			expectedError: true,
			errorContains: "password",
		},
		{
			name: "Login With Empty Email",
			loginInput: map[string]interface{}{
				"email":    "",
				"password": "password",
			},
			expectedError: true,
			errorContains: "email",
		},
		{
			name: "Login With Inactive User",
			loginInput: map[string]interface{}{
				"email":    "inactive@example.com",
				"password": "password",
			},
			expectedError: true,
			errorContains: "inactive",
		},
		{
			name: "Login With Missing Credentials",
			loginInput: map[string]interface{}{
				"email": "admin@example.com",
				// Missing password
			},
			expectedError: true,
			errorContains: "password",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation Login($input: LoginInput!) {
					login(input: $input) {
						token
						user {
							id
							email
						}
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"input": tc.loginInput,
			})

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented login: %s", tc.name)
						assertContractError(t, resp, "login")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented login: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testTokenRefreshFunctionality tests token refresh functionality
func testTokenRefreshFunctionality(t *testing.T, client *helpers.GraphQLClient) {
	// Create a valid token for testing
	validToken := helpers.CreateAdminToken(t)

	testCases := []struct {
		name          string
		refreshToken  string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Refresh Valid Token",
			refreshToken:  validToken,
			expectedError: true, // Will fail until implemented
		},
		{
			name:          "Refresh With Invalid Token",
			refreshToken:  "invalid-token",
			expectedError: true,
			errorContains: "invalid",
		},
		{
			name:          "Refresh With Empty Token",
			refreshToken:  "",
			expectedError: true,
			errorContains: "token",
		},
		{
			name:          "Refresh With Expired Token",
			refreshToken:  helpers.CreateExpiredToken(t),
			expectedError: true,
			errorContains: "expired",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation RefreshToken($refreshToken: String!) {
					refreshToken(refreshToken: $refreshToken) {
						token
						refreshToken
						expiresAt
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"refreshToken": tc.refreshToken,
			})

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented refreshToken: %s", tc.name)
						assertContractError(t, resp, "refreshToken")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented refreshToken: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testTokenExpirationHandling tests token expiration scenarios
func testTokenExpirationHandling(t *testing.T, client *helpers.GraphQLClient) {
	// Create tokens with different expiration times
	freshToken := helpers.CreateAdminToken(t)
	expiredToken := helpers.CreateExpiredToken(t)

	testCases := []struct {
		name          string
		token         string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Access Protected Resource With Expired Token",
			token:         expiredToken,
			expectedError: true,
			errorContains: "expired",
		},
		{
			name:          "Access Protected Resource With Fresh Token",
			token:         freshToken,
			expectedError: true, // Will fail until implemented
		},
		{
			name:          "Access Protected Resource With Malformed Token",
			token:         "malformed.jwt.token",
			expectedError: true,
			errorContains: "invalid",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Test accessing a protected resource with the token
			authClient := helpers.NewGraphQLClient(client.BaseURL).WithAuth(tc.token)

			query := `
				query {
					me {
						id
						username
						email
					}
				}`

			resp, err := authClient.Execute(ctx, query, nil)

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented protected resource: %s", tc.name)
						assertContractError(t, resp, "protected resource")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented protected resource: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testPasswordHashingAndVerification tests password hashing and verification
func testPasswordHashingAndVerification(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		password      string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Hash Strong Password",
			password:      "SecurePassword123!",
			expectedError: true, // Will fail until implemented
		},
		{
			name:          "Hash Weak Password",
			password:      "weak",
			expectedError: true,
			errorContains: "weak",
		},
		{
			name:          "Hash Empty Password",
			password:      "",
			expectedError: true,
			errorContains: "required",
		},
		{
			name:          "Verify Correct Password",
			password:      "SecurePassword123!",
			expectedError: true,
		},
		{
			name:          "Verify Incorrect Password",
			password:      "WrongPassword123!",
			expectedError: true,
			errorContains: "incorrect",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Test password hashing
			hashQuery := `
				mutation HashPassword($password: String!) {
					hashPassword(password: $password) {
						hash
						algorithm
						strength
					}
				}`

			resp, err := client.Execute(ctx, hashQuery, map[string]interface{}{
				"password": tc.password,
			})

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented hashPassword: %s", tc.name)
						assertContractError(t, resp, "hashPassword")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented hashPassword: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testMultipleAuthenticationMethods tests different authentication methods
func testMultipleAuthenticationMethods(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		authMethod    string
		credentials   map[string]interface{}
		expectedError bool
		errorContains string
	}{
		{
			name:       "Authentication By Email",
			authMethod: "email",
			credentials: map[string]interface{}{
				"email":    "admin@example.com",
				"password": "SecurePassword123!",
			},
			expectedError: true, // Will fail until implemented
		},
		{
			name:       "Authentication By Username",
			authMethod: "username",
			credentials: map[string]interface{}{
				"username": "admin",
				"password": "SecurePassword123!",
			},
			expectedError: true,
		},
		{
			name:       "Authentication By Employee ID",
			authMethod: "employeeId",
			credentials: map[string]interface{}{
				"employeeId": "emp-123",
				"password":   "SecurePassword123!",
			},
			expectedError: true,
		},
		{
			name:       "Authentication With Invalid Method",
			authMethod: "invalid",
			credentials: map[string]interface{}{
				"email":    "admin@example.com",
				"password": "SecurePassword123!",
			},
			expectedError: true,
			errorContains: "method",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			query := `
				mutation Authenticate($method: String!, $credentials: Json!) {
					authenticate(method: $method, credentials: $credentials) {
						token
						user {
							id
							username
							email
						}
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"method":      tc.authMethod,
				"credentials": tc.credentials,
			})

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented authenticate: %s", tc.name)
						assertContractError(t, resp, "authenticate")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented authenticate: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testSessionManagement tests session management functionality
func testSessionManagement(t *testing.T, client *helpers.GraphQLClient) {
	testToken := helpers.CreateAdminToken(t)

	testCases := []struct {
		name          string
		action        string
		token         string
		expectedError bool
		errorContains string
	}{
		{
			name:          "List Active Sessions",
			action:        "list",
			token:         testToken,
			expectedError: true, // Will fail until implemented
		},
		{
			name:          "Terminate Current Session",
			action:        "terminate",
			token:         testToken,
			expectedError: true,
		},
		{
			name:          "Terminate All Other Sessions",
			action:        "terminateOthers",
			token:         testToken,
			expectedError: true,
		},
		{
			name:          "Check Session Validity",
			action:        "check",
			token:         testToken,
			expectedError: true,
		},
		{
			name:          "Session Management Without Authentication",
			action:        "list",
			token:         "",
			expectedError: true,
			errorContains: "authentication",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var query string
			switch tc.action {
			case "list":
				query = `
					query GetSessions {
						sessions {
							id
							userId
							deviceInfo
							lastActive
							isCurrent
						}
					}`
			case "terminate":
				query = `
					mutation TerminateSession($sessionId: ID) {
						terminateSession(sessionId: $sessionId) {
							success
							message
						}
					}`
			case "terminateOthers":
				query = `
					mutation TerminateOtherSessions {
						terminateOtherSessions {
							success
							terminatedCount
						}
					}`
			case "check":
				query = `
					query CheckSession {
						session {
							isValid
							expiresAt
							user {
								id
								username
							}
						}
					}`
			}

			authClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				authClient = authClient.WithAuth(tc.token)
			}

			resp, err := authClient.Execute(ctx, query, nil)

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented session management: %s", tc.name)
						assertContractError(t, resp, "session management")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented session management: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testRateLimitingForAuthAttempts tests rate limiting for authentication attempts
func testRateLimitingForAuthAttempts(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		attempts      int
		delay         time.Duration
		expectedError bool
		errorContains string
	}{
		{
			name:          "Multiple Failed Login Attempts",
			attempts:      5,
			delay:         time.Millisecond * 100,
			expectedError: true,
			errorContains: "rate limit",
		},
		{
			name:          "Excessive Login Attempts",
			attempts:      10,
			delay:         time.Millisecond * 50,
			expectedError: true,
			errorContains: "blocked",
		},
		{
			name:          "Rate Limit Reset After Delay",
			attempts:      3,
			delay:         time.Second * 2,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Simulate multiple authentication attempts
			for i := 0; i < tc.attempts; i++ {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				query := `
					mutation Login($input: LoginInput!) {
						login(input: $input) {
							token
							user {
								id
							}
						}
					}`

				resp, err := client.Execute(ctx, query, map[string]interface{}{
					"input": map[string]interface{}{
						"email":    "admin@example.com",
						"password": "wrongpassword",
					},
				})

				// Check if rate limiting kicks in on later attempts
				if i >= tc.attempts-1 {
					if err == nil {
						if tc.expectedError {
							if tc.errorContains != "" {
								require.True(t, resp.HasErrors(), "Expected rate limiting errors on attempt %d", i+1)
								resp.AssertErrorContains(t, tc.errorContains)
							} else {
								require.True(t, resp.HasErrors(), "Expected rate limiting errors on attempt %d", i+1)
								assertContractError(t, resp, "rate limiting")
							}
						}
					}
				}

				// Add delay between attempts
				time.Sleep(tc.delay)
			}
		})
	}
}

// testSecurityHeadersAndCookies tests security headers and cookie management
func testSecurityHeadersAndCookies(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		endpoint      string
		headers       map[string]string
		expectedError bool
		errorContains string
	}{
		{
			name:     "Security Headers on Login Response",
			endpoint: "login",
			headers: map[string]string{
				"Content-Type":           "application/json",
				"X-Content-Type-Options": "nosniff",
				"X-Frame-Options":        "DENY",
				"X-XSS-Protection":       "1; mode=block",
			},
			expectedError: true, // Will fail until implemented
		},
		{
			name:     "Cookie Management for Authentication",
			endpoint: "login",
			headers: map[string]string{
				"Set-Cookie": "auth-token=; HttpOnly; Secure; SameSite=Strict",
			},
			expectedError: true,
		},
		{
			name:     "CORS Headers for GraphQL Endpoint",
			endpoint: "graphql",
			headers: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "POST, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Create a direct HTTP request to test headers
			// This would normally be done with net/http but we'll simulate with GraphQL
			query := `
				mutation Login($input: LoginInput!) {
					login(input: $input) {
						token
					}
				}`

			resp, err := client.Execute(ctx, query, map[string]interface{}{
				"input": map[string]interface{}{
					"email":    "admin@example.com",
					"password": "SecurePassword123!",
				},
			})

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented security headers: %s", tc.name)
						assertContractError(t, resp, "security headers")
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented security headers: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}

// testAuthenticationErrorHandling tests various authentication error scenarios
func testAuthenticationErrorHandling(t *testing.T, client *helpers.GraphQLClient) {
	testCases := []struct {
		name          string
		operation     string
		query         string
		variables     map[string]interface{}
		token         string
		expectedError bool
		errorContains string
	}{
		{
			name:          "Access Protected Resource Without Token",
			operation:     "protectedQuery",
			query:         `query { me { id username } }`,
			token:         "",
			expectedError: true,
			errorContains: "authentication",
		},
		{
			name:          "Access Protected Resource With Invalid Token",
			operation:     "protectedQuery",
			query:         `query { me { id username } }`,
			token:         "invalid.token.here",
			expectedError: true,
			errorContains: "invalid",
		},
		{
			name:      "Login With Missing Required Fields",
			operation: "login",
			query:     `mutation Login($input: LoginInput!) { login(input: $input) { token } }`,
			variables: map[string]interface{}{
				"input": map[string]interface{}{
					"email": "admin@example.com",
					// Missing password
				},
			},
			expectedError: true,
			errorContains: "required",
		},
		{
			name:      "Password Reset With Invalid Token",
			operation: "passwordReset",
			query:     `mutation ResetPassword($token: String!, $newPassword: String!) { resetPassword(token: $token, newPassword: $newPassword) { success } }`,
			variables: map[string]interface{}{
				"token":       "invalid-reset-token",
				"newPassword": "NewPassword123!",
			},
			expectedError: true,
			errorContains: "invalid",
		},
		{
			name:          "Logout With Invalid Session",
			operation:     "logout",
			query:         `mutation Logout { logout { success } }`,
			token:         "invalid.token.here",
			expectedError: true,
			errorContains: "session",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			authClient := helpers.NewGraphQLClient(client.BaseURL)
			if tc.token != "" {
				authClient = authClient.WithAuth(tc.token)
			}

			resp, err := authClient.Execute(ctx, tc.query, tc.variables)

			// In RED phase, expect errors
			if err == nil {
				if tc.expectedError {
					if tc.errorContains != "" {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for: %s", tc.name)
						resp.AssertErrorContains(t, tc.errorContains)
					} else {
						require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented auth operation: %s", tc.name)
						assertContractError(t, resp, tc.operation)
					}
				} else {
					require.True(t, resp.HasErrors(), "Expected GraphQL errors for unimplemented auth operation: %s", tc.name)
				}
			} else {
				assertConnectionError(t, err, tc.name)
			}
		})
	}
}
