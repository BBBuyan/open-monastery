package models

import (
	"errors"
	"testing"
)

func TestNewUserValidatesContact(t *testing.T) {
	tests := []struct {
		name  string
		email *string
		phone *string
		want  error
	}{
		{name: "missing contact", want: ErrContactRequired},
		{name: "invalid email with valid phone", email: new("invalid"), phone: new("+97699112233"), want: ErrEmailNotParseable},
		{name: "invalid phone with valid email", email: new("a@example.com"), phone: new("+15551234567"), want: ErrInvalidPhoneNumber},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := NewUser("A", "B", tt.email, tt.phone)
			if user != nil || !errors.Is(err, tt.want) {
				t.Fatalf("NewUser returned (%v, %v), want (nil, %v)", user, err, tt.want)
			}
		})
	}
}
