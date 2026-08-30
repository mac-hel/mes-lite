package csvimport

import (
	"errors"

	"github.com/go-fuego/fuego"

	"github.com/mac-hel/mes-lite/internal/platform/jobs"
)

// Handler holds HTTP handlers for CSV import operations.
type Handler struct {
	service      *Service
	asyncService *AsyncService
}

// NewHandler creates a CSV import handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// NewHandlerWithAsync creates a CSV import handler with synchronous and async behavior.
func NewHandlerWithAsync(service *Service, asyncService *AsyncService) *Handler {
	return &Handler{service: service, asyncService: asyncService}
}

// AsyncEnabled reports whether this handler can enqueue background import jobs.
func (h *Handler) AsyncEnabled() bool {
	return h.asyncService != nil
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

// ImportProductionEntriesAsync handles POST /imports/production-entries/jobs.
func (h *Handler) ImportProductionEntriesAsync(c fuego.ContextNoBody) (jobs.JobResponse, error) {
	job, err := h.asyncService.EnqueueProductionEntries(c.Context(), c.Request().Body)
	if err != nil {
		if errors.Is(err, jobs.ErrQueueFull) {
			return jobs.JobResponse{}, fuego.ConflictError{Err: err, Detail: err.Error()}
		}
		return jobs.JobResponse{}, err
	}

	return jobs.NewJobResponse(job), nil
}
