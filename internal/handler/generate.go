package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/solumD/auth-test-task/internal/logger"
	"github.com/solumD/auth-test-task/internal/utils/ip"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type generateTokensResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ErrorMsg     string `json:"error_message,omitempty"`
}

var (
	// ErrUserIPFailure ...
	ErrUserIPFailure = errors.New("failed to get user ip")
)

// GenerateTokens get's user's guid from params, generates access and
// refresh tokens and returns them in json
func (h *Handler) GenerateTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := r.URL.Query().Get("guid")
		userIP, err := ip.GetIP(r)
		if err != nil {
			logger.Error(ErrUserIPFailure.Error(), zap.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, generateTokensResponse{
				ErrorMsg: ErrUserIPFailure.Error(),
			})
			return
		}

		tokens, err := h.authService.GenerateTokens(ctx, guid, userIP)
		if err != nil {
			logger.Error(err.Error())

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, generateTokensResponse{
				ErrorMsg: err.Error(),
			})
			return
		}

		logger.Info("generated tokens for user", zap.String("guid", guid), zap.String("ip", userIP))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, generateTokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
