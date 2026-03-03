package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type FilesClient struct {
	client *Client
	token  string
}

func NewFilesClient(cfg *Config) *FilesClient {
	return &FilesClient{
		client: NewClient(cfg),
	}
}

func (c *FilesClient) SetToken(token string) {
	c.token = token
}

type FileItem struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Type       string    `json:"type"`
	Size       int64     `json:"size"`
	MimeType   string    `json:"mime_type"`
	Modified   time.Time `json:"modified"`
	Created    time.Time `json:"created"`
	Owner      string    `json:"owner"`
	Shared     bool      `json:"shared"`
	SharedLink string    `json:"shared_link,omitempty"`
}

type Folder struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	ParentID string    `json:"parent_id"`
	Owner    string    `json:"owner"`
	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (c *FilesClient) ListFiles(path string) ([]FileItem, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "files.list",
		"token":  c.token,
		"path":   path,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Files []FileItem `json:"files"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Files, nil
}

func (c *FilesClient) ListFolders(path string) ([]Folder, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "files.folders.list",
		"token":  c.token,
		"path":   path,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Folders []Folder `json:"folders"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Folders, nil
}

func (c *FilesClient) CreateFolder(name, parentPath string) (*Folder, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "files.folders.create",
		"token":       c.token,
		"name":        name,
		"parent_path": parentPath,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Folder
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *FilesClient) DeleteFile(fileID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "files.delete",
		"token":   c.token,
		"file_id": fileID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *FilesClient) DeleteFolder(folderID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "files.folders.delete",
		"token":     c.token,
		"folder_id": folderID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *FilesClient) UploadFile(name, parentPath string, data []byte) (*FileItem, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "files.upload",
		"token":       c.token,
		"name":        name,
		"parent_path": parentPath,
		"data":        data,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result FileItem
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *FilesClient) DownloadFile(fileID string) ([]byte, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "files.download",
		"token":   c.token,
		"file_id": fileID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *FilesClient) ShareFile(fileID string) (string, error) {
	if c.token == "" {
		return "", fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "files.share",
		"token":   c.token,
		"file_id": fileID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return "", err
	}

	var result struct {
		Link string `json:"link"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Link, nil
}

func (c *FilesClient) MoveFile(fileID, newPath string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":   "files.move",
		"token":    c.token,
		"file_id":  fileID,
		"new_path": newPath,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *FilesClient) CopyFile(fileID, newPath string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":   "files.copy",
		"token":    c.token,
		"file_id":  fileID,
		"new_path": newPath,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}
