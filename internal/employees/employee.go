package employees

import "errors"

// ErrNotFound is returned when an employee cannot be found by ID.
var ErrNotFound = errors.New("employee not found")
var ErrAlreadyExists = errors.New("employee already exists")

// Employee represents a person working in the company.
type Employee struct {
	ID        string
	FirstName string
	LastName  string
	Email     string
	IsActive  bool
}

// NewEmployee creates a new Employee with the given details and sets IsActive to true.
func NewEmployee(id, firstName, lastName, email string) Employee {
	return Employee{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		IsActive:  true,
	}
}
