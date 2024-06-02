package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/yageunpro/owl-backend-go/store"
)

type Service interface {
	Info(ctx context.Context, userId uuid.UUID) (*resInfo, error)
}

type service struct {
	store *store.Store
}

func New(sto *store.Store) Service {
	return &service{store: sto}
}

func (s *service) Info(ctx context.Context, userId uuid.UUID) (*resInfo, error) {
	out, err := s.store.User.Info(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &resInfo{
		Id:       out.Id,
		Username: out.Username,
		Email:    out.Email,
	}, nil
}
