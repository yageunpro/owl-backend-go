// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: schedule.sql

package query

import (
	"context"

	uuid "github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createSchedule = `-- name: CreateSchedule :exec
INSERT INTO calendar.schedule (id, user_id, title, period)
VALUES ($1, $2, $3, $4)
`

type CreateScheduleParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Title  string
	Period pgtype.Range[pgtype.Timestamptz]
}

func (q *Queries) CreateSchedule(ctx context.Context, arg CreateScheduleParams) error {
	_, err := q.db.Exec(ctx, createSchedule,
		arg.ID,
		arg.UserID,
		arg.Title,
		arg.Period,
	)
	return err
}

const deleteSchedule = `-- name: DeleteSchedule :exec
UPDATE calendar.schedule
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND user_id = $2
`

type DeleteScheduleParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

func (q *Queries) DeleteSchedule(ctx context.Context, arg DeleteScheduleParams) error {
	_, err := q.db.Exec(ctx, deleteSchedule, arg.ID, arg.UserID)
	return err
}

const findSchedule = `-- name: FindSchedule :many
SELECT id, title, period
FROM calendar.schedule
WHERE deleted_at IS NULL
  AND user_id = $1
  AND period && TSTZRANGE($2::timestamptz, $3::timestamptz, '[]')
ORDER BY LOWER(period), UPPER(period)
OFFSET $4 LIMIT $5
`

type FindScheduleParams struct {
	UserID      uuid.UUID
	StartTime   pgtype.Timestamptz
	EndTime     pgtype.Timestamptz
	OffsetCount int32
	LimitCount  int32
}

type FindScheduleRow struct {
	ID     uuid.UUID
	Title  string
	Period pgtype.Range[pgtype.Timestamptz]
}

func (q *Queries) FindSchedule(ctx context.Context, arg FindScheduleParams) ([]FindScheduleRow, error) {
	rows, err := q.db.Query(ctx, findSchedule,
		arg.UserID,
		arg.StartTime,
		arg.EndTime,
		arg.OffsetCount,
		arg.LimitCount,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FindScheduleRow
	for rows.Next() {
		var i FindScheduleRow
		if err := rows.Scan(&i.ID, &i.Title, &i.Period); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSchedule = `-- name: GetSchedule :one
SELECT id, title, period, deleted_at
FROM calendar.schedule
WHERE id = $1
  AND user_id = $2
`

type GetScheduleParams struct {
	ID     uuid.UUID
	UserID uuid.UUID
}

type GetScheduleRow struct {
	ID        uuid.UUID
	Title     string
	Period    pgtype.Range[pgtype.Timestamptz]
	DeletedAt pgtype.Timestamptz
}

func (q *Queries) GetSchedule(ctx context.Context, arg GetScheduleParams) (GetScheduleRow, error) {
	row := q.db.QueryRow(ctx, getSchedule, arg.ID, arg.UserID)
	var i GetScheduleRow
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Period,
		&i.DeletedAt,
	)
	return i, err
}
