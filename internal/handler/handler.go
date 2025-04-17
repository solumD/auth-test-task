package handler

import (
	"github.com/solumD/auth-test-task/internal/service"
)

// Handler contains handlers for all routes
type Handler struct {
	authService service.AuthService
}

// New returns new Handler object
func New(authService service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}
