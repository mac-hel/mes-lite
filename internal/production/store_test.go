package production

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestInMemoryStore_SaveAndFindByID(t *testing.T) {
	store := NewInMemoryStore()
	entry, err := NewEntry("entry-1", "emp-1", "sku-1", 12, "ws-1", time.Now(), "")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	got, err := store.FindByID(context.Background(), entry.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got != entry {
		t.Errorf("expected %#v, got %#v", entry, got)
	}
}

func TestInMemoryStore_SaveDuplicate(t *testing.T) {
	store := NewInMemoryStore()
	entry, err := NewEntry("entry-1", "emp-1", "sku-1", 12, "ws-1", time.Now(), "")
	if err != nil {
		t.Fatal(err)
	}

	if err := store.Save(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	err = store.Save(context.Background(), entry)
	if !errors.Is(err, ErrAlreadyExists) {
		t.Fatalf("expected ErrAlreadyExists, got %v", err)
	}
}

func TestInMemoryStore_FindByID_NotFound(t *testing.T) {
	store := NewInMemoryStore()

	_, err := store.FindByID(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
