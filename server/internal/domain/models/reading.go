package models

import (
	"time"

	"example.com/open-monastery/internal/domain/enums"
)

type Reading struct {
	Id         int64               `json:"id"`
	BookId     int64               `json:"book_id"`
	MonkId     *int64              `json:"monk_id"`
	CustomerId int64               `json:"customer_id"`
	Status     enums.ReadingStatus `json:"status"`
	CreatedAt  time.Time           `json:"created_at"`
	AssignetAt time.Time           `json:"assigned_at"`
}

func NewReading(
	bookId int64,
	customerId int64,
) *Reading {
	return &Reading{
		BookId:     bookId,
		CustomerId: customerId,
		CreatedAt:  time.Now(),
		Status:     enums.Pending,
	}
}

func (r *Reading) AssignToMonk(
	monkId int64,
) {
	r.MonkId = &monkId
	r.AssignetAt = time.Now()
	r.Status = enums.Scheduled
}
