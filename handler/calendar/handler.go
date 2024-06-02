package calendar

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/service"
	"github.com/yageunpro/owl-backend-go/service/calendar"
	"net/http"
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
	req := new(reqScheduleAdd)

	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	scheduleId, err := h.svc.Calendar.ScheduleAdd(c.Request().Context(), calendar.ScheduleAddParam{
		UserId:    userId,
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resSchedule{
		Id:        scheduleId,
		Title:     req.Title,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
}

func (h *handler) ScheduleInfo(c echo.Context) error {
	req := new(reqScheduleInfo)

	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	out, err := h.svc.Calendar.ScheduleInfo(c.Request().Context(), req.Id, userId)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resSchedule{
		Id:        out.Id,
		Title:     out.Title,
		StartTime: out.StartTime,
		EndTime:   out.EndTime,
	})
}

func (h *handler) ScheduleDelete(c echo.Context) error {
	req := new(reqScheduleDelete)

	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}

	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	err = h.svc.Calendar.ScheduleDelete(c.Request().Context(), req.Id, userId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) ScheduleList(c echo.Context) error {
	req := new(reqScheduleList)

	err := c.Bind(req)
	if err != nil {
		// TODO:: add validation logic
		return err
	}
	if req.Limit == 0 {
		req.Limit = 10
	}

	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	out, err := h.svc.Calendar.ScheduleList(c.Request().Context(), calendar.ScheduleListParam{
		UserId:    userId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		PageToken: req.PageToken,
		Limit:     req.Limit,
	})
	if err != nil {
		return err
	}

	res := resScheduleList{
		Data:      make([]resSchedule, len(out.ScheduleList)),
		NextToken: out.NextToken,
	}

	for i := range out.ScheduleList {
		res.Data[i].Id = out.ScheduleList[i].Id
		res.Data[i].Title = out.ScheduleList[i].Title
		res.Data[i].StartTime = out.ScheduleList[i].StartTime
		res.Data[i].EndTime = out.ScheduleList[i].EndTime
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) Sync(c echo.Context) error {
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	err := h.svc.Calendar.Sync(c.Request().Context(), userId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
