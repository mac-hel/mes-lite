package machines

import (
	"context"
	"errors"
	"time"

	"github.com/go-fuego/fuego"
)

// Handler exposes fake machine integration endpoints.
type Handler struct {
	receiver EventReceiver
	stats    IntakeStatsProvider
}

// EventReceiver defines machine event intake behavior needed by HTTP.
type EventReceiver interface {
	ReceiveEvent(ctx context.Context, cmd ReceiveEventCommand) (Event, error)
}

// IntakeStatsProvider defines machine event intake metrics needed by HTTP.
type IntakeStatsProvider interface {
	Stats() IntakeStats
}

type machineService interface {
	EventReceiver
	IntakeStatsProvider
}

// NewHandler creates a machine handler.
func NewHandler(service machineService) *Handler {
	return &Handler{receiver: service, stats: service}
}

// CreateMachineEventRequest is the fake machine event JSON body.
type CreateMachineEventRequest struct {
	ExternalEventID string    `json:"externalEventId" validate:"required"`
	Type            EventType `json:"type"            validate:"required"`
	OccurredAt      time.Time `json:"occurredAt"      validate:"required"`
	ProductSKU      string    `json:"productSku"`
	Quantity        int       `json:"quantity"`
	Workstation     string    `json:"workstation"      validate:"required"`
	Message         string    `json:"message"`
}

// EventResponse is the HTTP representation of a machine event.
type EventResponse struct {
	ID              string    `json:"id"`
	MachineID       string    `json:"machineId"`
	ExternalEventID string    `json:"externalEventId"`
	Type            EventType `json:"type"`
	OccurredAt      time.Time `json:"occurredAt"`
	ProductSKU      string    `json:"productSku"`
	Quantity        int       `json:"quantity"`
	Workstation     string    `json:"workstation"`
	Message         string    `json:"message"`
}

// IntakeStatsResponse is the HTTP representation of fake machine intake counters.
type IntakeStatsResponse struct {
	Received         uint64 `json:"received"`
	Accepted         uint64 `json:"accepted"`
	DuplicateRetries uint64 `json:"duplicateRetries"`
	Conflicts        uint64 `json:"conflicts"`
	Invalid          uint64 `json:"invalid"`
}

// CreateEvent handles POST /machines/{machineId}/events.
func (h *Handler) CreateEvent(c fuego.ContextWithBody[CreateMachineEventRequest]) (EventResponse, error) {
	body, err := c.Body()
	if err != nil {
		return EventResponse{}, err
	}

	event, err := h.receiver.ReceiveEvent(c.Context(), ReceiveEventCommand{
		MachineID:       c.PathParam("machineId"),
		ExternalEventID: body.ExternalEventID,
		Type:            body.Type,
		OccurredAt:      body.OccurredAt,
		ProductSKU:      body.ProductSKU,
		Quantity:        body.Quantity,
		Workstation:     body.Workstation,
		Message:         body.Message,
	})
	if err != nil {
		if errors.Is(err, ErrInvalidEvent) {
			return EventResponse{}, machineEventError(err)
		}
		if errors.Is(err, ErrEventConflict) {
			return EventResponse{}, fuego.ConflictError{Err: err, Detail: err.Error()}
		}
		return EventResponse{}, err
	}

	return eventResponse(event), nil
}

// Stats handles GET /machines/events/stats.
func (h *Handler) Stats(c fuego.ContextNoBody) (IntakeStatsResponse, error) {
	return intakeStatsResponse(h.stats.Stats()), nil
}

func machineEventError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{Err: err, Detail: err.Error()}
}

func eventResponse(event Event) EventResponse {
	return EventResponse(event)
}

func intakeStatsResponse(stats IntakeStats) IntakeStatsResponse {
	return IntakeStatsResponse(stats)
}
