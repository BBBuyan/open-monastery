package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/apperror"
	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/category"
	"github.com/gin-gonic/gin"
)

type categoryHandlerRepo struct {
	category.Repository
	categoryID int64
	bookID     int64
	name       string
	err        error
}

func (r *categoryHandlerRepo) GetAll(ctx context.Context) ([]category.Category, error) {
	return []category.Category{{Id: 7, Name: "Category"}}, r.err
}

func (r *categoryHandlerRepo) Add(ctx context.Context, c *category.Category) error {
	r.name = c.Name
	c.Id = 7
	return r.err
}

func (r *categoryHandlerRepo) AddBook(ctx context.Context, categoryID, bookID int64) error {
	r.categoryID, r.bookID = categoryID, bookID
	return r.err
}

func (r *categoryHandlerRepo) GetBooksByCategoryId(ctx context.Context, id int64, options book.ListOptions) ([]book.Book, error) {
	r.categoryID = id
	return []book.Book{{Id: 9, Name: "Prayer"}}, r.err
}

func (r *categoryHandlerRepo) GetCategoriesByBookId(ctx context.Context, id int64) ([]category.Category, error) {
	r.bookID = id
	return []category.Category{{Id: 7, Name: "Category"}}, r.err
}

func TestCategoryEndpoints(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, method, path, body string
		status                   int
		categoryID, bookID       int64
		createdName              string
		response                 string
		repoErr                  error
	}{
		{name: "list categories", method: "GET", path: "/categories", status: 200, response: `[{"id":7,"name":"Category"}]`},
		{name: "create", method: "POST", path: "/categories", body: `{"name":" Category "}`, status: 201, createdName: "Category", response: `{"id":7,"name":"Category"}`},
		{name: "blank name", method: "POST", path: "/categories", body: `{"name":" "}`, status: 400},
		{name: "invalid JSON", method: "POST", path: "/categories", body: `{`, status: 400},
		{name: "add book", method: "POST", path: "/categories/7/books", body: `{"book_id":9}`, status: 204, categoryID: 7, bookID: 9},
		{name: "invalid book", method: "POST", path: "/categories/7/books", body: `{"book_id":0}`, status: 400},
		{name: "invalid category", method: "POST", path: "/categories/no/books", body: `{"book_id":9}`, status: 400},
		{name: "category books", method: "GET", path: "/categories/7/books", status: 200, categoryID: 7, response: `{"books":[{"id":9,"name":"Prayer","description":"","price":0}],"next_cursor":null,"has_more":false}`},
		{name: "invalid lookup category", method: "GET", path: "/categories/0/books", status: 400},
		{name: "book categories", method: "GET", path: "/books/9/categories", status: 200, bookID: 9, response: `[{"id":7,"name":"Category"}]`},
		{name: "invalid lookup book", method: "GET", path: "/books/-1/categories", status: 400},
		{name: "missing reference", method: "POST", path: "/categories/7/books", body: `{"book_id":9}`, status: 404, categoryID: 7, bookID: 9, repoErr: apperror.New(apperror.NotFound, "referenced resource not found", nil)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &categoryHandlerRepo{err: tc.repoErr}
			h := NewCategoryHandler(category.NewService(repo))
			router := gin.New()
			router.Use(ErrorHandler())
			router.POST("/categories", h.AddCategory)
			router.GET("/categories", h.GetCategories)
			router.POST("/categories/:id/books", h.AddBook)
			router.GET("/categories/:id/books", h.GetBooksByCategoryId)
			router.GET("/books/:id/categories", h.GetCategoriesByBookId)
			req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			router.ServeHTTP(response, req)
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d: %s", response.Code, tc.status, response.Body.String())
			}
			if repo.categoryID != tc.categoryID || repo.bookID != tc.bookID || repo.name != tc.createdName {
				t.Fatalf("unexpected repository arguments: %+v", repo)
			}
			if tc.response != "" && response.Body.String() != tc.response {
				t.Fatalf("response = %s, want %s", response.Body.String(), tc.response)
			}
			if tc.status == http.StatusNoContent && response.Body.Len() != 0 {
				t.Fatal("204 response should have no body")
			}
		})
	}
}
