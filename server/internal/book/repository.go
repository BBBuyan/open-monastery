package book

import "context"

type Repository interface {
	FindById(ctx context.Context, id int64) (*Book, error)
	GetBooks(ctx context.Context, options ListOptions) ([]Book, error)
	Add(ctx context.Context, book *Book) error
}
