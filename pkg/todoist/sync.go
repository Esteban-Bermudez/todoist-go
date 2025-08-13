package todoist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Sync struct {
	// A special string, used to allow the client to perform incremental sync. Pass * to retrieve
	// all active resource data. More details about this below.
	SyncToken string `json:"sync_token"`

	// Used to specify what resources to fetch from the server. Here is a list of available resource
	// types: labels, projects, items, notes, sections, filters, reminders, reminders_location,
	// locations, user, live_notifications, collaborators, user_settings, notification_settings,
	// user_plan_limits, completed_info, stats, workspaces, workspace_users. You may use all to
	// include all the resource types. Resources can also be excluded by prefixing a - prior to the
	// name, for example, -projects
	ResourceTypes []string `json:"resource_types"`

	Commands []Command `json:"commands"`
	APIKey   string    `json:"-"`
}

// type String	The type of the command.
// args Object	The parameters of the command.
// uuid String	Command UUID. More details about this below.
// temp_id String	Temporary resource ID, Optional. Only specified for commands
// that create a new
// resource (e.g. item_add command). More details about this belowV
type Command struct {
	Type   string         `json:"type"`
	Args   map[string]any `json:"args"`
	UUID   string         `json:"uuid"`
	TempID string         `json:"temp_id,omitempty"`
}

type SyncReadResponse struct {
	SyncToken          string              `json:"sync_token"`
	FullSync           bool                `json:"full_sync"`
	User               User                `json:"user,omitempty"`
	Projects           []Project           `json:"projects,omitempty"`
	Items              []Task              `json:"items,omitempty"`
	Notes              []Comment           `json:"notes,omitempty"`
	ProjectNotes       []Comment           `json:"project_notes,omitempty"`
	Sections           []Section           `json:"sections,omitempty"`
	Labels             []Label             `json:"labels,omitempty"`
	Filters            []Filter            `json:"filters,omitempty"`
	DayOrders          map[string]int      `json:"day_orders,omitempty"`
	Reminders          []Reminder          `json:"reminders,omitempty"`
	Collaborators      []Collaborator      `json:"collaborators,omitempty"`
	CollaboratorStates []CollaboratorState `json:"collaborator_states,omitempty"`
	CompletedInfo      []CompletedInfo     `json:"completed_info,omitempty"`

	// this has a weird structure where the array contains different types
	LiveNotifications []map[string]any `json:"live_notifications,omitempty"`

	LiveNotificationsLastRead string          `json:"live_notifications_last_read,omitempty"`
	UserSettings              map[string]any  `json:"user_settings,omitempty"`
	UserPlanLimits            UserPlanLimits  `json:"user_plan_limits,omitempty"`
	Workspaces                *[]Workspace    `json:"workspaces,omitempty"`
	WorkspaceUsers            []WorkspaceUser `json:"workspace_users,omitempty"` // only included in incremental sync
}

type SyncWriteResponse struct {
	SyncToken     string            `json:"sync_token"`
	SyncStatus    map[string]string `json:"sync_status"`
	TempIDMapping map[string]string `json:"temp_id_mapping,omitempty"`
}

func (s *Sync) ReadResources(
	ctx context.Context,
	resourceTypes []string,
) (*SyncReadResponse, error) {
	if resourceTypes != nil {
		s.ResourceTypes = resourceTypes
	} else {
		return nil, fmt.Errorf("resourceTypes cannot be nil")
	}

	data := url.Values{}
	data.Set("sync_token", s.SyncToken)
	data.Set("resource_types", `["`+strings.Join(s.ResourceTypes, `","`)+`"]`)

	resp, err := s.request(ctx, strings.NewReader((data.Encode())))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to sync: %s", resp.Status)
	}

	// parse the respone.body and return a map string to any to show the json
	var result SyncReadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &result, nil
}

func (s *Sync) request(
	ctx context.Context,
	body io.Reader,
) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.todoist.com/api/v1/sync", body)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	return client.Do(req)
}

func (s *Sync) AddCommand(command Command) {
	s.Commands = append(s.Commands, command)
}
