package services

import (
	"github.com/egsam98/users-todos/todos/db"
)

// Сервис управления задачами
type TodoService struct {
	q *db.Queries
}

func NewTodoService(q *db.Queries) *TodoService {
	return &TodoService{q: q}
}

// Создать новую задачу
func (ts *TodoService) CreateTodo() {
}
