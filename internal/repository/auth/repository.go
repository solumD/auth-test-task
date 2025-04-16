package somereponame

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"
	"github.com/solumD/auth-test-task/internal/client/db"
	"github.com/solumD/auth-test-task/internal/repository"

	sq "github.com/Masterminds/squirrel"
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
	ErrInvalidQuery         = errors.New("failed to create query")
	ErrExecFailure          = errors.New("failed to exec query")
	ErrRefreshTokenNotExist = errors.New("refresh token not exists or was already used")
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

func (r *repo) SaveTokensInfo(ctx context.Context, refreshToken string, accessTokenUID string) error {
	query, args, err := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(refreshTokenCol, accessTokenUIDCol).
		Values(refreshToken, accessTokenUID).ToSql()

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

func (r *repo) GetAcccessTokenUID(ctx context.Context, refreshToken string) (string, error) {
	query, args, err := sq.Select(accessTokenUIDCol).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.And{sq.Eq{refreshTokenCol: refreshToken}, sq.Eq{isUsedCol: 0}}).
		ToSql()

	if err != nil {
		return "", ErrInvalidQuery
	}

	q := db.Query{
		Name:     "auth_repository.GetAccessTokenUID",
		QueryRaw: query,
	}

	var accessTokenUID string
	err = r.db.DB().ScanOneContext(ctx, &accessTokenUID, q, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrRefreshTokenNotExist
		}

		return "", ErrExecFailure
	}

	return accessTokenUID, nil
}

func (r *repo) UpdateTokensInfo(ctx context.Context, refreshToken string) error {
	query, args, err := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(isUsedCol, 1).
		Where(sq.Eq{refreshTokenCol: refreshToken}).
		ToSql()

	if err != nil {
		return ErrInvalidQuery
	}

	q := db.Query{
		Name:     "auth_repository.UpdateTokensInfo",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return ErrExecFailure
	}

	return nil
}
