package models

import (
	"time"

	"example.com/open-monastery/internal/domain/enums"
)

type User struct {
	Id           int64
	FirstName    string
	LastName     string
	PasswordHash string
	UserType     enums.UserType
	Email        Email
	PhoneNumber  *PhoneNumber
	MonasteryId  *int64
	CreatedAt    time.Time
}

func NewUser(
	firstName string,
	lastName string,
	email Email,
	phone *PhoneNumber,
	passHash string,
) (*User, error) {
	user := &User{
		FirstName:    firstName,
		LastName:     lastName,
		UserType:     enums.Customer,
		Email:        email,
		PhoneNumber:  phone,
		PasswordHash: passHash,
		CreatedAt: time.Now().UTC().Truncate(time.Microsecond),
	}
	return user, nil
}

func (u *User) PromoteToMonk(monasteryId int64) {
	u.MonasteryId = &monasteryId
	u.UserType = enums.Monk
}
