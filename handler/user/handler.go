package user

import (
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/service"
)

type Handler interface {
	Me(c echo.Context) error
	ListAccount(c echo.Context) error
	AddAccount(c echo.Context) error
	VerifyAccount(c echo.Context) error
	DeleteAccount(c echo.Context) error
}

type handler struct {
	svc *service.Service
}

func New(svc *service.Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) Me(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) ListAccount(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) AddAccount(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) VerifyAccount(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) DeleteAccount(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
