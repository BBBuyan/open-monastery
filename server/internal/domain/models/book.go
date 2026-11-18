package models

import (
	"time"
)

type Book struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint32    `json:"price"`
	MonasteryId int64     `json:"monastery_id"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
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
