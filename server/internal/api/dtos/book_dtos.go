package dtos

type CreateBookRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       uint32 `json:"price"`
}

type BooksURI struct {
	MonasteryId int64 `uri:"id" binding:"required,min=1"`
}

type BookURI struct {
	BooksURI
	BookId int64 `uri:"bookId" binding:"required,min=1"`
}
