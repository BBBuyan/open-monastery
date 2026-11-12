package abstractions

import (
	"context"

	"example.com/open-monastery/internal/app/appmodels"
	"example.com/open-monastery/internal/domain/models"
)

type BookRepository interface {
	FindById(ctx context.Context, id int64) (*models.Book, error)
	GetBooks(
		ctx context.Context,
		options appmodels.ListOptions,
	) ([]*models.Book, error)
	Add(ctx context.Context, book *models.Book) error
	GetCategories(ctx context.Context, bookId int64) ([]*models.Category, error)
}
