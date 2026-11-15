package abstractions

import (
	"context"

	"example.com/open-monastery/internal/domain/models"
)

type AuthRepository interface {
	FindCredentialsByEmail(
		ctx context.Context,
		email string,
	) (userID int64, passwordHash string, userType models.UserType, err error)
	AddUser(ctx context.Context, user *models.User, passwordHash string) error
}
