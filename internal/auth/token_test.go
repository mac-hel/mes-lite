package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestTokenManagerIssueAndVerify(t *testing.T) {
	manager := testTokenManager(t)
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	token, err := manager.Issue(user)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(token, ".") != 2 {
		t.Fatalf("expected compact jwt, got %q", token)
	}

	principal, err := manager.Verify(token)
	if err != nil {
		t.Fatal(err)
	}
	if principal.UserID != user.ID || principal.Email != user.Email || principal.Role != user.Role {
		t.Fatalf("unexpected principal: %+v", principal)
	}
}

func TestTokenManagerRejectsWrongSecret(t *testing.T) {
	manager := testTokenManager(t)
	other, err := NewTokenManager("different-test-secret-with-32-characters")
	if err != nil {
		t.Fatal(err)
	}
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	token, err := manager.Issue(user)
	if err != nil {
		t.Fatal(err)
	}

	_, err = other.Verify(token)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token, got %v", err)
	}
}

func TestTokenManagerRejectsShortSecret(t *testing.T) {
	_, err := NewTokenManager("short")
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected invalid token, got %v", err)
	}
}

func testTokenManager(t *testing.T) *TokenManager {
	t.Helper()
	manager, err := NewTokenManager("test-secret-with-at-least-32-characters")
	if err != nil {
		t.Fatal(err)
	}
	return manager
}
