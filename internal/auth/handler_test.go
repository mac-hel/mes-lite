package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-fuego/fuego"
)

func TestLoginHandler(t *testing.T) {
	store := NewInMemoryStore()
	user, err := NewUser("user-1", "admin@example.com", "secret", RoleAdmin)
	if err != nil {
		t.Fatal(err)
	}
	if err := store.Save(t.Context(), user); err != nil {
		t.Fatal(err)
	}
	tokens := testTokenManager(t)

	s := fuego.NewServer()
	fuego.Post(s, "/auth/login", NewHandler(NewService(store, tokens)).Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(`{"email":"admin@example.com","password":"secret"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var result LoginResult
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
}

func TestLoginHandlerRejectsInvalidCredentials(t *testing.T) {
	store := NewInMemoryStore()
	s := fuego.NewServer()
	fuego.Post(s, "/auth/login", NewHandler(NewService(store, testTokenManager(t))).Login)

	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(`{"email":"admin@example.com","password":"wrong"}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	s.Mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", w.Code, w.Body.String())
	}
}
