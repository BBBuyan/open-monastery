package dtos

import (
	"time"

	"example.com/open-monastery/internal/domain/enums"
	"example.com/open-monastery/internal/domain/models"
)

type CreateUserResponse struct {
	Id          int64          `json:"id"`
	FirstName   string         `json:"first_name"`
	LastName    string         `json:"last_name"`
	UserType    enums.UserType `json:"user_type"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	MonasteryId *int64         `json:"monastery_id"`
	CreatedAt   time.Time      `json:"created_at"`
}

func NewCreateUserResponse(user *models.User) CreateUserResponse {
	response := CreateUserResponse{
		Id:          user.Id,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		UserType:    user.UserType,
		Email:       user.Email.String(),
		MonasteryId: user.MonasteryId,
		CreatedAt:   user.CreatedAt,
	}
	if user.PhoneNumber != nil {
		response.PhoneNumber = user.PhoneNumber.String()
	}
	return response
}

type CreateUserRequest struct {
	Password    string  `json:"password" binding:"required"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Email       string  `json:"email" binding:"required" format:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateMonkRequest struct {
	Password    string  `json:"password" binding:"required"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Email       *string `json:"email" format:"email"`
	PhoneNumber *string `json:"phone_number"`
	MonasteryId int64   `json:"monastery_id" binding:"required,gt=0"`
}
