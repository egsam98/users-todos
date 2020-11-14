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

// Запустить тест - HTTP запрос к gin-роутеру
func RunGinTest(r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}
