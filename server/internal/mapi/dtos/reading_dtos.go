package dtos

type CreateReadingRequest struct {
	BookId int64 `json:"book_id" binding:"gt=0"`
	MonasteryId int64 `json:"monastery_id" binding:"gt=0"`
}
