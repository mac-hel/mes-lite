package orders

import (
	"errors"
	"testing"

	"github.com/mac-hel/mes-lite/internal/employees"
	"github.com/mac-hel/mes-lite/internal/products"
)

func TestService_CreateValidatesProductAndEmployeeReferences(t *testing.T) {
	service := testService(t)

	order, err := service.Create(t.Context(), CreateCommand{
		Lines:               []CreateLineCommand{{ProductSKU: "shaft-1", PlannedQuantity: 2}},
		AssignedEmployeeIDs: []string{"emp-1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if order.ID() == "" {
		t.Fatal("expected generated order ID")
	}
	if len(order.AssignedEmployees()) != 1 {
		t.Fatalf("assigned employee count = %d, want 1", len(order.AssignedEmployees()))
	}
}

func TestService_CreateRejectsMissingProduct(t *testing.T) {
	service := testService(t)

	_, err := service.Create(t.Context(), CreateCommand{Lines: []CreateLineCommand{{ProductSKU: "missing", PlannedQuantity: 2}}})
	if !errors.Is(err, ErrProductNotFound) {
		t.Fatalf("Create() error = %v, want ErrProductNotFound", err)
	}
}

func TestService_CreateRejectsInactiveProduct(t *testing.T) {
	ordersStore, empStore, prodStore := seededOrderStores(t)
	prod, err := prodStore.FindBySKU(t.Context(), "shaft-1")
	if err != nil {
		t.Fatal(err)
	}
	prod.IsActive = false
	if _, err := prodStore.Update(t.Context(), prod); err != nil {
		t.Fatal(err)
	}
	service := NewService(ordersStore, empStore, prodStore)

	_, err = service.Create(t.Context(), CreateCommand{Lines: []CreateLineCommand{{ProductSKU: "shaft-1", PlannedQuantity: 2}}})
	if !errors.Is(err, ErrProductInactive) {
		t.Fatalf("Create() error = %v, want ErrProductInactive", err)
	}
}

func TestService_AssignEmployee(t *testing.T) {
	service := testService(t)
	order := saveDraftOrder(t, service)

	updated, err := service.AssignEmployee(t.Context(), AssignEmployeeCommand{OrderID: order.ID(), EmployeeID: "emp-1"})
	if err != nil {
		t.Fatal(err)
	}
	if len(updated.AssignedEmployees()) != 1 || updated.AssignedEmployees()[0] != "emp-1" {
		t.Fatalf("AssignedEmployees = %#v, want [emp-1]", updated.AssignedEmployees())
	}
	if updated.Version() != order.Version()+1 {
		t.Fatalf("Version = %d, want %d", updated.Version(), order.Version()+1)
	}
}

func TestService_AssignEmployeeRejectsMissingEmployee(t *testing.T) {
	service := testService(t)
	order := saveDraftOrder(t, service)

	_, err := service.AssignEmployee(t.Context(), AssignEmployeeCommand{OrderID: order.ID(), EmployeeID: "missing"})
	if !errors.Is(err, ErrEmployeeNotFound) {
		t.Fatalf("AssignEmployee() error = %v, want ErrEmployeeNotFound", err)
	}
}

func TestService_AssignEmployeeRejectsInactiveEmployee(t *testing.T) {
	ordersStore, empStore, prodStore := seededOrderStores(t)
	emp, err := empStore.FindByID(t.Context(), "emp-1")
	if err != nil {
		t.Fatal(err)
	}
	emp.IsActive = false
	if _, err := empStore.Update(t.Context(), emp); err != nil {
		t.Fatal(err)
	}
	service := NewService(ordersStore, empStore, prodStore)
	order := saveDraftOrder(t, service)

	_, err = service.AssignEmployee(t.Context(), AssignEmployeeCommand{OrderID: order.ID(), EmployeeID: "emp-1"})
	if !errors.Is(err, ErrEmployeeInactive) {
		t.Fatalf("AssignEmployee() error = %v, want ErrEmployeeInactive", err)
	}
}

func TestService_StatusTransitions(t *testing.T) {
	service := testService(t)
	order := saveDraftOrder(t, service)
	order, err := service.AssignEmployee(t.Context(), AssignEmployeeCommand{OrderID: order.ID(), EmployeeID: "emp-1"})
	if err != nil {
		t.Fatal(err)
	}

	order, err = service.Release(t.Context(), order.ID())
	if err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusReleased {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusReleased)
	}
	if order.Version() != 3 {
		t.Fatalf("Version after release = %d, want 3", order.Version())
	}
	order, err = service.Start(t.Context(), order.ID())
	if err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusInProgress {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusInProgress)
	}
	order, err = service.Complete(t.Context(), order.ID())
	if err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusCompleted {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusCompleted)
	}
}

func TestService_InvalidStatusTransition(t *testing.T) {
	service := testService(t)
	order := saveDraftOrder(t, service)

	_, err := service.Start(t.Context(), order.ID())
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("Start() error = %v, want ErrInvalidTransition", err)
	}
}

func testService(t *testing.T) *Service {
	t.Helper()
	ordersStore, empStore, prodStore := seededOrderStores(t)
	return NewService(ordersStore, empStore, prodStore)
}

func seededOrderStores(t *testing.T) (*InMemoryStore, *employees.InMemoryStore, *products.InMemoryStore) {
	t.Helper()
	return NewInMemoryStore(), seededEmployeeStore(t), seededProductStore(t)
}

func saveDraftOrder(t *testing.T, service *Service) Order {
	t.Helper()
	order, err := service.Create(t.Context(), CreateCommand{Lines: []CreateLineCommand{{ProductSKU: "shaft-1", PlannedQuantity: 2}}})
	if err != nil {
		t.Fatal(err)
	}
	return order
}
