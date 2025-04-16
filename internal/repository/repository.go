package repository

import "context"

// AuthRepository interface of auth repository
type AuthRepository interface {
	SaveTokensInfo(ctx context.Context, refreshToken string, accessTokenUID string) error
	GetAcccessTokenUID(ctx context.Context, refreshToken string) (string, error)
	UpdateTokensInfo(ctx context.Context, refreshToken string) error
}
