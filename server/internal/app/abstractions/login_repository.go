package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type AuthRepository interface {
	FindCredentialsByEmail(
		ctx context.Context,
		email string,
	) (*appmodels.UserCreds, error)
	AddUser(ctx context.Context, user *models.User) error
	PromoteToMonk(ctx context.Context, userId, monasteryId int64) error
	PromoteToManager(ctx context.Context, userId, monasteryId int64) error
}
