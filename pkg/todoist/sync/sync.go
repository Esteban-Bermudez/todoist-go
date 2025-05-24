package sync

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// sync_token String	Yes	A special string, used to allow the client to perform incremental sync. Pass * to retrieve all active resource data. More details about this below.
// resource_types JSON array of strings	Yes	Used to specify what resources to fetch from the server. It should be a JSON-encoded array of strings. Here is a list of available resource types: labels, projects, items, notes, sections, filters, reminders, reminders_location, locations, user, live_notifications, collaborators, user_settings, notification_settings, user_plan_limits, completed_info, stats, workspaces, workspace_users. You may use all to include all the resource types. Resources can also be excluded by prefixing a - prior to the name, for example, -projects
type Sync struct {
	SyncToken     string    `json:"sync_token"`
	ResourceTypes []string  `json:"resource_types"`
	Commands      []Command `json:"commands"`
	APIKey        string    `json:"-"`
}

// type String	The type of the command.
// args Object	The parameters of the command.
// uuid String	Command UUID. More details about this below.
// temp_id String	Temporary resource ID, Optional. Only specified for commands that create a new resource (e.g. item_add command). More details about this belowV
type Command struct {
	Type   string         `json:"type"`
	Args   map[string]any `json:"args"`
	UUID   string         `json:"uuid"`
	TempID string         `json:"temp_id,omitempty"`
}

func (s *Sync) ReadResources(resourceTypes []string) (any, error) {
	if resourceTypes != nil {
		s.ResourceTypes = resourceTypes
	} else {
		return nil, fmt.Errorf("resourceTypes cannot be nil")
	}

	data := url.Values{}
	data.Set("sync_token", s.SyncToken)
	data.Set("resource_types", `["`+strings.Join(s.ResourceTypes, `","`)+`"]`)

	resp, err := s.request(strings.NewReader((data.Encode())))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to sync: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return body, nil
}

func (s *Sync) request(body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", "https://api.todoist.com/api/v1/sync", body)
	if err != nil {
		panic(err)
	}
	// fmt.Println("Request Body:", body)

	req.Header.Set("Authorization", "Bearer "+s.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	return client.Do(req)
}

func (s *Sync) AddCommand(command Command) {
	s.Commands = append(s.Commands, command)
}
