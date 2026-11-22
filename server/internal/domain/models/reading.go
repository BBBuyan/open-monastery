package models

import (
	"time"

	"example.com/open-monastery/internal/domain/enums"
)

type Reading struct {
	Id          int64
	BookId      int64
	MonkId      *int64
	CustomerId  int64
	MonasteryId int64

	Status enums.ReadingStatus
	// CreatedAt is the UTC creation time, when recorded.
	CreatedAt *time.Time
	// AssignedAt is the UTC time when a monk claimed the reading, when recorded.
	AssignedAt *time.Time
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
