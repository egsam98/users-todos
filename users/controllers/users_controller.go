package controllers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	errors2 "github.com/pkg/errors"

	"github.com/egsam98/users-todos/pkg/contract"
	"github.com/egsam98/users-todos/pkg/errors"
	"github.com/egsam98/users-todos/pkg/responses"
	"github.com/egsam98/users-todos/users/controllers/requests"
	responses2 "github.com/egsam98/users-todos/users/controllers/responses"
	"github.com/egsam98/users-todos/users/db"
	"github.com/egsam98/users-todos/users/services"
)

type UsersController struct {
	service *services.UserService
}

func NewUsersController(q *db.Queries) *UsersController {
	return &UsersController{service: services.NewUserService(q)}
}

// Signup godoc
// @Summary Регистрация пользователя в системе
// @Tags users
// @Accept json
// @Param user body requests.Signup true "Зарегистрировать пользователя"
// @Success 201
// @Failure 400 {object} responses.httpError
// @Router /signup [post]
func (uc *UsersController) Signup(ctx *gin.Context) {
	var req requests.Signup
	errs, ok := contract.Validate(ctx, &req)
	if !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	if err := uc.service.Register(ctx, req); err != nil {
		if errors.IsPgError(err, errors.PgErrUniqueViolated) {
			responses.RespondError(ctx, http.StatusBadRequest, "user with this username already exists")
			return
		}
		responses.RespondInternalError(ctx, err)
	}

	ctx.Status(http.StatusCreated)
}

// Signin godoc
// @Summary Вход в систему
// @Tags users
// @Accept json
// @Param user body requests.Signin true "Зарегистрировать пользователя"
// @Success 200 {object} responses.Token
// @Failure 400 {object} responses.httpError
// @Failure 401 {object} responses.httpError
// @Router /signin [post]
func (uc *UsersController) Signin(ctx *gin.Context) {
	var req requests.Signin
	errs, ok := contract.Validate(ctx, &req)
	if !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	token, err := uc.service.Login(ctx, req)
	if err != nil {
		if errors2.Cause(err) == sql.ErrNoRows {
			responses.RespondError(ctx, http.StatusUnauthorized, "username or/and password is incorrect")
		} else {
			responses.RespondInternalError(ctx, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, responses2.Token{Value: token})
}

// FetchUser godoc
// @Summary Запрос пользователя в системе по ID
// @Tags users
// @Accept json
// @Param id path int true "ID пользователя"
// @Success 200 {object} responses.User
// @Failure 400 {object} responses.httpError
// @Failure 403 {object} responses.httpError
// @Failure 404 {object} responses.httpError
// @Router /users/:id [get]
func (uc *UsersController) FetchUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		responses.RespondError(ctx, http.StatusBadRequest, "user ID must be integer")
		return
	}
	user, err := uc.service.FindUser(ctx, id)
	if err != nil {
		if errors2.Cause(err) == sql.ErrNoRows {
			responses.RespondError(ctx, http.StatusNotFound, fmt.Sprintf("user ID=%d is not found", id))
		} else {
			responses.RespondInternalError(ctx, err)
		}
		return
	}

	ctx.JSON(200, responses2.User{
		ID:       user.ID,
		Username: user.Username,
	})
}
