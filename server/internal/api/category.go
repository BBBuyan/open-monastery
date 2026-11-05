package api

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

type AddBookToCategoryRequest struct {
	BookId int64 `json:"book_id"`
}
