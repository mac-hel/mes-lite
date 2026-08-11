package auth

import (
	"errors"

	"github.com/go-fuego/fuego"
)

// Handler holds HTTP handlers for authentication endpoints.
type Handler struct {
	service *Service
}

// NewHandler creates a new authentication handler.
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// LoginRequest is the expected JSON body for logging in.
type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Login handles POST /auth/login.
func (h *Handler) Login(c fuego.ContextWithBody[LoginRequest]) (LoginResult, error) {
	body, err := c.Body()
	if err != nil {
		return LoginResult{}, err
	}

	result, err := h.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) || errors.Is(err, ErrInactiveUser) {
			return LoginResult{}, fuego.UnauthorizedError{
				Err:    err,
				Detail: "invalid email or password",
			}
		}
		return LoginResult{}, err
	}

	return result, nil
}
