package reading

import "time"

type Status string

const (
	Pending   Status = "pending"
	Scheduled Status = "scheduled"
	Cancelled Status = "cancelled"
	Refunded  Status = "refunded"
)

type Reading struct {
	Id         int64  `json:"id"`
	BookId     int64  `json:"book_id"`
	MonkId     *int64 `json:"monk_id"`
	CustomerId int64  `json:"customer_id"`
	Status     Status `json:"status"`
	CreatedAt  time.Time
}

func NewReading(
	bookId int64,  
	customerId int64,
) *Reading {
	return &Reading{
		BookId: bookId,
		CustomerId: customerId,
		CreatedAt: time.Now(),
		Status: Pending,
	}
}
