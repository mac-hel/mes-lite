package production

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-fuego/fuego"

	"github.com/mac-hel/mes-lite/internal/auth"
)

// Handler holds the HTTP handler methods for the production resource.
type Handler struct {
	entries     EntryRegistrar
	corrections CorrectionRegistrar
}

// EntryRegistrar defines production-entry registration and review behavior needed by HTTP.
type EntryRegistrar interface {
	Register(ctx context.Context, cmd RegisterCommand) (Entry, error)
	List(ctx context.Context, opts ListOptions) ([]Entry, error)
}

// CorrectionRegistrar defines production-entry correction behavior needed by HTTP.
type CorrectionRegistrar interface {
	CorrectEntry(ctx context.Context, cmd CorrectEntryCommand) (Correction, error)
	ListCorrections(ctx context.Context, entryID string) ([]Correction, error)
}

// NewHandler creates a new Handler with production entry and correction services.
func NewHandler(entries EntryRegistrar, corrections CorrectionRegistrar) *Handler {
	return &Handler{entries: entries, corrections: corrections}
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

	entries, err := h.entries.List(c.Context(), opts)
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

// CorrectProductionEntryRequest is the expected body for appending a production-entry correction.
type CorrectProductionEntryRequest struct {
	Reason      string    `json:"reason"      validate:"required"`
	EmployeeID  string    `json:"employeeId"  validate:"required"`
	ProductSKU  string    `json:"productSku"  validate:"required"`
	Quantity    int       `json:"quantity"    validate:"required,min=1"`
	Workstation string    `json:"workstation" validate:"required"`
	Timestamp   time.Time `json:"timestamp"   validate:"required"`
	Comment     string    `json:"comment"`
}

// ListCorrectionsResponse wraps correction history for JSON serialization.
type ListCorrectionsResponse struct {
	Corrections []Correction `json:"corrections"`
}

// Register handles POST /production-entries and stores a completed production entry.
func (h *Handler) Register(c fuego.ContextWithBody[RegisterProductionRequest]) (Entry, error) {
	body, err := c.Body()
	if err != nil {
		return Entry{}, err
	}

	entry, err := h.entries.Register(c.Context(), RegisterCommand(body))
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

// Correct appends an audit correction for a production entry without changing the original row.
func (h *Handler) Correct(c fuego.ContextWithBody[CorrectProductionEntryRequest]) (Correction, error) {
	principal, ok := auth.PrincipalFromContext(c.Context())
	if !ok {
		return Correction{}, fuego.UnauthorizedError{Detail: "missing authenticated principal"}
	}
	body, err := c.Body()
	if err != nil {
		return Correction{}, err
	}

	correction, err := h.corrections.CorrectEntry(c.Context(), CorrectEntryCommand{
		EntryID:     c.PathParam("id"),
		ActorUserID: principal.UserID,
		Reason:      body.Reason,
		EmployeeID:  body.EmployeeID,
		ProductSKU:  body.ProductSKU,
		Quantity:    body.Quantity,
		Workstation: body.Workstation,
		Timestamp:   body.Timestamp,
		Comment:     body.Comment,
	})
	if err != nil {
		return Correction{}, correctionError(c.PathParam("id"), err)
	}

	return correction, nil
}

// ListCorrections returns the append-only correction history for one production entry.
func (h *Handler) ListCorrections(c fuego.ContextNoBody) (ListCorrectionsResponse, error) {
	corrections, err := h.corrections.ListCorrections(c.Context(), c.PathParam("id"))
	if err != nil {
		return ListCorrectionsResponse{}, correctionError(c.PathParam("id"), err)
	}
	if corrections == nil {
		corrections = []Correction{}
	}

	return ListCorrectionsResponse{Corrections: corrections}, nil
}

func correctionError(entryID string, err error) error {
	if errors.Is(err, ErrNotFound) {
		return fuego.NotFoundError{Err: err, Detail: fmt.Sprintf("production entry %q not found", entryID)}
	}
	if errors.Is(err, ErrInvalidCorrection) || errors.Is(err, ErrInvalidEntry) || errors.Is(err, ErrEmployeeInactive) || errors.Is(err, ErrProductInactive) {
		return fuego.BadRequestError{Err: err, Detail: err.Error()}
	}
	if errors.Is(err, ErrEmployeeNotFound) || errors.Is(err, ErrProductNotFound) {
		return fuego.NotFoundError{Err: err, Detail: err.Error()}
	}
	if errors.Is(err, ErrAlreadyExists) {
		return fuego.ConflictError{Err: err, Detail: err.Error()}
	}
	return err
}

func invalidEntryError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{
		Err:    err,
		Detail: err.Error(),
	}
}
