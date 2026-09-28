package monastery

import "context"

type Repository interface {
	FindById(ctx context.Context, id int64) (*Monastery, error)
	Add(ctx context.Context, monastery *Monastery) error
}
