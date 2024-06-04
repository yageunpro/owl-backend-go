package service

import (
	"github.com/yageunpro/owl-backend-go/service/appointment"
	"github.com/yageunpro/owl-backend-go/service/auth"
	"github.com/yageunpro/owl-backend-go/service/calendar"
	"github.com/yageunpro/owl-backend-go/service/location"
	"github.com/yageunpro/owl-backend-go/service/user"
	"github.com/yageunpro/owl-backend-go/store"
)

type Service struct {
	Appointment appointment.Service
	Auth        auth.Service
	Calendar    calendar.Service
	Location    location.Service
	User        user.Service
}

func New(sto *store.Store) (*Service, error) {
	return &Service{
		Appointment: appointment.New(sto),
		Auth:        auth.New(sto),
		Calendar:    calendar.New(sto),
		Location:    location.New(sto),
		User:        user.New(sto),
	}, nil
}
