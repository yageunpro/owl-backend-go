package calendar

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/internal/query"
)

type Store interface {
	CreateSchedule(ctx context.Context, arg CreateScheduleParam) error
	CreateGoogleSchedule(ctx context.Context, arg CreateGoogleScheduleParam) error
	GetSchedule(ctx context.Context, id, userId uuid.UUID) (*resSchedule, error)
	DeleteSchedule(ctx context.Context, id, userId uuid.UUID) error
	FindSchedule(ctx context.Context, arg FindScheduleParam) ([]resSchedule, error)
	GetSyncToken(ctx context.Context, userId uuid.UUID) (string, error)
	UpdateSyncToken(ctx context.Context, userId uuid.UUID, syncToken string) error
}

type store struct {
	pool *db.Pool
}

func New(pool *db.Pool) Store {
	return &store{pool: pool}
}

func (s *store) CreateSchedule(ctx context.Context, arg CreateScheduleParam) error {
	qry := query.New(s.pool)

	err := qry.CreateSchedule(ctx, query.CreateScheduleParams{
		ID:     arg.Id,
		UserID: arg.UserId,
		Title:  arg.Title,
		Period: pgtype.Range[pgtype.Timestamptz]{
			Lower:     pgtype.Timestamptz{Time: arg.StartTime.UTC(), Valid: true},
			Upper:     pgtype.Timestamptz{Time: arg.EndTime.UTC(), Valid: true},
			LowerType: pgtype.Inclusive,
			UpperType: pgtype.Inclusive,
			Valid:     true,
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to create calendar schedule"), err)
	}

	return nil
}

func (s *store) CreateGoogleSchedule(ctx context.Context, arg CreateGoogleScheduleParam) error {
	qry := query.New(s.pool)

	err := qry.CreateGoogleSchedule(ctx, query.CreateGoogleScheduleParams{
		ID:     arg.Id,
		UserID: arg.UserId,
		Title:  arg.Title,
		Period: pgtype.Range[pgtype.Timestamptz]{
			Lower:     pgtype.Timestamptz{Time: arg.StartTime.UTC(), Valid: true},
			Upper:     pgtype.Timestamptz{Time: arg.EndTime.UTC(), Valid: true},
			LowerType: pgtype.Inclusive,
			UpperType: pgtype.Inclusive,
			Valid:     true,
		},
		GoogleCalcID: pgtype.Text{
			String: arg.GoogleCalcId,
			Valid:  true,
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to create calendar schedule"), err)
	}

	return nil
}

func (s *store) GetSchedule(ctx context.Context, id, userId uuid.UUID) (*resSchedule, error) {
	qry := query.New(s.pool)

	row, err := qry.GetSchedule(ctx, query.GetScheduleParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return nil, errors.Join(errors.New("failed to get calendar schedule"), err)
	}

	return &resSchedule{
		Id:        row.ID,
		Title:     row.Title,
		StartTime: row.Period.Lower.Time,
		EndTime:   row.Period.Upper.Time,
	}, nil
}

func (s *store) DeleteSchedule(ctx context.Context, id, userId uuid.UUID) error {
	qry := query.New(s.pool)

	err := qry.DeleteSchedule(ctx, query.DeleteScheduleParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		return errors.Join(errors.New("failed to delete calendar schedule"), err)
	}

	return nil
}

func (s *store) FindSchedule(ctx context.Context, arg FindScheduleParam) ([]resSchedule, error) {
	qry := query.New(s.pool)

	rows, err := qry.FindSchedule(ctx, query.FindScheduleParams{
		UserID:      arg.UserId,
		StartTime:   pgtype.Timestamptz{Time: arg.StartTime.UTC(), Valid: true},
		EndTime:     pgtype.Timestamptz{Time: arg.EndTime.UTC(), Valid: true},
		OffsetCount: int32(arg.Offset),
		LimitCount:  int32(arg.Limit),
	})
	if err != nil {
		return nil, errors.Join(errors.New("failed to find calendar schedule"), err)
	}

	res := make([]resSchedule, len(rows))
	for i := range rows {
		res[i].Id = rows[i].ID
		res[i].Title = rows[i].Title
		res[i].StartTime = rows[i].Period.Lower.Time
		res[i].EndTime = rows[i].Period.Upper.Time
	}

	return res, nil
}

func (s *store) GetSyncToken(ctx context.Context, userId uuid.UUID) (string, error) {
	qry := query.New(s.pool)

	row, err := qry.GetSync(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", errors.Join(errors.New("failed to get sync token"), err)
	}

	return row.SyncToken, nil
}

func (s *store) UpdateSyncToken(ctx context.Context, userId uuid.UUID, syncToken string) error {
	qry := query.New(s.pool)
	return qry.CreateSync(ctx, query.CreateSyncParams{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    userId,
		SyncToken: syncToken,
	})
}
