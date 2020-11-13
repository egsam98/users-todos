package services

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"

	"github.com/egsam98/users-todos/users/db"
	"github.com/egsam98/users-todos/users/utils/env"
)

const tokenExpiresIn = 5 * 24 * time.Hour

type JwtService struct {
	environment env.Environment
	q           *db.Queries
}

func NewJwtService(environment env.Environment, q *db.Queries) *JwtService {
	return &JwtService{environment: environment, q: q}
}

func (js *JwtService) Generate(user db.User) (string, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().UTC().Add(tokenExpiresIn).Unix(),
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(js.environment.Signature))
	return tokenString, errors.WithStack(err)
}

func (js *JwtService) Parse(ctx context.Context, tokenString string) (*db.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong signing method")
		}
		return []byte(js.environment.Signature), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	if err := claims.Valid(); err != nil {
		return nil, err
		//err := err.(*jwt.ValidationError)
	}

	user, err := js.q.FindUserById(ctx, int32(claims["sub"].(float64))) // TODO: check casting
	return &user, errors.WithStack(err)
}
