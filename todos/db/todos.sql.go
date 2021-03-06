// Code generated by sqlc. DO NOT EDIT.
// source: todos.sql

package db

import (
	"context"
	"database/sql"
)

const createTodo = `-- name: CreateTodo :one
insert into todos (title, description, deadline, user_id) values ($1, $2, $3, $4) returning id, title, description, deadline, user_id
`

type CreateTodoParams struct {
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Deadline    sql.NullTime   `json:"deadline"`
	UserID      int32          `json:"user_id"`
}

func (q *Queries) CreateTodo(ctx context.Context, arg CreateTodoParams) (Todo, error) {
	row := q.db.QueryRowContext(ctx, createTodo,
		arg.Title,
		arg.Description,
		arg.Deadline,
		arg.UserID,
	)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Deadline,
		&i.UserID,
	)
	return i, err
}

const deleteTodo = `-- name: DeleteTodo :execrows
delete from todos where id = $1 and user_id = $2
`

type DeleteTodoParams struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"`
}

func (q *Queries) DeleteTodo(ctx context.Context, arg DeleteTodoParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteTodo, arg.ID, arg.UserID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const findAll = `-- name: FindAll :many
select id, title, description, deadline, user_id from todos where user_id = $1 order by deadline
`

func (q *Queries) FindAll(ctx context.Context, userID int32) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, findAll, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Deadline,
			&i.UserID,
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

const findBeforeDeadline = `-- name: FindBeforeDeadline :many
select id, title, description, deadline, user_id from todos where deadline < $1 and user_id = $2 order by deadline
`

type FindBeforeDeadlineParams struct {
	Deadline sql.NullTime `json:"deadline"`
	UserID   int32        `json:"user_id"`
}

func (q *Queries) FindBeforeDeadline(ctx context.Context, arg FindBeforeDeadlineParams) ([]Todo, error) {
	rows, err := q.db.QueryContext(ctx, findBeforeDeadline, arg.Deadline, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Todo
	for rows.Next() {
		var i Todo
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Deadline,
			&i.UserID,
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

const findTodoById = `-- name: FindTodoById :one
select id, title, description, deadline, user_id from todos where id = $1
`

func (q *Queries) FindTodoById(ctx context.Context, id int32) (Todo, error) {
	row := q.db.QueryRowContext(ctx, findTodoById, id)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Deadline,
		&i.UserID,
	)
	return i, err
}

const updateTodo = `-- name: UpdateTodo :one
update todos set title=$1, description=$2, deadline=$3 where id=$4 returning id, title, description, deadline, user_id
`

type UpdateTodoParams struct {
	Title       string         `json:"title"`
	Description sql.NullString `json:"description"`
	Deadline    sql.NullTime   `json:"deadline"`
	ID          int32          `json:"id"`
}

func (q *Queries) UpdateTodo(ctx context.Context, arg UpdateTodoParams) (Todo, error) {
	row := q.db.QueryRowContext(ctx, updateTodo,
		arg.Title,
		arg.Description,
		arg.Deadline,
		arg.ID,
	)
	var i Todo
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Deadline,
		&i.UserID,
	)
	return i, err
}
