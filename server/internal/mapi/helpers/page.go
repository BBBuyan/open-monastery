package helpers

type PageInfo struct {
	HasMore    bool    `json:"hasMore"`
	NextCursor *string `json:"nextCursor"`
}
 
type Paginated[T any] struct {
	Data       []T      `json:"data"`
	Pagination PageInfo `json:"pagination"`
}
