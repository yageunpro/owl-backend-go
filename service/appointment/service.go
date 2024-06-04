package appointment

import (
	"context"
	"github.com/google/uuid"
	"github.com/yageunpro/owl-backend-go/store"
	"github.com/yageunpro/owl-backend-go/store/appointment"
	"strconv"
	"time"
)

type Service interface {
	Add(ctx context.Context, arg AddParam) error
	List(ctx context.Context, arg ListParam) (*resInfoList, error)
	Info(ctx context.Context, appointmentId uuid.UUID) (*resInfo, error)
	Edit(ctx context.Context, arg EditParam) (*resInfo, error)
	Delete(ctx context.Context, appointmentId, userId uuid.UUID) error
	Share(ctx context.Context) error
	Join(ctx context.Context, appointmentId, userId uuid.UUID) (*resInfo, error)
	JoinNonmember(ctx context.Context, arg JoinNonmemberParam) (*resInfo, error)
	RecommendTime(ctx context.Context, appointmentId uuid.UUID) error
	Confirm(ctx context.Context, appointmentId, userId uuid.UUID, confirmTime time.Time) error
}
type service struct {
	store *store.Store
}

func New(sto *store.Store) Service {
	return &service{store: sto}
}

func (s *service) Add(ctx context.Context, arg AddParam) error {
	return s.store.Appointment.Create(ctx, appointment.CreateAppointmentParam{
		Id:           uuid.Must(uuid.NewV7()),
		UserId:       arg.UserId,
		Title:        arg.Title,
		Description:  arg.Description,
		LocationId:   arg.LocationId,
		CategoryList: arg.CategoryList,
		Deadline:     arg.Deadline,
	})
}

func (s *service) List(ctx context.Context, arg ListParam) (*resInfoList, error) {
	var offset int
	if arg.PageToken == nil || *arg.PageToken == "" {
		offset = 0
	} else {
		i, err := strconv.ParseInt(*arg.PageToken, 10, 32)
		if err != nil {
			offset = 0
		} else {
			offset = int(i)
		}
	}

	out, err := s.store.Appointment.List(ctx, appointment.Status(arg.Status), arg.UserId, offset, arg.Limit)
	if err != nil {
		return nil, err
	}

	offset += len(out)

	res := resInfoList{
		InfoList:  make([]absInfo, len(out)),
		NextToken: strconv.Itoa(offset),
	}

	for i := range out {
		var location *resLocation
		if out[i].LocationId.Valid {
			l, err := s.store.Location.GetLocation(ctx, out[i].LocationId.UUID)
			if err != nil {
				location = nil
			} else {
				location = &resLocation{
					Id:       l.Id,
					Title:    l.Title,
					Address:  l.Address,
					Category: l.Category,
					Position: [2]int{l.MapX, l.MapY},
				}
			}
		}

		res.InfoList[i] = absInfo{
			Id:          out[i].Id,
			OrganizerId: out[i].OrganizerId,
			Title:       out[i].Title,
			Location:    location,
			Status:      string(out[i].Status),
			ConfirmTime: out[i].ConfirmTime,
			HeadCount:   out[i].HeadCount,
		}
	}

	return &res, nil
}

func (s *service) Info(ctx context.Context, appointmentId uuid.UUID) (*resInfo, error) {
	out, err := s.store.Appointment.Get(ctx, appointmentId)
	if err != nil {
		return nil, err
	}

	participants := make([]resParticipant, len(out.Participants))
	for i := range out.Participants {
		participants[i].UserId = out.Participants[i].UserId
		participants[i].Username = out.Participants[i].Username
	}

	var location *resLocation
	if out.LocationId.Valid {
		l, err := s.store.Location.GetLocation(ctx, out.LocationId.UUID)
		if err != nil {
			location = nil
		} else {
			location = &resLocation{
				Id:       l.Id,
				Title:    l.Title,
				Address:  l.Address,
				Category: l.Category,
				Position: [2]int{l.MapX, l.MapY},
			}
		}
	}

	return &resInfo{
		Id:           out.Id,
		OrganizerId:  out.OrganizerId,
		Status:       string(out.Status),
		Title:        out.Title,
		Description:  out.Description,
		Deadline:     out.Deadline,
		ConfirmTime:  out.ConfirmTime,
		Location:     location,
		CategoryList: out.CategoryList,
		Participants: participants,
	}, nil
}

func (s *service) Edit(ctx context.Context, arg EditParam) (*resInfo, error) {
	locationId := uuid.NullUUID{Valid: false}
	if arg.LocationId != nil {
		locationId.UUID = *arg.LocationId
		locationId.Valid = true
	}
	err := s.store.Appointment.Update(ctx, appointment.UpdateAppointmentParam{
		Id:           arg.Id,
		UserId:       arg.UserId,
		Title:        arg.Title,
		Description:  arg.Description,
		LocationId:   locationId,
		CategoryList: arg.CategoryList,
	})
	if err != nil {
		return nil, err
	}

	return s.Info(ctx, arg.Id)
}

func (s *service) Delete(ctx context.Context, appointmentId, userId uuid.UUID) error {
	return s.store.Appointment.Delete(ctx, appointmentId, userId)
}

func (s *service) Share(ctx context.Context) error {
	return nil
}

func (s *service) Join(ctx context.Context, appointmentId, userId uuid.UUID) (*resInfo, error) {
	err := s.store.Appointment.AddParticipant(ctx, appointmentId, userId)
	if err != nil {
		return nil, err
	}

	return s.Info(ctx, appointmentId)
}

func (s *service) JoinNonmember(ctx context.Context, arg JoinNonmemberParam) (*resInfo, error) {
	param := appointment.AddParticipantNonMemberParam{
		Id:           arg.Id,
		Username:     arg.Username,
		ScheduleList: nil,
	}
	if arg.JoinSchedule == nil {
		param.ScheduleList = make([]appointment.Schedule, 0)
	} else {
		param.ScheduleList = make([]appointment.Schedule, len(arg.JoinSchedule))
		for i := range arg.JoinSchedule {
			param.ScheduleList[i] = appointment.Schedule{
				Title:     arg.JoinSchedule[i].Title,
				StartTime: arg.JoinSchedule[i].StartTime,
				EndTime:   arg.JoinSchedule[i].EndTime,
			}
		}
	}

	err := s.store.Appointment.AddParticipantNonMember(ctx, param)
	if err != nil {
		return nil, err
	}

	return s.Info(ctx, param.Id)
}

func (s *service) RecommendTime(ctx context.Context, appointmentId uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (s *service) Confirm(ctx context.Context, appointmentId, userId uuid.UUID, confirmTime time.Time) error {
	return s.store.Appointment.Confirm(ctx, appointmentId, userId, confirmTime)
}
