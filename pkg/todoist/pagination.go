package todoist

type PaginationFilters struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

type PaginationResponse[T any] struct {
	NextCursor *string `json:"next_cursor"`
	Results    []T     `json:"results"`
}
