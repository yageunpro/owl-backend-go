-- name: CreateSchedule :exec
INSERT INTO calendar.schedule (id, user_id, title, period)
VALUES ($1, $2, $3, $4);

-- name: GetSchedule :one
SELECT id, title, period, deleted_at
FROM calendar.schedule
WHERE id = $1
  AND user_id = @user_id;

-- name: DeleteSchedule :exec
UPDATE calendar.schedule
SET deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND user_id = @user_id;

-- name: FindSchedule :many
SELECT id, title, period
FROM calendar.schedule
WHERE deleted_at IS NULL
  AND user_id = @user_id
  AND period && TSTZRANGE(@start_time::timestamptz, @end_time::timestamptz, '[]')
ORDER BY LOWER(period), UPPER(period)
OFFSET @offset_count LIMIT @limit_count;
