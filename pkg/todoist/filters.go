package todoist

type Filter struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Query      string `json:"query"`
	Color      int    `json:"color"`
	ItemOrder  int    `json:"item_order"`
	IsDeleted  bool   `json:"is_deleted"`
	IsFavorite bool   `json:"is_favorite"`
	IsFrozen   bool   `json:"is_frozen"`
}
