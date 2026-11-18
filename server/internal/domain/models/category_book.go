package models

import "time"

type CategoryBook struct {
	Id          int64
	CategoryId  int64
	BookId      int64
	MonasteryId int64
	CreatedBy   int64
	CreatedAt   time.Time
}

func NewCategoryBook(
	categoryId int64,
	bookId int64,
	userId int64,
) (*CategoryBook, error) {
	return &CategoryBook{
		CategoryId: categoryId,
		BookId:     bookId,
		CreatedBy:  userId,
		CreatedAt:  time.Now().UTC().Truncate(time.Microsecond),
	}, nil
}
