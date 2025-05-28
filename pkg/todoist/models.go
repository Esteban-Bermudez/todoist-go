package todoist

import (
	_ "encoding/json"
)

// I can use the same models from the rest api in the sync api they are now the
// same

type SyncReadResponse struct {
	SyncToken                 string              `json:"sync_token"`
	FullSync                  bool                `json:"full_sync"`
	User                      User                `json:"user,omitempty"`
	Projects                  []Project           `json:"projects,omitempty"`
	Items                     []Task              `json:"items,omitempty"`
	Notes                     []Comment           `json:"notes,omitempty"`
	ProjectNotes              []Comment           `json:"project_notes,omitempty"`
	Sections                  []Section           `json:"sections,omitempty"`
	Labels                    []Label             `json:"labels,omitempty"`
	Filters                   []Filter            `json:"filters,omitempty"`
	DayOrders                 map[string]int      `json:"day_orders,omitempty"`
	Reminders                 []Reminder          `json:"reminders,omitempty"`
	Collaborators             []Collaborator      `json:"collaborators,omitempty"`
	CollaboratorStates        []CollaboratorState `json:"collaborator_states,omitempty"`
	CompletedInfo             []CompletedInfo     `json:"completed_info,omitempty"`
	LiveNotifications         []LiveNotification  `json:"live_notifications,omitempty"`
	LiveNotificationsLastRead string              `json:"live_notifications_last_read,omitempty"`
	UserSettings              UserSettings        `json:"user_settings,omitempty"`
	UserPlanLimits            UserPlanLimits      `json:"user_plan_limits,omitempty"`
	Workspaces                []Workspace         `json:"workspaces,omitempty"`
	WorkspaceUsers            []WorkspaceUser     `json:"workspace_users,omitempty"` // only included in incremental sync
}

type SyncWriteResponse struct {
	SyncToken     string            `json:"sync_token"`
	SyncStatus    map[string]string `json:"sync_status"`
	TempIDMapping map[string]string `json:"temp_id_mapping,omitempty"`
}
