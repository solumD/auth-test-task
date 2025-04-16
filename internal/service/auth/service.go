package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/solumD/auth-test-task/internal/client/db"
	"github.com/solumD/auth-test-task/internal/logger"
	"github.com/solumD/auth-test-task/internal/model"
	"github.com/solumD/auth-test-task/internal/repository"
	"github.com/solumD/auth-test-task/internal/service"
	"github.com/solumD/auth-test-task/internal/utils/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type srv struct {
	authRepository repository.AuthRepository
	txManager      db.TxManager
}

var (
	jwtKey      = os.Getenv("JWT_KEY")
	jwtDuration = time.Minute * 30
)

var (
	ErrJwtGenerationFailure     = errors.New("failed to generate access token")
	ErrJwtVerificationFailure   = errors.New("failed to verify access token")
	ErrIPsNotMatch              = errors.New("old and curr user's ip do not match")
	ErrAccessTokensUIDsNotMatch = errors.New("old and curr access tokens's uid do not match")
)

// New returns new auth service object
func New(authRepository repository.AuthRepository, txManager db.TxManager) service.AuthService {
	return &srv{
		authRepository: authRepository,
		txManager:      txManager,
	}
}

func (s *srv) GenerateTokens(ctx context.Context, guid string, userIP string) (*model.Tokens, error) {
	accessTokenUID := uuid.NewString()

	info := &jwt.Info{
		UserIP:         userIP,
		UserGUID:       guid,
		AccessTokenUID: accessTokenUID,
	}

	accessToken, err := jwt.GenerateToken(info, []byte(jwtKey), jwtDuration)
	if err != nil {
		logger.Error(ErrJwtGenerationFailure.Error(), zap.Error(err))
		return nil, ErrJwtGenerationFailure
	}

	refreshToken := base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))

	err = s.authRepository.SaveTokensInfo(ctx, refreshToken, accessTokenUID)
	if err != nil {
		logger.Error("failed to save tokens info", zap.Error(err))
		return nil, err
	}

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *srv) RefreshTokens(ctx context.Context, tokens *model.Tokens, userIP string) (*model.Tokens, error) {
	claims, err := jwt.VerifyToken(tokens.AccessToken, []byte(jwtKey))
	if err != nil {
		logger.Error(ErrJwtVerificationFailure.Error(), zap.Error(err))
		return nil, ErrJwtVerificationFailure
	}

	if claims.UserIP != userIP {
		logger.Error(ErrIPsNotMatch.Error())
		return nil, ErrIPsNotMatch
	}

	oldAccessTokenUID, err := s.authRepository.GetAcccessTokenUID(ctx, tokens.RefreshToken)
	if err != nil {
		logger.Error("failed to get old access token uid", zap.Error(err))
		return nil, err
	}

	if claims.AccessTokenUID != oldAccessTokenUID {
		logger.Error(ErrAccessTokensUIDsNotMatch.Error())
		return nil, ErrAccessTokensUIDsNotMatch
	}

	err = s.authRepository.SetRefreshTokenUsed(ctx, tokens.RefreshToken)
	if err != nil {
		logger.Error("failed to update tokens info")
		return nil, err
	}

	return s.GenerateTokens(ctx, claims.UserGUID, userIP)
}
