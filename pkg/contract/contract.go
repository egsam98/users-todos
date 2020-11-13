package contract

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Валидация входных данных в контроллерах
// Используется custom-тэг "error" для большей читаемости JSON-ответа сервера с ошибками валидации
// Возвращает map ошибок валидация и boolean - является объект валидным ?
func Validate(ctx *gin.Context, obj interface{}) (gin.H, bool) {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		panic("object must be a pointer")
	}

	err := ctx.ShouldBindJSON(obj)
	if err == nil {
		return nil, true
	}

	h := gin.H{}

	if errs, ok := err.(validator.ValidationErrors); ok {
		t := reflect.TypeOf(obj).Elem()
		for _, err := range errs {
			field, _ := t.FieldByName(err.Field())
			h[field.Tag.Get("json")] = field.Tag.Get("error")
		}
	} else {
		h["body"] = "request body must be JSON"
	}
	return h, false
}
