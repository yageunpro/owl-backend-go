package calendar

import (
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/service"
)

type Handler interface {
	ScheduleAdd(c echo.Context) error
	ScheduleInfo(c echo.Context) error
	ScheduleDelete(c echo.Context) error
	ScheduleList(c echo.Context) error
	Sync(c echo.Context) error
}

type handler struct {
	svc *service.Service
}

func New(svc *service.Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) ScheduleAdd(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) ScheduleInfo(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) ScheduleDelete(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) ScheduleList(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Sync(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
