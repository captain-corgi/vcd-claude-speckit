package domain

import (
	"errors"
	"fmt"
)

// EmployeeStatus represents the employment status of an employee
type EmployeeStatus string

const (
	EmployeeStatusActive     EmployeeStatus = "ACTIVE"
	EmployeeStatusTerminated EmployeeStatus = "TERMINATED"
	EmployeeStatusOnLeave    EmployeeStatus = "ON_LEAVE"
)

// AllEmployeeStatuses returns all valid employee statuses
func AllEmployeeStatuses() []EmployeeStatus {
	return []EmployeeStatus{
		EmployeeStatusActive,
		EmployeeStatusTerminated,
		EmployeeStatusOnLeave,
	}
}

// IsValid checks if the employee status is valid
func (s EmployeeStatus) IsValid() bool {
	switch s {
	case EmployeeStatusActive, EmployeeStatusTerminated, EmployeeStatusOnLeave:
		return true
	default:
		return false
	}
}

// String returns the string representation of the employee status
func (s EmployeeStatus) String() string {
	return string(s)
}

// CanChangeTo checks if the current status can be changed to the target status
func (s EmployeeStatus) CanChangeTo(target EmployeeStatus) bool {
	// TERMINATED employees cannot change status
	if s == EmployeeStatusTerminated {
		return false
	}

	// All other transitions are allowed
	return target.IsValid()
}

// ValidateTransition validates the status transition and returns an error if invalid
func (s EmployeeStatus) ValidateTransition(target EmployeeStatus) error {
	if !s.CanChangeTo(target) {
		return fmt.Errorf("invalid status transition from %s to %s", s, target)
	}
	return nil
}

// ParseEmployeeStatus parses a string into an EmployeeStatus
func ParseEmployeeStatus(status string) (EmployeeStatus, error) {
	es := EmployeeStatus(status)
	if !es.IsValid() {
		return "", fmt.Errorf("invalid employee status: %s", status)
	}
	return es, nil
}

// MustParseEmployeeStatus parses a string into an EmployeeStatus or panics
func MustParseEmployeeStatus(status string) EmployeeStatus {
	es, err := ParseEmployeeStatus(status)
	if err != nil {
		panic(err)
	}
	return es
}

var (
	ErrInvalidEmployeeStatus = errors.New("invalid employee status")
	ErrInvalidStatusTransition = errors.New("invalid employee status transition")
)