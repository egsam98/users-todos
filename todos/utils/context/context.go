package context

import (
	"context"
	"errors"
)

const userIDKey = "userID"

func GetUserID(ctx context.Context) (int32, error) {
	i := ctx.Value(userIDKey)
	if i == nil {
		return 0, errors.New("key '" + userIDKey + "' is not found")
	}
	userID, ok := i.(int32)
	if !ok {
		return 0, errors.New("must be int32")
	}
	return userID, nil
}
