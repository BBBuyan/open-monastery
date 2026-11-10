package helpers

import (
	"net/mail"
	"strings"

	"example.com/open-monastery/internal/app/apperr"
)

func NewEmail(email string) (string, error) {
	email = strings.TrimSpace(email)
	if len(email) > 254 {
		return "", apperr.NewInvalid(
			"email is too long",
			"email must be less than 254 characters",
			nil,
		)
	}

	address, err := mail.ParseAddress(email)
	if err != nil {
		return "", apperr.NewInvalid(
			"email not parseable",
			"email not parseable",
			err,
		)
	}
	if address.Address != email {
		return "", apperr.NewInvalid(
			"invalid email",
			"invalid email",
			nil,
		)
	}

	return address.Address, nil
}
