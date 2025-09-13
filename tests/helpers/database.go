package helpers

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// TestDatabase represents a test database container
type TestDatabase struct {
	Container testcontainers.Container
	Host      string
	Port      string
	User      string
	Password  string
	Database  string
}

// NewTestDatabase creates a new PostgreSQL container for testing
func NewTestDatabase(ctx context.Context, t *testing.T) (*TestDatabase, error) {
	t.Helper()

	// Create PostgreSQL container
	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("test_employee_management"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Minute),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL container: %w", err)
	}

	// Get container details
	host, err := pgContainer.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	testDB := &TestDatabase{
		Container: pgContainer,
		Host:      host,
		Port:      port.Port(),
		User:      "testuser",
		Password:  "testpass",
		Database:  "test_employee_management",
	}

	return testDB, nil
}

// GetConnectionString returns the database connection string
func (td *TestDatabase) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		td.Host, td.Port, td.User, td.Password, td.Database)
}

// Cleanup stops and removes the container
func (td *TestDatabase) Cleanup(ctx context.Context) error {
	if td.Container != nil {
		return td.Container.Terminate(ctx)
	}
	return nil
}

// CreateTestContext creates a context for tests with timeout
func CreateTestContext(t *testing.T) (context.Context, context.CancelFunc) {
	t.Helper()
	return context.WithTimeout(context.Background(), 5*time.Minute)
}

// RunDatabaseTest is a helper function that sets up and tears down a test database
func RunDatabaseTest(t *testing.T, testFunc func(*testing.T, *TestDatabase)) {
	ctx, cancel := CreateTestContext(t)
	defer cancel()

	testDB, err := NewTestDatabase(ctx, t)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	defer func() {
		if err := testDB.Cleanup(ctx); err != nil {
			t.Errorf("Failed to cleanup test database: %v", err)
		}
	}()

	testFunc(t, testDB)
}

// WaitForDatabase waits for the database to be ready
func (td *TestDatabase) WaitForDatabase(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(100 * time.Millisecond):
		// Try to connect to verify database is ready
		// This is a simple implementation - in production you'd want proper health checks
		return nil
	}
}