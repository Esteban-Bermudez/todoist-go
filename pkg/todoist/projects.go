package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Project struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Color        string  `json:"color"`
	ParentID     *string `json:"parent_id"` // Nullable
	Order        int     `json:"order"`
	CommentCount int     `json:"comment_count"`
	IsShared     bool    `json:"is_shared"`
	IsFavorite   bool    `json:"is_favorite"`
	IsInbox      bool    `json:"is_inbox_project"`
	IsTeamInbox  bool    `json:"is_team_inbox"`
	ViewStyle    string  `json:"view_style"`
	URL          string  `json:"url"`
}

type ProjectOptions struct {
	Name       string  `json:"name,omitempty"`
	ParentID   *string `json:"parent_id,omitempty"` // Nullable
	Color      string  `json:"color,omitempty"`
	IsFavorite bool    `json:"is_favorite,omitempty"`
	ViewStyle  string  `json:"view_style,omitempty"`
}

type Collaborator struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (c *Client) GetAllProjects() ([]Project, error) {
	res, err := c.request("GET", "/projects", nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var projects []Project
	err = json.NewDecoder(res.Body).Decode(&projects)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

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

func (c *Client) GetAllCollaborators(projectId string) ([]Collaborator, error) {
	res, err := c.request("GET", fmt.Sprintf("/projects/%s/collaborators", projectId), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var collaborators []Collaborator
	err = json.NewDecoder(res.Body).Decode(&collaborators)
	if err != nil {
		return nil, err
	}
	return collaborators, nil
}
