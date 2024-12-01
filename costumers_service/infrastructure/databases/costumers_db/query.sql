-- name: GetCostumer :one
SELECT * FROM costumers
WHERE (id = $1 AND deleted = false) 
LIMIT 1; 

-- name: AddCostumer :one
INSERT INTO costumers (name, cpf, birthdate, phone, email, phone_validated, email_validated, deleted,created_at,updated_at) VALUES
($1,$2,$3,$4,$5, false,false,false, NOW(), NULL) RETURNING *;

-- name: UpdateCostumer :exec
UPDATE costumers 
SET name = $2, cpf = $3, birthdate =$4, email=$5, phone=$6, phone_validated=$7, email_validated=$8, updated_at=NOW()
WHERE id = $1 AND deleted = false;

-- name: DeleteCostumer :exec
UPDATE costumers 
SET deleted = true
WHERE id = $1;

-- name: SearchCostumers :many
SELECT * 
FROM costumers
WHERE 1 = 1
    AND (name ILIKE '%' || sqlc.arg(name)::text || '%' OR sqlc.arg(name)::text = '')
    AND (cpf ILIKE '%' || sqlc.arg(cpf)::text || '%' OR sqlc.arg(cpf)::text = '')
    AND (email ILIKE '%' || sqlc.arg(email)::text || '%' OR sqlc.arg(email)::text = '')
    AND (phone ILIKE '%' || sqlc.arg(phone)::text || '%' OR sqlc.arg(phone)::text = '')
    AND deleted = false
LIMIT sqlc.arg(pageCount)::int OFFSET (sqlc.arg(page)::int - 1) * sqlc.arg(pageCount)::int;

-- name: AddEvent :exec
INSERT INTO outbox (event_name, message,created_at,processed_at) VALUES
($1,$2,NOW(),NULL);

-- name: GetEvents :many
SELECT * FROM outbox WHERE processed_at IS NULL LIMIT sqlc.arg(length)::int;

-- name: ProcessEvent :exec
UPDATE outbox SET processed_at = NOW() WHERE id = $1;