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
	monasteryId int64,
	userId int64,
) *Category {
	return &Category{
		Name:        name,
		MonasteryId: monasteryId,
		CreatedBy:   userId,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}
}
