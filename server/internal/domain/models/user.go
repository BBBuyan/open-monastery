package models

import (
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var (
	ErrContactRequired    = errors.New("at least one contact method is required")
	ErrEmailTooLong       = errors.New("email is too long")
	ErrEmailNotParseable  = errors.New("email not parseable")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrPhoneNumberTooLong = errors.New("phone number is too long")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)

var mongolianPhoneRegex = regexp.MustCompile(`^\+976[0-9]{8}$`)

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
	if email == nil && phone == nil {
		return nil, ErrContactRequired
	}
	if email != nil {
		normalized, err := normalizeEmail(*email)
		if err != nil {
			return nil, err
		}
		email = &normalized
	}
	if phone != nil {
		number := *phone
		if err := validatePhoneNumber(number); err != nil {
			return nil, err
		}
		phone = &number
	}

	user := &User{
		FirstName:   firstName,
		LastName:    lastName,
		UserType:    Customer,
		Email:       email,
		PhoneNumber: phone,
	}

	return user, nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return "", ErrEmailTooLong
	}
	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrEmailNotParseable, err)
	}
	if address.Address != email {
		return "", ErrInvalidEmail
	}
	return address.Address, nil
}

func validatePhoneNumber(number string) error {
	if len(number) > 12 {
		return ErrPhoneNumberTooLong
	}
	if !mongolianPhoneRegex.MatchString(number) {
		return ErrInvalidPhoneNumber
	}
	return nil
}
