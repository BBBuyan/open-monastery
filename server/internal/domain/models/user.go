package models

type UserType string

const (
	Monk     UserType = "monk"
	Customer UserType = "customer"
)

type User struct {
	Id          int64    `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	UserType    UserType `json:"user_type"`
	MonasteryId *int64   `json:"monastery_id"`
}

