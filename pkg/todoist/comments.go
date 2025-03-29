package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Attachment represents a file attached to a comment.
type Attachment struct {
	FileName     string `json:"file_name"`
	FileType     string `json:"file_type"`
	FileURL      string `json:"file_url"`
	ResourceType string `json:"resource_type"` // e.g., "file", "image", "audio", etc.
}

// Comment represents a Todoist comment on a task or project.
type Comment struct {
	ID         string      `json:"id"`
	TaskID     *string     `json:"task_id"`    // Pointer for nullability
	ProjectID  *string     `json:"project_id"` // Pointer for nullability
	PostedAt   string      `json:"posted_at"`  // RFC3339 format UTC
	Content    string      `json:"content"`
	Attachment *Attachment `json:"attachment"` // Pointer for nullability
}

// CommentFilters holds the required filter parameters for retrieving comments.
// Exactly one of TaskID or ProjectID must be provided.
type CommentFilters struct {
	TaskID    string `json:"task_id,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
}

// CommentOptions holds the parameters for creating a new comment.
// Exactly one of TaskID or ProjectID must be provided. Content is handled separately.
type CommentOptions struct {
	TaskID     string      `json:"task_id,omitempty"`
	ProjectID  string      `json:"project_id,omitempty"`
	Content    string      `json:"content"`              // Required content for the comment
	Attachment *Attachment `json:"attachment,omitempty"` // Pointer for optional attachment
}

// GetAllComments retrieves comments for a specific task or project.
// Exactly one of options.TaskID or options.ProjectID must be non-empty.
func (c *Client) GetAllComments(options *CommentOptions) ([]Comment, error) {
	res, err := c.request("GET", "/comments", nil, options)
	if err != nil {
		return nil, fmt.Errorf("failed to make get all comments request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, _ := io.ReadAll(res.Body)
		fmt.Printf("Response body: %s\n", resBody) // Print the response body for debugging
		return nil, fmt.Errorf("get all comments unexpected status code: %d", res.StatusCode)
	}

	var comments []Comment
	err = json.NewDecoder(res.Body).Decode(&comments)
	if err != nil {
		return nil, fmt.Errorf("failed to decode get all comments response: %w", err)
	}
	return comments, nil // Return the list of comments
}

// CreateComment creates a new comment on a task or project.
// content is required. Exactly one of options.TaskID or options.ProjectID must be non-empty.
func (c *Client) CreateComment(content string, options *CommentOptions) (*Comment, error) {
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

	res, err := c.request("POST", "/comments", body, nil)
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
func (c *Client) GetComment(commentID string) (*Comment, error) {
	if commentID == "" {
		return nil, fmt.Errorf("comment ID cannot be empty")
	}

	res, err := c.request("GET", fmt.Sprintf("/comments/%s", commentID), nil, nil)
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
	return &comment, nil // Return the comment object
}

func (c *Client) UpdateComment(commentID string, content string) (*Comment, error) {
	body := CommentOptions{
		Content: content,
	}

	res, err := c.request("POST", fmt.Sprintf("/comments/%s", commentID), body, nil)
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
	return &comment, nil // Return the updated comment object
}

// DeleteComment deletes a comment by its ID.
func (c *Client) DeleteComment(commentID string) error {
	endpoint := fmt.Sprintf("/comments/%s", commentID)

	res, err := c.request("DELETE", endpoint, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to make delete comment request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete comment unexpected status code: %d", res.StatusCode)
	}

	return nil // Success (204 No Content)
}
