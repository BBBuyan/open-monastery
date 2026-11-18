package models

import (
	"time"

	"example.com/open-monastery/internal/domain/enums"
)

type Reading struct {
	Id          int64  `json:"id"`
	BookId      int64  `json:"book_id"`
	MonkId      *int64 `json:"monk_id"`
	CustomerId  int64  `json:"customer_id"`
	MonasteryId int64  `json:"monastery_id"`

	Status enums.ReadingStatus `json:"status"`
	// CreatedAt is the UTC creation time, when recorded.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// AssignedAt is the UTC time when a monk claimed the reading, when recorded.
	AssignedAt *time.Time `json:"assigned_at,omitempty"`
}

func NewReading(
	bookId int64,
	customerId int64,
	monasteryId int64,
) *Reading {
	now := time.Now().UTC().Truncate(time.Microsecond)
	return &Reading{
		BookId:      bookId,
		CustomerId:  customerId,
		MonasteryId: monasteryId,
		CreatedAt:   &now,
		Status:      enums.Pending,
	}
}

func (r *Reading) AssignToMonk(
	userId int64,
) error {
	r.MonkId = &userId
	r.Status = enums.Scheduled

	now := time.Now().UTC().Truncate(time.Microsecond)
	r.AssignedAt = &now
	return nil
}
