package context

import (
	"context"
)

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("userID").(int)
	return userID, ok
}
