package middlewares

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/egsam98/users-todos/pkg/responses"
)

// Миддлвар авторизации через сторонний сервис
// В случае успеха gin.Context добавляет ID пользователя по ключу "userID"
func CheckAuth(authUrl string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req, _ := http.NewRequestWithContext(ctx, "POST", authUrl, nil)
		req.Header.Set("Authorization", ctx.GetHeader("Authorization"))

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			responses.RespondInternalError(ctx, err)
			return
		}
		defer res.Body.Close()

		jsonDecoder := json.NewDecoder(res.Body)

		switch res.StatusCode {
		case http.StatusOK:
			userID := struct {
				Value int `json:"id"`
			}{}
			if err := jsonDecoder.Decode(&userID); err != nil {
				responses.RespondInternalError(ctx, err)
				return
			}
			ctx.Set("userID", userID.Value)
		case http.StatusForbidden:
			var body map[string]interface{}
			if err := jsonDecoder.Decode(&body); err != nil {
				responses.RespondInternalError(ctx, err)
				return
			}
			ctx.AbortWithStatusJSON(http.StatusForbidden, body)
		default:
			responses.RespondInternalError(ctx, err)
		}
	}
}
