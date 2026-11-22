package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type CategoryRepository interface {
	FindById(ctx context.Context, monasteryID, id int64) (*models.Category, error)
	GetAll(ctx context.Context, monasteryID int64) ([]*models.Category, error)
	GetBooksByCategoryId(
		ctx context.Context,
		monasteryID, categoryId int64,
		options appmodels.ListOptions,
	) ([]*models.Book, error)
	Add(ctx context.Context, category *models.Category) error
	RemoveBook(ctx context.Context, monasteryID, categoryId, bookId int64) error
	AddCategoryBook(ctx context.Context, categoryBook *models.CategoryBook) error
}
