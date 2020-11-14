package responses

import (
	"github.com/egsam98/users-todos/todos/db"
)

// Ответ сервера с задачей в теле
type Todo struct {
	ID          int32   `json:"id"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	Deadline    *int64  `json:"deadline"`
	UserID      int32   `json:"user_id"`
}

func NewTodo(todo db.Todo) Todo {
	var desc *string
	if todo.Description.Valid {
		desc = &todo.Description.String
	}

	var deadline *int64
	if todo.Deadline.Valid {
		deadline = new(int64)
		*deadline = todo.Deadline.Time.Unix()
	}

	return Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: desc,
		Deadline:    deadline,
		UserID:      todo.UserID,
	}
}
