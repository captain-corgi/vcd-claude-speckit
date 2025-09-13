package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Database migration tool for employee management system",
	}

	var dbHost, dbUser, dbPassword, dbName, dbPort string

	rootCmd.PersistentFlags().StringVar(&dbHost, "host", "localhost", "Database host")
	rootCmd.PersistentFlags().StringVar(&dbPort, "port", "5432", "Database port")
	rootCmd.PersistentFlags().StringVar(&dbUser, "user", "employeemgmt", "Database user")
	rootCmd.PersistentFlags().StringVar(&dbPassword, "password", "employeemgmt_password", "Database password")
	rootCmd.PersistentFlags().StringVar(&dbName, "database", "employee_management", "Database name")

	// Add commands
	rootCmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Run all pending migrations",
		Run: func(cmd *cobra.Command, args []string) {
			m, err := createMigrateInstance(dbHost, dbPort, dbUser, dbPassword, dbName)
			if err != nil {
				log.Fatalf("Failed to create migrate instance: %v", err)
			}
			defer m.Close()

			if err := m.Up(); err != nil {
				if err == migrate.ErrNoChange {
					fmt.Println("Database is already up to date")
					return
				}
				log.Fatalf("Failed to run up migrations: %v", err)
			}
			fmt.Println("Migrations completed successfully")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback the last migration",
		Run: func(cmd *cobra.Command, args []string) {
			m, err := createMigrateInstance(dbHost, dbPort, dbUser, dbPassword, dbName)
			if err != nil {
				log.Fatalf("Failed to create migrate instance: %v", err)
			}
			defer m.Close()

			if err := m.Steps(-1); err != nil {
				if err == migrate.ErrNoChange {
					fmt.Println("No migrations to rollback")
					return
				}
				log.Fatalf("Failed to rollback migration: %v", err)
			}
			fmt.Println("Migration rolled back successfully")
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new migration file",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			if err := createMigrationFile(name); err != nil {
				log.Fatalf("Failed to create migration file: %v", err)
			}
			fmt.Printf("Migration file created: %s\n", name)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show current migration version",
		Run: func(cmd *cobra.Command, args []string) {
			m, err := createMigrateInstance(dbHost, dbPort, dbUser, dbPassword, dbName)
			if err != nil {
				log.Fatalf("Failed to create migrate instance: %v", err)
			}
			defer m.Close()

			version, dirty, err := m.Version()
			if err != nil {
				if err == migrate.ErrNilVersion {
					fmt.Println("No migrations have been applied")
					return
				}
				log.Fatalf("Failed to get version: %v", err)
			}

			fmt.Printf("Current version: %d (dirty: %v)\n", version, dirty)
		},
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "force",
		Short: "Force set migration version",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			version, err := strconv.ParseUint(args[0], 10, 32)
			if err != nil {
				log.Fatalf("Invalid version number: %v", err)
			}

			m, err := createMigrateInstance(dbHost, dbPort, dbUser, dbPassword, dbName)
			if err != nil {
				log.Fatalf("Failed to create migrate instance: %v", err)
			}
			defer m.Close()

			if err := m.Force(int(version)); err != nil {
				log.Fatalf("Failed to force version: %v", err)
			}
			fmt.Printf("Forced version to: %d\n", version)
		},
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}

func createMigrateInstance(host, port, user, password, dbName string) (*migrate.Migrate, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres driver: %w", err)
	}

	// Use file system source for migrations
	m, err := migrate.NewWithDatabaseInstance("file://migrations", dbName, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

func createMigrationFile(name string) error {
	timestamp := time.Now().Format("20060102150405")
	version := timestamp

	// Create up migration file
	upFileName := fmt.Sprintf("%s_%s.up.sql", version, name)
	upPath := filepath.Join("migrations", upFileName)
	upContent := fmt.Sprintf(`-- Migration: %s
-- Created: %s
-- Description: %s

-- Write your UP migration SQL here

`, name, time.Now().Format("2006-01-02 15:04:05"), name)

	if err := os.WriteFile(upPath, []byte(upContent), 0644); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	// Create down migration file
	downFileName := fmt.Sprintf("%s_%s.down.sql", version, name)
	downPath := filepath.Join("migrations", downFileName)
	downContent := fmt.Sprintf(`-- Migration: %s (DOWN)
-- Created: %s
-- Description: %s

-- Write your DOWN migration SQL here

`, name, time.Now().Format("2006-01-02 15:04:05"), name)

	if err := os.WriteFile(downPath, []byte(downContent), 0644); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	return nil
}

