package appmodels

type PageInfo struct {
	HasMore    bool    `json:"hasMore"`
	NextCursor *string `json:"nextCursor"`
}

type Paginated[T any] struct {
	Data       []T      `json:"data"`
	Pagination PageInfo `json:"pagination"`
}

func ToPage[T any](
	data []T,
	hasMore bool,
	cursor *Cursor,
) *Paginated[T] {
	pageInfo := PageInfo{
		HasMore: hasMore,
	}
	if cursor != nil {
		encodedCursor, err := EncodeCursor(*cursor)
		if err == nil {
			pageInfo.NextCursor = &encodedCursor
		}
	}

	return &Paginated[T]{
		Data:       data,
		Pagination: pageInfo,
	}
}
