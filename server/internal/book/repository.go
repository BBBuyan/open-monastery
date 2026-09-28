package book

import "context"

type Repository interface {
	FindById(ctx context.Context, id int64) (*Book, error)
	Add(ctx context.Context, monastery *Book) error
}
