package abstractions

import (
	"context"

	"example.com/open-monastery/internal/domain/enums"
	"example.com/open-monastery/internal/domain/models"
)

type AuthRepository interface {
	FindCredentialsByEmail(
		ctx context.Context,
		email string,
	) (userID int64, passwordHash string, userType enums.UserType, monasteryId *int64, err error)
	AddUser(ctx context.Context, user *models.User) error
	PromoteToMonk(ctx context.Context, userId, monasteryId int64) error
	PromoteToManager(ctx context.Context, userId, monasteryId int64) error
}
