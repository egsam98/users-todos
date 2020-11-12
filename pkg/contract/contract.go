package contract

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Валидация входных данных в контроллерах
// Используется custom-тэг "error" для большей читаемости JSON-ответа сервера с ошибками валидации
func Validate(ctx *gin.Context, obj interface{}) gin.H {
	h := gin.H{}

	err := ctx.ShouldBindJSON(obj)
	if err == nil {
		return h
	}

	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		panic(err)
	}

	t := reflect.TypeOf(obj).Elem()

	for _, err := range errs {
		field, _ := t.FieldByName(err.Field())
		h[field.Tag.Get("json")] = field.Tag.Get("error")
	}
	return h
}
