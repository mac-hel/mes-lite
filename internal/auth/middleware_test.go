package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareAuthenticate(t *testing.T) {
	tokens := testTokenManager(t)
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	token, err := tokens.Issue(user)
	if err != nil {
		t.Fatal(err)
	}

	called := false
	handler := NewMiddleware(tokens).Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		principal, ok := PrincipalFromContext(r.Context())
		if !ok {
			t.Fatal("expected principal in context")
		}
		if principal.UserID != user.ID || principal.Role != RoleAdmin {
			t.Fatalf("unexpected principal: %+v", principal)
		}
		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if !called {
		t.Fatal("expected next handler to be called")
	}
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

func TestMiddlewareAuthenticateRejectsMissingToken(t *testing.T) {
	handler := NewMiddleware(testTokenManager(t)).Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestMiddlewareAuthenticateRejectsInvalidToken(t *testing.T) {
	handler := NewMiddleware(testTokenManager(t)).Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	}))

	req := httptest.NewRequest(http.MethodGet, "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-jwt")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}
