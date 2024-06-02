package user

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/internal/query"
)

type Store interface {
	Info(ctx context.Context, userId uuid.UUID) (*resInfo, error)
	GetAccount(ctx context.Context, userId uuid.UUID) error
	AddAccount(ctx context.Context, userId uuid.UUID) error
	DeleteAccount(ctx context.Context, userId, accountId uuid.UUID) error
}

type store struct {
	pool *db.Pool
}

func New(pool *db.Pool) Store {
	return &store{pool: pool}
}

func (s *store) Info(ctx context.Context, userId uuid.UUID) (*resInfo, error) {
	qry := query.New(s.pool)

	row, err := qry.GetUser(ctx, userId)
	if err != nil {
		return nil, errors.Join(errors.New("fail to get user info"), err)
	}

	return &resInfo{
		Id:       row.ID,
		Username: row.Username,
		Email:    row.Email,
	}, nil
}

func (s *store) GetAccount(ctx context.Context, userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *store) AddAccount(ctx context.Context, userId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *store) DeleteAccount(ctx context.Context, userId, accountId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}
