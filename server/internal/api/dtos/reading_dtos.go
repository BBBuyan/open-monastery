package dtos

type CreateReadingRequest struct {
	BookId int64 `json:"book_id" binding:"gt=0"`
}

type AssignReadingRequest struct {
	MonkId int64 `json:"monk_id" binding:"gt=0"`
}

type ReadingsURI struct {
	MonasteryId int64 `uri:"id" binding:"required,min=1"`
}

type ReadingURI struct {
	ReadingsURI
	ReadingId int64 `uri:"readingId" binding:"required,min=1"`
}
