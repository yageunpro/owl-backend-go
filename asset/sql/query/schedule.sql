-- name: CreateSchedule :exec
INSERT INTO calendar.schedule (id, user_id, title, period, google_calc_id)
VALUES ($1, $2, $3, $4, NULL);

-- name: CreateGoogleSchedule :exec
INSERT INTO calendar.schedule (id, user_id, title, period, google_calc_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (google_calc_id)
    DO UPDATE
    SET user_id    = $2,
        title      = $3,
        period     = $4,
        updated_at = NOW();

-- name: GetSchedule :one
SELECT id,
       title,
       period,
       deleted_at
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

-- name: CreateSync :exec
INSERT INTO calendar.sync (id, user_id, sync_token)
VALUES ($1, $2, $3);

-- name: GetSync :one
SELECT id, user_id, sync_token
FROM calendar.sync
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: GetAllSchedule :many
SELECT id, user_id, period
FROM calendar.schedule
WHERE deleted_at IS NULL
  AND user_id = ANY (@user_ids::uuid[])
  AND period && TSTZRANGE(@start_time::timestamptz, @end_time::timestamptz, '[]');