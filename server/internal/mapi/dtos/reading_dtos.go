package dtos

type CreateReadingRequest struct {
	BookId int64 `json:"book_id" binding:"gt=0"`
}
