package csvimport

import (
	"errors"

	"github.com/go-fuego/fuego"
)

// Handler holds HTTP handlers for CSV import operations.
type Handler struct {
	service *Service
}

// NewHandler creates a CSV import handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ImportProductionEntries handles POST /imports/production-entries.
func (h *Handler) ImportProductionEntries(c fuego.ContextNoBody) (ImportSummary, error) {
	summary, err := h.service.ImportProductionEntries(c.Context(), c.Request().Body)
	if err != nil {
		if errors.Is(err, ErrInvalidHeader) || errors.Is(err, ErrInvalidRecord) {
			return ImportSummary{}, fuego.BadRequestError{Err: err, Detail: err.Error()}
		}
		return ImportSummary{}, err
	}

	return summary, nil
}
