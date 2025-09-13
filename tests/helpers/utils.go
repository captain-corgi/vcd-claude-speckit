package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestEmployee represents test employee data
type TestEmployee struct {
	ID         string                 `json:"id"`
	FirstName  string                 `json:"firstName"`
	LastName   string                 `json:"lastName"`
	Email      string                 `json:"email"`
	Department string                 `json:"department"`
	Position   string                 `json:"position"`
	Salary     float64                `json:"salary"`
	Status     string                 `json:"status"`
	Address    map[string]interface{} `json:"address,omitempty"`
}

// TestUser represents test user data
type TestUser struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Role     string   `json:"role"`
	IsActive bool     `json:"isActive"`
}

// GenerateTestEmployee creates a test employee with random data
func GenerateTestEmployee() *TestEmployee {
	return &TestEmployee{
		FirstName:  fmt.Sprintf("Test%s", RandomString(8)),
		LastName:   fmt.Sprintf("User%s", RandomString(8)),
		Email:      fmt.Sprintf("test.%s@example.com", RandomString(10)),
		Department: RandomDepartment(),
		Position:   RandomPosition(),
		Salary:     RandomSalary(),
		Status:     "ACTIVE",
		Address: map[string]interface{}{
			"street":     fmt.Sprintf("%d Test St", rand.Intn(9999)+1),
			"city":      RandomCity(),
			"state":     RandomState(),
			"postalCode": fmt.Sprintf("%05d", rand.Intn(99999)),
			"country":   "US",
		},
	}
}

// GenerateTestUser creates a test user with random data
func GenerateTestUser() *TestUser {
	return &TestUser{
		Username: fmt.Sprintf("testuser%s", RandomString(8)),
		Email:    fmt.Sprintf("user.%s@example.com", RandomString(10)),
		Role:     RandomRole(),
		IsActive: true,
	}
}

// RandomString generates a random string of specified length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// RandomDepartment returns a random department name
func RandomDepartment() string {
	departments := []string{
		"Engineering",
		"Marketing",
		"Sales",
		"HR",
		"Finance",
		"Operations",
		"Customer Support",
		"Product",
	}
	return departments[rand.Intn(len(departments))]
}

// RandomPosition returns a random job position
func RandomPosition() string {
	positions := []string{
		"Software Engineer",
		"Senior Software Engineer",
		"Lead Developer",
		"Engineering Manager",
		"Product Manager",
		"Marketing Manager",
		"Sales Representative",
		"HR Specialist",
		"Financial Analyst",
		"Operations Manager",
	}
	return positions[rand.Intn(len(positions))]
}

// RandomSalary returns a random salary between 40000 and 200000
func RandomSalary() float64 {
	return float64(rand.Intn(160000)+40000) + rand.Float64()
}

// RandomRole returns a random user role
func RandomRole() string {
	roles := []string{"ADMIN", "MANAGER", "VIEWER"}
	return roles[rand.Intn(len(roles))]
}

// RandomCity returns a random city name
func RandomCity() string {
	cities := []string{
		"New York",
		"San Francisco",
		"Los Angeles",
		"Chicago",
		"Houston",
		"Phoenix",
		"Philadelphia",
		"San Antonio",
		"San Diego",
		"Dallas",
	}
	return cities[rand.Intn(len(cities))]
}

// RandomState returns a random state code
func RandomState() string {
	states := []string{"CA", "NY", "TX", "FL", "IL", "PA", "OH", "GA", "NC", "MI"}
	return states[rand.Intn(len(states))]
}

// JSONString converts an interface to JSON string
func JSONString(t *testing.T, v interface{}) string {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}
	return string(data)
}

// FromJSONString parses JSON string to interface
func FromJSONString(t *testing.T, s string, v interface{}) {
	t.Helper()
	err := json.Unmarshal([]byte(s), v)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
}

// AssertEqual asserts that two values are equal
func AssertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", msg, expected, actual)
	}
}

// AssertNotEqual asserts that two values are not equal
func AssertNotEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if expected == actual {
		t.Errorf("%s: expected values to be different, but both are %v", msg, expected)
	}
}

// AssertNotNil asserts that a value is not nil
func AssertNotNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value == nil {
		t.Errorf("%s: expected non-nil value, got nil", msg)
	}
}

// AssertNil asserts that a value is nil
func AssertNil(t *testing.T, value interface{}, msg string) {
	t.Helper()
	if value != nil {
		t.Errorf("%s: expected nil value, got %v", msg, value)
	}
}

// AssertContains asserts that a string contains a substring
func AssertContains(t *testing.T, s, substr string, msg string) {
	t.Helper()
	if !strings.Contains(s, substr) {
		t.Errorf("%s: expected '%s' to contain '%s'", msg, s, substr)
	}
}

// AssertNotContains asserts that a string does not contain a substring
func AssertNotContains(t *testing.T, s, substr string, msg string) {
	t.Helper()
	if strings.Contains(s, substr) {
		t.Errorf("%s: expected '%s' to not contain '%s'", msg, s, substr)
	}
}

// AssertHTTPStatus asserts that an HTTP response has the expected status code
func AssertHTTPStatus(t *testing.T, resp *http.Response, expected int, msg string) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Errorf("%s: expected status %d, got %d", msg, expected, resp.StatusCode)
	}
}

// AssertJSONHeader asserts that a response has JSON content type
func AssertJSONHeader(t *testing.T, resp *http.Response) {
	t.Helper()
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected JSON content type, got: %s", contentType)
	}
}

// ReadResponseBody reads and returns the response body as string
func ReadResponseBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	return string(bodyBytes)
}

// CreateTestHTTPRequest creates an HTTP request for testing
func CreateTestHTTPRequest(t *testing.T, method, url string, body interface{}, headers map[string]string) *http.Request {
	t.Helper()

	var reqBody *strings.Reader
	if body != nil {
		jsonBody := JSONString(t, body)
		reqBody = strings.NewReader(jsonBody)
	} else {
		reqBody = strings.NewReader("")
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		t.Fatalf("Failed to create HTTP request: %v", err)
	}

	// Set headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	return req
}

// SendTestRequest sends an HTTP request and returns the response
func SendTestRequest(t *testing.T, handler http.Handler, req *http.Request) *http.Response {
	t.Helper()

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr.Result()
}

// MustParseTime parses a time string or fails the test
func MustParseTime(t *testing.T, layout, value string) time.Time {
	t.Helper()
	parsed, err := time.Parse(layout, value)
	if err != nil {
		t.Fatalf("Failed to parse time '%s' with layout '%s': %v", value, layout, err)
	}
	return parsed
}

// SeedRandom seeds the random number generator for consistent tests
func SeedRandom(t *testing.T) {
	t.Helper()
	rand.Seed(time.Now().UnixNano())
}

// Initialize the random seed when package is loaded
func init() {
	rand.Seed(time.Now().UnixNano())
}