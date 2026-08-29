package orders

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/ids"
	"github.com/mac-hel/mes-lite/internal/products"
)

var (
	// ErrEmployeeNotFound is returned when an order references an unknown employee.
	ErrEmployeeNotFound = errors.New("order employee not found")
	// ErrEmployeeInactive is returned when an order references an inactive employee.
	ErrEmployeeInactive = errors.New("order employee inactive")
	// ErrProductNotFound is returned when an order line references an unknown product.
	ErrProductNotFound = errors.New("order product not found")
	// ErrProductInactive is returned when an order line references an inactive product.
	ErrProductInactive = errors.New("order product inactive")
)

// EmployeeLookup defines the employee data needed by order workflows.
type EmployeeLookup interface {
	FindByID(ctx context.Context, id string) (employees.Employee, error)
}

// ProductLookup defines the product data needed by order workflows.
type ProductLookup interface {
	FindBySKU(ctx context.Context, sku string) (products.Product, error)
}

// CreateCommand contains the business input needed to create an order.
type CreateCommand struct {
	Lines               []CreateLineCommand
	AssignedEmployeeIDs []string
}

// CreateLineCommand contains one planned product line.
type CreateLineCommand struct {
	ProductSKU      string
	PlannedQuantity int
}

// AssignEmployeeCommand contains the business input needed to assign an employee.
type AssignEmployeeCommand struct {
	OrderID    string
	EmployeeID string
}

// Service coordinates production order business workflows.
type Service struct {
	orders    Store
	employees EmployeeLookup
	products  ProductLookup
}

// NewService creates an order application service.
func NewService(orders Store, employees EmployeeLookup, products ProductLookup) *Service {
	return &Service{orders: orders, employees: employees, products: products}
}

// Create validates references and persists a new draft order.
func (s *Service) Create(ctx context.Context, cmd CreateCommand) (Order, error) {
	lines := make([]OrderLine, 0, len(cmd.Lines))
	for _, lineCmd := range cmd.Lines {
		line, err := NewOrderLine(lineCmd.ProductSKU, lineCmd.PlannedQuantity)
		if err != nil {
			return Order{}, err
		}
		if err := s.validateProduct(ctx, line.ProductSKU()); err != nil {
			return Order{}, err
		}
		lines = append(lines, line)
	}
	orderLines, err := NewOrderLines(lines...)
	if err != nil {
		return Order{}, err
	}
	order, err := NewOrder(ids.New(), orderLines, time.Now())
	if err != nil {
		return Order{}, err
	}
	for _, employeeID := range cmd.AssignedEmployeeIDs {
		if err := s.validateEmployee(ctx, employeeID); err != nil {
			return Order{}, err
		}
		if err := order.AssignEmployee(employeeID, time.Now()); err != nil {
			return Order{}, err
		}
	}
	if err := s.orders.Save(ctx, order); err != nil {
		return Order{}, err
	}

	return order, nil
}

// Get returns a production order by ID.
func (s *Service) Get(ctx context.Context, id string) (Order, error) {
	return s.orders.FindByID(ctx, id)
}

// AssignEmployee assigns an active employee to an existing order.
func (s *Service) AssignEmployee(ctx context.Context, cmd AssignEmployeeCommand) (Order, error) {
	if err := s.validateEmployee(ctx, cmd.EmployeeID); err != nil {
		return Order{}, err
	}
	order, err := s.orders.FindByID(ctx, cmd.OrderID)
	if err != nil {
		return Order{}, err
	}
	if err := order.AssignEmployee(cmd.EmployeeID, time.Now()); err != nil {
		return Order{}, err
	}
	if err := s.orders.Update(ctx, order); err != nil {
		return Order{}, err
	}

	return s.orders.FindByID(ctx, order.ID())
}

// Release moves a draft order into released status.
func (s *Service) Release(ctx context.Context, id string) (Order, error) {
	return s.transition(ctx, id, func(order *Order) error { return order.Release(time.Now()) })
}

// Start moves a released order into in-progress status.
func (s *Service) Start(ctx context.Context, id string) (Order, error) {
	return s.transition(ctx, id, func(order *Order) error { return order.Start(time.Now()) })
}

// Complete moves an in-progress order into completed status.
func (s *Service) Complete(ctx context.Context, id string) (Order, error) {
	return s.transition(ctx, id, func(order *Order) error { return order.Complete(time.Now()) })
}

// Cancel moves a non-completed order into cancelled status.
func (s *Service) Cancel(ctx context.Context, id string) (Order, error) {
	return s.transition(ctx, id, func(order *Order) error { return order.Cancel(time.Now()) })
}

func (s *Service) transition(ctx context.Context, id string, change func(*Order) error) (Order, error) {
	order, err := s.orders.FindByID(ctx, id)
	if err != nil {
		return Order{}, err
	}
	if err := change(&order); err != nil {
		return Order{}, err
	}
	if err := s.orders.Update(ctx, order); err != nil {
		return Order{}, err
	}

	return s.orders.FindByID(ctx, order.ID())
}

func (s *Service) validateEmployee(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return fmt.Errorf("employee id is required: %w", ErrInvalidOrder)
	}
	emp, err := s.employees.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, employees.ErrNotFound) {
			return fmt.Errorf("employee %q: %w", id, ErrEmployeeNotFound)
		}
		return err
	}
	if !emp.IsActive {
		return fmt.Errorf("employee %q: %w", id, ErrEmployeeInactive)
	}
	return nil
}

func (s *Service) validateProduct(ctx context.Context, sku string) error {
	prod, err := s.products.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, products.ErrNotFound) {
			return fmt.Errorf("product %q: %w", sku, ErrProductNotFound)
		}
		return err
	}
	if !prod.IsActive {
		return fmt.Errorf("product %q: %w", sku, ErrProductInactive)
	}
	return nil
}
