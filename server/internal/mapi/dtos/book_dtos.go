package dtos

type CreateBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
