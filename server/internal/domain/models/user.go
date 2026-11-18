package models

import (
	"regexp"

	"example.com/open-monastery/internal/domain/enums"
)

var mongolianPhoneRegex = regexp.MustCompile(`^\+976[0-9]{8}$`)

type User struct {
	Id           int64
	FirstName    string
	LastName     string
	PasswordHash string
	UserType     enums.UserType
	Email        Email
	PhoneNumber  *PhoneNumber
	MonasteryId  *int64
}

func NewUser(
	firstName string,
	lastName string,
	email Email,
	phone *PhoneNumber,
) (*User, error) {
	user := &User{
		FirstName:   firstName,
		LastName:    lastName,
		UserType:    enums.Customer,
		Email:       email,
		PhoneNumber: phone,
	}
	return user, nil
}

func (u *User) PromoteToMonk(monasteryId int64) {
	u.MonasteryId = &monasteryId
	u.UserType = enums.Monk
}
