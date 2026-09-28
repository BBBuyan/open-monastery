package user

import "context"

type Repository interface {
	FindById(ctx context.Context, id int64) (*User, error)
	Add(ctx context.Context, user *User) error
}
