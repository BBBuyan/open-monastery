package reading

type Reading struct {
	Id         int64 `json:"id"`
	BookId     int64 `json:"book_id"`
	MonkId     int64 `json:"monk_id"`
	CustomerId int64 `json:"customer_id"`
}
