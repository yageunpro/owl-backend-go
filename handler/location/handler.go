package location

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/service"
	"net/http"
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
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqSearch)
	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}
	if req.Query == nil && len(req.Query) == 0 {
		return c.JSON(http.StatusOK, make([]resLocation, 0))
	}

	out, err := h.svc.Location.Query(c.Request().Context(), req.Query)
	if err != nil {
		return err
	}

	res := make([]resLocation, len(out.Locations))
	for i := range out.Locations {
		res[i].Id = out.Locations[i].Id
		res[i].Title = out.Locations[i].Title
		res[i].Address = out.Locations[i].Address
		res[i].Category = out.Locations[i].Category
		res[i].Position = [2]int{out.Locations[i].MapX, out.Locations[i].MapY}
	}

	return c.JSON(http.StatusOK, res)
}
