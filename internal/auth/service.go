package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrInvalidCredentials means the submitted email/password pair is not valid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInactiveUser means credentials were valid but the user is disabled.
	ErrInactiveUser = errors.New("inactive user")
)

// Service coordinates authentication use cases.
type Service struct {
	store Store
}

// NewService creates an authentication service.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// Login verifies credentials and returns an opaque access token for the caller.
func (s *Service) Login(ctx context.Context, email, password string) (LoginResult, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(password) == "" {
		return LoginResult{}, ErrInvalidCredentials
	}

	user, err := s.store.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return LoginResult{}, ErrInvalidCredentials
		}
		return LoginResult{}, err
	}

	if !user.IsActive {
		return LoginResult{}, ErrInactiveUser
	}
	if !user.VerifyPassword(password) {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := newAccessToken()
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{AccessToken: token, User: user}, nil
}

func newAccessToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// LoginResult is returned after successful credential verification.
type LoginResult struct {
	AccessToken string `json:"accessToken"`
	User        User   `json:"user"`
}
