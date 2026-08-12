package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const defaultAccessTokenTTL = 15 * time.Minute

// ErrInvalidToken means a JWT cannot be trusted or used for authentication.
var ErrInvalidToken = errors.New("invalid token")

// TokenManager signs and verifies JWT access tokens.
type TokenManager struct {
	secret []byte
	ttl    time.Duration
}

// NewTokenManager creates a JWT token manager using an HMAC secret.
func NewTokenManager(secret string) (*TokenManager, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("jwt secret must be at least 32 characters: %w", ErrInvalidToken)
	}
	return &TokenManager{secret: []byte(secret), ttl: defaultAccessTokenTTL}, nil
}

type accessClaims struct {
	Email string `json:"email"`
	Role  Role   `json:"role"`
	jwt.RegisteredClaims
}

// Issue creates a signed JWT for the given user.
func (m *TokenManager) Issue(user User) (string, error) {
	now := time.Now().UTC()
	claims := accessClaims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(m.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}
	return signed, nil
}

// Verify validates a JWT and returns its authenticated principal.
func (m *TokenManager) Verify(raw string) (Principal, error) {
	claims := accessClaims{}
	token, err := jwt.ParseWithClaims(raw, &claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected jwt signing method: %w", ErrInvalidToken)
		}
		return m.secret, nil
	})
	if err != nil {
		return Principal{}, fmt.Errorf("parse jwt: %w", ErrInvalidToken)
	}
	if !token.Valid {
		return Principal{}, ErrInvalidToken
	}
	if claims.Subject == "" || claims.Email == "" || !claims.Role.Valid() {
		return Principal{}, fmt.Errorf("missing required jwt claims: %w", ErrInvalidToken)
	}

	return Principal{UserID: claims.Subject, Email: claims.Email, Role: claims.Role}, nil
}

// Principal is the authenticated identity extracted from an access token.
type Principal struct {
	UserID string
	Email  string
	Role   Role
}
