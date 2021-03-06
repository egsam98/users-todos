package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"

	"github.com/egsam98/users-todos/pkg/contract"
	errors2 "github.com/egsam98/users-todos/pkg/errors"
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
// @Tags todos
// @Param Authorization header string true "JWT-токен"
// @Param todo body requests.NewTodo true "Новая задача"
// @Success 201 {object} responses.Todo
// @Failure 400 {object} responses.httpError
// @Failure 401 {object} responses.httpError
// @Router /todos [post]
func (tc *TodosController) CreateTodo(ctx *gin.Context) {
	var req requests.NewTodo
	if errs, ok := contract.ValidateJSON(ctx, &req); !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	todo, err := tc.service.CreateTodo(ctx, req)
	if err != nil {
		if errors2.IsPgError(err, errors2.PgErrDatetimeFieldOverflow) {
			responses.RespondError(ctx, http.StatusBadRequest, err)
		} else {
			responses.RespondInternalError(ctx, err)
		}
		return
	}

	ctx.JSON(201, responses2.NewTodo(*todo))
}

// UpdateTodo godoc
// @Summary Обновить существующую задачу
// @Tags todos
// @Param Authorization header string true "JWT-токен"
// @Param id path int true "ID задачи"
// @Param todo body requests.NewTodo true "Новые значения задачи"
// @Success 200 {object} responses.Todo
// @Failure 400 {object} responses.httpError
// @Failure 401 {object} responses.httpError
// @Failure 403 {object} responses.httpError
// @Failure 404 {object} responses.httpError
// @Router /todos/:id [put]
func (tc *TodosController) UpdateTodo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		responses.RespondError(ctx, http.StatusBadRequest, "id must be integer")
		return
	}

	var req requests.NewTodo
	if errs, ok := contract.ValidateJSON(ctx, &req); !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	if err := tc.validateRequiredKeysForUpdate(ctx); err != nil {
		responses.RespondError(ctx, http.StatusBadRequest, err)
		return
	}

	todo, err := tc.service.UpdateTodo(ctx, id, req)
	if err != nil {
		switch cause := errors.Cause(err); cause {
		case services.ErrNoTodoFound:
			responses.RespondError(ctx, http.StatusNotFound, cause)
		case services.ErrNoAccessToTodo:
			responses.RespondError(ctx, http.StatusForbidden, cause)
		default:
			if errors2.IsPgError(cause, errors2.PgErrDatetimeFieldOverflow) {
				responses.RespondError(ctx, http.StatusBadRequest, cause)
			} else {
				responses.RespondInternalError(ctx, cause)
			}
		}
		return
	}

	ctx.JSON(200, responses2.NewTodo(*todo))
}

// DeleteTodo godoc
// @Summary Удалить существующую задачу
// @Tags todos
// @Param Authorization header string true "JWT-токен"
// @Param id path int true "ID задачи"
// @Success 200
// @Failure 400 {object} responses.httpError
// @Failure 401 {object} responses.httpError
// @Failure 404 {object} responses.httpError
// @Router /todos/:id [delete]
func (tc *TodosController) DeleteTodo(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		responses.RespondError(ctx, http.StatusBadRequest, "id must be integer")
		return
	}

	if err := tc.service.DeleteTodo(ctx, id); err != nil {
		if errors.Cause(err) == services.ErrNoTodoFound {
			responses.RespondError(ctx, http.StatusNotFound, err)
		} else {
			responses.RespondInternalError(ctx, err)
		}
	}
}

// FetchAll godoc
// @Summary Все задачи пользователя
// @Tags todos
// @Param Authorization header string true "JWT-токен"
// @Success 200 {array} responses.Todo
// @Failure 401 {object} responses.httpError
// @Router /todos [get]
func (tc *TodosController) FetchAll(ctx *gin.Context) {
	todos, err := tc.service.All(ctx)
	if err != nil {
		responses.RespondInternalError(ctx, err)
		return
	}
	ctx.JSON(200, responses2.NewTodos(todos))
}

// FetchBeforeDeadline godoc
// @Summary Все задачи пользователя до определенного времени deadline
// @Tags todos
// @Param Authorization header string true "JWT-токен"
// @Param deadline body requests.Deadline true "Deadline"
// @Success 200 {array} responses.Todo
// @Failure 400 {object} responses.httpError
// @Failure 401 {object} responses.httpError
// @Router /todos/before [post]
func (tc *TodosController) FetchBeforeDeadline(ctx *gin.Context) {
	var req requests.Deadline
	if errs, ok := contract.ValidateJSON(ctx, &req); !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	todos, err := tc.service.BeforeDeadline(ctx, req)
	if err != nil {
		if errors2.IsPgError(err, errors2.PgErrDatetimeFieldOverflow) {
			responses.RespondError(ctx, http.StatusBadRequest, err)
		} else {
			responses.RespondInternalError(ctx, err)
		}
		return
	}

	ctx.JSON(200, responses2.NewTodos(todos))
}

// Проверка наличия ключей "description" и "deadline" в JSON-запросе для PUT /todos/:id
func (_ *TodosController) validateRequiredKeysForUpdate(ctx *gin.Context) error {
	reqMap := map[string]interface{}{}
	if err := ctx.ShouldBindBodyWith(&reqMap, binding.JSON); err != nil {
		panic(err)
	}

	missingKeys := make([]string, 0)
	for _, key := range []string{"description", "deadline"} {
		if _, ok := reqMap[key]; !ok {
			missingKeys = append(missingKeys, "'"+key+"'")
		}
	}

	if len(missingKeys) > 0 {
		return fmt.Errorf("%s must present", strings.Join(missingKeys, ", "))
	}

	return nil
}
