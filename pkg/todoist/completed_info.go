package todoist

type CompletedInfo struct {
	ArchivedSections int    `json:"archived_sections,omitempty"`
	CompletedItems   int    `json:"completed_items,omitempty"`
	ProjectID        string `json:"project_id,omitempty"`
}
