package models

import (
	"time"
)

type Category struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	MonasteryId int64     `json:"monastery_id"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
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
