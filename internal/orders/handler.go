package orders

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-fuego/fuego"
)

// Handler holds the HTTP handler methods for production orders.
type Handler struct {
	store Store
}

// NewHandler creates a new Handler with the given Store.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
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

	lines := make([]OrderLine, 0, len(body.Lines))
	for _, requestLine := range body.Lines {
		line, err := NewOrderLine(requestLine.ProductSKU, requestLine.PlannedQuantity)
		if err != nil {
			return OrderResponse{}, invalidOrderError(err)
		}
		lines = append(lines, line)
	}
	orderLines, err := NewOrderLines(lines...)
	if err != nil {
		return OrderResponse{}, invalidOrderError(err)
	}
	id, err := NewOrderID()
	if err != nil {
		return OrderResponse{}, err
	}
	order, err := NewOrder(id, orderLines, time.Now())
	if err != nil {
		return OrderResponse{}, invalidOrderError(err)
	}
	for _, employeeID := range body.AssignedEmployeeIDs {
		if err := order.AssignEmployee(employeeID, time.Now()); err != nil {
			return OrderResponse{}, invalidOrderError(err)
		}
	}

	if err := h.store.Save(c.Context(), order); err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			return OrderResponse{}, fuego.ConflictError{Err: err, Detail: fmt.Sprintf("production order %q already exists", order.ID())}
		}
		if errors.Is(err, ErrInvalidOrder) {
			return OrderResponse{}, invalidOrderError(err)
		}
		return OrderResponse{}, err
	}

	return orderResponse(order), nil
}

// Get handles GET /production-orders/{id} and returns one production order.
func (h *Handler) Get(c fuego.ContextNoBody) (OrderResponse, error) {
	id := c.PathParam("id")
	order, err := h.store.FindByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return OrderResponse{}, fuego.NotFoundError{Err: err, Detail: fmt.Sprintf("production order %q not found", id)}
		}
		return OrderResponse{}, err
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
