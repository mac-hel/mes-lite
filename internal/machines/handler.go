package machines

import (
	"errors"
	"time"

	"github.com/go-fuego/fuego"
)

// Handler exposes fake machine integration endpoints.
type Handler struct {
	store Store
}

// NewHandler creates a machine handler.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
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

// CreateEvent handles POST /machines/{machineId}/events.
func (h *Handler) CreateEvent(c fuego.ContextWithBody[CreateMachineEventRequest]) (EventResponse, error) {
	body, err := c.Body()
	if err != nil {
		return EventResponse{}, err
	}

	event, err := NewEvent(c.PathParam("machineId"), body.ExternalEventID, body.Type, body.OccurredAt, body.ProductSKU, body.Quantity, body.Workstation, body.Message)
	if err != nil {
		return EventResponse{}, machineEventError(err)
	}

	if err := h.store.Save(c.Context(), event); err != nil {
		if errors.Is(err, ErrInvalidEvent) {
			return EventResponse{}, machineEventError(err)
		}
		return EventResponse{}, err
	}

	return eventResponse(event), nil
}

func machineEventError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{Err: err, Detail: err.Error()}
}

func eventResponse(event Event) EventResponse {
	return EventResponse(event)
}
