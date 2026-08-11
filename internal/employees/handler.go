package employees

import (
	"errors"
	"fmt"
	"strconv"

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

	emp, err := NewEmployee(body.ID, body.FirstName, body.LastName, body.Email)
	if err != nil {
		return Employee{}, invalidEmployeeError(err)
	}

	if err := h.store.Save(c.Context(), emp); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return Employee{}, fuego.ConflictError{
				Err:    err,
				Detail: fmt.Sprintf("employee %q already exists", emp.ID),
			}
		}
		if errors.Is(err, ErrInvalidEmployee) {
			return Employee{}, invalidEmployeeError(err)
		}
		return Employee{}, err
	}

	return emp, nil
}

// ListEmployeesResponse wraps a slice of employees for JSON serialization.
type ListEmployeesResponse struct {
	Employees  []Employee `json:"employees"`
	Pagination Page       `json:"pagination"`
}

// List handles GET /employees and returns all employees.
func (h *Handler) List(c fuego.ContextNoBody) (ListEmployeesResponse, error) {
	opts, err := listOptionsFromQuery(c)
	if err != nil {
		return ListEmployeesResponse{}, invalidListOptionsError(err)
	}

	emps, err := h.store.List(c.Context(), opts)
	if err != nil {
		if errors.Is(err, ErrInvalidListOptions) {
			return ListEmployeesResponse{}, invalidListOptionsError(err)
		}
		return ListEmployeesResponse{}, err
	}

	if emps == nil {
		emps = []Employee{}
	}

	return ListEmployeesResponse{
		Employees: emps,
		Pagination: Page{
			Limit:  opts.Limit,
			Offset: opts.Offset,
			Count:  len(emps),
		},
	}, nil
}

func listOptionsFromQuery(c fuego.ContextNoBody) (ListOptions, error) {
	limit, err := parseIntQuery(c.QueryParam("limit"), "limit")
	if err != nil {
		return ListOptions{}, err
	}
	offset, err := parseIntQuery(c.QueryParam("offset"), "offset")
	if err != nil {
		return ListOptions{}, err
	}
	active, err := parseBoolQuery(c.QueryParam("active"), "active")
	if err != nil {
		return ListOptions{}, err
	}

	return ListOptions{
		Limit:  limit,
		Offset: offset,
		Sort:   c.QueryParam("sort"),
		Query:  c.QueryParam("q"),
		Active: active,
	}.normalize()
}

func parseIntQuery(raw, name string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, ErrInvalidListOptions)
	}
	return value, nil
}

func parseBoolQuery(raw, name string) (*bool, error) {
	if raw == "" {
		return nil, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return nil, fmt.Errorf("%s must be true or false: %w", name, ErrInvalidListOptions)
	}
	return &value, nil
}

func invalidListOptionsError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}

// UpdateEmployeeRequest is the expected JSON body for updating an employee.
// Validation tags enforce required fields and email format.
type UpdateEmployeeRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName"  validate:"required"`
	Email     string `json:"email"     validate:"required,email"`
	Version   int    `json:"version"   validate:"min=1"`
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

	if err := emp.UpdateDetails(body.FirstName, body.LastName, body.Email); err != nil {
		return Employee{}, invalidEmployeeError(err)
	}
	emp.Version = body.Version

	updated, err := h.store.Update(c.Context(), emp)
	if err != nil {
		if errors.Is(err, ErrInvalidEmployee) {
			return Employee{}, invalidEmployeeError(err)
		}
		if errors.Is(err, ErrVersionConflict) {
			return Employee{}, versionConflictError(err)
		}
		return Employee{}, err
	}

	return updated, nil
}

func invalidEmployeeError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}

func versionConflictError(err error) fuego.ConflictError {
	return fuego.ConflictError{
		Err:    err,
		Detail: err.Error(),
	}
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

	updated, err := h.store.Update(c.Context(), emp)
	if err != nil {
		if errors.Is(err, ErrVersionConflict) {
			return Employee{}, versionConflictError(err)
		}
		return Employee{}, err
	}

	return updated, nil
}
