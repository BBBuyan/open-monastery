package api

import "example.com/open-monastery/internal/user"

type CreateUserRequest struct {
	FirstName   string        `json:"first_name"`
	LastName    string        `json:"last_name"`
	UserType    user.UserType `json:"user_type"`
	MonasteryId *int64        `json:"monastery_id,omitempty"`
}
