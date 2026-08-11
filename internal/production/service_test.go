package production

import (
	"errors"
	"testing"
	"time"

	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/products"
)

func newTestService(t *testing.T) (*Service, *InMemoryStore, *employees.InMemoryStore, *products.InMemoryStore) {
	t.Helper()

	entryStore := NewInMemoryStore()
	empStore := employees.NewInMemoryStore()
	prodStore := products.NewInMemoryStore()

	emp, err := employees.NewEmployee("emp-1", "Ana", "Worker", "ana@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if err := empStore.Save(t.Context(), emp); err != nil {
		t.Fatal(err)
	}
	prod, err := products.NewProduct("sku-1", "Ventilation Unit", "piece", products.CategoryVentilation)
	if err != nil {
		t.Fatal(err)
	}
	if err := prodStore.Save(t.Context(), prod); err != nil {
		t.Fatal(err)
	}

	return NewService(entryStore, empStore, prodStore), entryStore, empStore, prodStore
}

func validRegisterCommand() RegisterCommand {
	return RegisterCommand{
		EmployeeID:  "emp-1",
		ProductSKU:  "sku-1",
		Quantity:    12,
		Workstation: "ws-1",
		Timestamp:   time.Date(2026, 8, 8, 10, 30, 0, 0, time.UTC),
		Comment:     "batch finished",
	}
}

func TestService_Register(t *testing.T) {
	service, store, _, _ := newTestService(t)

	entry, err := service.Register(t.Context(), validRegisterCommand())
	if err != nil {
		t.Fatal(err)
	}

	if entry.ID == "" {
		t.Fatal("expected generated ID")
	}
	if entry.EmployeeID != "emp-1" {
		t.Errorf("expected employee emp-1, got %q", entry.EmployeeID)
	}
	if entry.ProductSKU != "sku-1" {
		t.Errorf("expected product sku-1, got %q", entry.ProductSKU)
	}

	stored, err := store.FindByID(t.Context(), entry.ID)
	if err != nil {
		t.Fatal(err)
	}
	if stored != entry {
		t.Errorf("expected stored entry %#v, got %#v", entry, stored)
	}
}

func TestService_Register_MissingEmployee(t *testing.T) {
	service, _, _, _ := newTestService(t)
	cmd := validRegisterCommand()
	cmd.EmployeeID = "missing"

	_, err := service.Register(t.Context(), cmd)
	if !errors.Is(err, ErrEmployeeNotFound) {
		t.Fatalf("expected ErrEmployeeNotFound, got %v", err)
	}
}

func TestService_Register_InactiveEmployee(t *testing.T) {
	service, _, empStore, _ := newTestService(t)
	emp, err := empStore.FindByID(t.Context(), "emp-1")
	if err != nil {
		t.Fatal(err)
	}
	emp.IsActive = false
	if _, err := empStore.Update(t.Context(), emp); err != nil {
		t.Fatal(err)
	}

	_, err = service.Register(t.Context(), validRegisterCommand())
	if !errors.Is(err, ErrEmployeeInactive) {
		t.Fatalf("expected ErrEmployeeInactive, got %v", err)
	}
}

func TestService_Register_MissingProduct(t *testing.T) {
	service, _, _, _ := newTestService(t)
	cmd := validRegisterCommand()
	cmd.ProductSKU = "missing"

	_, err := service.Register(t.Context(), cmd)
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("expected ErrProductNotFound, got %v", err)
	}
}

func TestService_Register_InactiveProduct(t *testing.T) {
	service, _, _, prodStore := newTestService(t)
	prod, err := prodStore.FindBySKU(t.Context(), "sku-1")
	if err != nil {
		t.Fatal(err)
	}
	prod.IsActive = false
	if _, err := prodStore.Update(t.Context(), prod); err != nil {
		t.Fatal(err)
	}

	_, err = service.Register(t.Context(), validRegisterCommand())
	if !errors.Is(err, ErrProductInactive) {
		t.Fatalf("expected ErrProductInactive, got %v", err)
	}
}

func TestService_Register_InvalidEntry(t *testing.T) {
	service, _, _, _ := newTestService(t)
	cmd := validRegisterCommand()
	cmd.Quantity = 0

	_, err := service.Register(t.Context(), cmd)
	if !errors.Is(err, ErrInvalidEntry) {
		t.Fatalf("expected ErrInvalidEntry, got %v", err)
	}
}
