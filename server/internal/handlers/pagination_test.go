package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/open-monastery/internal/book"
	"example.com/open-monastery/internal/category"
	"github.com/gin-gonic/gin"
)

type paginatedRepo struct {
	book.Repository
	options book.ListOptions
}

func (r *paginatedRepo) GetBooks(ctx context.Context, options book.ListOptions) ([]book.Book, error) {
	r.options = options
	all := []book.Book{{Id: 2, Name: "Prayer"}, {Id: 5, Name: "Other"}, {Id: 9, Name: "Prayer"}, {Id: 12, Name: "Prayer"}}
	var result []book.Book
	for _, b := range all {
		if b.Id > options.AfterId && strings.Contains(b.Name, options.Search) {
			result = append(result, b)
			if len(result) == options.Limit+1 {
				break
			}
		}
	}
	return result, nil
}

type paginatedCategoryRepo struct {
	category.Repository
	books *paginatedRepo
}

func (r *paginatedCategoryRepo) GetBooksByCategoryId(ctx context.Context, id int64, options book.ListOptions) ([]book.Book, error) {
	if id != 7 {
		return []book.Book{}, nil
	}
	return r.books.GetBooks(ctx, options)
}

func TestBookPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, path := range []string{"/books", "/categories/7/books"} {
		for _, tc := range []struct {
			query  string
			status int
			ids    []int64
			cursor int64
			limit  int
		}{
			{query: "?limit=2", status: 200, ids: []int64{2, 5}, cursor: 5, limit: 2},
			{query: "?limit=2&after_id=5", status: 200, ids: []int64{9, 12}, limit: 2},
			{query: "?limit=2&after_id=12", status: 200, ids: []int64{}, limit: 2},
			{query: "?limit=1&q=Prayer&after_id=2", status: 200, ids: []int64{9}, cursor: 9, limit: 1},
			{query: "?q=absent", status: 200, ids: []int64{}, limit: 20},
			{query: "", status: 200, ids: []int64{2, 5, 9, 12}, limit: 20},
			{query: "?limit=999", status: 200, ids: []int64{2, 5, 9, 12}, limit: 100},
			{query: "?limit=0", status: 400},
			{query: "?limit=-1", status: 400},
			{query: "?limit=no", status: 400},
			{query: "?limit=", status: 400},
			{query: "?after_id=-1", status: 400},
			{query: "?after_id=no", status: 400},
			{query: "?after_id=9223372036854775808", status: 400},
		} {
			t.Run(path+tc.query, func(t *testing.T) {
				repo := &paginatedRepo{}
				router := gin.New()
				router.Use(ErrorHandler())
				router.GET("/books", NewBookHandler(book.NewService(repo)).GetBooks)
				handler := NewCategoryHandler(category.NewService(&paginatedCategoryRepo{books: repo}))
				router.GET("/categories/:id/books", handler.GetBooksByCategoryId)
				response := httptest.NewRecorder()
				router.ServeHTTP(response, httptest.NewRequest("GET", path+tc.query, nil))
				if response.Code != tc.status {
					t.Fatalf("status %d: %s", response.Code, response.Body.String())
				}
				if tc.status != 200 {
					return
				}
				var page book.Page
				if err := json.Unmarshal(response.Body.Bytes(), &page); err != nil {
					t.Fatal(err)
				}
				if page.Books == nil || len(page.Books) != len(tc.ids) {
					t.Fatalf("unexpected books: %+v", page)
				}
				for i, id := range tc.ids {
					if page.Books[i].Id != id {
						t.Fatalf("unexpected book ID: %+v", page)
					}
				}
				if page.HasMore != (tc.cursor != 0) {
					t.Fatalf("unexpected has_more: %+v", page)
				}
				if tc.cursor == 0 && page.NextCursor != nil {
					t.Fatal("final page cursor must be null")
				}
				if tc.cursor != 0 && (page.NextCursor == nil || *page.NextCursor != tc.cursor) {
					t.Fatal("incorrect next cursor")
				}
				if repo.options.Limit != tc.limit {
					t.Fatalf("limit = %d", repo.options.Limit)
				}
			})
		}
	}
}
