package appointment

import (
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/service"
)

type Handler interface {
	Add(c echo.Context) error
	List(c echo.Context) error
	Info(c echo.Context) error
	Edit(c echo.Context) error
	Delete(c echo.Context) error
	Share(c echo.Context) error
	Join(c echo.Context) error
	JoinNonmember(c echo.Context) error
	RecommendTime(c echo.Context) error
	Confirm(c echo.Context) error
}
type handler struct {
	svc *service.Service
}

func New(svc *service.Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) Add(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) List(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Info(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Edit(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Delete(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Share(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Join(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) JoinNonmember(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) RecommendTime(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Confirm(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}
