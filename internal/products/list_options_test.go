package products

import "testing"

func TestInMemoryStore_ListOptions(t *testing.T) {
	store := NewInMemoryStore()
	seedProduct(t, store, "VX-100", "Ventilation X100", CategoryVentilation, true)
	seedProduct(t, store, "FT-200", "Filter T200", CategoryFilter, false)
	seedProduct(t, store, "DC-300", "Duct C300", CategoryDuct, true)

	active := true
	got, err := store.List(t.Context(), ListOptions{Limit: 1, Offset: 1, Sort: "name", Query: "0", Active: &active})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 product, got %d", len(got))
	}
	if got[0].SKU != "VX-100" {
		t.Fatalf("expected second active product sorted by name, got %q", got[0].SKU)
	}
}

func seedProduct(t *testing.T, store *InMemoryStore, sku, name string, category ProductCategory, active bool) {
	t.Helper()
	prod, err := NewProduct(sku, name, "piece", category)
	if err != nil {
		t.Fatal(err)
	}
	prod.IsActive = active
	if err := store.Save(t.Context(), prod); err != nil {
		t.Fatal(err)
	}
}
