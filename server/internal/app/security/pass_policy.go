package security

import "errors"

const MinPasswordLength = 8

var ErrPasswordTooShort = errors.New(
	"password must contain at least 8 characters",
)

func ValidatePassword(password string) error {
	if len([]rune(password)) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	return nil
}
