package abstractions

import "context"

type LoginRepository interface {
	FindCredentialsByEmail(ctx context.Context, email string) (userID int64, passwordHash string, err error)
}
