package handler

import (
	"github.com/solumD/auth-test-task/internal/service"
)

type Handler struct {
	authService service.AuthService
}

// New возвращает объект Handler
func New(authService service.AuthService) *Handler {
	return &Handler{
		authService: authService,
	}
}
