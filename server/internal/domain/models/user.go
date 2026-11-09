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
	Email       *string  `json:"email"`
	PhoneNumber *string  `json:"phone_number"`
	MonasteryId *int64   `json:"monastery_id"`
}

func NewMonk(
	user *User,
	monasteryId int64,
) (*User, error) {
	user.MonasteryId = &monasteryId
	user.UserType = Monk

	return user, nil
}

func NewUser(
	firstName string,
	lastName string,
	email *string,
	phone *string,
) (*User, error) {
	user := &User{
		FirstName:   firstName,
		LastName:    lastName,
		UserType:    Customer,
		Email:       email,
		PhoneNumber: phone,
	}

	return user, nil
}
