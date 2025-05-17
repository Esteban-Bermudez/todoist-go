package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Task represents a task in Todoist.
// The Task struct contains all the fields returned by the API.
type Task struct {
	UserID         string          `json:"user_id"`
	ID             string          `json:"id"`
	ProjectID      string          `json:"project_id"`
	SectionID      *string         `json:"section_id"`
	ParentID       *string         `json:"parent_id"`
	AddedByUID     *string         `json:"added_by_uid"`
	AssignedByUID  *string         `json:"assigned_by_uid"`
	ResponsibleUID *string         `json:"responsible_uid"`
	Labels         []string        `json:"labels"`
	Deadline       *map[string]any `json:"deadline"`
	Duration       *map[string]any `json:"duration"`
	Checked        bool            `json:"checked"`
	IsDeleted      bool            `json:"is_deleted"`
	AddedAt        *string         `json:"added_at"`
	CompletedAt    *string         `json:"completed_at"`
	UpdatedAt      *string         `json:"updated_at"`
	Due            *map[string]any `json:"due"`
	Priority       int             `json:"priority"`
	ChildOrder     int             `json:"child_order"`
	Content        string          `json:"content"`
	Description    string          `json:"description"`
	NoteCount      int             `json:"note_count"`
	DayOrder       int             `json:"day_order"`
	IsCollapsed    bool            `json:"is_collapsed"`
}

// TaskOptions represents the body parameters for creating or updating a task.
type TaskOptions struct {
	Content      string   `json:"content,omitempty"`
	Description  string   `json:"description,omitempty"`
	ProjectID    string   `json:"project_id,omitempty"`
	SectionID    string   `json:"section_id,omitempty"`
	ParentID     string   `json:"parent_id,omitempty"`
	Order        int      `json:"order,omitempty"`
	Labels       []string `json:"labels,omitempty"`
	Priority     int      `json:"priority,omitempty"`
	AssigneeID   string   `json:"assignee_id,omitempty"`
	DueString    string   `json:"due_string,omitempty"`
	DueDate      string   `json:"due_date,omitempty"`
	DueDateTime  string   `json:"due_datetime,omitempty"`
	DueLang      string   `json:"due_lang,omitempty"`
	Duration     int      `json:"duration,omitempty"` // If duration is set, duration_unit must also be set
	DurationUnit string   `json:"duration_unit,omitempty"` // The unit of the duration has to be "days" or "minutes".
	DeadlineDate string   `json:"deadline_date,omitempty"`
	DeadlineLang string   `json:"deadline_lang,omitempty"`
}

// TaskFilters represents the query parameters for filtering tasks.
type TaskFilters struct {
	ProjectID string `json:"project_id,omitempty"`
	SectionID string `json:"section_id,omitempty"`
	ParentID  string `json:"parent_id,omitempty"`
	Label     string `json:"label,omitempty"`
	Filter    string `json:"filter,omitempty"`
	Lang      string `json:"lang,omitempty"`
	IDs       string `json:"ids,omitempty"` //A list of the task IDs to retrieve, this should be a comma separated list
	PaginationFilters
}

func (c *Client) CreateTask(content string, options *TaskOptions) (*Task, error) {
	body := options
	if body == nil {
		body = &TaskOptions{}
	}
	body.Content = content

	res, err := c.request("POST", "/tasks", body, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var task Task
	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// GetTasks returns a list containing all active user tasks and a cursor for
// pagination. The cursor is nil if there are no more pages to return.
func (c *Client) GetTasks(filters *TaskFilters) ([]Task, *string, error) {
	res, err := c.request("GET", "/tasks", nil, filters)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var pagiResp PaginationResponse[Task]
	err = json.NewDecoder(res.Body).Decode(&pagiResp)
	if err != nil {
		return nil, nil, err
	}
	return pagiResp.Results, pagiResp.NextCursor, nil
}

// Tasks Completed By Completion Date

// Tasks Completed By Due Date

// Get Tasks By Filter

// QuickAddTask creates a new task using the quick add feature. This is what
// Todoist uses to create tasks with natural language processing. The text
// parameter is the text of the task to create. The options parameter is
// optional if you want to set additional options for the task.
func (c *Client) QuickAddTask(text string, options *TaskOptions) (*Task, error) {
	body := struct {
		Text string `json:"text"`
	}{
		Text: text,
	}
	res, err := c.request("POST", "/tasks/quick", body, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	var task Task
	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// ReopenTask reopens a task that has been completed. The taskID parameter is
// the ID of the task to reopen.
func (c *Client) ReopenTask(taskID string) error {
	res, err := c.request("POST", fmt.Sprintf("/tasks/%s/reopen", taskID), nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

// CloseTask closes a task that has been completed. The taskID parameter is
// the ID of the task to close.
func (c *Client) CloseTask(taskID string) error {
	res, err := c.request("POST", fmt.Sprintf("/tasks/%s/close", taskID), nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

// Move Task

// GetTask returns a task related to the given taskID. The taskID parameter
// is the ID of the task to get. The taskID parameter is required.
func (c *Client) GetTask(taskID string) (*Task, error) {
	res, err := c.request("GET", fmt.Sprintf("/tasks/%s", taskID), nil, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var task Task
	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// UpdateTask updates a task with the given taskID with the given TaskOptions.
func (c *Client) UpdateTask(taskID string, options *TaskOptions) (*Task, error) {
	res, err := c.request("POST", fmt.Sprintf("/tasks/%s", taskID), options, nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var task Task
	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

// DeleteTask deletes a task with the given taskID. The taskID parameter is the
// ID of the task to delete. The taskID parameter is required.
func (c *Client) DeleteTask(taskID string) error {
	res, err := c.request("DELETE", fmt.Sprintf("/tasks/%s", taskID), nil, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}
