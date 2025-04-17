package repository

import "context"

// AuthRepository interface of auth repository
type AuthRepository interface {
	SaveTokensInfo(ctx context.Context, refreshToken string, accessTokenUID string) error
	GetRefreshTokenByAccessUID(ctx context.Context, accessTokenUID string) (string, error)
	SetRefreshTokenUsed(ctx context.Context, accessTokenUID string) error
}
