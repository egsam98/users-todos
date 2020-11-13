package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/egsam98/users-todos/pkg/responses"
	"github.com/egsam98/users-todos/users/db"
	"github.com/egsam98/users-todos/users/services"
	"github.com/egsam98/users-todos/users/utils/env"
)

// Миддлвар для проверки присутствия и валидности JWT-токена
type JwtMiddleware struct {
	service *services.JwtService
}

func NewJwtMiddleware(environment env.Environment, q *db.Queries) *JwtMiddleware {
	return &JwtMiddleware{service: services.NewJwtService(environment, q)}
}

// Функция вызова миддлвара
func (jm *JwtMiddleware) Process(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		responses.RespondError(ctx, http.StatusForbidden, "Authorization header is not provided")
		return
	}

	user, err := jm.service.Parse(ctx, authHeader)
	if err != nil {
		responses.RespondError(ctx, http.StatusForbidden, gin.H{"jwt": err.Error()})
		return
	}

	ctx.Set("user", user)
}
