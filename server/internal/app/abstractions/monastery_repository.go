package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type MonasteryRepository interface {
	FindById(ctx context.Context, id int64) (*models.Monastery, error)
	Add(ctx context.Context, monastery *models.Monastery) error
	GetMonasteries(ctx context.Context, opts appmodels.ListOptions) ([]*models.Monastery, error)
	CanWrite(ctx context.Context, monasteryID, userID int64) (bool, error)
	GrantManager(ctx context.Context, monasteryID, userID int64) error
	RevokeManager(ctx context.Context, monasteryID, userID int64) error
}
