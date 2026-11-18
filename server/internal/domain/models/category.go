package models

import (
	"time"
)

type Category struct {
	Id          int64
	Name        string
	MonasteryId int64
	CreatedBy   int64
	CreatedAt   time.Time
}

func NewCategory(
	name string,
	userId int64,
	monasteryId int64,
) (*Category, error) {
	return &Category{
		Name:        name,
		MonasteryId: monasteryId,
		CreatedBy:   userId,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}, nil
}
