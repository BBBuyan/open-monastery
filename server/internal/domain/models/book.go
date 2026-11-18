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
	CreatedBy   int64
	CreatedAt   time.Time
}

func NewBook(
	name string,
	desc string,
	price uint32,
	monasteryId int64,
	userId int64,
) *Book {
	return &Book{
		Name:        name,
		Description: desc,
		Price:       price,
		MonasteryId: monasteryId,
		CreatedBy:   userId,
		CreatedAt:   time.Now().UTC().Truncate(time.Microsecond),
	}
}
