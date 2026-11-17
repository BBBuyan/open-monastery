package dtos

type CreateReadingRequest struct {
	BookId int64 `json:"book_id" binding:"gt=0"`
}

type AssignReadingRequest struct {
	MonkId int64 `json:"monk_id" binding:"gt=0"`
}
