package appointment

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/yageunpro/owl-backend-go/internal/db"
	"github.com/yageunpro/owl-backend-go/store/internal/query"
	"time"
)

type Store interface {
	Create(ctx context.Context, arg CreateAppointmentParam) error
	Get(ctx context.Context, appointmentId uuid.UUID) (*resAppointment, error)
	List(ctx context.Context, status Status, userId uuid.UUID, offset, limit int) ([]resAppointmentList, error)
	Delete(ctx context.Context, appointmentId, userId uuid.UUID) error
	Update(ctx context.Context, arg UpdateAppointmentParam) error
	Confirm(ctx context.Context, appointmentId, userId uuid.UUID, confirmTime time.Time) error
	AddParticipant(ctx context.Context, appointmentId, userId uuid.UUID) error
	AddParticipantNonMember(ctx context.Context, arg AddParticipantNonMemberParam) error
}

type store struct {
	pool *db.Pool
}

func New(pool *db.Pool) Store {
	return &store{pool: pool}
}

func (s *store) Create(ctx context.Context, arg CreateAppointmentParam) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to begin transaction"), err)
	}
	defer tx.Rollback(ctx)

	qry := query.New(s.pool).WithTx(tx)

	err = qry.CreateAppointment(ctx, query.CreateAppointmentParams{
		ID:          arg.Id,
		OrganizerID: arg.UserId,
		Title:       arg.Title,
		Description: arg.Description,
		Category:    arg.CategoryList,
		LocationID:  arg.LocationId,
		Deadline:    pgtype.Timestamptz{Time: arg.Deadline.UTC(), Valid: true},
	})
	if err != nil {
		return errors.Join(errors.New("failed to create appointment"), err)
	}

	err = qry.CreateParticipant(ctx, query.CreateParticipantParams{
		ID:            uuid.Must(uuid.NewV7()),
		AppointmentID: arg.Id,
		UserID:        arg.UserId,
	})
	if err != nil {
		return errors.Join(errors.New("failed to create participant"), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}

	return nil
}

func (s *store) Get(ctx context.Context, appointmentId uuid.UUID) (*resAppointment, error) {
	qry := query.New(s.pool)

	row, err := qry.GetAppointment(ctx, appointmentId)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get appointment"), err)
	}

	participantRows, err := qry.GetParticipants(ctx, appointmentId)
	if err != nil {
		return nil, errors.Join(errors.New("failed to get participants"), err)
	}

	var confirmTime *time.Time
	if row.ConfirmTime.Valid {
		confirmTime = &row.ConfirmTime.Time
	}

	res := resAppointment{
		Id:           appointmentId,
		OrganizerId:  row.OrganizerID,
		Status:       Status(row.Status),
		Title:        row.Title,
		Description:  row.Description,
		Deadline:     row.Deadline.Time,
		ConfirmTime:  confirmTime,
		LocationId:   row.LocationID,
		CategoryList: row.Category,
		Participants: make([]resParticipant, len(participantRows)),
	}
	if row.DeletedAt.Valid {
		res.Status = DELETE
	}

	for i := range participantRows {
		res.Participants[i].UserId = participantRows[i].UserID
		res.Participants[i].Username = participantRows[i].Username
	}

	return &res, nil
}

func (s *store) List(ctx context.Context, status Status, userId uuid.UUID, offset, limit int) ([]resAppointmentList, error) {
	qry := query.New(s.pool)

	switch status {
	case DRAFT, CANCEL, DELETE:
		rows, err := qry.GetStatusAppointment(ctx, query.GetStatusAppointmentParams{
			UserID: userId,
			Status: query.AppointmentStatus(status),
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			return nil, errors.Join(errors.New("failed to get appointment"), err)
		}
		res := make([]resAppointmentList, len(rows))
		for i := range rows {
			res[i] = resAppointmentList{
				Id:          rows[i].ID,
				OrganizerId: rows[i].OrganizerID,
				Title:       rows[i].Title,
				LocationId:  rows[i].LocationID,
				Status:      Status(rows[i].Status),
				ConfirmTime: nil,
				HeadCount:   int(rows[i].HeadCount),
			}
			if rows[i].ConfirmTime.Valid {
				res[i].ConfirmTime = &rows[i].ConfirmTime.Time
			}
		}
		return res, nil
	case CONFIRM:
		rows, err := qry.GetConfirmAppointment(ctx, query.GetConfirmAppointmentParams{
			UserID: userId,
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			return nil, errors.Join(errors.New("failed to get appointment"), err)
		}
		res := make([]resAppointmentList, len(rows))
		for i := range rows {
			res[i] = resAppointmentList{
				Id:          rows[i].ID,
				OrganizerId: rows[i].OrganizerID,
				Title:       rows[i].Title,
				LocationId:  rows[i].LocationID,
				Status:      Status(rows[i].Status),
				ConfirmTime: nil,
				HeadCount:   int(rows[i].HeadCount),
			}
			if rows[i].ConfirmTime.Valid {
				res[i].ConfirmTime = &rows[i].ConfirmTime.Time
			}
		}
		return res, nil
	case DONE:
		rows, err := qry.GetDoneAppointment(ctx, query.GetDoneAppointmentParams{
			UserID: userId,
			Offset: int32(offset),
			Limit:  int32(limit),
		})
		if err != nil {
			return nil, errors.Join(errors.New("failed to get appointment"), err)
		}
		res := make([]resAppointmentList, len(rows))
		for i := range rows {
			res[i] = resAppointmentList{
				Id:          rows[i].ID,
				OrganizerId: rows[i].OrganizerID,
				Title:       rows[i].Title,
				LocationId:  rows[i].LocationID,
				Status:      Status(rows[i].Status),
				ConfirmTime: nil,
				HeadCount:   int(rows[i].HeadCount),
			}
			if rows[i].ConfirmTime.Valid {
				res[i].ConfirmTime = &rows[i].ConfirmTime.Time
			}
		}
		return res, nil
	}

	return nil, errors.New("invalid status")
}

func (s *store) Delete(ctx context.Context, appointmentId, userId uuid.UUID) error {
	qry := query.New(s.pool)

	err := qry.DeleteAppointment(ctx, query.DeleteAppointmentParams{
		ID:          appointmentId,
		OrganizerID: userId,
	})
	if err != nil {
		return errors.Join(errors.New("failed to delete appointment"), err)
	}

	return nil
}

func (s *store) Update(ctx context.Context, arg UpdateAppointmentParam) error {
	qry := query.New(s.pool)

	param := query.UpdateAppointmentParams{
		ID:          arg.Id,
		OrganizerID: arg.UserId,
		UTitle:      pgtype.Text{Valid: false},
		UDesc:       pgtype.Text{Valid: false},
		ULocationID: arg.LocationId,
		UCategory:   arg.CategoryList,
	}
	if arg.Title != nil {
		param.UTitle.String = *arg.Title
		param.UTitle.Valid = true
	}
	if arg.Description != nil {
		param.UDesc.String = *arg.Description
		param.UDesc.Valid = true
	}

	err := qry.UpdateAppointment(ctx, param)
	if err != nil {
		return errors.Join(errors.New("failed to update appointment"), err)
	}

	return nil
}

func (s *store) Confirm(ctx context.Context, appointmentId, userId uuid.UUID, confirmTime time.Time) error {
	qry := query.New(s.pool)

	err := qry.ConfirmAppointment(ctx, query.ConfirmAppointmentParams{
		ID:          appointmentId,
		OrganizerID: userId,
		ConfirmTime: pgtype.Timestamptz{
			Time:  confirmTime.UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return errors.Join(errors.New("failed to confirm appointment"), err)
	}

	return nil
}

func (s *store) AddParticipant(ctx context.Context, appointmentId, userId uuid.UUID) error {
	qry := query.New(s.pool)

	err := qry.CreateParticipant(ctx, query.CreateParticipantParams{
		ID:            uuid.Must(uuid.NewV7()),
		AppointmentID: appointmentId,
		UserID:        userId,
	})
	if err != nil {
		return errors.Join(errors.New("failed to add appointment"), err)
	}

	return nil
}

func (s *store) AddParticipantNonMember(ctx context.Context, arg AddParticipantNonMemberParam) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to begin transaction"), err)
	}
	defer tx.Rollback(ctx)

	qry := query.New(s.pool).WithTx(tx)

	userId := uuid.Must(uuid.NewV7())

	err = qry.CreateUser(ctx, query.CreateUserParams{
		ID:       userId,
		Email:    userId.String() + "temp.user",
		Username: arg.Username,
	})
	if err != nil {
		return errors.Join(errors.New("failed to create user"), err)
	}

	for i := range arg.ScheduleList {
		err = qry.CreateSchedule(ctx, query.CreateScheduleParams{
			ID:     uuid.Must(uuid.NewV7()),
			UserID: userId,
			Title:  arg.ScheduleList[i].Title,
			Period: pgtype.Range[pgtype.Timestamptz]{
				Lower:     pgtype.Timestamptz{Time: arg.ScheduleList[i].StartTime.UTC(), Valid: true},
				Upper:     pgtype.Timestamptz{Time: arg.ScheduleList[i].EndTime.UTC(), Valid: true},
				LowerType: pgtype.Inclusive,
				UpperType: pgtype.Inclusive,
				Valid:     true,
			},
		})
		if err != nil {
			return errors.Join(errors.New("failed to create schedule"), err)
		}
	}

	err = qry.CreateParticipant(ctx, query.CreateParticipantParams{
		ID:            uuid.Must(uuid.NewV7()),
		AppointmentID: arg.Id,
		UserID:        userId,
	})
	if err != nil {
		return errors.Join(errors.New("failed to add appointment"), err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.Join(errors.New("failed to commit transaction"), err)
	}
	return nil
}
