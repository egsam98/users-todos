package contract

import (
	"encoding/json"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Валидация входных данных в формате JSON в контроллерах
// Используется custom-тэг "error" для большей читаемости JSON-ответа сервера с ошибками валидации
// Возвращает map ошибок валидация и boolean - является объект валидным ?
func ValidateJSON(ctx *gin.Context, obj interface{}) (gin.H, bool) {
	if reflect.ValueOf(obj).Kind() != reflect.Ptr {
		panic("object must be a pointer")
	}

	err := ctx.ShouldBindBodyWith(obj, binding.JSON)
	if err == nil {
		return nil, true
	}

	h := gin.H{}

	switch err := err.(type) {
	case validator.ValidationErrors:
		t := reflect.TypeOf(obj).Elem()
		for _, err := range err {
			field, _ := t.FieldByName(err.Field())
			h[field.Tag.Get("json")] = field.Tag.Get("error")
		}
	case *json.UnmarshalTypeError:
		if err.Field != "" {
			h[err.Field] = "must be " + err.Type.String()
		}
	}

	// Дефолтная ошибка
	if len(h) == 0 {
		h["body"] = err.Error()
	}
	return h, false
}
