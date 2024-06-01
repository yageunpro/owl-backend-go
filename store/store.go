package store

import (
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/auth"
)

type Store struct {
	Auth auth.Store
}

func New(pool *db.Pool) (*Store, error) {
	return &Store{Auth: auth.New(pool)}, nil
}
