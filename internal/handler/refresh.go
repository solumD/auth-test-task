package handler

import (
	"context"
	"net/http"

	"github.com/solumD/auth-test-task/internal/logger"
	"github.com/solumD/auth-test-task/internal/model"
	"github.com/solumD/auth-test-task/internal/utils/ip"

	"github.com/go-chi/render"
	"go.uber.org/zap"
)

type refreshTokensRequest struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type refreshTokensResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ErrorMsg     string `json:"error_message,omitempty"`
}

// RefreshTokens get's access and refresh tokens from body, validates them
// and returns new refreshed pair of access and refresh tokens
func (h *Handler) RefreshTokens(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userIP, err := ip.GetIP(r)
		if err != nil {
			logger.Error(ErrUserIPFailure.Error(), zap.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, refreshTokensResponse{
				ErrorMsg: ErrUserIPFailure.Error(),
			})
			return
		}

		var req refreshTokensRequest
		err = render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("failed to decode refresh tokens request", zap.Error(err))

			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, refreshTokensResponse{
				ErrorMsg: "failed to decode refresh tokens request",
			})
			return
		}

		logger.Info("RefreshTokens request recieved", zap.Any("request", req))

		t := &model.Tokens{
			AccessToken:  req.AccessToken,
			RefreshToken: req.RefreshToken,
		}

		tokens, err := h.authService.RefreshTokens(ctx, t, userIP)
		if err != nil {
			logger.Error(err.Error())

			render.Status(r, http.StatusUnauthorized)
			render.JSON(w, r, refreshTokensResponse{
				ErrorMsg: err.Error(),
			})
			return
		}

		logger.Info("refreshed tokens for user")
		render.Status(r, http.StatusOK)
		render.JSON(w, r, refreshTokensResponse{
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
		})
	}
}
