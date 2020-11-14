package testutils

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

// Проверка, что маршрут требует JWT-авторизацию
func JwtAuthRequired(t *testing.T, h http.Handler, method, route string) {
	res := RunHTTPTest(h, httptest.NewRequest(method, route, nil))
	assert.Equal(t, http.StatusForbidden, res.Code)
}
