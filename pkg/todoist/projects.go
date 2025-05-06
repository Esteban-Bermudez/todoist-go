package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Both the WorkspaceProject and PersonalProject are
// combined into one as I felt there is not a need to have them as separate
// since they have so many similar keys. This also avoids having to add
// more logic to identify the types and have them as different types in one
// array.
type Project struct {
	ID             string  `json:"id"`
	CanAssignTasks bool    `json:"can_assign_tasks"`
	ChildOrder     int     `json:"child_order"`
	Color          string  `json:"color"`
	CreatedAt      *string `json:"created_at"`
	IsArchived     bool    `json:"is_archived"`
	IsDeleted      bool    `json:"is_deleted"`
	IsFavorite     bool    `json:"is_favorite"`
	IsFrozen       bool    `json:"is_frozen"`
	Name           string  `json:"name"`
	UpdatedAt      *string `json:"updated_at"`
	ViewStyle      string  `json:"view_style"`
	DefaultOrder   int     `json:"default_order"`
	Description    string  `json:"description"`
	Access         *struct {
		Visibility    string `json:"visibility"`
		Configuration any    `json:"configuration"`
	} `json:"access"`
	CollaboratorRoleDefault string  `json:"collaborator_role_default"`
	FolderID                *string `json:"folder_id,omitempty"`
	IsInviteOnly            *bool   `json:"is_invite_only"`
	IsLinkSharingEnabled    bool    `json:"is_link_sharing_enabled,omitempty"`
	Role                    string  `json:"role,omitempty"`
	Status                  string  `json:"status,omitempty"`
	WorkspaceID             string  `json:"workspace_id,omitempty"`

	ParentID     *string `json:"parent_id,omitempty"`
	InboxProject bool    `json:"inbox_project,omitempty"`
	IsCollapsed  bool    `json:"is_collapsed"`
	IsShared     bool    `json:"is_shared"`
}

type ProjectResponse struct {
	NextCursor *string   `json:"next_cursor"`
	Results    []Project `json:"results"`
}

type ProjectOptions struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	ParentID    *string `json:"parent_id,omitempty"` // Nullable
	Color       string  `json:"color,omitempty"`
	IsFavorite  bool    `json:"is_favorite,omitempty"`
	ViewStyle   *string `json:"view_style,omitempty"`
}

type ProjectCollaboratorResponse struct {
	NextCursor *string        `json:"next_cursor"`
	Results    []Collaborator `json:"results"`
}

type Collaborator struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *Client) GetProjects(pagination PaginationOptions) ([]Project, *string, error) {
	res, err := c.request("GET", "/projects/", nil, pagination)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var response ProjectResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return nil, nil, err
	}
	return response.Results, response.NextCursor, nil
}

func (c *Client) GetArchived(pagination PaginationOptions) ([]Project, *string, error) {
	res, err := c.request("GET", "/projects/archived", nil, pagination)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var pr ProjectResponse
	err = json.NewDecoder(res.Body).Decode(&pr)
	if err != nil {
		return nil, nil, err
	}
	return pr.Results, pr.NextCursor, nil
}

func (c *Client) CreateProject(name string, options *ProjectOptions) (*Project, error) {
	body := options
	body.Name = &name

	res, err := c.request("POST", "/projects", body, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var project Project
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) GetProject(projectId string) (*Project, error) {
	res, err := c.request("GET", fmt.Sprintf("/projects/%s", projectId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var project Project
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) UpdateProject(projectId string, options *ProjectOptions) (*Project, error) {
	res, err := c.request("POST", fmt.Sprintf("/projects/%s", projectId), options, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var project Project
	err = json.NewDecoder(res.Body).Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) ArchiveProject(projectId string) error {
	res, err := c.request("POST", fmt.Sprintf("/projects/%s/archive", projectId), nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) UnarchiveProject(projectId string) error {
	res, err := c.request("POST", fmt.Sprintf("/projects/%s/unarchive", projectId), nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return nil
}

func (c *Client) DeleteProject(projectId string) error {
	res, err := c.request("DELETE", fmt.Sprintf("/projects/%s", projectId), nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

func (c *Client) GetProjectCollaborators(
	projectId string,
	pagination PaginationOptions,
) ([]Collaborator, *string, error) {
	res, err := c.request(
		"GET",
		fmt.Sprintf("/projects/%s/collaborators", projectId),
		nil,
		pagination,
	)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var pcr ProjectCollaboratorResponse
	err = json.NewDecoder(res.Body).Decode(&pcr)
	if err != nil {
		return nil, nil, err
	}
	return pcr.Results, pcr.NextCursor, nil
}

// TODO: Implement Projects Join Endpoint
// https://developer.todoist.com/api/v1#tag/Projects/operation/join_api_v1_projects__project_id__join_post
