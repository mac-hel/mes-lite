package orders

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestNewOrder_NormalizesTextAndTimestamps(t *testing.T) {
	now := time.Date(2026, 8, 12, 10, 30, 0, 0, time.FixedZone("CEST", 2*60*60))

	order, err := NewOrder(" order-1 ", mustOrderLines(t, mustOrderLine(t, " VX-100 ", 100), mustOrderLine(t, " FILTER-1 ", 4)), now)
	if err != nil {
		t.Fatal(err)
	}

	if order.ID() != "order-1" {
		t.Errorf("ID = %q, want %q", order.ID(), "order-1")
	}
	lines := order.Lines().Values()
	if len(lines) != 2 {
		t.Fatalf("line count = %d, want 2", len(lines))
	}
	if lines[0].ProductSKU() != "VX-100" {
		t.Errorf("first line ProductSKU = %q, want %q", lines[0].ProductSKU(), "VX-100")
	}
	if lines[1].ProductSKU() != "FILTER-1" {
		t.Errorf("second line ProductSKU = %q, want %q", lines[1].ProductSKU(), "FILTER-1")
	}
	if order.Status() != StatusDraft {
		t.Errorf("Status = %q, want %q", order.Status(), StatusDraft)
	}
	if order.Version() != 1 {
		t.Errorf("Version = %d, want 1", order.Version())
	}
	if order.CreatedAt().Location() != time.UTC {
		t.Errorf("CreatedAt location = %s, want UTC", order.CreatedAt().Location())
	}
	if order.UpdatedAt().Location() != time.UTC {
		t.Errorf("UpdatedAt location = %s, want UTC", order.UpdatedAt().Location())
	}
}

func TestNewOrderLines_CopiesLines(t *testing.T) {
	lines := []OrderLine{mustOrderLine(t, "VX-100", 100)}
	orderLines, err := NewOrderLines(lines...)
	if err != nil {
		t.Fatal(err)
	}

	lines[0].productSKU = "FILTER-1"

	if orderLines.Values()[0].ProductSKU() != "VX-100" {
		t.Errorf("order line was mutated through caller-owned slice, got %q", orderLines.Values()[0].ProductSKU())
	}
}

func TestOrderLines_ValuesReturnsCopy(t *testing.T) {
	orderLines := mustOrderLines(t, mustOrderLine(t, "VX-100", 100))
	lines := orderLines.Values()
	lines[0].productSKU = "FILTER-1"

	if orderLines.Values()[0].ProductSKU() != "VX-100" {
		t.Errorf("order line was mutated through returned slice, got %q", orderLines.Values()[0].ProductSKU())
	}
}

func TestNewOrderLines(t *testing.T) {
	orderLines, err := NewOrderLines(mustOrderLine(t, "VX-100", 100), mustOrderLine(t, "FILTER-1", 4))
	if err != nil {
		t.Fatal(err)
	}

	if orderLines.Len() != 2 {
		t.Fatalf("Len() = %d, want 2", orderLines.Len())
	}
}

func TestOrder_AssignedEmployeesReturnsCopy(t *testing.T) {
	order := mustOrder(t)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}

	assignedEmployees := order.AssignedEmployees()
	assignedEmployees[0] = "emp-2"

	if order.AssignedEmployees()[0] != "emp-1" {
		t.Errorf("assigned employee was mutated through returned slice, got %q", order.AssignedEmployees()[0])
	}
}

func TestNewOrderLine(t *testing.T) {
	line, err := NewOrderLine(" VX-100 ", 100)
	if err != nil {
		t.Fatal(err)
	}

	if line.ProductSKU() != "VX-100" {
		t.Errorf("ProductSKU = %q, want %q", line.ProductSKU(), "VX-100")
	}
	if line.PlannedQuantity() != 100 {
		t.Errorf("PlannedQuantity = %d, want 100", line.PlannedQuantity())
	}
}

func TestOrderLine_Validate(t *testing.T) {
	tests := []struct {
		name string
		line OrderLine
	}{
		{"valid", mustOrderLine(t, "VX-100", 100)},
		{"blank product sku", OrderLine{plannedQuantity: 100}},
		{"zero planned quantity", OrderLine{productSKU: "VX-100"}},
		{"negative planned quantity", OrderLine{productSKU: "VX-100", plannedQuantity: -1}},
		{"too large planned quantity", OrderLine{productSKU: "VX-100", plannedQuantity: math.MaxInt32 + 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.line.Validate()
			if tt.name == "valid" && err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
			if tt.name != "valid" && !errors.Is(err, ErrInvalidOrder) {
				t.Fatalf("Validate() error = %v, want ErrInvalidOrder", err)
			}
		})
	}
}

func TestOrderLines_Validate(t *testing.T) {
	tests := []struct {
		name       string
		orderLines OrderLines
	}{
		{"valid", mustOrderLines(t, mustOrderLine(t, "VX-100", 100))},
		{"missing lines", OrderLines{}},
		{"invalid line", OrderLines{values: []OrderLine{{productSKU: "VX-100"}}}},
		{"duplicate product sku", OrderLines{values: []OrderLine{mustOrderLine(t, "VX-100", 100), mustOrderLine(t, "VX-100", 50)}}},
		{"duplicate product sku with spaces", OrderLines{values: []OrderLine{mustOrderLine(t, "VX-100", 100), {productSKU: " VX-100 ", plannedQuantity: 50}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.orderLines.Validate()
			if tt.name == "valid" && err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
			if tt.name != "valid" && !errors.Is(err, ErrInvalidOrder) {
				t.Fatalf("Validate() error = %v, want ErrInvalidOrder", err)
			}
		})
	}
}

func TestNewOrderLine_RejectsInvalidState(t *testing.T) {
	_, err := NewOrderLine("", 100)
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("NewOrderLine() error = %v, want ErrInvalidOrder", err)
	}
}

func TestOrder_Validate(t *testing.T) {
	now := time.Now()
	valid, err := NewOrder("order-1", mustOrderLines(t, mustOrderLine(t, "VX-100", 100), mustOrderLine(t, "FILTER-1", 4)), now)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name  string
		order Order
	}{
		{"valid", valid},
		{"missing id", Order{lines: mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), status: StatusDraft, createdAt: now, updatedAt: now}},
		{"missing lines", Order{id: "order-1", status: StatusDraft, createdAt: now, updatedAt: now}},
		{"invalid line", Order{id: "order-1", lines: OrderLines{values: []OrderLine{{productSKU: "VX-100"}}}, status: StatusDraft, createdAt: now, updatedAt: now}},
		{"invalid status", Order{id: "order-1", lines: mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), status: Status("unknown"), createdAt: now, updatedAt: now}},
		{"missing created at", Order{id: "order-1", lines: mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), status: StatusDraft, updatedAt: now}},
		{"missing updated at", Order{id: "order-1", lines: mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), status: StatusDraft, createdAt: now}},
		{"blank assigned employee", Order{id: "order-1", lines: mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), status: StatusDraft, createdAt: now, updatedAt: now, assignedEmployees: []string{"emp-1", " "}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.order.Validate()
			if tt.name == "valid" && err != nil {
				t.Fatalf("Validate() error = %v, want nil", err)
			}
			if tt.name != "valid" && !errors.Is(err, ErrInvalidOrder) {
				t.Fatalf("Validate() error = %v, want ErrInvalidOrder", err)
			}
		})
	}
}

func TestNewOrder_RejectsInvalidState(t *testing.T) {
	_, err := NewOrder("", mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), time.Now())
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("NewOrder() error = %v, want ErrInvalidOrder", err)
	}
}

func TestStatus_Valid(t *testing.T) {
	tests := []struct {
		name   string
		status Status
		want   bool
	}{
		{"draft", StatusDraft, true},
		{"released", StatusReleased, true},
		{"in progress", StatusInProgress, true},
		{"completed", StatusCompleted, true},
		{"cancelled", StatusCancelled, true},
		{"empty", Status(""), false},
		{"unknown", Status("unknown"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.status.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOrder_AssignEmployee(t *testing.T) {
	order := mustOrder(t)
	now := time.Date(2026, 8, 12, 11, 0, 0, 0, time.UTC)

	if err := order.AssignEmployee(" emp-1 ", now); err != nil {
		t.Fatal(err)
	}

	assignedEmployees := order.AssignedEmployees()
	if len(assignedEmployees) != 1 || assignedEmployees[0] != "emp-1" {
		t.Fatalf("AssignedEmployees = %#v, want [emp-1]", assignedEmployees)
	}
	if !order.UpdatedAt().Equal(now) {
		t.Errorf("UpdatedAt = %s, want %s", order.UpdatedAt(), now)
	}
}

func TestOrder_AssignEmployeeIsIdempotent(t *testing.T) {
	order := mustOrder(t)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}

	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}

	if len(order.AssignedEmployees()) != 1 {
		t.Fatalf("assigned employee count = %d, want 1", len(order.AssignedEmployees()))
	}
}

func TestOrder_AssignEmployeeRejectsBlankEmployee(t *testing.T) {
	order := mustOrder(t)

	err := order.AssignEmployee(" ", time.Now())
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("AssignEmployee() error = %v, want ErrInvalidOrder", err)
	}
}

func TestOrder_AssignEmployeeRejectsClosedOrder(t *testing.T) {
	order := mustOrder(t)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := order.Release(time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := order.Start(time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := order.Complete(time.Now()); err != nil {
		t.Fatal(err)
	}

	err := order.AssignEmployee("emp-2", time.Now())
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("AssignEmployee() error = %v, want ErrInvalidTransition", err)
	}
}

func TestOrder_StatusTransitions(t *testing.T) {
	order := mustOrder(t)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}

	if err := order.Release(time.Now()); err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusReleased {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusReleased)
	}

	if err := order.Start(time.Now()); err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusInProgress {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusInProgress)
	}

	if err := order.Complete(time.Now()); err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusCompleted {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusCompleted)
	}
}

func TestOrder_ReleaseRequiresAssignedEmployee(t *testing.T) {
	order := mustOrder(t)

	err := order.Release(time.Now())
	if !errors.Is(err, ErrInvalidOrder) {
		t.Fatalf("Release() error = %v, want ErrInvalidOrder", err)
	}
	if order.Status() != StatusDraft {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusDraft)
	}
}

func TestOrder_InvalidStatusTransitions(t *testing.T) {
	tests := []struct {
		name string
		act  func(*Order) error
	}{
		{"start draft", func(o *Order) error { return o.Start(time.Now()) }},
		{"complete draft", func(o *Order) error { return o.Complete(time.Now()) }},
		{"release cancelled", func(o *Order) error {
			if err := o.Cancel(time.Now()); err != nil {
				return err
			}
			return o.Release(time.Now())
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			order := mustOrder(t)
			err := tt.act(&order)
			if !errors.Is(err, ErrInvalidTransition) {
				t.Fatalf("transition error = %v, want ErrInvalidTransition", err)
			}
		})
	}
}

func TestOrder_Cancel(t *testing.T) {
	order := mustOrder(t)

	if err := order.Cancel(time.Now()); err != nil {
		t.Fatal(err)
	}
	if order.Status() != StatusCancelled {
		t.Fatalf("Status = %q, want %q", order.Status(), StatusCancelled)
	}
}

func TestOrder_CancelRejectsCompletedOrder(t *testing.T) {
	order := mustOrder(t)
	if err := order.AssignEmployee("emp-1", time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := order.Release(time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := order.Start(time.Now()); err != nil {
		t.Fatal(err)
	}
	if err := order.Complete(time.Now()); err != nil {
		t.Fatal(err)
	}

	err := order.Cancel(time.Now())
	if !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("Cancel() error = %v, want ErrInvalidTransition", err)
	}
}

func TestSentinelErrors(t *testing.T) {
	if ErrInvalidOrder == nil {
		t.Fatal("ErrInvalidOrder should not be nil")
	}
	if ErrInvalidOrder.Error() != "invalid production order" {
		t.Errorf("ErrInvalidOrder = %q, want %q", ErrInvalidOrder.Error(), "invalid production order")
	}
	if ErrInvalidTransition == nil {
		t.Fatal("ErrInvalidTransition should not be nil")
	}
	if ErrInvalidTransition.Error() != "invalid production order status transition" {
		t.Errorf("ErrInvalidTransition = %q, want %q", ErrInvalidTransition.Error(), "invalid production order status transition")
	}
}

func mustOrder(t *testing.T) Order {
	t.Helper()

	order, err := NewOrder("order-1", mustOrderLines(t, mustOrderLine(t, "VX-100", 100)), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	return order
}

func mustOrderLine(t *testing.T, productSKU string, plannedQuantity int) OrderLine {
	t.Helper()

	line, err := NewOrderLine(productSKU, plannedQuantity)
	if err != nil {
		t.Fatal(err)
	}
	return line
}

func mustOrderLines(t *testing.T, lines ...OrderLine) OrderLines {
	t.Helper()

	orderLines, err := NewOrderLines(lines...)
	if err != nil {
		t.Fatal(err)
	}
	return orderLines
}
