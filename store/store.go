package store

import (
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/appointment"
	"github.com/yageunpro/owl-backend-go/store/auth"
	"github.com/yageunpro/owl-backend-go/store/calendar"
	"github.com/yageunpro/owl-backend-go/store/location"
	"github.com/yageunpro/owl-backend-go/store/user"
)

type Store struct {
	Appointment appointment.Store
	Auth        auth.Store
	Calendar    calendar.Store
	Location    location.Store
	User        user.Store
}

func New(pool *db.Pool) (*Store, error) {
	return &Store{
		Appointment: appointment.New(pool),
		Auth:        auth.New(pool),
		Calendar:    calendar.New(pool),
		Location:    location.New(pool),
		User:        user.New(pool),
	}, nil
}
