package auth

import (
	"errors"
	"testing"
)

func TestServiceLogin(t *testing.T) {
	store := NewInMemoryStore()
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}

	result, err := NewService(store).Login(t.Context(), "ADMIN@example.com", "secret")
	if err != nil {
		t.Fatal(err)
	}

	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
	if result.User.Email != "admin@example.com" {
		t.Fatalf("expected logged-in user, got %q", result.User.Email)
	}
}

func TestServiceLoginRejectsWrongPassword(t *testing.T) {
	store := NewInMemoryStore()
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}

	_, err = NewService(store).Login(t.Context(), "admin@example.com", "wrong")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected invalid credentials, got %v", err)
	}
}

func TestServiceLoginRejectsInactiveUser(t *testing.T) {
	store := NewInMemoryStore()
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	user.IsActive = false
	if err := store.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}

	_, err = NewService(store).Login(t.Context(), "admin@example.com", "secret")
	if !errors.Is(err, ErrInactiveUser) {
		t.Fatalf("expected inactive user, got %v", err)
	}
}
