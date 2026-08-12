package auth

import (
	"context"
	"errors"
)

// ErrNotFound means no auth user matched the lookup.
var ErrNotFound = errors.New("user not found")

// Store defines the user lookup behavior needed by authentication.
type Store interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	Save(ctx context.Context, user User) error
}

// BootstrapStore defines startup behavior for creating the initial admin user.
type BootstrapStore interface {
	Store
	EnsureBootstrapAdmin(ctx context.Context, user User) error
}
