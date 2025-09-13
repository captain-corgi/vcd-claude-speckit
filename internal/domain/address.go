package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Address represents an employee's address
type Address struct {
	Street     string
	City       string
	State      string
	PostalCode string
	Country    string
}

// NewAddress creates a new Address with validation
func NewAddress(street, city, state, postalCode, country string) (*Address, error) {
	addr := &Address{
		Street:     strings.TrimSpace(street),
		City:       strings.TrimSpace(city),
		State:      strings.TrimSpace(state),
		PostalCode: strings.TrimSpace(postalCode),
		Country:    strings.TrimSpace(country),
	}

	if err := addr.Validate(); err != nil {
		return nil, err
	}

	return addr, nil
}

// Validate validates the address fields
func (a *Address) Validate() error {
	if a == nil {
		return nil // Nil address is valid (optional field)
	}

	// If any field is provided, validate required fields
	if a.Street != "" || a.City != "" || a.State != "" || a.PostalCode != "" || a.Country != "" {
		if a.Street == "" {
			return errors.New("street is required when other address fields are provided")
		}
		if a.City == "" {
			return errors.New("city is required when other address fields are provided")
		}
		if a.State == "" {
			return errors.New("state is required when other address fields are provided")
		}
		if a.PostalCode == "" {
			return errors.New("postal code is required when other address fields are provided")
		}
		if a.Country == "" {
			return errors.New("country is required when other address fields are provided")
		}
	}

	// Validate street format
	if a.Street != "" && len(a.Street) < 5 {
		return errors.New("street must be at least 5 characters long")
	}
	if len(a.Street) > 200 {
		return errors.New("street cannot exceed 200 characters")
	}

	// Validate city format
	if a.City != "" && len(a.City) < 2 {
		return errors.New("city must be at least 2 characters long")
	}
	if len(a.City) > 100 {
		return errors.New("city cannot exceed 100 characters")
	}

	// Validate state format (US states, Canadian provinces, or general format)
	if a.State != "" && len(a.State) < 2 {
		return errors.New("state must be at least 2 characters long")
	}
	if len(a.State) > 50 {
		return errors.New("state cannot exceed 50 characters")
	}

	// Validate postal code format
	if a.PostalCode != "" {
		if err := validatePostalCode(a.PostalCode); err != nil {
			return fmt.Errorf("invalid postal code: %w", err)
		}
	}

	// Validate country format
	if a.Country != "" && len(a.Country) < 2 {
		return errors.New("country must be at least 2 characters long")
	}
	if len(a.Country) > 100 {
		return errors.New("country cannot exceed 100 characters")
	}

	return nil
}

// validatePostalCode validates various postal code formats
func validatePostalCode(postalCode string) error {
	// US ZIP code format (12345 or 12345-6789)
	usZip := regexp.MustCompile(`^\d{5}(-\d{4})?$`)
	// Canadian postal code format (A1A 1A1)
	caPostal := regexp.MustCompile(`^[A-Z]\d[A-Z][ ]?\d[A-Z]\d$`)
	// General postal code format (alphanumeric, 3-10 characters)
	general := regexp.MustCompile(`^[A-Z0-9]{3,10}$`)

	normalized := strings.ToUpper(strings.ReplaceAll(postalCode, " ", ""))

	if usZip.MatchString(postalCode) || caPostal.MatchString(postalCode) || general.MatchString(normalized) {
		return nil
	}

	return errors.New("postal code format is invalid")
}

// IsEmpty checks if the address is empty (all fields are empty)
func (a *Address) IsEmpty() bool {
	if a == nil {
		return true
	}
	return a.Street == "" && a.City == "" && a.State == "" && a.PostalCode == "" && a.Country == ""
}

// Format returns a formatted address string
func (a *Address) Format() string {
	if a == nil || a.IsEmpty() {
		return ""
	}

	var parts []string
	if a.Street != "" {
		parts = append(parts, a.Street)
	}
	if a.City != "" || a.State != "" || a.PostalCode != "" {
		cityStateZip := fmt.Sprintf("%s, %s %s", a.City, a.State, a.PostalCode)
		parts = append(parts, strings.TrimSpace(cityStateZip))
	}
	if a.Country != "" {
		parts = append(parts, a.Country)
	}

	return strings.Join(parts, "\n")
}

// Equal checks if two addresses are equal
func (a *Address) Equal(other *Address) bool {
	if a == nil || other == nil {
		return a == other
	}
	return a.Street == other.Street &&
		a.City == other.City &&
		a.State == other.State &&
		a.PostalCode == other.PostalCode &&
		a.Country == other.Country
}

// Clone returns a copy of the address
func (a *Address) Clone() *Address {
	if a == nil {
		return nil
	}
	return &Address{
		Street:     a.Street,
		City:       a.City,
		State:      a.State,
		PostalCode: a.PostalCode,
		Country:    a.Country,
	}
}

var (
	ErrInvalidAddress = errors.New("invalid address")
	ErrInvalidPostalCode = errors.New("invalid postal code")
	ErrAddressFieldRequired = errors.New("address field is required")
)