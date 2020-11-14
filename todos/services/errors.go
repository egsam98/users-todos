package services

import (
	"errors"
)

var (
	ErrNoAccessToTodo = errors.New("you have no access to this todo")
	ErrNoTodoFound    = errors.New("todo is not found")
)
