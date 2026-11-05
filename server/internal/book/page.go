package book

import (
	"strings"

	"example.com/open-monastery/internal/apperror"
)

type ListOptions struct {
	Limit   int
	AfterId int64
	Search  string
}

func (o ListOptions) Normalize() (ListOptions, error) {
	if o.Limit == 0 {
		o.Limit = 20
	}
	if o.Limit < 1 || o.AfterId < 0 {
		return o, apperror.New(apperror.BadRequest, "limit must be positive and after_id must be nonnegative", nil)
	}
	if o.Limit > 100 {
		o.Limit = 100
	}
	o.Search = strings.TrimSpace(o.Search)
	return o, nil
}

type Page struct {
	Books      []Book `json:"books"`
	NextCursor *int64 `json:"next_cursor"`
	HasMore    bool   `json:"has_more"`
}

// Repositories return up to limit+1 rows to detect the next page.
func NewPage(books []Book, limit int) *Page {
	if books == nil {
		books = make([]Book, 0)
	}
	page := &Page{Books: books}
	if len(books) > limit {
		page.Books = books[:limit]
		page.HasMore = true
		cursor := page.Books[limit-1].Id
		page.NextCursor = &cursor
	}
	return page
}
