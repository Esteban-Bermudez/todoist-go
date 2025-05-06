package todoist

type PaginationOptions struct {
  Cursor string `json:"cursor,omitempty"`
  Limit  int    `json:"limit,omitempty"`
}
