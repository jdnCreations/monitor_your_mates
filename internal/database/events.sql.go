// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: events.sql

package database

import (
	"context"
	"database/sql"
	"time"
)

const createEvent = `-- name: CreateEvent :one
INSERT INTO events (id, message, created_at, severity)
VALUES(
  $1,
  $2,
  $3,
  $4
) RETURNING id, message, created_at, severity
`

type CreateEventParams struct {
	ID        int32
	Message   sql.NullString
	CreatedAt time.Time
	Severity  sql.NullString
}

func (q *Queries) CreateEvent(ctx context.Context, arg CreateEventParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, createEvent,
		arg.ID,
		arg.Message,
		arg.CreatedAt,
		arg.Severity,
	)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.CreatedAt,
		&i.Severity,
	)
	return i, err
}

const getCriticalEvents = `-- name: GetCriticalEvents :many
SELECT id, message, created_at, severity from events
where severity = "Critical"
`

func (q *Queries) GetCriticalEvents(ctx context.Context) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getCriticalEvents)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Message,
			&i.CreatedAt,
			&i.Severity,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEventById = `-- name: GetEventById :one
SELECT id, message, created_at, severity FROM events
WHERE id = $1 and created_at = $2
`

type GetEventByIdParams struct {
	ID        int32
	CreatedAt time.Time
}

func (q *Queries) GetEventById(ctx context.Context, arg GetEventByIdParams) (Event, error) {
	row := q.db.QueryRowContext(ctx, getEventById, arg.ID, arg.CreatedAt)
	var i Event
	err := row.Scan(
		&i.ID,
		&i.Message,
		&i.CreatedAt,
		&i.Severity,
	)
	return i, err
}

const getEvents = `-- name: GetEvents :many
SELECT id, message, created_at, severity FROM events
LIMIT $1
`

func (q *Queries) GetEvents(ctx context.Context, limit int32) ([]Event, error) {
	rows, err := q.db.QueryContext(ctx, getEvents, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Event
	for rows.Next() {
		var i Event
		if err := rows.Scan(
			&i.ID,
			&i.Message,
			&i.CreatedAt,
			&i.Severity,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
