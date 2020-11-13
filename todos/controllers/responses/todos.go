package responses

import (
	"time"
)

// Ответ сервера с задачей в теле
type Todo struct {
	ID          int32      `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	UserID      int32      `json:"user_id"`
}
