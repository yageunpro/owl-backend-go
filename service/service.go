package service

import (
	"github.com/yageunpro/owl-backend-go/service/auth"
	"github.com/yageunpro/owl-backend-go/store"
)

type Service struct {
	Auth auth.Service
}

func New(sto *store.Store) (*Service, error) {
	return &Service{Auth: auth.New(sto)}, nil
}
