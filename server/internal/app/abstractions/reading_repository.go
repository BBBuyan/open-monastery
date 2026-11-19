package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingRepository interface {
	Add(ctx context.Context, reading *models.Reading) error
	GetReadingsByUserId(
		ctx context.Context,
		userID int64,
		opts appmodels.ListOptions,
	) ([]*models.Reading, error)
}
