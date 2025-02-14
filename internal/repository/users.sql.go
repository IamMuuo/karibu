// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (
  firstname, othernames, email, organization,
  role, phone, ssh_key, created_at, updated_at
) VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW()) 
RETURNING id, firstname, othernames, email, organization, role, phone, ssh_key, created_at, updated_at
`

type CreateUserParams struct {
	Firstname    string      `json:"firstname"`
	Othernames   string      `json:"othernames"`
	Email        string      `json:"email"`
	Organization pgtype.Text `json:"organization"`
	Role         pgtype.Text `json:"role"`
	Phone        pgtype.Text `json:"phone"`
	SshKey       pgtype.Text `json:"ssh_key"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Firstname,
		arg.Othernames,
		arg.Email,
		arg.Organization,
		arg.Role,
		arg.Phone,
		arg.SshKey,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Firstname,
		&i.Othernames,
		&i.Email,
		&i.Organization,
		&i.Role,
		&i.Phone,
		&i.SshKey,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllUsers = `-- name: GetAllUsers :many
select id, firstname, othernames, email, organization, role, phone, ssh_key, created_at, updated_at
from users
limit $1
offset $2
`

type GetAllUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, getAllUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Firstname,
			&i.Othernames,
			&i.Email,
			&i.Organization,
			&i.Role,
			&i.Phone,
			&i.SshKey,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
