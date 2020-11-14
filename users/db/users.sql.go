// Code generated by sqlc. DO NOT EDIT.
// source: users.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
insert into users (username, password) values ($1, $2) returning id, username, password
`

type CreateUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Username, arg.Password)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}

const findUser = `-- name: FindUser :one
select id, username, password from users where username = $1 and password = $2 limit 1
`

type FindUserParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (q *Queries) FindUser(ctx context.Context, arg FindUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, findUser, arg.Username, arg.Password)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}

const findUserById = `-- name: FindUserById :one
select id, username, password from users where id = $1
`

func (q *Queries) FindUserById(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRowContext(ctx, findUserById, id)
	var i User
	err := row.Scan(&i.ID, &i.Username, &i.Password)
	return i, err
}
