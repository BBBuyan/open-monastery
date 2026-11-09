package helpers

import (
	"regexp"

	"example.com/open-monastery/internal/app/apperr"
)

var mongolianPhoneRegex = regexp.MustCompile(`^\+976[0-9]{8}$`)

func NewPhoneNumber(val string) (string, error) {
	if len(val) > 12 {
		return "", apperr.NewInvalid("phone number is too long", nil)
	}
	if !mongolianPhoneRegex.MatchString(val) {
		return "", apperr.NewInvalid("invalid phone number", nil)
	}

	return val, nil
}
