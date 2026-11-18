package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingRepository interface {
	Add(ctx context.Context, reading *models.Reading) error
	Claim(ctx context.Context, monasteryID, readingID, monkID int64) error
	Assign(ctx context.Context, monasteryID, readingID, monkID int64) error
	GetReadingsByUserId(
		ctx context.Context,
		monasteryID int64,
		userID int64,
		opts appmodels.ListOptions,
	) ([]*models.Reading, error)
}
