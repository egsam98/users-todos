package errors

import (
	"github.com/lib/pq"
)

const PgErrUniqueViolated pq.ErrorCode = "23505"

// Проверка ошибки PostgreSQL по коду
func IsPgError(err error, code pq.ErrorCode) bool {
	if err, ok := err.(*pq.Error); ok {
		return code == err.Code
	}
	return false
}
