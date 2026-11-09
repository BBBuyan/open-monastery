package dtos

type CreateUserRequest struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

type CreateMonkRequest struct {
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
	MonasteryId int64   `json:"monastery_id" binding:"required,gt=0"`
}
