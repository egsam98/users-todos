package testutils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func NewRequestJSON(method, target string, jsonBody interface{}) *http.Request {
	var data []byte
	if jsonBody != nil {
		var err error
		data, err = json.Marshal(jsonBody)
		if err != nil {
			panic(err)
		}
	}
	return httptest.NewRequest(method, target, bytes.NewBuffer(data))
}

// Запустить тест - HTTP запрос
func RunHTTPTest(h http.Handler, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

// Поиск gin.HandlerFunc по методу и маршруту
// Используется для установления нового тестового маршрута во избежание миддлваров
func GinHandler(router *gin.Engine, method, route string) gin.HandlerFunc {
	for _, info := range router.Routes() {
		if info.Path == route && info.Method == method {
			return info.HandlerFunc
		}
	}
	panic("handler is not found")
}

// Преобразовать JSON-ответ сервера в obj
func DecodeBody(res *httptest.ResponseRecorder, obj interface{}) {
	if err := json.NewDecoder(res.Body).Decode(obj); err != nil {
		panic(err)
	}
}
