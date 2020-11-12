package services

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"hash"

	"github.com/egsam98/users-todos/users/controllers/requests"
	"github.com/egsam98/users-todos/users/db"
)

// Сервис взаимодействия с пользователями
type UserService struct {
	q    *db.Queries
	hash hash.Hash
}

func NewUserService(q *db.Queries) *UserService {
	return &UserService{q: q, hash: sha1.New()}
}

// Зарегистрировать пользователя в системе
func (us *UserService) Register(ctx context.Context, req requests.RegisterUser) error {
	_, _ = us.hash.Write([]byte(req.Password))
	encodedPassword := hex.EncodeToString(us.hash.Sum(nil))
	return us.q.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Password: encodedPassword,
	})
}
