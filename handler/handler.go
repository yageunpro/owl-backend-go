package handler

import (
	"github.com/yageunpro/owl-backend-go/handler/appointment"
	"github.com/yageunpro/owl-backend-go/handler/auth"
	"github.com/yageunpro/owl-backend-go/handler/calendar"
	"github.com/yageunpro/owl-backend-go/handler/location"
	"github.com/yageunpro/owl-backend-go/handler/user"
	"github.com/yageunpro/owl-backend-go/service"
)

type Handler struct {
	Auth        auth.Handler
	Appointment appointment.Handler
	Calendar    calendar.Handler
	Location    location.Handler
	User        user.Handler
}

func New(svc *service.Service) (*Handler, error) {
	return &Handler{
		Auth:        auth.New(svc),
		Appointment: appointment.New(svc),
		Calendar:    calendar.New(svc),
		Location:    location.New(svc),
		User:        user.New(svc),
	}, nil
}
