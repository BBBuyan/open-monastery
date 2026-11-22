package models

import (
	"errors"
	"fmt"
	"net/mail"
	"strings"
)

var (
	ErrEmailTooLong      = errors.New("email is too long")
	ErrEmailNotParseable = errors.New("email not parseable")
	ErrInvalidEmail      = errors.New("invalid email")
)

type Email struct {
	value string
}

func NewEmail(value string) (*Email, error) {
	addr, err := NormalizeEmail(value)
	if err != nil {
		return nil, err
	}
	return &Email{value: addr}, nil
}

func NormalizeEmail(email string) (string, error) {
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

func (e Email) String() string {
	return e.value
}
