package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Order      *int   `json:"order"`
	IsFavorite bool   `json:"is_favorite"`
}

type SharedLabelFilters struct {
	OmitPersonal bool `json:"omit_personal,omitempty"`
	PaginationFilters
}

type SharedLabelOptions struct {
	NewName string `json:"new_name,omitempty"`
}

type LabelOptions struct {
	Name       string `json:"name,omitempty"`
	Order      int    `json:"order,omitempty"`
	Color      string `json:"color,omitempty"`
	IsFavorite bool   `json:"is_favorite,omitempty"`
}

// SharedLabels returns a set of unique strings containing labels from active tasks.
// By default, the names of a user's personal labels will also be included.
// These can be excluded by passing the OmitPersonal field in the
// SharedLabelFilters.
func (c *Client) SharedLabels(filters *SharedLabelFilters) ([]string, *string, error) {
	res, err := c.request("GET", "/labels/shared", filters, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get shared labels: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get shared labels unexpected status code: %d", res.StatusCode)
	}

	var pagiResp PaginationResponse[string]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode shared labels response: %w", err)
	}

	return pagiResp.Results, pagiResp.NextCursor, nil
}

// GetLabels returns a list of all user labels.
func (c *Client) GetLabels(filters *PaginationFilters) ([]Label, *string, error) {
	res, err := c.request("GET", "/labels", nil, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get labels: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get labels unexpected status code: %d", res.StatusCode)
	}

	var pagiResp PaginationResponse[Label]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode labels response: %w", err)
	}

	return pagiResp.Results, pagiResp.NextCursor, nil
}

// CreateLabel creates a new personal label with the given name.
// The name is required and will override the value in the LabelOptions.
func (c *Client) CreateLabel(name string, options *LabelOptions) (*Label, error) {
	if name == "" {
		return nil, fmt.Errorf("label name is required")
	}

	if options == nil {
		options = &LabelOptions{}
	}

	options.Name = name

	res, err := c.request("POST", "/labels", options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create label unexpected status code: %d", res.StatusCode)
	}

	var label Label
	err = json.NewDecoder(res.Body).Decode(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to decode label response: %w", err)
	}

	return &label, nil
}

// SharedLabelsRemove removes the given shared label from all active tasks. If
// no instances of the label name are found, the request will still be
// considered successful.
func (c *Client) SharedLabelsRemove(name string) error {
	if name == "" {
		return fmt.Errorf("label name is required")
	}

	res, err := c.request("POST", "/labels/shared/remove", LabelOptions{Name: name}, nil)
	if err != nil {
		return fmt.Errorf("failed to remove shared label: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("remove shared label unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// SharedLabelsRename renames the given shared label from all active tasks.
func (c *Client) SharedLabelsRename(name string, newName string) error {
	if name == "" {
		return fmt.Errorf("label name is required")
	}
	if newName == "" {
		return fmt.Errorf("new label name is required")
	}

	res, err := c.request(
		"POST",
		"/labels/shared/rename",
		SharedLabelOptions{NewName: newName},
		map[string]string{"name": name},
	)
	if err != nil {
		return fmt.Errorf("failed to rename shared label: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("rename shared label unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// DeleteLabel deletes a personal label by its ID. All instances of the label
// will be removed from tasks.
func (c *Client) DeleteLabel(id string) error {
	if id == "" {
		return fmt.Errorf("label ID is required")
	}

	res, err := c.request("DELETE", fmt.Sprintf("/labels/%s", id), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("delete label unexpected status code: %d", res.StatusCode)
	}

	return nil
}

// GetLabel returns a label by its ID.
func (c *Client) GetLabel(id string) (*Label, error) {
	if id == "" {
		return nil, fmt.Errorf("label ID is required")
	}

	res, err := c.request("GET", fmt.Sprintf("/labels/%s", id), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get label unexpected status code: %d", res.StatusCode)
	}

	var label Label
	err = json.NewDecoder(res.Body).Decode(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to decode label response: %w", err)
	}

	return &label, nil
}

// UpdateLabel updates a label by its ID. The ID is required, and the
// LabelOptions are required to specify the fields to update.
func (c *Client) UpdateLabel(id string, options *LabelOptions) (*Label, error) {
	if id == "" {
		return nil, fmt.Errorf("label ID is required")
	}
	if options == nil {
		return nil, fmt.Errorf("options are required")
	}

	res, err := c.request("POST", fmt.Sprintf("/labels/%s", id), options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to update label: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update label unexpected status code: %d", res.StatusCode)
	}

	var label Label
	err = json.NewDecoder(res.Body).Decode(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to decode label response: %w", err)
	}

	return &label, nil
}
