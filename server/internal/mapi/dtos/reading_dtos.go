package dtos

type CreateReadingRequest struct {
	BookId     int64 `json:"book_id" binding:"gt=0"`
	MonkId     int64 `json:"monk_id" binding:"gt=0"`
	CustomerId int64 `json:"customer_id" binding:"gt=0"`
}
