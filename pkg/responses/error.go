package responses

import (
	"log"

	"github.com/gin-gonic/gin"
)

type httpError struct {
	Error interface{} `json:"error"`
}

// Ответ сервера с ошибкой и статусом code
func RespondError(ctx *gin.Context, code int, value interface{}) {
	if err, ok := value.(error); ok {
		value = err.Error()
	}
	ctx.JSON(code, httpError{Error: value})
}

// Для ошибок со статусом 500 (предусматривается логирование)
func RespondInternalError(ctx *gin.Context, err error) {
	log.Printf("%+v\n", err)
	ctx.Status(500)
}
