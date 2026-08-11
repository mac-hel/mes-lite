package auth

import (
	"context"
	"fmt"
	"strings"
)

// InMemoryStore stores users in memory for development and fast handler tests.
type InMemoryStore struct {
	byEmail map[string]User
}

// NewInMemoryStore creates an empty in-memory user store.
func NewInMemoryStore() *InMemoryStore {
	return &InMemoryStore{byEmail: make(map[string]User)}
}

// Save stores a user by normalized email.
func (s *InMemoryStore) Save(ctx context.Context, user User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	s.byEmail[normalizeEmail(user.Email)] = user
	return nil
}

// FindByEmail returns a user by normalized email.
func (s *InMemoryStore) FindByEmail(ctx context.Context, email string) (User, error) {
	user, ok := s.byEmail[normalizeEmail(email)]
	if !ok {
		return User{}, fmt.Errorf("email %q: %w", email, ErrNotFound)
	}
	return user, nil
}

func normalizeEmail(email string) string {
	return strings.TrimSpace(strings.ToLower(email))
}
