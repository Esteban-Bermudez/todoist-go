package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	ID          string   `json:"id"`
	ProjectID   string   `json:"project_id"`
	SectionID   string   `json:"section_id"`
	Content     string   `json:"content"`
	Description string   `json:"description"`
	IsCompleted bool     `json:"is_completed"`
	Labels      []string `json:"labels"`
	ParentID    string   `json:"parent_id"`
	Order       int      `json:"order"`
	Priority    int      `json:"priority"`
	Due         struct {
		Date        string `json:"date"`
		IsRecurring bool   `json:"is_recurring"`
		DateTime    string `json:"datetime"`
		TimeZone    string `json:"timezone"`
		Lang        string `json:"lang"`
	} `json:"due"`
	Deadline struct {
		Date string `json:"date"`
		Lang string `json:"lang"`
	} `json:"deadline"`
	URL          string `json:"url"`
	CommentCount int    `json:"comment_count"`
	CreatedAt    string `json:"created_at"`
	CreatorID    string `json:"creator_id"`
	AssigneeID   string `json:"assignee_id"`
	AssignerID   string `json:"assigner_id"`
	Duration     struct {
		Amount int    `json:"amount"`
		Unit   string `json:"unit"`
	} `json:"duration"`
}

type TaskOptions struct {
	Content      string   `json:"content,omitempty"`
	Description  string   `json:"description,omitempty"`
	ProjectID    string   `json:"project_id,omitempty"`
	SectionID    string   `json:"section_id,omitempty"`
	ParentID     string   `json:"parent_id,omitempty"`
	Order        int      `json:"order,omitempty"`
	Labels       []string `json:"labels,omitempty"`
	Priority     int      `json:"priority,omitempty"`
	DueString    string   `json:"due_string,omitempty"`
	DueDate      string   `json:"due_date,omitempty"`
	DueDateTime  string   `json:"due_datetime,omitempty"`
	DueLang      string   `json:"due_lang,omitempty"`
	AssigneeID   string   `json:"assignee_id,omitempty"`
	Duration     int      `json:"duration,omitempty"`
	DurationUnit string   `json:"duration_unit,omitempty"`
	DeadlineDate string   `json:"deadline_date,omitempty"`
	DeadlineLang string   `json:"deadline_lang,omitempty"`
}

type TaskFilters struct {
	ProjectID string `json:"project_id,omitempty"`
	SectionID string `json:"section_id,omitempty"`
	Label     string `json:"label,omitempty"`
	Filter    string `json:"filter,omitempty"`
	Lang      string `json:"lang,omitempty"`
	IDs       []int  `json:"ids,omitempty"`
}

func (c *Client) GetActiveTasks(filters *TaskFilters) ([]Task, error) {
	res, err := c.request("GET", "/tasks", nil, filters)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var tasks []Task
	err = json.NewDecoder(res.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

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

func (c *Client) CreateTask(content string, options *TaskOptions) (*Task, error) {
	body := options
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
