package handler

import (
	"context"
	"net/http"

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

func (h *Handler) GenerateTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := r.URL.Query().Get("GUID")

	}
}

func (h *Handler) RefreshTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
