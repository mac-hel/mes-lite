package auth

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Role describes what a user is allowed to do in the system.
type Role string

const (
	// RoleAdmin can administer MES Lite.
	RoleAdmin Role = "admin"
	// RoleManager can manage production and reports.
	RoleManager Role = "manager"
	// RoleLeader can monitor and correct production work.
	RoleLeader Role = "leader"
	// RoleWorker can register production work.
	RoleWorker Role = "worker"
)

// ErrInvalidUser means a user violates an auth domain invariant.
var ErrInvalidUser = errors.New("invalid user")

// ErrAlreadyExists means an auth user already exists.
var ErrAlreadyExists = errors.New("user already exists")

// User is an authenticated person who can access MES Lite.
type User struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"-"`
	Role         Role   `json:"role"`
	IsActive     bool   `json:"isActive"`
}

// NewUser creates a valid user and stores only a password hash, never the raw password.
func NewUser(id, email, password string, role Role) (User, error) {
	if strings.TrimSpace(password) == "" {
		return User{}, fmt.Errorf("password is required: %w", ErrInvalidUser)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hash password: %w", err)
	}

	user := User{
		ID:           strings.TrimSpace(id),
		Email:        strings.TrimSpace(strings.ToLower(email)),
		PasswordHash: hash,
		Role:         role,
		IsActive:     true,
	}
	if err := user.Validate(); err != nil {
		return User{}, err
	}

	return user, nil
}

// Validate checks user invariants that must hold regardless of transport or storage.
func (u User) Validate() error {
	if strings.TrimSpace(u.ID) == "" {
		return fmt.Errorf("id is required: %w", ErrInvalidUser)
	}
	if strings.TrimSpace(u.Email) == "" {
		return fmt.Errorf("email is required: %w", ErrInvalidUser)
	}
	if len(u.PasswordHash) == 0 {
		return fmt.Errorf("password hash is required: %w", ErrInvalidUser)
	}
	if !u.Role.Valid() {
		return fmt.Errorf("role %q is invalid: %w", u.Role, ErrInvalidUser)
	}
	return nil
}

// VerifyPassword compares a raw login password with the stored password hash.
func (u User) VerifyPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password)) == nil
}

// Valid reports whether the role is one MES Lite understands.
func (r Role) Valid() bool {
	switch r {
	case RoleAdmin, RoleManager, RoleLeader, RoleWorker:
		return true
	default:
		return false
	}
}
