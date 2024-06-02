package user

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/service"
	"net/http"
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
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	out, err := h.svc.User.Info(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	res := resUserInfo{
		Id:       out.Id,
		Username: out.Username,
		Email:    out.Email,
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) ListAccount(c echo.Context) error {
	return echo.ErrNotFound
}

func (h *handler) AddAccount(c echo.Context) error {
	return echo.ErrNotFound
}

func (h *handler) VerifyAccount(c echo.Context) error {
	return echo.ErrNotFound
}

func (h *handler) DeleteAccount(c echo.Context) error {
	return echo.ErrNotFound
}
