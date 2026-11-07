package abstractions

import (
	"context"

	"example.com/open-monastery/internal/domain/models"
)

type UserRepository interface {
	FindById(ctx context.Context, id int64) (*models.User, error)
	Add(ctx context.Context, user *models.User) error
}
