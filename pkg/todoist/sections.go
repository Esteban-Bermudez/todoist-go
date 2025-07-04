package todoist

import (
	"context"
	"encoding/json"
	"fmt"
)

// Section represents a section in Todoist. A section will always belong to a
// project. Sections are used to organize tasks within a project.
type Section struct {
	ID           string  `json:"id"`
	UserID       string  `json:"user_id"`
	ProjectID    string  `json:"project_id"` // ID of the project section belongs to
	AddedAt      string  `json:"added_at"`
	UpdatedAt    *string `json:"updated_at"`
	ArchivedAt   *string `json:"archived_at"`
	Name         string  `json:"name"`
	SectionOrder int     `json:"section_order"` // Section position among other sections from the same project
	IsArchived   bool    `json:"is_archived"`
	IsDeleted    bool    `json:"is_deleted"`
	IsCollapsed  bool    `json:"is_collapsed"`
}

// SectionFilters holds the required filter parameters for retrieving sections.
type SectionFilters struct {
	ProjectID string `json:"project_id,omitempty"`
	PaginationFilters
}

// SectionOptions holds the parameters for creating or updating a section.
type SectionOptions struct {
	Name      string `json:"name,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	Order     int    `json:"order,omitempty"`
}

// CreateSection creates a new section in the specified project. The parameters
// name and projectID are required. They will override the values in options.
func (c *Client) CreateSection(
	ctx context.Context,
	name string,
	projectID string,
	options *SectionOptions,
) (*Section, error) {
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if projectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}

	if options == nil {
		options = &SectionOptions{}
	}

	options.Name = name
	options.ProjectID = projectID

	res, err := c.request(ctx, "POST", "/sections", options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create section: %w", err)
	}
	defer res.Body.Close()

	var section Section
	err = json.NewDecoder(res.Body).Decode(&section)
	if err != nil {
		return nil, fmt.Errorf("failed to decode section response: %w", err)
	}

	return &section, nil
}

// GetSections returns a list of all active sections for the user
// or a specific project if filters.ProjectID is provided.
func (c *Client) GetSections(
	ctx context.Context,
	filters *SectionFilters,
) ([]Section, *string, error) {
	res, err := c.request(ctx, "GET", "/sections", nil, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get sections: %w", err)
	}
	defer res.Body.Close()

	var pagiResp PaginationResponse[Section]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode sections response: %w", err)
	}

	return pagiResp.Results, pagiResp.NextCursor, nil
}

// GetSection returns the section for the given section ID
func (c *Client) GetSection(ctx context.Context, id string) (*Section, error) {
	if id == "" {
		return nil, fmt.Errorf("section ID is required")
	}

	res, err := c.request(ctx, "GET", fmt.Sprintf("/sections/%s", id), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get section: %w", err)
	}
	defer res.Body.Close()

	var section Section
	err = json.NewDecoder(res.Body).Decode(&section)
	if err != nil {
		return nil, fmt.Errorf("failed to decode section response: %w", err)
	}

	return &section, nil
}

// UpdateSection updates the section name with the given ID.
func (c *Client) UpdateSection(ctx context.Context, id string, name string) (*Section, error) {
	if id == "" {
		return nil, fmt.Errorf("section ID is required")
	}
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}

	options := &SectionOptions{
		Name: name,
	}

	res, err := c.request(ctx, "POST", fmt.Sprintf("/sections/%s", id), options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update section: %w", err)
	}
	defer res.Body.Close()

	var section Section
	err = json.NewDecoder(res.Body).Decode(&section)
	if err != nil {
		return nil, fmt.Errorf("failed to decode section response: %w", err)
	}

	return &section, nil
}

// DeleteSection deletes the section with the given ID and all of its tasks.
func (c *Client) DeleteSection(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("section ID is required")
	}

	res, err := c.request(ctx, "DELETE", fmt.Sprintf("/sections/%s", id), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete section: %w", err)
	}
	defer res.Body.Close()

	return nil
}
