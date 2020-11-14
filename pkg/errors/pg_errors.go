package errors

import (
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	PgErrUniqueViolated        pq.ErrorCode = "23505"
	PgErrDatetimeFieldOverflow pq.ErrorCode = "22008"
)

// Проверка ошибки PostgreSQL по коду
func IsPgError(err error, code pq.ErrorCode) bool {
	if err, ok := errors.Cause(err).(*pq.Error); ok {
		return code == err.Code
	}
	return false
}
