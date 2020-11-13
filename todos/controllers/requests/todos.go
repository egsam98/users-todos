package requests

import (
	"time"
)

// Тело запроса на создание новой задачи
type NewTodo struct {
	Title       string     `json:"title" binding:"required" error:"must be non empty"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}
