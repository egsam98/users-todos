package requests

// Тело запроса на создание новой задачи
type NewTodo struct {
	Title       string  `json:"title" binding:"required" error:"must be non empty"`
	Description *string `json:"description"`
	Deadline    *int64  `json:"deadline" swaggertype:"integer" format:"int64"`
}

// Тело запроса для фильтрации задач < deadline
type Deadline struct {
	Value int64 `json:"deadline" format:"int64" binding:"required" error:"must be unix time"`
}
