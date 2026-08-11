package production

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/products"
)

var (
	// ErrEmployeeNotFound is returned when a production entry references an unknown employee.
	ErrEmployeeNotFound = errors.New("production employee not found")
	// ErrEmployeeInactive is returned when a production entry references an inactive employee.
	ErrEmployeeInactive = errors.New("production employee inactive")
	// ErrProductNotFound is returned when a production entry references an unknown product.
	ErrProductNotFound = errors.New("production product not found")
	// ErrProductInactive is returned when a production entry references an inactive product.
	ErrProductInactive = errors.New("production product inactive")
)

// EmployeeLookup defines the employee data needed to register production.
type EmployeeLookup interface {
	FindByID(ctx context.Context, id string) (employees.Employee, error)
}

// ProductLookup defines the product data needed to register production.
type ProductLookup interface {
	FindBySKU(ctx context.Context, sku string) (products.Product, error)
}

// RegisterCommand contains the business input needed to register completed production.
type RegisterCommand struct {
	EmployeeID  string
	ProductSKU  string
	Quantity    int
	Workstation string
	Timestamp   time.Time
	Comment     string
}

// Service coordinates production registration business rules.
type Service struct {
	entries   Store
	employees EmployeeLookup
	products  ProductLookup
}

// NewService creates a production application service.
func NewService(entries Store, employees EmployeeLookup, products ProductLookup) *Service {
	return &Service{
		entries:   entries,
		employees: employees,
		products:  products,
	}
}

// Register validates business references and persists a production entry.
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (Entry, error) {
	id, err := NewEntryID()
	if err != nil {
		return Entry{}, err
	}

	entry, err := NewEntry(
		id,
		cmd.EmployeeID,
		cmd.ProductSKU,
		cmd.Quantity,
		cmd.Workstation,
		cmd.Timestamp,
		cmd.Comment,
	)
	if err != nil {
		return Entry{}, err
	}

	if err := s.validateEmployee(ctx, entry.EmployeeID); err != nil {
		return Entry{}, err
	}
	if err := s.validateProduct(ctx, entry.ProductSKU); err != nil {
		return Entry{}, err
	}

	if err := s.entries.Save(ctx, entry); err != nil {
		return Entry{}, err
	}

	return entry, nil
}

func (s *Service) validateEmployee(ctx context.Context, id string) error {
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
