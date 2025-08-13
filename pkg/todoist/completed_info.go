package todoist

type CompletedInfo []struct {
	ArchivedSections int    `json:"archived_sections"`
	CompletedItems   int    `json:"completed_items"`
	ProjectID        string `json:"project_id"`
}
