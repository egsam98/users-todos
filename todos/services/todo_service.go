package services

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	time2 "github.com/egsam98/users-todos/pkg/time"
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
		deadline.Time = time2.UtcFromUnix(*req.Deadline)
		deadline.Valid = true
	}

	userID, err := context2.GetUserID(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	todo, err := ts.q.CreateTodo(ctx, db.CreateTodoParams{
		Title:       req.Title,
		Description: desc,
		Deadline:    deadline,
		UserID:      userID,
	})
	return &todo, errors.WithStack(err)
}

// Обновить текущую задачу с опр. id новыми значения из requests.NewTodo
func (ts *TodoService) UpdateTodo(ctx context.Context, id int, req requests.NewTodo) (*db.Todo, error) {
	todo, err := ts.q.FindTodoById(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.WithStack(ErrNoTodoFound)
		}
		return nil, errors.WithStack(err)
	}

	currentUserID, err := context2.GetUserID(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if todo.UserID != currentUserID {
		return nil, ErrNoAccessToTodo
	}

	params := db.UpdateTodoParams{ID: int32(id), Title: req.Title}

	if req.Description != nil {
		params.Description.String = *req.Description
		params.Description.Valid = true
	} else {
		params.Description.Valid = false
	}
	if req.Deadline != nil {
		params.Deadline.Time = time2.UtcFromUnix(*req.Deadline)
		params.Deadline.Valid = true
	} else {
		params.Deadline.Valid = false
	}

	todo, err = ts.q.UpdateTodo(ctx, params)
	return &todo, errors.WithStack(err)
}

// Удалить задачу
func (ts *TodoService) DeleteTodo(ctx context.Context, id int) error {
	userID, err := context2.GetUserID(ctx)
	if err != nil {
		return errors.WithStack(err)
	}
	rowsAffected, err := ts.q.DeleteTodo(ctx, db.DeleteTodoParams{
		ID:     int32(id),
		UserID: userID,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	if rowsAffected == 0 {
		return errors.WithStack(ErrNoTodoFound)
	}

	return nil
}

// Все задачи пользователя, отсортированные в порядке возрастания deadline
func (ts *TodoService) All(ctx context.Context) ([]db.Todo, error) {
	userID, err := context2.GetUserID(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return ts.q.FindAll(ctx, userID)
}

// Задачи, deadline которых меньше переданного значения. Задачи отсортированы в порядке возрастания deadline
func (ts *TodoService) BeforeDeadline(ctx context.Context, req requests.Deadline) ([]db.Todo, error) {
	userID, err := context2.GetUserID(ctx)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return ts.q.FindBeforeDeadline(ctx, db.FindBeforeDeadlineParams{
		Deadline: sql.NullTime{Time: time2.UtcFromUnix(req.Value), Valid: true},
		UserID:   userID,
	})
}
