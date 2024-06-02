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
