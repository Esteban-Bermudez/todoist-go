package todoist

import (
	"context"
	"encoding/json"
	"fmt"
)

// Label represents a Todoist label. We do not need to have a struct for
// SharedLabel as they are represented as a []string in the API.
type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Order      *int   `json:"order"`
	IsFavorite bool   `json:"is_favorite"`
}

// SharedLabelFilters holds the required filter parameters for retrieving shared
// labels.
type SharedLabelFilters struct {
	OmitPersonal bool `json:"omit_personal,omitempty"`
	PaginationFilters
}

// SharedLabelOptions holds the parameters for renaming or removing a shared
// label. This is only used for the SharedLabelsRename.
type SharedLabelOptions struct {
	NewName string `json:"new_name,omitempty"`
}

// LabelOptions holds the parameters for creating or updating a label.
type LabelOptions struct {
	Name       string `json:"name,omitempty"`
	Order      int    `json:"order,omitempty"`
	Color      string `json:"color,omitempty"`
	IsFavorite bool   `json:"is_favorite,omitempty"`
}

// SharedLabels returns a set of unique strings containing labels from active
// tasks.
// By default, the names of a user's personal labels will also be included.
// These can be excluded by passing the OmitPersonal field in the
// SharedLabelFilters.
func (c *Client) SharedLabels(
	ctx context.Context,
	filters *SharedLabelFilters,
) ([]string, *string, error) {
	res, err := c.request(ctx, "GET", "/labels/shared", filters, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get shared labels: %w", err)
	}
	defer res.Body.Close()

	var pagiResp PaginationResponse[string]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to decode shared labels response: %w",
			err,
		)
	}

	return pagiResp.Results, pagiResp.NextCursor, nil
}

// GetLabels returns a list of all user labels.
func (c *Client) GetLabels(
	ctx context.Context,
	filters *PaginationFilters,
) ([]Label, *string, error) {
	res, err := c.request(ctx, "GET", "/labels", nil, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get labels: %w", err)
	}
	defer res.Body.Close()

	var pagiResp PaginationResponse[Label]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode labels response: %w", err)
	}

	return pagiResp.Results, pagiResp.NextCursor, nil
}

// CreateLabel creates a new personal label with the given name.
// The name is required and will override the value in the LabelOptions.
func (c *Client) CreateLabel(
	ctx context.Context,
	name string,
	options *LabelOptions,
) (*Label, error) {
	if name == "" {
		return nil, fmt.Errorf("label name is required")
	}

	if options == nil {
		options = &LabelOptions{}
	}

	options.Name = name

	res, err := c.request(ctx, "POST", "/labels", options, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create label: %w", err)
	}
	defer res.Body.Close()

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
func (c *Client) SharedLabelsRemove(ctx context.Context, name string) error {
	if name == "" {
		return fmt.Errorf("label name is required")
	}

	res, err := c.request(
		ctx,
		"POST",
		"/labels/shared/remove",
		LabelOptions{Name: name},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to remove shared label: %w", err)
	}
	defer res.Body.Close()

	return nil
}

// SharedLabelsRename renames the given shared label from all active tasks.
func (c *Client) SharedLabelsRename(
	ctx context.Context,
	name string,
	newName string,
) error {
	if name == "" {
		return fmt.Errorf("label name is required")
	}
	if newName == "" {
		return fmt.Errorf("new label name is required")
	}

	res, err := c.request(
		ctx,
		"POST",
		"/labels/shared/rename",
		SharedLabelOptions{NewName: newName},
		map[string]string{"name": name},
	)
	if err != nil {
		return fmt.Errorf("failed to rename shared label: %w", err)
	}
	defer res.Body.Close()

	return nil
}

// DeleteLabel deletes a personal label by its ID. All instances of the label
// will be removed from tasks.
func (c *Client) DeleteLabel(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("label ID is required")
	}

	res, err := c.request(
		ctx,
		"DELETE",
		fmt.Sprintf("/labels/%s", id),
		nil,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to delete label: %w", err)
	}
	defer res.Body.Close()

	return nil
}

// GetLabel returns a label by its ID.
func (c *Client) GetLabel(ctx context.Context, id string) (*Label, error) {
	if id == "" {
		return nil, fmt.Errorf("label ID is required")
	}

	res, err := c.request(ctx, "GET", fmt.Sprintf("/labels/%s", id), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get label: %w", err)
	}
	defer res.Body.Close()

	var label Label
	err = json.NewDecoder(res.Body).Decode(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to decode label response: %w", err)
	}

	return &label, nil
}

// UpdateLabel updates a label by its ID. The ID is required, and the
// LabelOptions are required to specify the fields to update.
func (c *Client) UpdateLabel(
	ctx context.Context,
	id string,
	options *LabelOptions,
) (*Label, error) {
	if id == "" {
		return nil, fmt.Errorf("label ID is required")
	}
	if options == nil {
		return nil, fmt.Errorf("options are required")
	}

	res, err := c.request(
		ctx,
		"POST",
		fmt.Sprintf("/labels/%s", id),
		options,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update label: %w", err)
	}
	defer res.Body.Close()

	var label Label
	err = json.NewDecoder(res.Body).Decode(&label)
	if err != nil {
		return nil, fmt.Errorf("failed to decode label response: %w", err)
	}

	return &label, nil
}
