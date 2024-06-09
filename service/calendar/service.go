package calendar

import (
	"context"
	"github.com/google/uuid"
	calc "github.com/yageunpro/owl-backend-go/internal/google/calendar"
	"github.com/yageunpro/owl-backend-go/internal/google/oauth"
	"github.com/yageunpro/owl-backend-go/store"
	"github.com/yageunpro/owl-backend-go/store/calendar"
	"strconv"
)

type Service interface {
	ScheduleAdd(ctx context.Context, arg ScheduleAddParam) (uuid.UUID, error)
	ScheduleInfo(ctx context.Context, id uuid.UUID, userId uuid.UUID) (*resSchedule, error)
	ScheduleDelete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error
	ScheduleList(ctx context.Context, arg ScheduleListParam) (*resScheduleList, error)
	Sync(ctx context.Context, userId uuid.UUID) error
}

type service struct {
	store *store.Store
}

func New(sto *store.Store) Service {
	return &service{store: sto}
}

func (s *service) ScheduleAdd(ctx context.Context, arg ScheduleAddParam) (uuid.UUID, error) {
	scheduleId := uuid.Must(uuid.NewV7())
	err := s.store.Calendar.CreateSchedule(ctx, calendar.CreateScheduleParam{
		Id:        scheduleId,
		UserId:    arg.UserId,
		Title:     arg.Title,
		StartTime: arg.StartTime,
		EndTime:   arg.EndTime,
	})
	if err != nil {
		return uuid.Nil, err
	}

	return scheduleId, nil
}

func (s *service) ScheduleInfo(ctx context.Context, id uuid.UUID, userId uuid.UUID) (*resSchedule, error) {
	out, err := s.store.Calendar.GetSchedule(ctx, id, userId)
	if err != nil {
		return nil, err
	}

	return &resSchedule{
		Id:        out.Id,
		Title:     out.Title,
		StartTime: out.StartTime,
		EndTime:   out.EndTime,
	}, nil
}

func (s *service) ScheduleDelete(ctx context.Context, id uuid.UUID, userId uuid.UUID) error {
	return s.store.Calendar.DeleteSchedule(ctx, id, userId)
}

func (s *service) ScheduleList(ctx context.Context, arg ScheduleListParam) (*resScheduleList, error) {
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

	out, err := s.store.Calendar.FindSchedule(ctx, calendar.FindScheduleParam{
		UserId:    arg.UserId,
		StartTime: arg.StartTime,
		EndTime:   arg.EndTime,
		Offset:    offset,
		Limit:     arg.Limit,
	})
	if err != nil {
		return nil, err
	}

	res := resScheduleList{
		ScheduleList: make([]resSchedule, len(out)),
		NextToken:    "",
	}

	for i := range out {
		res.ScheduleList[i].Id = out[i].Id
		res.ScheduleList[i].Title = out[i].Title
		res.ScheduleList[i].StartTime = out[i].StartTime
		res.ScheduleList[i].EndTime = out[i].EndTime
	}

	if len(res.ScheduleList) != 0 {
		offset += len(res.ScheduleList)
		res.NextToken = strconv.FormatInt(int64(offset), 10)
	}

	return &res, nil
}

func (s *service) Sync(ctx context.Context, userId uuid.UUID) error {
	out, err := s.store.Auth.GetOAuthToken(ctx, userId)
	if err != nil {
		return err
	}

	tok := oauth.ToToken(out.AccessToken, out.RefreshToken, out.ExpireTime)
	source, err := oauth.TokenSource(ctx, tok)
	if err != nil {
		return err
	}

	syncToken, err := s.store.Calendar.GetSyncToken(ctx, userId)
	if err != nil {
		return err
	}

	res, err := calc.GetAllEvents(ctx, calc.Param{
		Source:    source,
		SyncToken: syncToken,
	})
	if err != nil {
		return err
	}

	for _, event := range res.Events {
		err = s.store.Calendar.CreateGoogleSchedule(ctx, calendar.CreateGoogleScheduleParam{
			Id:           uuid.Must(uuid.NewV7()),
			UserId:       userId,
			Title:        event.Title,
			StartTime:    event.StartTime,
			EndTime:      event.EndTime,
			GoogleCalcId: event.CalcId,
		})
		if err != nil {
			return err
		}
	}

	err = s.store.Calendar.UpdateSyncToken(ctx, userId, res.SyncToken)
	if err != nil {
		return err
	}

	return nil
}
