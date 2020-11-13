package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/egsam98/users-todos/pkg/contract"
	"github.com/egsam98/users-todos/pkg/responses"
	"github.com/egsam98/users-todos/todos/controllers/requests"
	responses2 "github.com/egsam98/users-todos/todos/controllers/responses"
	"github.com/egsam98/users-todos/todos/db"
	"github.com/egsam98/users-todos/todos/services"
)

type TodosController struct {
	service *services.TodoService
}

func NewTodosController(q *db.Queries) *TodosController {
	return &TodosController{service: services.NewTodoService(q)}
}

// CreateTodo godoc
// @Summary Создать новую задачу
// @Tags auth
// @Param Authorization header string true "JWT-токен"
// @Param todo body requests.NewTodo true "Новая задача"
// @Success 201 {object} responses.Todo
// @Failure 400 {object} responses.httpError
// @Router /todos [post]
func (tc *TodosController) CreateTodo(ctx *gin.Context) {
	var req requests.NewTodo
	if errs, ok := contract.ValidateJSON(ctx, &req); !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	todo, err := tc.service.CreateTodo(ctx, req)
	if err != nil {
		responses.RespondInternalError(ctx, err)
		return
	}

	var desc *string
	if todo.Description.Valid {
		desc = &todo.Description.String
	}

	var deadline *time.Time
	if todo.Deadline.Valid {
		deadline = &todo.Deadline.Time
	}

	ctx.JSON(201, responses2.Todo{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: desc,
		Deadline:    deadline,
		UserID:      todo.UserID,
	})
}
