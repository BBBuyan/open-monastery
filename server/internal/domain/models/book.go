package models

type Book struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint32 `json:"price"`
}

