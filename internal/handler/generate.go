package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/render"
	"github.com/solumD/auth-test-task/internal/logger"
	"github.com/solumD/auth-test-task/internal/utils/ip"
	"go.uber.org/zap"
)

type GenerateTokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ErrorMsg     string `json:"error_message,omitempty"`
}

var (
	ErrUserIPFailure = errors.New("failed to get user ip")
)

func (h *Handler) GenerateTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guid := r.URL.Query().Get("GUID")
		userIP, err := ip.GetIP(r)
		if err != nil {
			logger.Error(ErrUserIPFailure.Error(), zap.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, GenerateTokenResponse{
				ErrorMsg: ErrUserIPFailure.Error(),
			})
			return
		}

		tokens, err := h.authService.GenerateTokens(ctx, guid, userIP)
		if err != nil {
			logger.Error(err.Error())

			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, GenerateTokenResponse{
				ErrorMsg: err.Error(),
			})
			return
		}

		logger.Info("generated tokens for user", zap.String("guid", guid), zap.String("ip", userIP))
		render.Status(r, http.StatusOK)
		render.JSON(w, r, GenerateTokenResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
