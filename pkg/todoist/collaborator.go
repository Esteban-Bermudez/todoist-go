package todoist

// Collaborator represents a collaborator on a project.
type Collaborator struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email"`
	ImageID  string `json:"image_id,omitempty"`
	Timezone string `json:"timezone,omitempty"`
}

type CollaboratorState struct {
	ProjectID string `json:"project_id"`
	UserID    string `json:"user_id"`
	State     string `json:"state"` // "active", "invited"
	IsDeleted bool   `json:"is_deleted,omitempty"`
	Role      string `json:"role,omitempty"` // only available for teams
}
