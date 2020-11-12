package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/egsam98/users-todos/pkg/contract"
	"github.com/egsam98/users-todos/pkg/errors"
	"github.com/egsam98/users-todos/pkg/responses"
	"github.com/egsam98/users-todos/users/controllers/requests"
	"github.com/egsam98/users-todos/users/db"
	"github.com/egsam98/users-todos/users/services"
)

type UsersController struct {
	service *services.UserService
}

func NewUsersController(q *db.Queries) *UsersController {
	return &UsersController{service: services.NewUserService(q)}
}

// RegisterUser godoc
// @Summary Регистрация пользователя в системе
// @Tags users
// @Accept json
// @Param user body requests.RegisterUser true "Зарегистрировать пользователя"
// @Success 201
// @Failure 400 {object} responses.httpError
// @Router /users [post]
func (uc *UsersController) RegisterUser(ctx *gin.Context) {
	var req requests.RegisterUser
	errs := contract.Validate(ctx, &req)
	if len(errs) > 0 {
		responses.RespondError(ctx, http.StatusBadRequest, errs)
		return
	}

	if err := uc.service.Register(ctx, req); err != nil {
		if errors.IsPgError(err, errors.PgErrUniqueViolated) {
			responses.RespondError(ctx, 400, "user with this username already exists")
			return
		}
		responses.RespondInternalError(ctx, err)
	}

	ctx.Status(http.StatusCreated)
}
