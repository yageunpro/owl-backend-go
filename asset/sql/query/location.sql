-- name: GetLocation :one
SELECT id, title, address, category, map_x, map_y
FROM location.location
WHERE id = $1;

-- name: GetLocationWithQueryId :many
SELECT id, title, address, category, map_x, map_y
FROM location.location
WHERE deleted_at IS NULL
  AND query_id = $1;

-- name: FindQueryId :one
SELECT id, updated_at
FROM location.query
WHERE data = $1;

-- name: DeprecateQuery :exec
UPDATE location.location
SET deleted_at = NOW()
WHERE query_id = $1;

-- name: UpdateQueryTime :exec
UPDATE location.query
SET updated_at = NOW()
WHERE id = $1;

-- name: CreateQuery :exec
INSERT INTO location.query (id, data)
VALUES ($1, $2);

-- name: CreateLocation :exec
INSERT INTO location.location (id, query_id, title, address, category, map_x, map_y)
VALUES ($1, $2, $3, $4, $5, $6, $7);
