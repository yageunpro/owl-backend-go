package location

import (
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/service"
)

type Handler interface {
	Search(c echo.Context) error
}

type handler struct {
	svc *service.Service
}

func New(svc *service.Service) Handler {
	return &handler{svc: svc}
}

func (h *handler) Search(c echo.Context) error {
	return nil
}
