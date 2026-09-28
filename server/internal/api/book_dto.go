package api

type CreateBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
