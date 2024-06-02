package service

import (
	"github.com/yageunpro/owl-backend-go/service/auth"
	"github.com/yageunpro/owl-backend-go/service/calendar"
	"github.com/yageunpro/owl-backend-go/service/user"
	"github.com/yageunpro/owl-backend-go/store"
)

type Service struct {
	Auth     auth.Service
	Calendar calendar.Service
	User     user.Service
}

func New(sto *store.Store) (*Service, error) {
	return &Service{
		Auth:     auth.New(sto),
		Calendar: calendar.New(sto),
		User:     user.New(sto),
	}, nil
}
