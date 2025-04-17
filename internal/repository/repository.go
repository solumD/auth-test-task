package repository

import "context"

// AuthRepository interface of auth repository
type AuthRepository interface {
	SaveTokensInfo(ctx context.Context, refreshToken string, accessTokenUID string) error
	GetAccessTokenUID(ctx context.Context, refreshToken string) (string, error)
	SetRefreshTokenUsed(ctx context.Context, refreshToken string) error
}
