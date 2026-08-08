package production

import (
	"errors"
	"math"
	"testing"
	"time"
)

func TestNewEntry_NormalizesTextFieldsAndTimestamp(t *testing.T) {
	timestamp := time.Date(2026, 8, 8, 10, 30, 0, 0, time.FixedZone("CEST", 2*60*60))

	entry := NewEntry(" entry-1 ", " emp-1 ", " sku-1 ", 12, " ws-1 ", timestamp, " done ")

	if entry.ID != "entry-1" {
		t.Errorf("expected trimmed ID, got %q", entry.ID)
	}
	if entry.EmployeeID != "emp-1" {
		t.Errorf("expected trimmed employee ID, got %q", entry.EmployeeID)
	}
	if entry.ProductSKU != "sku-1" {
		t.Errorf("expected trimmed product SKU, got %q", entry.ProductSKU)
	}
	if entry.Workstation != "ws-1" {
		t.Errorf("expected trimmed workstation, got %q", entry.Workstation)
	}
	if entry.Comment != "done" {
		t.Errorf("expected trimmed comment, got %q", entry.Comment)
	}
	if entry.Timestamp.Location() != time.UTC {
		t.Errorf("expected UTC timestamp, got %s", entry.Timestamp.Location())
	}
}

func TestEntry_Validate(t *testing.T) {
	valid := NewEntry("entry-1", "emp-1", "sku-1", 12, "ws-1", time.Now(), "")

	tests := []struct {
		name  string
		entry Entry
	}{
		{"valid", valid},
		{"missing id", Entry{EmployeeID: "emp-1", ProductSKU: "sku-1", Quantity: 12, Workstation: "ws-1", Timestamp: time.Now()}},
		{"missing employee", Entry{ID: "entry-1", ProductSKU: "sku-1", Quantity: 12, Workstation: "ws-1", Timestamp: time.Now()}},
		{"missing product", Entry{ID: "entry-1", EmployeeID: "emp-1", Quantity: 12, Workstation: "ws-1", Timestamp: time.Now()}},
		{"zero quantity", Entry{ID: "entry-1", EmployeeID: "emp-1", ProductSKU: "sku-1", Workstation: "ws-1", Timestamp: time.Now()}},
		{"negative quantity", Entry{ID: "entry-1", EmployeeID: "emp-1", ProductSKU: "sku-1", Quantity: -1, Workstation: "ws-1", Timestamp: time.Now()}},
		{"quantity too large", Entry{ID: "entry-1", EmployeeID: "emp-1", ProductSKU: "sku-1", Quantity: math.MaxInt32 + 1, Workstation: "ws-1", Timestamp: time.Now()}},
		{"missing workstation", Entry{ID: "entry-1", EmployeeID: "emp-1", ProductSKU: "sku-1", Quantity: 12, Timestamp: time.Now()}},
		{"missing timestamp", Entry{ID: "entry-1", EmployeeID: "emp-1", ProductSKU: "sku-1", Quantity: 12, Workstation: "ws-1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if tt.name == "valid" && err != nil {
				t.Fatalf("expected valid entry, got %v", err)
			}
			if tt.name != "valid" && !errors.Is(err, ErrInvalidEntry) {
				t.Fatalf("expected ErrInvalidEntry, got %v", err)
			}
		})
	}
}

func TestNewEntryID(t *testing.T) {
	id, err := NewEntryID()
	if err != nil {
		t.Fatal(err)
	}

	if len(id) != 36 {
		t.Fatalf("expected UUID-shaped ID length 36, got %d", len(id))
	}
	if id[8] != '-' || id[13] != '-' || id[18] != '-' || id[23] != '-' {
		t.Fatalf("expected UUID-shaped ID, got %q", id)
	}
}
