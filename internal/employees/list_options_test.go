package employees

import "testing"

func TestInMemoryStore_ListOptions(t *testing.T) {
	store := NewInMemoryStore()
	seedEmployee(t, store, "003", "Carla", "Zulu", "carla@example.com", true)
	seedEmployee(t, store, "001", "Ana", "Alpha", "ana@example.com", true)
	seedEmployee(t, store, "002", "Bob", "Beta", "bob@example.com", false)

	active := true
	got, err := store.List(t.Context(), ListOptions{Limit: 1, Offset: 1, Sort: "name", Query: "a", Active: &active})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 1 {
		t.Fatalf("expected 1 employee, got %d", len(got))
	}
	if got[0].ID != "003" {
		t.Fatalf("expected second active employee sorted by name, got %q", got[0].ID)
	}
}

func seedEmployee(t *testing.T, store *InMemoryStore, id, firstName, lastName, email string, active bool) {
	t.Helper()
	emp, err := NewEmployee(id, firstName, lastName, email)
	if err != nil {
		t.Fatal(err)
	}
	emp.IsActive = active
	if err := store.Save(t.Context(), emp); err != nil {
		t.Fatal(err)
	}
}
