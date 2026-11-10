package helpers

import (
	"regexp"

	"example.com/open-monastery/internal/app/apperr"
)

var mongolianPhoneRegex = regexp.MustCompile(`^\+976[0-9]{8}$`)

func NewPhoneNumber(val string) (string, error) {
	if len(val) > 12 {
		return "", apperr.NewInvalidWithMsg(
			"phone number is too long",
			"phone number must be less than 12 characters",
			nil,
		)
	}
	if !mongolianPhoneRegex.MatchString(val) {
		return "", apperr.NewInvalidWithMsg(
			"invalid phone number",
			"phone number is not mongolian",
			nil,
		)
	}

	return val, nil
}
