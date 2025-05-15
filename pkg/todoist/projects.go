package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Project represents a Todoist project.
// Both the WorkspaceProject and PersonalProject types are
// combined into one as I felt there is not a need to have them as separate
// since they have so many similar keys. This also avoids having to add
// more logic to identify which types was returned and allows for storing an
// array of projects without having to worry about the type.
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

// ProjectOptions represents the body parameters for creating or updating a
// project.
type ProjectOptions struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parent_id,omitempty"`
	Color       string `json:"color,omitempty"`
	IsFavorite  bool   `json:"is_favorite,omitempty"`
	ViewStyle   string `json:"view_style,omitempty"`
}

// Collaborator represents a collaborator on a project.
type Collaborator struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetProjects returns a list containing all active user projects and a cursor
// for pagination. The cursor is nil if there are no more pages to return.
func (c *Client) GetProjects(pagination *PaginationFilters) ([]Project, *string, error) {
	res, err := c.request("GET", "/projects/", nil, pagination)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var pagiResp PaginationResponse[Project]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, err
	}
	return pagiResp.Results, pagiResp.NextCursor, nil
}

// GetArchived returns a list containing all archived user projects and a cursor
// for pagination. The cursor is nil if there are no more pages to return.
func (c *Client) GetArchived(pagination *PaginationFilters) ([]Project, *string, error) {
	res, err := c.request("GET", "/projects/archived", nil, pagination)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var pagiResp PaginationResponse[Project]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, err
	}
	return pagiResp.Results, pagiResp.NextCursor, nil
}

// CreateProject creates a new project with the given name and options.
// The name is required, and any additional options can be set in the
// ProjectOptions
func (c *Client) CreateProject(name string, options *ProjectOptions) (*Project, error) {
	body := options
	body.Name = name

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

// GetProject returns a project related to the given projectId.
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

// UpdateProject updates a project with the given projectId with the given
// ProjectOptions.
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

// ArchiveProject archives a project with the given projectId.
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

// UnarchiveProject unarchives a project with the given projectId.
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

// DeleteProject deletes a project with the given projectId.
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

// GetProjectCollaborators returns a list of collaborators for the given
// projectId and a cursor for pagination. The cursor is nil if there are no
// more pages to return.
func (c *Client) GetProjectCollaborators(
	projectId string,
	pagination *PaginationFilters,
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

	var pagiResp PaginationResponse[Collaborator]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, err
	}
	return pagiResp.Results, pagiResp.NextCursor, nil
}

// TODO: Implement Projects Join Endpoint
// https://developer.todoist.com/api/v1#tag/Projects/operation/join_api_v1_projects__project_id__join_post
