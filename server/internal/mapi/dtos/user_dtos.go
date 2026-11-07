package dtos

import (
	"example.com/open-monastery/internal/domain/models"
)

type CreateUserRequest struct {
	FirstName   string        `json:"first_name"`
	LastName    string        `json:"last_name"`
	UserType    models.UserType `json:"user_type"`
	MonasteryId *int64        `json:"monastery_id,omitempty"`
}
