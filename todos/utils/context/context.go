package context

import (
	"context"
	"errors"
)

var ErrNoUserID = errors.New("userID does not exist in context")

func GetUserID(ctx context.Context) (int32, error) {
	userID, ok := ctx.Value("userID").(int32)
	if !ok {
		return 0, ErrNoUserID
	}
	return userID, nil
}
