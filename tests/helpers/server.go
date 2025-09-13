package helpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// TestServer represents a test HTTP server
type TestServer struct {
	Server   *httptest.Server
	Router   *gin.Engine
	BaseURL  string
}

// NewTestServer creates a new test server with Gin router
func NewTestServer(t *testing.T) *TestServer {
	t.Helper()

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	server := httptest.NewServer(router)

	return &TestServer{
		Server:  server,
		Router:  router,
		BaseURL: server.URL,
	}
}

// Close shuts down the test server
func (ts *TestServer) Close() {
	if ts.Server != nil {
		ts.Server.Close()
	}
}

// CreateTestRequest creates a test HTTP request
func (ts *TestServer) CreateTestRequest(method, path string, body interface{}) *http.Request {
	// Implementation will depend on your specific request needs
	req, _ := http.NewRequest(method, ts.BaseURL+path, nil)
	return req
}

// RunServerTest is a helper function that sets up and tears down a test server
func RunServerTest(t *testing.T, testFunc func(*testing.T, *TestServer)) {
	t.Helper()

	testServer := NewTestServer(t)
	defer testServer.Close()

	testFunc(t, testServer)
}

// MockHandler creates a simple mock handler for testing
func MockHandler(message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": message,
		})
	}
}

// MockHandlerWithError creates a handler that returns an error
func MockHandlerWithError(statusCode int, message string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(statusCode, gin.H{
			"error": message,
		})
	}
}

// CreateTestContextWithTimeout creates a context with timeout for tests
func CreateTestContextWithTimeout(t *testing.T, timeout time.Duration) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), timeout)
}

// DefaultTestTimeout returns the default timeout for tests
func DefaultTestTimeout() time.Duration {
	return 30 * time.Second
}