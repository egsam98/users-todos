package services

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/egsam98/users-todos/todos/controllers/requests"
	"github.com/egsam98/users-todos/todos/db"
	context2 "github.com/egsam98/users-todos/todos/utils/context"
)

// Сервис управления задачами
type TodoService struct {
	q *db.Queries
}

func NewTodoService(q *db.Queries) *TodoService {
	return &TodoService{q: q}
}

// Создать новую задачу
func (ts *TodoService) CreateTodo(ctx context.Context, req requests.NewTodo) (*db.Todo, error) {
	desc := sql.NullString{}
	if req.Description != nil {
		desc.String = *req.Description
		desc.Valid = true
	}

	deadline := sql.NullTime{}
	if req.Deadline != nil {
		deadline.Time = *req.Deadline
		deadline.Valid = true
	}

	userID, ok := context2.GetUserID(ctx)
	if !ok {
		return nil, errors.New("userID doesn't exist in context")
	}

	todo, err := ts.q.CreateTodo(ctx, db.CreateTodoParams{
		Title:       req.Title,
		Description: desc,
		Deadline:    deadline,
		UserID:      int32(userID),
	})
	return &todo, errors.WithStack(err)
}
