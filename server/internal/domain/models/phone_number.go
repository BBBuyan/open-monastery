package models

import (
	"errors"
	"regexp"
)

var (
	ErrPhoneNumberTooLong = errors.New("phone number is too long")
	ErrInvalidPhoneNumber = errors.New("invalid phone number")
)

var mongolianPhoneRegex = regexp.MustCompile(`^\+976[0-9]{8}$`)

type PhoneNumber struct {
	value string
}

func NewPhoneNumber(number string) (*PhoneNumber, error) {
	if len(number) > 12 {
		return nil, ErrPhoneNumberTooLong
	}
	if !mongolianPhoneRegex.MatchString(number) {
		return nil, ErrInvalidPhoneNumber
	}
	return &PhoneNumber{value: number}, nil
}

func (p PhoneNumber) String() string {
	return p.value
}
