package orders

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// ErrInvalidOrder is returned when production order data breaks domain rules.
var ErrInvalidOrder = errors.New("invalid production order")

// ErrInvalidTransition is returned when an order status change is not allowed.
var ErrInvalidTransition = errors.New("invalid production order status transition")

// Status describes the production order lifecycle.
type Status string

// Production order statuses.
const (
	StatusDraft      Status = "draft"
	StatusReleased   Status = "released"
	StatusInProgress Status = "in_progress"
	StatusCompleted  Status = "completed"
	StatusCancelled  Status = "cancelled"
)

// Valid reports whether the status is one of the supported order statuses.
func (s Status) Valid() bool {
	switch s {
	case StatusDraft, StatusReleased, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	default:
		return false
	}
}

// OrderLine describes one planned product inside a production order.
type OrderLine struct {
	ProductSKU      string
	PlannedQuantity int
}

// NewOrderLine creates a valid order line and normalizes text fields.
func NewOrderLine(productSKU string, plannedQuantity int) (OrderLine, error) {
	line := OrderLine{
		ProductSKU:      strings.TrimSpace(productSKU),
		PlannedQuantity: plannedQuantity,
	}
	if err := line.Validate(); err != nil {
		return OrderLine{}, err
	}

	return line, nil
}

// Validate checks the order line invariants.
func (l OrderLine) Validate() error {
	if strings.TrimSpace(l.ProductSKU) == "" {
		return fmt.Errorf("product sku is required: %w", ErrInvalidOrder)
	}
	if l.PlannedQuantity <= 0 {
		return fmt.Errorf("planned quantity must be greater than zero: %w", ErrInvalidOrder)
	}
	if l.PlannedQuantity > math.MaxInt32 {
		return fmt.Errorf("planned quantity must fit PostgreSQL integer: %w", ErrInvalidOrder)
	}

	return nil
}

// Order is the aggregate root for planned production work.
type Order struct {
	ID                string
	Lines             []OrderLine
	Status            Status
	AssignedEmployees []string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// NewOrder creates a draft production order and normalizes text and timestamps.
func NewOrder(id string, lines []OrderLine, now time.Time) (Order, error) {
	createdAt := now.UTC()
	order := Order{
		ID:        strings.TrimSpace(id),
		Lines:     copyOrderLines(lines),
		Status:    StatusDraft,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
	if err := order.Validate(); err != nil {
		return Order{}, err
	}

	return order, nil
}

// Validate checks the order invariants that must hold in every entry point.
func (o Order) Validate() error {
	if strings.TrimSpace(o.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidOrder)
	}
	if len(o.Lines) == 0 {
		return fmt.Errorf("at least one order line is required: %w", ErrInvalidOrder)
	}
	if !o.Status.Valid() {
		return fmt.Errorf("status %q is not supported: %w", o.Status, ErrInvalidOrder)
	}
	if o.CreatedAt.IsZero() {
		return fmt.Errorf("created at is required: %w", ErrInvalidOrder)
	}
	if o.UpdatedAt.IsZero() {
		return fmt.Errorf("updated at is required: %w", ErrInvalidOrder)
	}
	seenProductSKUs := make(map[string]struct{}, len(o.Lines))
	for _, line := range o.Lines {
		if err := line.Validate(); err != nil {
			return err
		}
		productSKU := strings.TrimSpace(line.ProductSKU)
		if _, ok := seenProductSKUs[productSKU]; ok {
			return fmt.Errorf("duplicate product sku %q: %w", productSKU, ErrInvalidOrder)
		}
		seenProductSKUs[productSKU] = struct{}{}
	}
	for _, employeeID := range o.AssignedEmployees {
		if strings.TrimSpace(employeeID) == "" {
			return fmt.Errorf("assigned employee id is required: %w", ErrInvalidOrder)
		}
	}

	return nil
}

// AssignEmployee assigns an employee to the order if not already assigned.
func (o *Order) AssignEmployee(employeeID string, now time.Time) error {
	employeeID = strings.TrimSpace(employeeID)
	if employeeID == "" {
		return fmt.Errorf("employee id is required: %w", ErrInvalidOrder)
	}
	if o.Status == StatusCompleted || o.Status == StatusCancelled {
		return fmt.Errorf("cannot assign employee to %s order: %w", o.Status, ErrInvalidTransition)
	}
	if o.hasAssignedEmployee(employeeID) {
		return nil
	}

	updated := *o
	updated.AssignedEmployees = append(copyAssignedEmployees(o.AssignedEmployees), employeeID)
	updated.UpdatedAt = now.UTC()
	if err := updated.Validate(); err != nil {
		return err
	}

	*o = updated
	return nil
}

// Release moves a draft order into released status.
func (o *Order) Release(now time.Time) error {
	if o.Status != StatusDraft {
		return fmt.Errorf("cannot release %s order: %w", o.Status, ErrInvalidTransition)
	}
	if len(o.AssignedEmployees) == 0 {
		return fmt.Errorf("released order requires at least one assigned employee: %w", ErrInvalidOrder)
	}

	return o.changeStatus(StatusReleased, now)
}

// Start moves a released order into in-progress status.
func (o *Order) Start(now time.Time) error {
	if o.Status != StatusReleased {
		return fmt.Errorf("cannot start %s order: %w", o.Status, ErrInvalidTransition)
	}

	return o.changeStatus(StatusInProgress, now)
}

// Complete moves an in-progress order into completed status.
func (o *Order) Complete(now time.Time) error {
	if o.Status != StatusInProgress {
		return fmt.Errorf("cannot complete %s order: %w", o.Status, ErrInvalidTransition)
	}

	return o.changeStatus(StatusCompleted, now)
}

// Cancel moves a non-completed order into cancelled status.
func (o *Order) Cancel(now time.Time) error {
	if o.Status == StatusCompleted {
		return fmt.Errorf("cannot cancel completed order: %w", ErrInvalidTransition)
	}
	if o.Status == StatusCancelled {
		return nil
	}

	return o.changeStatus(StatusCancelled, now)
}

func (o *Order) changeStatus(status Status, now time.Time) error {
	updated := *o
	updated.Status = status
	updated.UpdatedAt = now.UTC()
	if err := updated.Validate(); err != nil {
		return err
	}

	*o = updated
	return nil
}

func (o Order) hasAssignedEmployee(employeeID string) bool {
	for _, assigned := range o.AssignedEmployees {
		if assigned == employeeID {
			return true
		}
	}

	return false
}

func copyAssignedEmployees(employeeIDs []string) []string {
	if len(employeeIDs) == 0 {
		return nil
	}

	copied := make([]string, len(employeeIDs))
	copy(copied, employeeIDs)
	return copied
}

func copyOrderLines(lines []OrderLine) []OrderLine {
	if len(lines) == 0 {
		return nil
	}

	copied := make([]OrderLine, len(lines))
	copy(copied, lines)
	return copied
}
