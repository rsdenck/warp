package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type NotesClient struct {
	client *Client
	token  string
}

func NewNotesClient(cfg *Config) *NotesClient {
	return &NotesClient{
		client: NewClient(cfg),
	}
}

func (c *NotesClient) SetToken(token string) {
	c.token = token
}

type Note struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Content  string    `json:"content"`
	Color    string    `json:"color"`
	FolderID string    `json:"folder_id"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
	Tags     []string  `json:"tags"`
}

type NoteFolder struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	ParentID string       `json:"parent_id"`
	Notes    []Note       `json:"notes"`
	Children []NoteFolder `json:"children,omitempty"`
}

func (c *NotesClient) ListNotes(folderID string) ([]Note, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "notes.list",
		"token":     c.token,
		"folder_id": folderID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Notes []Note `json:"notes"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Notes, nil
}

func (c *NotesClient) ListFolders() ([]NoteFolder, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "notes.folders.list",
		"token":  c.token,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Folders []NoteFolder `json:"folders"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Folders, nil
}

func (c *NotesClient) CreateNote(title, content, color, folderID string, tags []string) (*Note, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "notes.create",
		"token":     c.token,
		"title":     title,
		"content":   content,
		"color":     color,
		"folder_id": folderID,
		"tags":      tags,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Note
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *NotesClient) UpdateNote(noteID string, title, content, color string, tags []string) (*Note, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "notes.update",
		"token":   c.token,
		"note_id": noteID,
		"title":   title,
		"content": content,
		"color":   color,
		"tags":    tags,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Note
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *NotesClient) DeleteNote(noteID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "notes.delete",
		"token":   c.token,
		"note_id": noteID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *NotesClient) CreateFolder(name, parentID string) (*NoteFolder, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "notes.folders.create",
		"token":     c.token,
		"name":      name,
		"parent_id": parentID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result NoteFolder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *NotesClient) DeleteFolder(folderID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "notes.folders.delete",
		"token":     c.token,
		"folder_id": folderID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}
