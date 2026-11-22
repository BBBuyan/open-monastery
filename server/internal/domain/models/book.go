package models

import (
	"time"
)

type Book struct {
	Id          int64
	Name        string
	Description string
	Price       uint32
	MonasteryId int64
	CreatedBy   *string
	CreatedAt   time.Time
}

func NewBook(
	name string,
	desc string,
	price uint32,
	monasteryId int64,
	createdBy *string,
) *Book {
	return &Book{
		Name:        name,
		Description: desc,
		Price:       price,
		MonasteryId: monasteryId,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}
}
