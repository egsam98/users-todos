package testmocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/egsam98/users-todos/users/db"
	"github.com/egsam98/users-todos/users/services"
)

var _ services.TokenService = (*TokenServiceMock)(nil)

type TokenServiceMock struct {
	mock.Mock
	q *db.Queries
}

func NewTokenServiceMock(q *db.Queries) *TokenServiceMock {
	return &TokenServiceMock{q: q}
}

func (tsm *TokenServiceMock) Generate(db.User) (string, error) {
	panic("implement me")
}

func (tsm *TokenServiceMock) Parse(_ context.Context, tokenString string) (*db.User, error) {
	args := tsm.Called(tokenString)
	return args.Get(0).(*db.User), args.Error(1)
}
