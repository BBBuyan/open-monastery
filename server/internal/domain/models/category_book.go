package models

import "time"

type CategoryBook struct {
	Id          int64
	CategoryId  int64
	BookId      int64
	MonasteryId int64
	CreatedBy   *string
	CreatedAt   time.Time
}

func NewCategoryBook(
	categoryId int64,
	bookId int64,
	monasteryId int64,
	createdBy *string,
) *CategoryBook {
	return &CategoryBook{
		CategoryId: categoryId,
		BookId:     bookId,
		MonasteryId: monasteryId,
		CreatedBy:  createdBy,
		CreatedAt:  time.Now().UTC().Truncate(time.Microsecond),
	}
}
