package category

import (
	"context"

	"example.com/open-monastery/internal/book"
)

type Repository interface {
	FindById(ctx context.Context, id int64) (*Category, error)
	GetAll(ctx context.Context) ([]Category, error)
	GetBooksByCategoryId(ctx context.Context, categoryId int64, options book.ListOptions) ([]book.Book, error)
	GetCategoriesByBookId(ctx context.Context, bookId int64) ([]Category, error)
	Add(ctx context.Context, category *Category) error
	GetCategoryBooks(ctx context.Context) ([]CategoryBook, error)
	AddBook(ctx context.Context, categoryId, bookId int64) error
	RemoveBook(ctx context.Context, categoryId, bookId int64) error
}
