// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package costumers_db

import (
	"context"
	"encoding/json"
	"time"
)

const addCostumer = `-- name: AddCostumer :one
INSERT INTO costumers (name, cpf, birthdate, phone, email, phone_validated, email_validated, deleted,created_at,updated_at) VALUES
($1,$2,$3,$4,$5, false,false,false, NOW(), NULL) RETURNING id, name, cpf, birthdate, phone, email, phone_validated, email_validated, deleted, created_at, updated_at
`

type AddCostumerParams struct {
	Name      string
	Cpf       string
	Birthdate time.Time
	Phone     string
	Email     string
}

func (q *Queries) AddCostumer(ctx context.Context, arg AddCostumerParams) (Costumer, error) {
	row := q.db.QueryRowContext(ctx, addCostumer,
		arg.Name,
		arg.Cpf,
		arg.Birthdate,
		arg.Phone,
		arg.Email,
	)
	var i Costumer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Cpf,
		&i.Birthdate,
		&i.Phone,
		&i.Email,
		&i.PhoneValidated,
		&i.EmailValidated,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const addEvent = `-- name: AddEvent :exec
INSERT INTO outbox (event_name, message,created_at,processed_at) VALUES
($1,$2,NOW(),NULL)
`

type AddEventParams struct {
	EventName string
	Message   json.RawMessage
}

func (q *Queries) AddEvent(ctx context.Context, arg AddEventParams) error {
	_, err := q.db.ExecContext(ctx, addEvent, arg.EventName, arg.Message)
	return err
}

const deleteCostumer = `-- name: DeleteCostumer :exec
UPDATE costumers 
SET deleted = true
WHERE id = $1
`

func (q *Queries) DeleteCostumer(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteCostumer, id)
	return err
}

const getCostumer = `-- name: GetCostumer :one
SELECT id, name, cpf, birthdate, phone, email, phone_validated, email_validated, deleted, created_at, updated_at FROM costumers
WHERE (id = $1 AND deleted = false) 
LIMIT 1
`

func (q *Queries) GetCostumer(ctx context.Context, id int32) (Costumer, error) {
	row := q.db.QueryRowContext(ctx, getCostumer, id)
	var i Costumer
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Cpf,
		&i.Birthdate,
		&i.Phone,
		&i.Email,
		&i.PhoneValidated,
		&i.EmailValidated,
		&i.Deleted,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getEvents = `-- name: GetEvents :many
SELECT id, event_name, message, created_at, processed_at FROM outbox WHERE processed_at IS NULL LIMIT $1::int
`

func (q *Queries) GetEvents(ctx context.Context, length int32) ([]Outbox, error) {
	rows, err := q.db.QueryContext(ctx, getEvents, length)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Outbox
	for rows.Next() {
		var i Outbox
		if err := rows.Scan(
			&i.ID,
			&i.EventName,
			&i.Message,
			&i.CreatedAt,
			&i.ProcessedAt,
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

const processEvent = `-- name: ProcessEvent :exec
UPDATE outbox SET processed_at = NOW() WHERE id = $1
`

func (q *Queries) ProcessEvent(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, processEvent, id)
	return err
}

const searchCostumers = `-- name: SearchCostumers :many
SELECT id, name, cpf, birthdate, phone, email, phone_validated, email_validated, deleted, created_at, updated_at 
FROM costumers
WHERE 1 = 1
    AND (name ILIKE '%' || $1::text || '%' OR $1::text = '')
    AND (cpf ILIKE '%' || $2::text || '%' OR $2::text = '')
    AND (email ILIKE '%' || $3::text || '%' OR $3::text = '')
    AND (phone ILIKE '%' || $4::text || '%' OR $4::text = '')
    AND deleted = false
LIMIT $6::int OFFSET ($5::int - 1) * $6::int
`

type SearchCostumersParams struct {
	Name      string
	Cpf       string
	Email     string
	Phone     string
	Page      int32
	Pagecount int32
}

func (q *Queries) SearchCostumers(ctx context.Context, arg SearchCostumersParams) ([]Costumer, error) {
	rows, err := q.db.QueryContext(ctx, searchCostumers,
		arg.Name,
		arg.Cpf,
		arg.Email,
		arg.Phone,
		arg.Page,
		arg.Pagecount,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Costumer
	for rows.Next() {
		var i Costumer
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Cpf,
			&i.Birthdate,
			&i.Phone,
			&i.Email,
			&i.PhoneValidated,
			&i.EmailValidated,
			&i.Deleted,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updateCostumer = `-- name: UpdateCostumer :exec
UPDATE costumers 
SET name = $2, cpf = $3, birthdate =$4, email=$5, phone=$6, phone_validated=$7, email_validated=$8, updated_at=NOW()
WHERE id = $1 AND deleted = false
`

type UpdateCostumerParams struct {
	ID             int32
	Name           string
	Cpf            string
	Birthdate      time.Time
	Email          string
	Phone          string
	PhoneValidated bool
	EmailValidated bool
}

func (q *Queries) UpdateCostumer(ctx context.Context, arg UpdateCostumerParams) error {
	_, err := q.db.ExecContext(ctx, updateCostumer,
		arg.ID,
		arg.Name,
		arg.Cpf,
		arg.Birthdate,
		arg.Email,
		arg.Phone,
		arg.PhoneValidated,
		arg.EmailValidated,
	)
	return err
}
