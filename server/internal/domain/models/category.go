package models

type Category struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	MonasteryId int64  `json:"monastery_id"`
}
