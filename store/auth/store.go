package auth

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/internal/query"
)

type Store interface {
	CreateOAuthUser(ctx context.Context, arg CreateOAuthUserParam) error
	UpdateOAuthUser(ctx context.Context, arg UpdateOAuthUserParam) error
	GetOAuthUser(ctx context.Context, openId string) (*resGetOAuthUser, error)
	GetOAuthToken(ctx context.Context, userId uuid.UUID) (*resGetOAuthToken, error)
	CreateDevUser(ctx context.Context, arg CreateDevUserParam) error
	GetDevUser(ctx context.Context, email string) (*resGetDevUser, error)
	GetAllOAuthUserIds(ctx context.Context) ([]uuid.UUID, error)
}

type store struct {
	pool *db.Pool
}

func New(pool *db.Pool) Store {
	return &store{pool: pool}
}

func (s *store) CreateOAuthUser(ctx context.Context, arg CreateOAuthUserParam) error {
	tx, err := s.pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		errors.Join(errors.New("failed to begin transaction"), err)
	}

	qry := query.New(s.pool).WithTx(tx)
	err = qry.CreateUser(ctx, query.CreateUserParams{
		ID:       arg.UserId,
		Email:    arg.Email,
		Username: arg.UserName,
	})
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return ErrUserExisted
		}
		errors.Join(errors.New("failed to create user"), err)
	}

	param := query.CreateOAuthParams{
		ID:           arg.UserId,
		OpenID:       arg.OpenId,
		AccessToken:  arg.AccessToken,
		RefreshToken: pgtype.Text{},
		AllowSync:    arg.AllowSync,
		ValidUntil: pgtype.Timestamptz{
			Time:  arg.ValidUntil.UTC(),
			Valid: true,
		},
	}
	if arg.RefreshToken != nil {
		param.RefreshToken.String = *arg.RefreshToken
		param.RefreshToken.Valid = true
	}
	err = qry.CreateOAuth(ctx, param)
	if err != nil {
		return errors.Join(errors.New("failed to create oauth"), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}

	return nil
}

func (s *store) UpdateOAuthUser(ctx context.Context, arg UpdateOAuthUserParam) error {
	qry := query.New(s.pool)

	param := query.UpdateOAuthTokenParams{
		ID:          arg.UserId,
		UAccess:     pgtype.Text{Valid: false},
		URefresh:    pgtype.Text{Valid: false},
		UValidUntil: pgtype.Timestamptz{Time: arg.ValidUntil.UTC(), Valid: true},
		UAllowSync:  pgtype.Bool{Bool: arg.AllowSync, Valid: true},
	}

	if arg.AccessToken != "" {
		param.UAccess.String = arg.AccessToken
		param.UAccess.Valid = true
	}
	if arg.RefreshToken != nil {
		param.URefresh.String = *arg.RefreshToken
		param.URefresh.Valid = true
	}

	err := qry.UpdateOAuthToken(ctx, param)
	if err != nil {
		return errors.Join(errors.New("failed to update user"), err)
	}
	return nil
}

func (s *store) GetOAuthUser(ctx context.Context, openId string) (*resGetOAuthUser, error) {
	qry := query.New(s.pool)
	row, err := qry.FindOAuth(ctx, openId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Join(errors.New("failed to find user"), err)
	}

	return &resGetOAuthUser{UserId: row.ID}, nil
}

func (s *store) GetOAuthToken(ctx context.Context, userId uuid.UUID) (*resGetOAuthToken, error) {
	qry := query.New(s.pool)
	row, err := qry.GetUserOAuthToken(ctx, userId)
	if err != nil {
		return nil, errors.Join(errors.New("failed to find user"), err)
	}

	return &resGetOAuthToken{
		UserId:       userId,
		OpenId:       row.OpenID,
		AccessToken:  row.AccessToken,
		RefreshToken: row.RefreshToken.String,
		ExpireTime:   row.ValidUntil.Time,
	}, nil
}

func (s *store) CreateDevUser(ctx context.Context, arg CreateDevUserParam) error {
	tx, err := s.pool.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		errors.Join(errors.New("failed to begin transaction"), err)
	}

	qry := query.New(s.pool).WithTx(tx)
	err = qry.CreateUser(ctx, query.CreateUserParams{
		ID:       arg.UserId,
		Email:    arg.Email,
		Username: "USER_" + arg.UserId.String()[len(arg.UserId.String())-4:],
	})
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			return ErrUserExisted
		}
		errors.Join(errors.New("failed to create user"), err)
	}

	err = qry.CreatePassword(ctx, query.CreatePasswordParams{
		ID:           arg.UserId,
		PasswordHash: arg.PasswordHash,
	})
	if err != nil {
		errors.Join(errors.New("failed to create password"), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}

	return nil
}

func (s *store) GetDevUser(ctx context.Context, email string) (*resGetDevUser, error) {
	qry := query.New(s.pool)
	row, err := qry.FindUser(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, errors.Join(errors.New("failed to find user"), err)
	}

	return &resGetDevUser{
		UserId:       row.ID,
		PasswordHash: row.PasswordHash,
	}, nil
}

func (s *store) GetAllOAuthUserIds(ctx context.Context) ([]uuid.UUID, error) {
	qry := query.New(s.pool)
	rows, err := qry.GetAllOAuthUserIds(ctx)
	if err != nil {
		return nil, err
	}

	return rows, nil
}
