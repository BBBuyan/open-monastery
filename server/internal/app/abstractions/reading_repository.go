package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingRepository interface {
	Add(ctx context.Context, monastery *models.Reading) error
	Claim(ctx context.Context, readingID, monkID int64) error
	UserExists(ctx context.Context, id int64) (bool, error)
	BookExists(ctx context.Context, id int64) (bool, error)
	MonasteryExists(ctx context.Context, id int64) (bool, error)
	GetReadingsByUserId(
		ctx context.Context,
		userID int64,
		opts appmodels.ListOptions,
	) ([]*models.Reading, error)
}
