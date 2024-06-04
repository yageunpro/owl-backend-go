-- name: CreateAppointment :exec
INSERT INTO appointment.appointment (id, organizer_id, title, description, category, location_id, deadline)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateParticipant :exec
INSERT INTO appointment.participant(id, appointment_id, user_id)
VALUES ($1, $2, $3);

-- name: GetAppointment :one
SELECT id,
       organizer_id,
       status,
       title,
       description,
       category,
       location_id,
       deadline,
       confirm_time,
       deleted_at
FROM appointment.appointment
WHERE id = $1;

-- name: GetParticipants :many
SELECT u.id       AS user_id,
       u.username AS username
FROM appointment.participant AS p
         INNER JOIN auth.user AS u ON p.user_id = u.id
WHERE appointment_id = $1;

-- name: DeleteAppointment :exec
UPDATE appointment.appointment
SET deleted_at = NOW(),
    status     = 'DELETE'
WHERE id = $1
  AND organizer_id = $2;

-- name: ConfirmAppointment :exec
UPDATE appointment.appointment
SET confirm_time = $3,
    status       = 'CONFIRM',
    updated_at   = NOW()
WHERE id = $1
  AND organizer_id = $2
  AND confirm_time IS NULL;

-- name: UpdateAppointment :exec
UPDATE appointment.appointment
SET title       = COALESCE(sqlc.narg(u_title), title),
    description = COALESCE(sqlc.narg(u_desc), description),
    location_id = COALESCE(sqlc.narg(u_location_id), location_id),
    category    = COALESCE(sqlc.narg(u_category), category),
    updated_at  = NOW()
WHERE id = $1
  AND organizer_id = $2;

-- name: GetStatusAppointment :many
WITH ap_with_count AS (SELECT a.id           AS id,
                              a.organizer_id AS organizer_id,
                              a.status       AS status,
                              a.title        AS title,
                              a.location_id  AS location_id,
                              a.confirm_time AS confirm_time,
                              a.created_at   AS created_at,
                              a.deleted_at   AS deleted_at,
                              COUNT(p.id)    AS head_count
                       FROM appointment.appointment AS a
                                INNER JOIN appointment.participant p ON a.id = p.appointment_id
                       GROUP BY a.id)
SELECT a.id           AS id,
       a.organizer_id AS organizer_id,
       a.status       AS status,
       a.title        AS title,
       a.location_id  AS location_id,
       a.confirm_time AS confirm_time,
       a.deleted_at   AS deleted_at,
       a.head_count   AS head_count
FROM (SELECT appointment_id AS id
      FROM appointment.participant AS t
      WHERE t.user_id = $1) AS ids
         INNER JOIN ap_with_count AS a ON ids.id = a.id
WHERE status = $2
ORDER BY a.created_at DESC
OFFSET $3 LIMIT $4;


-- name: GetConfirmAppointment :many
WITH ap_with_count AS (SELECT a.id           AS id,
                              a.organizer_id AS organizer_id,
                              a.status       AS status,
                              a.title        AS title,
                              a.location_id  AS location_id,
                              a.confirm_time AS confirm_time,
                              a.created_at   AS created_at,
                              a.deleted_at   AS deleted_at,
                              COUNT(p.id)    AS head_count
                       FROM appointment.appointment AS a
                                INNER JOIN appointment.participant p ON a.id = p.appointment_id
                       GROUP BY a.id)
SELECT a.id           AS id,
       a.organizer_id AS organizer_id,
       a.status       AS status,
       a.title        AS title,
       a.location_id  AS location_id,
       a.confirm_time AS confirm_time,
       a.deleted_at   AS deleted_at,
       a.head_count   AS head_count
FROM (SELECT appointment_id AS id
      FROM appointment.participant AS t
      WHERE t.user_id = $1) AS ids
         INNER JOIN ap_with_count AS a ON ids.id = a.id
WHERE status = 'CONFIRM'
ORDER BY a.confirm_time
OFFSET $2 LIMIT $3;


-- name: GetDoneAppointment :many
WITH ap_with_count AS (SELECT a.id           AS id,
                              a.organizer_id AS organizer_id,
                              a.status       AS status,
                              a.title        AS title,
                              a.location_id  AS location_id,
                              a.confirm_time AS confirm_time,
                              a.created_at   AS created_at,
                              a.deleted_at   AS deleted_at,
                              COUNT(p.id)    AS head_count
                       FROM appointment.appointment AS a
                                INNER JOIN appointment.participant p ON a.id = p.appointment_id
                       GROUP BY a.id)
SELECT a.id           AS id,
       a.organizer_id AS organizer_id,
       a.status       AS status,
       a.title        AS title,
       a.location_id  AS location_id,
       a.confirm_time AS confirm_time,
       a.deleted_at   AS deleted_at,
       a.head_count   AS head_count
FROM (SELECT appointment_id AS id
      FROM appointment.participant AS t
      WHERE t.user_id = $1) AS ids
         INNER JOIN ap_with_count AS a ON ids.id = a.id
WHERE status = 'DONE'
ORDER BY a.confirm_time DESC
OFFSET $2 LIMIT $3;
