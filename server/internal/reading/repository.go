package reading

import "context"

type Repository interface {
	FindById(ctx context.Context, id int64) (*Reading, error)
	Add(ctx context.Context, monastery *Reading) error
}
