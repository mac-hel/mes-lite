package production

import (
	"context"
	"errors"
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
}

// NewHandler creates a new Handler with the given Registrar.
func NewHandler(registrar Registrar) *Handler {
	return &Handler{registrar: registrar}
}

// RegisterProductionRequest is the expected JSON body for registering completed production.
type RegisterProductionRequest struct {
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
