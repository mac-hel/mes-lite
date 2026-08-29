package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/mac-hel/mes-lite/internal/auth/authdb"
	"github.com/mac-hel/mes-lite/internal/platform/postgres"
)

// NewPostgresStore creates a PostgreSQL-backed auth user store.
func NewPostgresStore(db authdb.DBTX) *PostgresStore {
	return &PostgresStore{queries: authdb.New(db)}
}

// PostgresStore stores auth users in PostgreSQL through sqlc-generated queries.
type PostgresStore struct {
	queries *authdb.Queries
}

// Save stores an auth user keyed by ID and email.
func (s *PostgresStore) Save(ctx context.Context, user User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	_, err := s.queries.CreateUser(ctx, authdb.CreateUserParams{
		ID:           user.ID,
		Email:        normalizeEmail(user.Email),
		PasswordHash: user.PasswordHash,
		Role:         string(user.Role),
		IsActive:     user.IsActive,
	})
	if err != nil {
		return mapPostgresError(user.Email, err)
	}
	return nil
}

// FindByEmail returns an auth user by normalized email.
func (s *PostgresStore) FindByEmail(ctx context.Context, email string) (User, error) {
	user, err := s.queries.GetUserByEmail(ctx, normalizeEmail(email))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, fmt.Errorf("email %q: %w", email, ErrNotFound)
		}
		return User{}, err
	}
	return userFromDB(user)
}

// EnsureBootstrapAdmin creates the configured bootstrap admin only when it does not exist.
func (s *PostgresStore) EnsureBootstrapAdmin(ctx context.Context, user User) error {
	if err := user.Validate(); err != nil {
		return err
	}
	exists, err := s.queries.UserExistsByEmail(ctx, normalizeEmail(user.Email))
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	return s.Save(ctx, user)
}

func userFromDB(user authdb.AuthUser) (User, error) {
	role := Role(user.Role)
	domain := User{
		ID:           user.ID,
		Email:        normalizeEmail(user.Email),
		PasswordHash: user.PasswordHash,
		Role:         role,
		IsActive:     user.IsActive,
	}
	if err := domain.Validate(); err != nil {
		return User{}, err
	}
	return domain, nil
}

func mapPostgresError(email string, err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch postgres.SQLState(pgErr.Code) {
		case postgres.UniqueViolation:
			return fmt.Errorf("auth user %q: %w", email, ErrAlreadyExists)
		case postgres.CheckViolation, postgres.NotNullViolation:
			return fmt.Errorf("auth user %q: %w", email, ErrInvalidUser)
		}
	}
	return err
}
