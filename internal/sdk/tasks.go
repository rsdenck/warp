package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type TasksClient struct {
	client *Client
	token  string
}

func NewTasksClient(cfg *Config) *TasksClient {
	return &TasksClient{
		client: NewClient(cfg),
	}
}

func (c *TasksClient) SetToken(token string) {
	c.token = token
}

type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    int        `json:"priority"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date"`
	Completed   *time.Time `json:"completed"`
	Reminder    *time.Time `json:"reminder"`
	FolderID    string     `json:"folder_id"`
	Assignee    string     `json:"assignee"`
	Owner       string     `json:"owner"`
	Tags        []string   `json:"tags"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
}

type TaskFolder struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	ParentID string       `json:"parent_id"`
	Tasks    []Task       `json:"tasks"`
	Children []TaskFolder `json:"children,omitempty"`
}

func (c *TasksClient) ListTasks(folderID string) ([]Task, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "tasks.list",
		"token":     c.token,
		"folder_id": folderID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Tasks []Task `json:"tasks"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Tasks, nil
}

func (c *TasksClient) ListFolders() ([]TaskFolder, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "tasks.folders.list",
		"token":  c.token,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Folders []TaskFolder `json:"folders"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Folders, nil
}

func (c *TasksClient) CreateTask(title, description, folderID, assignee string, priority int, dueDate *time.Time, tags []string) (*Task, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "tasks.create",
		"token":       c.token,
		"title":       title,
		"description": description,
		"folder_id":   folderID,
		"assignee":    assignee,
		"priority":    priority,
		"tags":        tags,
	}

	if dueDate != nil {
		payload["due_date"] = dueDate.Format(time.RFC3339)
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Task
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *TasksClient) UpdateTask(taskID string, title, description, status string, priority int, dueDate *time.Time, tags []string) (*Task, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "tasks.update",
		"token":       c.token,
		"task_id":     taskID,
		"title":       title,
		"description": description,
		"status":      status,
		"priority":    priority,
		"tags":        tags,
	}

	if dueDate != nil {
		payload["due_date"] = dueDate.Format(time.RFC3339)
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Task
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *TasksClient) CompleteTask(taskID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "tasks.complete",
		"token":   c.token,
		"task_id": taskID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *TasksClient) DeleteTask(taskID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "tasks.delete",
		"token":   c.token,
		"task_id": taskID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *TasksClient) CreateFolder(name, parentID string) (*TaskFolder, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "tasks.folders.create",
		"token":     c.token,
		"name":      name,
		"parent_id": parentID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result TaskFolder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *TasksClient) DeleteFolder(folderID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "tasks.folders.delete",
		"token":     c.token,
		"folder_id": folderID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}
