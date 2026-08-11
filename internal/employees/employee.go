package employees

import (
	"errors"
	"fmt"
	"strings"
)

// ErrNotFound is returned when an employee cannot be found by ID.
var ErrNotFound = errors.New("employee not found")

// ErrAlreadyExists is returned when trying to create an employee with a duplicate ID.
var ErrAlreadyExists = errors.New("employee already exists")

// ErrInvalidEmployee is returned when employee data breaks domain rules.
var ErrInvalidEmployee = errors.New("invalid employee")

// ErrVersionConflict is returned when an update uses a stale employee version.
var ErrVersionConflict = errors.New("employee version conflict")

// Employee represents a person working in the company.
type Employee struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	IsActive  bool
	Version   int
}

// NewEmployee creates a valid Employee with IsActive set to true.
func NewEmployee(id, firstName, lastName, email string) (Employee, error) {
	emp := Employee{
		ID:        strings.TrimSpace(id),
		FirstName: strings.TrimSpace(firstName),
		LastName:  strings.TrimSpace(lastName),
		Email:     strings.TrimSpace(email),
		IsActive:  true,
		Version:   1,
	}
	if err := emp.Validate(); err != nil {
		return Employee{}, err
	}

	return emp, nil
}

// UpdateDetails replaces mutable employee fields and preserves employee invariants.
func (e *Employee) UpdateDetails(firstName, lastName, email string) error {
	updated := *e
	updated.FirstName = strings.TrimSpace(firstName)
	updated.LastName = strings.TrimSpace(lastName)
	updated.Email = strings.TrimSpace(email)
	if err := updated.Validate(); err != nil {
		return err
	}

	*e = updated
	return nil
}

// Validate checks the employee invariants that must hold in every entry point.
func (e Employee) Validate() error {
	if strings.TrimSpace(e.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidEmployee)
	}
	if strings.TrimSpace(e.FirstName) == "" {
		return fmt.Errorf("first name is required: %w", ErrInvalidEmployee)
	}
	if strings.TrimSpace(e.LastName) == "" {
		return fmt.Errorf("last name is required: %w", ErrInvalidEmployee)
	}
	if strings.TrimSpace(e.Email) == "" {
		return fmt.Errorf("email is required: %w", ErrInvalidEmployee)
	}
	if e.Version <= 0 {
		return fmt.Errorf("version must be greater than zero: %w", ErrInvalidEmployee)
	}

	return nil
}
