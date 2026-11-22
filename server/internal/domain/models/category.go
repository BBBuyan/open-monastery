package models

import (
	"time"
)

type Category struct {
	Id          int64
	Name        string
	MonasteryId int64
	CreatedBy   *string
	CreatedAt   time.Time
}

func NewCategory(
	name string,
	monasteryId int64,
	createdBy *string,
) *Category {
	return &Category{
		Name:        name,
		MonasteryId: monasteryId,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}
}
