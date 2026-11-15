package abstractions

import (
	"context"

	"example.com/open-monastery/internal/domain/models"
)

type LoginRepository interface {
	FindCredentialsByEmail(ctx context.Context, email string) (userID int64, passwordHash string, userType models.UserType, err error)
}
