package dtos

type CreateMonasteryRequest struct {
	Name string `json:"name" binding:"required"`
}
