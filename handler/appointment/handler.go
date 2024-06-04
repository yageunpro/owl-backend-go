package appointment

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/yageunpro/owl-backend-go/internal/jwt"
	"github.com/yageunpro/owl-backend-go/service"
	"github.com/yageunpro/owl-backend-go/service/appointment"
	"net/http"
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
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqAdd)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	param := appointment.AddParam{
		UserId:       userId,
		Title:        req.Title,
		Description:  req.Description,
		LocationId:   uuid.NullUUID{},
		CategoryList: nil,
		Deadline:     req.Deadline,
	}
	if req.LocationId == nil {
		param.LocationId = uuid.NullUUID{
			UUID:  uuid.Nil,
			Valid: false,
		}
	} else {
		param.LocationId = uuid.NullUUID{
			UUID:  *req.LocationId,
			Valid: true,
		}
	}

	if req.CategoryList == nil {
		param.CategoryList = make([]string, 0)
	} else {
		param.CategoryList = req.CategoryList
	}

	err = h.svc.Appointment.Add(c.Request().Context(), param)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) List(c echo.Context) error {
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqList)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	var limit int
	if req.Limit == nil {
		limit = 10
	} else {
		limit = *req.Limit
	}

	out, err := h.svc.Appointment.List(c.Request().Context(), appointment.ListParam{
		UserId:    userId,
		Status:    req.Type,
		PageToken: req.PageToken,
		Limit:     limit,
	})
	if err != nil {
		return err
	}

	res := resInfoList{
		Data:      make([]absInfo, len(out.InfoList)),
		NextToken: out.NextToken,
	}
	for i := range out.InfoList {
		res.Data[i] = absInfo{
			Id:          out.InfoList[i].Id,
			OrganizerId: out.InfoList[i].OrganizerId,
			Title:       out.InfoList[i].Title,
			Location:    nil,
			Status:      out.InfoList[i].Status,
			ConfirmTime: out.InfoList[i].ConfirmTime,
			HeadCount:   out.InfoList[i].HeadCount,
		}

		if out.InfoList[i].Location != nil {
			res.Data[i].Location = &resLocation{
				Id:       out.InfoList[i].Location.Id,
				Title:    out.InfoList[i].Location.Title,
				Address:  out.InfoList[i].Location.Address,
				Category: out.InfoList[i].Location.Category,
				Position: out.InfoList[i].Location.Position,
			}
		}
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) Info(c echo.Context) error {
	req := new(reqInfo)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	out, err := h.svc.Appointment.Info(c.Request().Context(), req.Id)
	if err != nil {
		return err
	}

	res := resInfo{
		Id:              out.Id,
		OrganizerId:     out.OrganizerId,
		Title:           out.Title,
		Location:        nil,
		Status:          out.Status,
		ConfirmTime:     out.ConfirmTime,
		Description:     out.Description,
		CategoryList:    out.CategoryList,
		ParticipantList: make([]resParticipant, len(out.Participants)),
		Deadline:        out.Deadline,
	}

	if out.Location != nil {
		res.Location = &resLocation{
			Id:       out.Location.Id,
			Title:    out.Location.Title,
			Address:  out.Location.Address,
			Category: out.Location.Category,
			Position: out.Location.Position,
		}
	}
	for i := range out.Participants {
		res.ParticipantList[i].Id = out.Participants[i].UserId
		res.ParticipantList[i].Name = out.Participants[i].Username
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) Edit(c echo.Context) error {
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqEdit)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	out, err := h.svc.Appointment.Edit(c.Request().Context(), appointment.EditParam{
		Id:           req.Id,
		UserId:       userId,
		Title:        req.Title,
		Description:  req.Description,
		LocationId:   req.LocationId,
		CategoryList: req.CategoryList,
	})
	if err != nil {
		return err
	}

	res := resInfo{
		Id:              out.Id,
		OrganizerId:     out.OrganizerId,
		Title:           out.Title,
		Location:        nil,
		Status:          out.Status,
		ConfirmTime:     out.ConfirmTime,
		Description:     out.Description,
		CategoryList:    out.CategoryList,
		ParticipantList: make([]resParticipant, len(out.Participants)),
		Deadline:        out.Deadline,
	}

	if out.Location != nil {
		res.Location = &resLocation{
			Id:       out.Location.Id,
			Title:    out.Location.Title,
			Address:  out.Location.Address,
			Category: out.Location.Category,
			Position: out.Location.Position,
		}
	}
	for i := range out.Participants {
		res.ParticipantList[i].Id = out.Participants[i].UserId
		res.ParticipantList[i].Name = out.Participants[i].Username
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) Delete(c echo.Context) error {
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqDelete)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = h.svc.Appointment.Delete(c.Request().Context(), req.Id, userId)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}

func (h *handler) Share(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *handler) Join(c echo.Context) error {
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqJoin)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	out, err := h.svc.Appointment.Join(c.Request().Context(), req.Id, userId)
	if err != nil {
		return err
	}

	res := resInfo{
		Id:              out.Id,
		OrganizerId:     out.OrganizerId,
		Title:           out.Title,
		Location:        nil,
		Status:          out.Status,
		ConfirmTime:     out.ConfirmTime,
		Description:     out.Description,
		CategoryList:    out.CategoryList,
		ParticipantList: make([]resParticipant, len(out.Participants)),
		Deadline:        out.Deadline,
	}

	if out.Location != nil {
		res.Location = &resLocation{
			Id:       out.Location.Id,
			Title:    out.Location.Title,
			Address:  out.Location.Address,
			Category: out.Location.Category,
			Position: out.Location.Position,
		}
	}
	for i := range out.Participants {
		res.ParticipantList[i].Id = out.Participants[i].UserId
		res.ParticipantList[i].Name = out.Participants[i].Username
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) JoinNonmember(c echo.Context) error {
	req := new(reqJoinNonmember)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	param := appointment.JoinNonmemberParam{
		Id:           req.Id,
		Username:     req.Username,
		JoinSchedule: nil,
	}
	if req.ScheduleList == nil {
		param.JoinSchedule = make([]appointment.JoinSchedule, 0)
	} else {
		param.JoinSchedule = make([]appointment.JoinSchedule, len(req.ScheduleList))
		for i := range req.ScheduleList {
			param.JoinSchedule[i] = appointment.JoinSchedule{
				Title:     req.ScheduleList[i].Title,
				StartTime: req.ScheduleList[i].StartTime,
				EndTime:   req.ScheduleList[i].EndTime,
			}
		}
	}

	out, err := h.svc.Appointment.JoinNonmember(c.Request().Context(), param)
	if err != nil {
		return err
	}

	res := resInfo{
		Id:              out.Id,
		OrganizerId:     out.OrganizerId,
		Title:           out.Title,
		Location:        nil,
		Status:          out.Status,
		ConfirmTime:     out.ConfirmTime,
		Description:     out.Description,
		CategoryList:    out.CategoryList,
		ParticipantList: make([]resParticipant, len(out.Participants)),
		Deadline:        out.Deadline,
	}

	if out.Location != nil {
		res.Location = &resLocation{
			Id:       out.Location.Id,
			Title:    out.Location.Title,
			Address:  out.Location.Address,
			Category: out.Location.Category,
			Position: out.Location.Position,
		}
	}
	for i := range out.Participants {
		res.ParticipantList[i].Id = out.Participants[i].UserId
		res.ParticipantList[i].Name = out.Participants[i].Username
	}

	return c.JSON(http.StatusOK, res)
}

func (h *handler) RecommendTime(c echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func (h *handler) Confirm(c echo.Context) error {
	userId := jwt.GetUserID(c.Request().Context())
	if userId == uuid.Nil {
		return echo.ErrUnauthorized
	}

	req := new(reqConfirm)
	err := c.Bind(req)
	if err != nil {
		return err
	}

	err = h.svc.Appointment.Confirm(c.Request().Context(), req.Id, userId, req.Time)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
