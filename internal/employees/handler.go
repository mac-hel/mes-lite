package employees

import (
	"errors"
	"fmt"

	"github.com/go-fuego/fuego"
)

// Handler holds the HTTP handler methods for the employees resource.
type Handler struct {
	store Store
}

// NewHandler creates a new Handler with the given Store.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// CreateEmployeeRequest is the expected JSON body for creating an employee.
// Validation tags enforce required fields and email format.
type CreateEmployeeRequest struct {
	ID        string `json:"id"        validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"  validate:"required"`
	Email     string `json:"email"     validate:"required,email"`
}

// Create handles POST /employees and stores a new employee.
func (h *Handler) Create(c fuego.ContextWithBody[CreateEmployeeRequest]) (Employee, error) {
	body, err := c.Body()
	if err != nil {
		return Employee{}, err
	}

	emp := NewEmployee(body.ID, body.FirstName, body.LastName, body.Email)

	if err := h.store.Save(c.Context(), emp); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return Employee{}, fuego.ConflictError{
				Err:    err,
				Detail: fmt.Sprintf("employee %q already exists", emp.ID),
			}
		}
	}

	return emp, nil
}

// ListEmployeesResponse wraps a slice of employees for JSON serialization.
type ListEmployeesResponse struct {
	Employees []Employee `json:"employees"`
}

// List handles GET /employees and returns all employees.
func (h *Handler) List(c fuego.ContextNoBody) (ListEmployeesResponse, error) {
	emps, err := h.store.List(c.Context())
	if err != nil {
		return ListEmployeesResponse{}, err
	}

	if emps == nil {
		emps = []Employee{}
	}

	return ListEmployeesResponse{Employees: emps}, nil
}

// UpdateEmployeeRequest is the expected JSON body for updating an employee.
// Validation tags enforce required fields and email format.
type UpdateEmployeeRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"  validate:"required"`
	Email     string `json:"email"     validate:"required,email"`
}

// Update handles PUT /employees/{id} and replaces the employee's mutable fields.
func (h *Handler) Update(c fuego.ContextWithBody[UpdateEmployeeRequest]) (Employee, error) {
	id := c.PathParam("id")

	emp, err := h.store.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Employee{}, fuego.NotFoundError{
				Err:    err,
				Detail: fmt.Sprintf("employee %q not found", id),
			}
		}
		return Employee{}, err
	}

	body, err := c.Body()
	if err != nil {
		return Employee{}, err
	}

	emp.FirstName = body.FirstName
	emp.LastName = body.LastName
	emp.Email = body.Email

	if err := h.store.Update(c.Context(), emp); err != nil {
		return Employee{}, err
	}

	return emp, nil
}

// Deactivate handles PUT /employees/{id}/deactivate and sets IsActive to false.
func (h *Handler) Deactivate(c fuego.ContextNoBody) (Employee, error) {
	id := c.PathParam("id")

	emp, err := h.store.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return Employee{}, fuego.NotFoundError{
				Err:    err,
				Detail: fmt.Sprintf("employee %q not found", id),
			}
		}
		return Employee{}, err
	}

	emp.IsActive = false

	if err := h.store.Update(c.Context(), emp); err != nil {
		return Employee{}, err
	}

	return emp, nil
}
