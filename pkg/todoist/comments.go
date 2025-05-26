package todoist

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Comment represents a Todoist comment on a task or project.
type Comment struct {
	ID             string          `json:"id"`
	PostedUID      *string         `json:"posted_uid"`
	Content        string          `json:"content"`
	FileAttachment *map[string]any `json:"file_attachment"`
	UIDsToNotify   *[]string       `json:"uids_to_notify"`
	IsDeleted      bool            `json:"is_deleted"`
	PostedAt       string          `json:"posted_at"`
	Reactions      *map[string]any `json:"reactions"`
}

// CommentFilters holds the required filter parameters for retrieving comments.
// Exactly one of TaskID or ProjectID must be provided.
type CommentFilters struct {
	TaskID    string `json:"task_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	PaginationFilters
}

// CommentOptions holds the parameters for creating a new comment.
// Exactly one of TaskID or ProjectID must be provided. Content is handled separately.
type CommentOptions struct {
	Content      string          `json:"content"` // Required content for the comment
	TaskID       string          `json:"task_id,omitempty"`
	ProjectID    string          `json:"project_id,omitempty"`
	Attachment   *map[string]any `json:"attachment,omitempty"` // Optional file attachment
	UIDsToNotify *[]int          `json:"uids_to_notify,omitempty"`
}

// GetComments returns a list of all comments for a given task_id or project_id
// Exactly one of filters.TaskID or filters.ProjectID must be non-empty.
func (c *Client) GetComments(
	ctx context.Context,
	filters CommentFilters,
) ([]Comment, *string, error) {
	if filters.TaskID == "" && filters.ProjectID == "" {
		return nil, nil, fmt.Errorf("either TaskID or ProjectID must be provided in filters")
	}
	if filters.TaskID != "" && filters.ProjectID != "" {
		return nil, nil, fmt.Errorf("provide either TaskID or ProjectID, not both")
	}

	res, err := c.request(ctx, "GET", "/comments", nil, filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make get all comments request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("get all comments unexpected status code: %d", res.StatusCode)
	}

	var pagiResp PaginationResponse[Comment]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, err
	}
	return pagiResp.Results, pagiResp.NextCursor, nil
}

// CreateComment creates a new comment on a task or project.
// content is required. Exactly one of options.TaskID or options.ProjectID must be non-empty.
func (c *Client) CreateComment(
	ctx context.Context,
	content string,
	options *CommentOptions,
) (*Comment, error) {
	if content == "" {
		return nil, fmt.Errorf("comment content cannot be empty")
	}
	if options == nil || (options.TaskID == "" && options.ProjectID == "") {
		return nil, fmt.Errorf("either TaskID or ProjectID must be provided in options")
	}
	if options.TaskID != "" && options.ProjectID != "" {
		return nil, fmt.Errorf("provide either TaskID or ProjectID, not both")
	}

	body := options
	body.Content = content // Set the content in the options struct

	res, err := c.request(ctx, "POST", "/comments", body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make create comment request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("create comment unexpected status code: %d", res.StatusCode)
	}

	var comment Comment
	err = json.NewDecoder(res.Body).Decode(&comment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode create comment response: %w", err)
	}
	return &comment, nil // Return the created comment object
}

// GetComment retrieves a single comment by its ID.
func (c *Client) GetComment(ctx context.Context, commentID string) (*Comment, error) {
	if commentID == "" {
		return nil, fmt.Errorf("comment ID cannot be empty")
	}

	res, err := c.request(ctx, "GET", fmt.Sprintf("/comments/%s", commentID), nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make get comment request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get comment unexpected status code: %d", res.StatusCode)
	}

	var comment Comment
	err = json.NewDecoder(res.Body).Decode(&comment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode get comment response: %w", err)
	}
	return &comment, nil
}

// UpdateComment updates an existing comment's content by its ID.
func (c *Client) UpdateComment(
	ctx context.Context,
	commentID string,
	content string,
) (*Comment, error) {
	body := CommentOptions{
		Content: content,
	}

	res, err := c.request(ctx, "POST", fmt.Sprintf("/comments/%s", commentID), body, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make update comment request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update comment unexpected status code: %d", res.StatusCode)
	}

	var comment Comment
	err = json.NewDecoder(res.Body).Decode(&comment)
	if err != nil {
		return nil, fmt.Errorf("failed to decode update comment response: %w", err)
	}
	return &comment, nil
}

// DeleteComment deletes a comment by its ID.
func (c *Client) DeleteComment(ctx context.Context, commentID string) error {
	endpoint := fmt.Sprintf("/comments/%s", commentID)

	res, err := c.request(ctx, "DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to make delete comment request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("delete comment unexpected status code: %d", res.StatusCode)
	}

	return nil
}
