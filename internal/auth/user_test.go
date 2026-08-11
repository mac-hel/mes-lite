package auth

import "testing"

func TestNewUserHashesPassword(t *testing.T) {
	user, err := NewUser("user-1", "ADMIN@EXAMPLE.COM", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}

	if user.Email != "admin@example.com" {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if string(user.PasswordHash) == "secret" {
		t.Fatal("password was stored as plaintext")
	}
	if !user.VerifyPassword("secret") {
		t.Fatal("expected password verification to succeed")
	}
	if user.VerifyPassword("wrong") {
		t.Fatal("expected wrong password verification to fail")
	}
}

func TestNewUserRejectsInvalidRole(t *testing.T) {
	_, err := NewUser("user-1", "admin@example.com", "secret", Role("owner"))
	if err == nil {
		t.Fatal("expected invalid role error")
	}
}
