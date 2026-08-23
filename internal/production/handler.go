package production

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-fuego/fuego"
)

// Handler holds the HTTP handler methods for the production resource.
type Handler struct {
	registrar Registrar
}

// Registrar defines the production registration behavior needed by HTTP.
type Registrar interface {
	Register(ctx context.Context, cmd RegisterCommand) (Entry, error)
	List(ctx context.Context, opts ListOptions) ([]Entry, error)
}

// NewHandler creates a new Handler with the given Registrar.
func NewHandler(registrar Registrar) *Handler {
	return &Handler{registrar: registrar}
}

// ListProductionEntriesResponse wraps production entries for review responses.
type ListProductionEntriesResponse struct {
	Entries    []Entry `json:"entries"`
	Pagination Page    `json:"pagination"`
}

// List handles GET /production-entries for production-entry review.
func (h *Handler) List(c fuego.ContextNoBody) (ListProductionEntriesResponse, error) {
	opts, err := productionListOptionsFromQuery(c)
	if err != nil {
		return ListProductionEntriesResponse{}, invalidProductionListOptionsError(err)
	}

	entries, err := h.registrar.List(c.Context(), opts)
	if err != nil {
		if errors.Is(err, ErrInvalidListOptions) {
			return ListProductionEntriesResponse{}, invalidProductionListOptionsError(err)
		}
		return ListProductionEntriesResponse{}, err
	}
	if entries == nil {
		entries = []Entry{}
	}

	return ListProductionEntriesResponse{
		Entries: entries,
		Pagination: Page{
			Limit:  opts.Limit,
			Offset: opts.Offset,
			Count:  len(entries),
		},
	}, nil
}

func productionListOptionsFromQuery(c fuego.ContextNoBody) (ListOptions, error) {
	limit, err := parseProductionIntQuery(c.QueryParam("limit"), "limit")
	if err != nil {
		return ListOptions{}, err
	}
	offset, err := parseProductionIntQuery(c.QueryParam("offset"), "offset")
	if err != nil {
		return ListOptions{}, err
	}
	from, err := parseProductionTimeQuery(c.QueryParam("from"), "from")
	if err != nil {
		return ListOptions{}, err
	}
	to, err := parseProductionTimeQuery(c.QueryParam("to"), "to")
	if err != nil {
		return ListOptions{}, err
	}

	return ListOptions{
		Limit:       limit,
		Offset:      offset,
		EmployeeID:  c.QueryParam("employeeId"),
		ProductSKU:  c.QueryParam("productSku"),
		Workstation: c.QueryParam("workstation"),
		From:        from,
		To:          to,
	}.normalize()
}

func parseProductionIntQuery(raw, name string) (int, error) {
	if raw == "" {
		return 0, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", name, ErrInvalidListOptions)
	}
	return value, nil
}

func parseProductionTimeQuery(raw, name string) (time.Time, error) {
	if raw == "" {
		return time.Time{}, nil
	}
	value, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, fmt.Errorf("%s must be RFC3339: %w", name, ErrInvalidListOptions)
	}
	return value, nil
}

func invalidProductionListOptionsError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}

// RegisterProductionRequest is the expected JSON body for registering completed production.
type RegisterProductionRequest struct {
	RequestID   string    `json:"requestId"   validate:"required"`
	EmployeeID  string    `json:"employeeId"  validate:"required"`
	ProductSKU  string    `json:"productSku"  validate:"required"`
	Quantity    int       `json:"quantity"    validate:"required,min=1"`
	Workstation string    `json:"workstation" validate:"required"`
	Timestamp   time.Time `json:"timestamp"   validate:"required"`
	Comment     string    `json:"comment"`
}

// Register handles POST /production-entries and stores a completed production entry.
func (h *Handler) Register(c fuego.ContextWithBody[RegisterProductionRequest]) (Entry, error) {
	body, err := c.Body()
	if err != nil {
		return Entry{}, err
	}

	entry, err := h.registrar.Register(c.Context(), RegisterCommand(body))
	if err != nil {
		if errors.Is(err, ErrRequestConflict) {
			return Entry{}, fuego.ConflictError{
				Err:    err,
				Detail: err.Error(),
			}
		}
		if errors.Is(err, ErrAlreadyExists) {
			return Entry{}, fuego.ConflictError{
				Err:    err,
				Detail: "production entry already exists",
			}
		}
		if errors.Is(err, ErrInvalidEntry) {
			return Entry{}, invalidEntryError(err)
		}
		if errors.Is(err, ErrEmployeeNotFound) || errors.Is(err, ErrProductNotFound) {
			return Entry{}, fuego.NotFoundError{
				Err:    err,
				Detail: err.Error(),
			}
		}
		if errors.Is(err, ErrEmployeeInactive) || errors.Is(err, ErrProductInactive) {
			return Entry{}, fuego.BadRequestError{
				Err:    err,
				Detail: err.Error(),
			}
		}
		return Entry{}, err
	}

	return entry, nil
}

func invalidEntryError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}
