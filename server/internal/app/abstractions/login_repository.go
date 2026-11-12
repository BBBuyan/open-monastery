package abstractions

import "context"

type LoginRepository interface {
	FindPasswordHash(ctx context.Context, userID int64) (string, error)
}
