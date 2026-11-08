package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type CategoryRepository interface {
	FindById(ctx context.Context, id int64) (*models.Category, error)
	GetAll(ctx context.Context) ([]*models.Category, error)
	GetBooksByCategoryId(
		ctx context.Context,
		categoryId int64,
		options appmodels.ListOptions,
	) ([]*models.Book, error)
	Add(ctx context.Context, category *models.Category) error
	AddBook(ctx context.Context, categoryId, bookId int64) error
	RemoveBook(ctx context.Context, categoryId, bookId int64) error
}
