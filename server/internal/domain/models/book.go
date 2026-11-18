package models

import "time"

type Book struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       uint32    `json:"price"`
	MonasteryId int64     `json:"monastery_id"`
	CreatedBy   int64     `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
}
