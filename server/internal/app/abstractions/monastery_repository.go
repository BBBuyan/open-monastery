package abstractions

import (
	"context"

	"example.com/open-monastery/internal/domain/models"
)

type MonasteryRepository interface {
	FindById(ctx context.Context, id int64) (*models.Monastery, error)
	Add(ctx context.Context, monastery *models.Monastery) error
}
