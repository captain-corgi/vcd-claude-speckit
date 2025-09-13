package helpers

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TestTokenConfig holds configuration for test JWT tokens
type TestTokenConfig struct {
	Secret     string
	Expiration time.Duration
}

// TestClaims represents JWT claims for testing
type TestClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Email    string   `json:"email"`
	jwt.RegisteredClaims
}

// DefaultTestTokenConfig returns default token configuration for tests
func DefaultTestTokenConfig() *TestTokenConfig {
	return &TestTokenConfig{
		Secret:     "test-secret-key-for-jwt-tokens",
		Expiration: time.Hour,
	}
}

// CreateTestToken creates a JWT token for testing
func CreateTestToken(t *testing.T, config *TestTokenConfig, claims *TestClaims) string {
	t.Helper()

	if config == nil {
		config = DefaultTestTokenConfig()
	}

	if claims == nil {
		claims = &TestClaims{
			UserID:   "test-user-id",
			Username: "testuser",
			Roles:    []string{"ADMIN"},
			Email:    "test@example.com",
		}
	}

	// Set expiration
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(config.Expiration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.Issuer = "employee-management-test"

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	return tokenString
}

// CreateAdminToken creates a token with admin role for testing
func CreateAdminToken(t *testing.T) string {
	t.Helper()
	return CreateTestToken(t, nil, &TestClaims{
		UserID:   "admin-user-id",
		Username: "admin",
		Roles:    []string{"ADMIN"},
		Email:    "admin@test.com",
	})
}

// CreateManagerToken creates a token with manager role for testing
func CreateManagerToken(t *testing.T) string {
	t.Helper()
	return CreateTestToken(t, nil, &TestClaims{
		UserID:   "manager-user-id",
		Username: "manager",
		Roles:    []string{"MANAGER"},
		Email:    "manager@test.com",
	})
}

// CreateViewerToken creates a token with viewer role for testing
func CreateViewerToken(t *testing.T) string {
	t.Helper()
	return CreateTestToken(t, nil, &TestClaims{
		UserID:   "viewer-user-id",
		Username: "viewer",
		Roles:    []string{"VIEWER"},
		Email:    "viewer@test.com",
	})
}

// CreateExpiredToken creates an expired token for testing
func CreateExpiredToken(t *testing.T) string {
	t.Helper()
	config := DefaultTestTokenConfig()
	config.Expiration = -time.Hour // Expired

	claims := &TestClaims{
		UserID:   "expired-user-id",
		Username: "expired",
		Roles:    []string{"VIEWER"},
		Email:    "expired@test.com",
	}

	return CreateTestToken(t, config, claims)
}

// ParseTestToken parses a JWT token for testing (returns claims or error)
func ParseTestToken(t *testing.T, tokenString string, secret string) *TestClaims {
	t.Helper()

	token, err := jwt.ParseWithClaims(tokenString, &TestClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse test token: %v", err)
		return nil
	}

	if claims, ok := token.Claims.(*TestClaims); ok && token.Valid {
		return claims
	}

	t.Fatalf("Invalid token claims")
	return nil
}

// GetAuthHeader returns the Authorization header value for HTTP requests
func GetAuthHeader(token string) string {
	return "Bearer " + token
}

// AuthHeaders returns HTTP headers with authentication
func AuthHeaders(token string) map[string]string {
	return map[string]string{
		"Authorization": GetAuthHeader(token),
		"Content-Type":  "application/json",
	}
}

// DefaultAuthHeaders returns headers with default admin token
func DefaultAuthHeaders(t *testing.T) map[string]string {
	return AuthHeaders(CreateAdminToken(t))
}