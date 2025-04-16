package service

import "context"

// AuthService interface of Auth service
type AuthService interface {
	GenerateTokens(ctx context.Context, GUID string) (string, string, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, error)
}
