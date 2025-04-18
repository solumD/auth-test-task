package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/solumD/auth-test-task/internal/logger"
	"github.com/solumD/auth-test-task/internal/model"
	"github.com/solumD/auth-test-task/internal/repository"
	"github.com/solumD/auth-test-task/internal/service"
	"github.com/solumD/auth-test-task/internal/utils/hash"
	"github.com/solumD/auth-test-task/internal/utils/jwt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type srv struct {
	authRepository repository.AuthRepository
	emailService   service.EmailService
}

var (
	jwtKey      = os.Getenv("JWT_KEY")
	jwtDuration = time.Minute * 30
)

var (
	// ErrJwtGenerationFailure ...
	ErrJwtGenerationFailure = errors.New("failed to generate access token")

	// ErrInvalidJwtToken  ...
	ErrInvalidJwtToken = errors.New("invalid access token")

	// ErrIPsNotMatch ...
	ErrIPsNotMatch = errors.New("old and curr user's ip do not match")

	// ErrRefreshTokensNotMatch ...
	ErrRefreshTokensNotMatch = errors.New("old and curr refresh tokens's do not match")

	// ErrTokensIsNil ...
	ErrTokensIsNil = errors.New(`"tokens" is nil`)
)

// New returns new auth service object
func New(authRepository repository.AuthRepository, emailService service.EmailService) service.AuthService {
	return &srv{
		authRepository: authRepository,
		emailService:   emailService,
	}
}

// GenerateTokens validates guid and ip, generates access and refresh token and makes a request
// in repository to save info
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

	logger.Info("saved tokens info in repo")

	return &model.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// RefreshTokens validates access and refresh tokens, makes a request in repository to get info
// and generates new pair of them
func (s *srv) RefreshTokens(ctx context.Context, tokens *model.Tokens, userIP string) (*model.Tokens, error) {
	if tokens == nil {
		return nil, ErrTokensIsNil
	}

	claims, err := jwt.VerifyToken(tokens.AccessToken, []byte(jwtKey))
	if err != nil {
		logger.Error(ErrInvalidJwtToken.Error(), zap.Error(err))
		return nil, ErrInvalidJwtToken
	}

	if claims.UserIP != userIP {
		logger.Error(ErrIPsNotMatch.Error())

		s.emailService.SendEmail(
			"medods@email.ru",
			"someUser@email.ru",
			"ip change warning",
			"Dear user, we noticed that your current IP address is different from the main one. Check the security of your account in profile page.",
		)

		return nil, ErrIPsNotMatch
	}

	refreshTokenHash, err := s.authRepository.GetRefreshTokenByAccessUID(ctx, claims.AccessTokenUID)
	if err != nil {
		logger.Error("failed to get refresh token hash", zap.Error(err))
		return nil, err
	}

	logger.Info("got refresh token hash from repo")

	if err := hash.CompareHashAndRaw(tokens.RefreshToken, refreshTokenHash); err != nil {
		return nil, ErrRefreshTokensNotMatch
	}

	err = s.authRepository.SetRefreshTokenUsed(ctx, claims.AccessTokenUID)
	if err != nil {
		logger.Error("failed to update tokens info")
		return nil, err
	}

	logger.Info("set refresh token used in repo")

	return s.GenerateTokens(ctx, claims.UserGUID, userIP)
}
