package todoist

// PaginationFilters is used to specify the page size and cursor for
// paginated requests. The cursor is used to get the next page of results.
// If there are no more pages, the cursor will be nil.
type PaginationFilters struct {
	Cursor string `json:"cursor,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// Paginated endpoints are marked by having the next_cursor attribute in the
// response. When a response comes back with next_cursor: null, it means the
// endpoint is paginated but there are no more pages to request data from. If
// the cursor is non-null, there are more objects to return, and a new request
// is necessary to get the next page. The next_cursor contains an opaque string
// that shouldn't be modified in any way. It should be sent as-is in the cursor
// parameter along with the same parameters used in the previous request.
//
// The NextCursor and Results fields will usually be returned as separate return
// values when an endpoint is called. The NextCursor will be nil if there are no
// more pages to return. The Results field will contain the results of the
// request.
//
// Example:
//
//	results, nextCursor, err := client.GetProjects(&todoist.PaginationFilters{
//	 Limit:  10,
//	})
type PaginationResponse[T any] struct {
	NextCursor *string `json:"next_cursor"`
	Results    []T     `json:"results"`
}
