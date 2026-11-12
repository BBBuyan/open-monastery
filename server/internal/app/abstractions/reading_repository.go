package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type ReadingRepository interface {
	FindByIdForCustomer(ctx context.Context, id, customerID int64) (*models.Reading, error)
	Add(ctx context.Context, monastery *models.Reading) error
	UserExists(ctx context.Context, id int64) (bool, error)
	BookExists(ctx context.Context, id int64) (bool, error)
	GetReadingsByUserId(
		ctx context.Context,
		id int64,
		opts appmodels.ListOptions,
	) ([]*models.Reading, error)
}
