package services

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"hash"

	"github.com/pkg/errors"

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
func (us *UserService) Register(ctx context.Context, req requests.Signup) error {
	_, err := us.q.CreateUser(ctx, db.CreateUserParams{
		Username: req.Username,
		Password: us.hashPassword(req.Password),
	})
	return errors.WithStack(err)
}

// Аутентификация пользователя по username и password
func (us *UserService) Authenticate(ctx context.Context, req requests.Signin) (db.User, error) {
	user, err := us.q.FindUser(ctx, db.FindUserParams{
		Username: req.Username,
		Password: us.hashPassword(req.Password),
	})
	return user, errors.WithStack(err)
}

// Поиск пользователя в БД по его ID
func (us *UserService) FindUser(ctx context.Context, id int) (*db.User, error) {
	user, err := us.q.FindUserById(ctx, int32(id))
	return &user, errors.WithStack(err)
}

// Хэширование пароля (алгоритм SHA1)
func (us *UserService) hashPassword(password string) string {
	_, _ = us.hash.Write([]byte(password))
	encoded := hex.EncodeToString(us.hash.Sum(nil))
	us.hash.Reset()
	return encoded
}
