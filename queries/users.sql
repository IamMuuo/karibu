-- name: GetAllUsers :many
select *
from users
limit $1
offset $2
;

-- name: CreateUser :one
INSERT INTO users (
  firstname, othernames, email, organization,
  role, phone, ssh_key, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW()) 
RETURNING *;

