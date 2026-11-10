package appmodels

import (
	"strings"

	"example.com/open-monastery/internal/app/apperr"
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
		return o, apperr.NewBadRequest(
			"request validation",
			"limit must be positive and after_id must be nonnegative",
			nil,
		)
	}
	if o.Limit > 100 {
		o.Limit = 100
	}
	o.Search = strings.TrimSpace(o.Search)
	return o, nil
}
