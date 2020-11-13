package controllers

import (
	"github.com/gin-gonic/gin"

	"github.com/egsam98/users-todos/todos/db"
	"github.com/egsam98/users-todos/todos/services"
)

type TodosController struct {
	service *services.TodoService
}

func NewTodosController(q *db.Queries) *TodosController {
	return &TodosController{service: services.NewTodoService(q)}
}

func (tc *TodosController) CreateTodo(ctx *gin.Context) {
	tc.service.CreateTodo()
}
