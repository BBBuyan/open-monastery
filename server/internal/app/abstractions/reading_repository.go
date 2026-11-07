package abstractions

import (
	"context"

	"example.com/open-monastery/internal/domain/models"
)

type ReadingRepository interface {
	FindById(ctx context.Context, id int64) (*models.Reading, error)
	Add(ctx context.Context, monastery *models.Reading) error
	UserExists(ctx context.Context, id int64) (bool, error)
	BookExists(ctx context.Context, id int64) (bool, error)
}
