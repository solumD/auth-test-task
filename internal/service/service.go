package service

import (
	"context"

	"github.com/solumD/auth-test-task/internal/model"
)

// AuthService interface of auth service
type AuthService interface {
	GenerateTokens(ctx context.Context, guid string, userIP string) (*model.Tokens, error)
	RefreshTokens(ctx context.Context, tokens *model.Tokens, userIP string) (*model.Tokens, error)
}

// EmailService interface of email service
type EmailService interface {
	SendEmail(from string, to string, theme string, text string)
}
