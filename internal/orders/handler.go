package orders

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-fuego/fuego"
)

// Handler holds the HTTP handler methods for production orders.
type Handler struct {
	service OrderService
}

// OrderService defines the order workflows needed by HTTP.
type OrderService interface {
	Create(ctx context.Context, cmd CreateCommand) (Order, error)
	Get(ctx context.Context, id string) (Order, error)
	AssignEmployee(ctx context.Context, cmd AssignEmployeeCommand) (Order, error)
	Release(ctx context.Context, id string) (Order, error)
	Start(ctx context.Context, id string) (Order, error)
	Complete(ctx context.Context, id string) (Order, error)
	Cancel(ctx context.Context, id string) (Order, error)
}

// NewHandler creates a new Handler with the given service.
func NewHandler(service OrderService) *Handler {
	return &Handler{service: service}
}

// CreateOrderRequest is the expected JSON body for creating a production order.
type CreateOrderRequest struct {
	Lines               []CreateOrderLineRequest `json:"lines"               validate:"required,min=1,dive"`
	AssignedEmployeeIDs []string                 `json:"assignedEmployeeIds"`
}

// CreateOrderLineRequest is one planned product line in a production order request.
type CreateOrderLineRequest struct {
	ProductSKU      string `json:"productSku"      validate:"required"`
	PlannedQuantity int    `json:"plannedQuantity" validate:"required,min=1"`
}

// AssignEmployeeRequest is the expected JSON body for assigning an employee to an order.
type AssignEmployeeRequest struct {
	EmployeeID string `json:"employeeId" validate:"required"`
}

// OrderResponse is the JSON representation of a production order.
type OrderResponse struct {
	ID                  string              `json:"id"`
	Lines               []OrderLineResponse `json:"lines"`
	Status              Status              `json:"status"`
	AssignedEmployeeIDs []string            `json:"assignedEmployeeIds"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
}

// OrderLineResponse is the JSON representation of one production-order line.
type OrderLineResponse struct {
	ProductSKU      string `json:"productSku"`
	PlannedQuantity int    `json:"plannedQuantity"`
}

// Create handles POST /production-orders and stores a new draft production order.
func (h *Handler) Create(c fuego.ContextWithBody[CreateOrderRequest]) (OrderResponse, error) {
	body, err := c.Body()
	if err != nil {
		return OrderResponse{}, err
	}

	lines := make([]CreateLineCommand, 0, len(body.Lines))
	for _, requestLine := range body.Lines {
		lines = append(lines, CreateLineCommand(requestLine))
	}

	order, err := h.service.Create(c.Context(), CreateCommand{Lines: lines, AssignedEmployeeIDs: body.AssignedEmployeeIDs})
	if err != nil {
		return OrderResponse{}, orderHTTPError(err, "")
	}

	return orderResponse(order), nil
}

// Get handles GET /production-orders/{id} and returns one production order.
func (h *Handler) Get(c fuego.ContextNoBody) (OrderResponse, error) {
	id := c.PathParam("id")
	order, err := h.service.Get(c.Context(), id)
	if err != nil {
		return OrderResponse{}, orderHTTPError(err, id)
	}

	return orderResponse(order), nil
}

// AssignEmployee handles POST /production-orders/{id}/assignments.
func (h *Handler) AssignEmployee(c fuego.ContextWithBody[AssignEmployeeRequest]) (OrderResponse, error) {
	body, err := c.Body()
	if err != nil {
		return OrderResponse{}, err
	}
	id := c.PathParam("id")
	order, err := h.service.AssignEmployee(c.Context(), AssignEmployeeCommand{OrderID: id, EmployeeID: body.EmployeeID})
	if err != nil {
		return OrderResponse{}, orderHTTPError(err, id)
	}

	return orderResponse(order), nil
}

// Release handles PUT /production-orders/{id}/release.
func (h *Handler) Release(c fuego.ContextNoBody) (OrderResponse, error) {
	return h.transition(c, h.service.Release)
}

// Start handles PUT /production-orders/{id}/start.
func (h *Handler) Start(c fuego.ContextNoBody) (OrderResponse, error) {
	return h.transition(c, h.service.Start)
}

// Complete handles PUT /production-orders/{id}/complete.
func (h *Handler) Complete(c fuego.ContextNoBody) (OrderResponse, error) {
	return h.transition(c, h.service.Complete)
}

// Cancel handles PUT /production-orders/{id}/cancel.
func (h *Handler) Cancel(c fuego.ContextNoBody) (OrderResponse, error) {
	return h.transition(c, h.service.Cancel)
}

func (h *Handler) transition(c fuego.ContextNoBody, change func(context.Context, string) (Order, error)) (OrderResponse, error) {
	id := c.PathParam("id")
	order, err := change(c.Context(), id)
	if err != nil {
		return OrderResponse{}, orderHTTPError(err, id)
	}

	return orderResponse(order), nil
}

func orderResponse(order Order) OrderResponse {
	lines := order.Lines().Values()
	lineResponses := make([]OrderLineResponse, 0, len(lines))
	for _, line := range lines {
		lineResponses = append(lineResponses, OrderLineResponse{ProductSKU: line.ProductSKU(), PlannedQuantity: line.PlannedQuantity()})
	}
	assignedEmployees := order.AssignedEmployees()
	if assignedEmployees == nil {
		assignedEmployees = []string{}
	}

	return OrderResponse{
		ID:                  order.ID(),
		Lines:               lineResponses,
		Status:              order.Status(),
		AssignedEmployeeIDs: assignedEmployees,
		CreatedAt:           order.CreatedAt(),
		UpdatedAt:           order.UpdatedAt(),
	}
}

func invalidOrderError(err error) fuego.BadRequestError {
	return fuego.BadRequestError{Err: err, Detail: err.Error()}
}

func orderHTTPError(err error, id string) error {
	if errors.Is(err, ErrAlreadyExists) {
		return fuego.ConflictError{Err: err, Detail: "production order already exists"}
	}
	if errors.Is(err, ErrNotFound) {
		return fuego.NotFoundError{Err: err, Detail: fmt.Sprintf("production order %q not found", id)}
	}
	if errors.Is(err, ErrInvalidOrder) || errors.Is(err, ErrInvalidTransition) || errors.Is(err, ErrEmployeeInactive) || errors.Is(err, ErrProductInactive) {
		return invalidOrderError(err)
	}
	if errors.Is(err, ErrEmployeeNotFound) || errors.Is(err, ErrProductNotFound) {
		return fuego.NotFoundError{Err: err, Detail: err.Error()}
	}
	return err
}
