package auth

import (
	"context"
	"errors"
	"log"

	"github.com/solumD/auth-test-task/internal/client/db"
	"github.com/solumD/auth-test-task/internal/repository"
	"github.com/solumD/auth-test-task/internal/utils/hash"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
)

const (
	tableName         = "auth"
	idCol             = "id"
	refreshTokenCol   = "refresh_token"
	accessTokenUIDCol = "access_token_uid"
	isUsedCol         = "is_used"
	createdAtCol      = "created_at"
)

var (
	// ErrInvalidQuery ...
	ErrInvalidQuery = errors.New("failed to create query")

	// ErrExecFailure ...
	ErrExecFailure = errors.New("failed to exec query")

	// ErrRefreshTokenNotExist ...
	ErrRefreshTokenNotExist = errors.New("refresh token not exists or was already used")

	// ErrEncryptingFailure ...
	ErrEncryptingFailure = errors.New("failed to encrypt a string")
)

type repo struct {
	db db.Client
}

// New returns new auth repository object
func New(db db.Client) repository.AuthRepository {
	return &repo{
		db: db,
	}
}

// SaveTokensInfo saves info about refresh and access tokens in storage
func (r *repo) SaveTokensInfo(ctx context.Context, refreshToken string, accessTokenUID string) error {
	hashedRefreshToken, err := hash.Encrypt(refreshToken)
	if err != nil {
		return ErrEncryptingFailure
	}

	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(refreshTokenCol, accessTokenUIDCol).
		Values(hashedRefreshToken, accessTokenUID).ToSql()

	if err != nil {
		return ErrInvalidQuery
	}

	q := db.Query{
		Name:     "auth_repository.SaveTokensInfo",
		QueryRaw: query,
	}

	if _, err = r.db.DB().ExecContext(ctx, q, args...); err != nil {
		return ErrExecFailure
	}

	return nil
}

// GetRefreshTokenByAccessUID gets refresh token's hash by access token uid if exist
func (r *repo) GetRefreshTokenByAccessUID(ctx context.Context, accessTokenUID string) (string, error) {
	query, args, err := sq.Select(refreshTokenCol).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.And{sq.Eq{accessTokenUIDCol: accessTokenUID}, sq.Eq{isUsedCol: 0}}).
		ToSql()

	log.Println("1231")

	if err != nil {
		return "", ErrInvalidQuery
	}

	q := db.Query{
		Name:     "auth_repository.GetRefreshTokenByAccessUID",
		QueryRaw: query,
	}

	var refreshTokenHash string
	err = r.db.DB().ScanOneContext(ctx, &refreshTokenHash, q, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrRefreshTokenNotExist
		}

		return "", ErrExecFailure
	}

	return refreshTokenHash, nil
}

// SetRefreshTokenUsed sets refresh token used by accessTokenUID
func (r *repo) SetRefreshTokenUsed(ctx context.Context, accessTokenUID string) error {
	query, args, err := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(isUsedCol, 1).
		Where(sq.Eq{accessTokenUIDCol: accessTokenUID}).
		ToSql()

	if err != nil {
		return ErrInvalidQuery
	}

	q := db.Query{
		Name:     "auth_repository.SetRefreshTokenUsed",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return ErrExecFailure
	}

	return nil
}
