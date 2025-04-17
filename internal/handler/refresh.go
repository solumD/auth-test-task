package handler

import (
	"context"
	"net/http"
)

func (h *Handler) RefreshTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
