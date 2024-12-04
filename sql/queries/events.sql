-- name: CreateEvent :one
INSERT INTO events (id, message, created_at, severity)
VALUES(
  $1,
  $2,
  $3,
  $4
) RETURNING *;

-- name: GetEventById :one
SELECT * FROM events
WHERE id = $1 and created_at = $2;

-- name: GetCriticalEvents :many
SELECT * from events
where severity = "Critical";

-- name: GetEvents :many
SELECT * FROM events
LIMIT $1;