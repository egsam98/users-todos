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
	"github.com/egsam98/users-todos/users/utils/env"
)

type UsersController struct {
	userService *services.UserService
	jwtService  *services.JwtService
}

func NewUsersController(environment env.Environment, q *db.Queries) *UsersController {
	return &UsersController{userService: services.NewUserService(q), jwtService: services.NewJwtService(environment, q)}
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
	errs, ok := contract.ValidateJSON(ctx, &req)
	if !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	if err := uc.userService.Register(ctx, req); err != nil {
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
	errs, ok := contract.ValidateJSON(ctx, &req)
	if !ok {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	user, err := uc.userService.Authenticate(ctx, req)
	if err != nil {
		if errors2.Cause(err) == sql.ErrNoRows {
			responses.RespondError(ctx, http.StatusUnauthorized, "username or/and password is incorrect")
		} else {
			responses.RespondInternalError(ctx, err)
		}
		return
	}

	token, err := uc.jwtService.Generate(user)
	if err != nil {
		responses.RespondInternalError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, responses2.Token{Value: token})
}

// FetchUser godoc
// @Summary Запрос пользователя в системе по ID
// @Tags users
// @Accept json
// @Param Authorization header string true "JWT-токен"
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
	user, err := uc.userService.FindUser(ctx, id)
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

// Auth godoc
// @Summary Аутентификация пользователя по JWT-токену
// @Tags auth
// @Param Authorization header string true "JWT-токен"
// @Success 200 {object} responses.User "Текущий пользователь в системе, определенный по JWT-токену"
// @Failure 401 {object} responses.httpError
// @Router /auth [post]
func (uc *UsersController) Auth(isMiddleware bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			responses.RespondError(ctx, http.StatusForbidden, "Authorization header is not provided")
			return
		}

		user, err := uc.jwtService.Parse(ctx, authHeader)
		if err != nil {
			responses.RespondError(ctx, http.StatusForbidden, gin.H{"jwt": err.Error()})
			return
		}

		if isMiddleware {
			ctx.Set("user", user)
			return
		}
		ctx.JSON(200, responses2.User{
			ID:       user.ID,
			Username: user.Username,
		})
	}
}
