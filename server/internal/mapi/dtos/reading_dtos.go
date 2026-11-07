package dtos

type CreateReadingRequest struct {
	BookId     int64 `json:"book_id"`
	MonkId     int64 `json:"monk_id"`
	CustomerId int64 `json:"customer_id"`
}
