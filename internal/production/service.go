package production

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	RequestID   string
	EmployeeID  string
	ProductSKU  string
	Quantity    int
	Workstation string
	Timestamp   time.Time
	Comment     string
}

// CorrectEntryCommand contains the replacement values and audit data for a correction.
type CorrectEntryCommand struct {
	EntryID     string
	ActorUserID string
	Reason      string
	EmployeeID  string
	ProductSKU  string
	Quantity    int
	Workstation string
	Timestamp   time.Time
	Comment     string
}

// Service coordinates production registration business rules.
type Service struct {
	entries     EntryStore
	corrections CorrectionStore
	employees   EmployeeLookup
	products    ProductLookup
}

// NewService creates a production application service.
func NewService(entries EntryStore, corrections CorrectionStore, employees EmployeeLookup, products ProductLookup) *Service {
	return &Service{
		entries:     entries,
		corrections: corrections,
		employees:   employees,
		products:    products,
	}
}

// Register validates business references and persists a production entry.
func (s *Service) Register(ctx context.Context, cmd RegisterCommand) (Entry, error) {
	if strings.TrimSpace(cmd.RequestID) == "" {
		return Entry{}, fmt.Errorf("request id is required: %w", ErrInvalidEntry)
	}

	id, err := NewEntryID()
	if err != nil {
		return Entry{}, err
	}

	entry, err := NewEntryWithRequestID(
		id,
		cmd.RequestID,
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
		if errors.Is(err, ErrRequestConflict) {
			return s.resolveIdempotentRegistration(ctx, entry)
		}
		return Entry{}, err
	}

	return entry, nil
}

func (s *Service) resolveIdempotentRegistration(ctx context.Context, entry Entry) (Entry, error) {
	existing, err := s.entries.FindByRequestID(ctx, entry.RequestID)
	if err != nil {
		return Entry{}, err
	}
	if sameProductionRequest(existing, entry) {
		return existing, nil
	}
	return Entry{}, fmt.Errorf("request id %q already used for different production entry: %w", entry.RequestID, ErrRequestConflict)
}

func sameProductionRequest(a, b Entry) bool {
	return a.RequestID == b.RequestID &&
		a.EmployeeID == b.EmployeeID &&
		a.ProductSKU == b.ProductSKU &&
		a.Quantity == b.Quantity &&
		a.Workstation == b.Workstation &&
		a.Timestamp.Equal(b.Timestamp) &&
		a.Comment == b.Comment
}

// List returns production entries for review workflows.
func (s *Service) List(ctx context.Context, opts ListOptions) ([]Entry, error) {
	return s.entries.List(ctx, opts)
}

// CorrectEntry appends a correction without overwriting the original production entry.
func (s *Service) CorrectEntry(ctx context.Context, cmd CorrectEntryCommand) (Correction, error) {
	if _, err := s.entries.FindByID(ctx, cmd.EntryID); err != nil {
		return Correction{}, err
	}

	id, err := NewEntryID()
	if err != nil {
		return Correction{}, err
	}
	correction, err := NewCorrection(id, cmd.EntryID, cmd.ActorUserID, cmd.Reason, cmd.EmployeeID, cmd.ProductSKU, cmd.Quantity, cmd.Workstation, cmd.Timestamp, cmd.Comment)
	if err != nil {
		return Correction{}, err
	}
	correction.CreatedAt = time.Now().UTC()

	if err := s.validateEmployee(ctx, correction.EmployeeID); err != nil {
		return Correction{}, err
	}
	if err := s.validateProduct(ctx, correction.ProductSKU); err != nil {
		return Correction{}, err
	}

	if err := s.corrections.SaveCorrection(ctx, correction); err != nil {
		return Correction{}, err
	}

	return correction, nil
}

// ListCorrections returns append-only correction history for a production entry.
func (s *Service) ListCorrections(ctx context.Context, entryID string) ([]Correction, error) {
	if _, err := s.entries.FindByID(ctx, entryID); err != nil {
		return nil, err
	}
	return s.corrections.ListCorrections(ctx, entryID)
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
