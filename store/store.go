package store

import (
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/auth"
	"github.com/yageunpro/owl-backend-go/store/calendar"
)

type Store struct {
	Auth     auth.Store
	Calendar calendar.Store
}

func New(pool *db.Pool) (*Store, error) {
	return &Store{
		Auth:     auth.New(pool),
		Calendar: calendar.New(pool),
	}, nil
}
