package sdk

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MailClient struct {
	client *Client
}

func NewMailClient(cfg *Config) *MailClient {
	return &MailClient{
		client: NewClient(cfg),
	}
}

type ShowVersionResponse struct {
	Version   string `json:"version"`
	Build     string `json:"build"`
	StartTime string `json:"start_time"`
}

func (c *MailClient) ShowVersion() (*ShowVersionResponse, error) {
	resp, err := c.client.postRawJSON("/mailapi/ShowVersion", nil)
	if err != nil {
		return nil, err
	}

	var result ShowVersionResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type UploadedItem struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UploadTime string `json:"upload_time"`
	Path       string `json:"path"`
}

func (c *MailClient) GetUploadedItems() ([]UploadedItem, error) {
	payload := map[string]interface{}{
		"method": "GetUploadedItems",
	}

	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return nil, err
	}

	var result []UploadedItem
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

type FolderInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Messages int    `json:"messages"`
	Unread   int    `json:"unread"`
	Size     int64  `json:"size"`
}

func (c *MailClient) ListFolders() ([]FolderInfo, error) {
	payload := map[string]interface{}{
		"method": "ListFolders",
	}

	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return nil, err
	}

	var result []FolderInfo
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *MailClient) CreateFolder(name, parentID string) error {
	payload := map[string]interface{}{
		"method":    "CreateFolder",
		"folder":    name,
		"parent_id": parentID,
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}

func (c *MailClient) DeleteFolder(folderID string) error {
	payload := map[string]interface{}{
		"method":   "DeleteFolder",
		"folderid": folderID,
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}

func (c *MailClient) RenameFolder(folderID, newName string) error {
	payload := map[string]interface{}{
		"method":    "RenameFolder",
		"folderid":  folderID,
		"newfolder": newName,
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}

func (c *MailClient) MoveFolder(folderID, newParentID string) error {
	payload := map[string]interface{}{
		"method":      "MoveFolder",
		"folderid":    folderID,
		"newparentid": newParentID,
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}

func (c *MailClient) GetFolderInfo(folderID string) (*FolderInfo, error) {
	payload := map[string]interface{}{
		"method":   "GetFolderInfo",
		"folderid": folderID,
	}

	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return nil, err
	}

	var result FolderInfo
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type GroupRoot struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Type     string      `json:"type"`
	Children []GroupRoot `json:"children,omitempty"`
}

func (c *MailClient) ListGroupRoots() ([]GroupRoot, error) {
	payload := map[string]interface{}{
		"method": "ListGroupRoots",
	}

	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return nil, err
	}

	var result []GroupRoot
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *MailClient) UploadFile(fileName string, data []byte) (string, error) {
	payload := map[string]interface{}{
		"method": "UploadFile",
		"file":   fileName,
		"data":   data,
	}

	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return "", err
	}

	var result map[string]string
	if err := json.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if id, ok := result["id"]; ok {
		return id, nil
	}
	return "", nil
}

func (c *MailClient) DeleteUploadedFile(fileID string) error {
	payload := map[string]interface{}{
		"method": "DeleteUploadedFile",
		"fileid": fileID,
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}

func (c *MailClient) parseMethodResponse(methodName string, payload map[string]interface{}) (map[string]interface{}, error) {
	payload["method"] = methodName
	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *MailClient) GetMessages(folderID string, limit, offset int) ([]map[string]interface{}, error) {
	payload := map[string]interface{}{
		"method":   "GetMessages",
		"folderid": folderID,
		"limit":    limit,
		"offset":   offset,
	}

	resp, err := c.client.postJSON("/mailapi/", payload)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(string(resp), "[") {
		var result []map[string]interface{}
		if err := json.Unmarshal(resp, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}
		return result, nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if msgs, ok := result["messages"]; ok {
		if messages, ok := msgs.([]interface{}); ok {
			resultSlice := make([]map[string]interface{}, len(messages))
			for i, m := range messages {
				if mm, ok := m.(map[string]interface{}); ok {
					resultSlice[i] = mm
				}
			}
			return resultSlice, nil
		}
	}

	return []map[string]interface{}{result}, nil
}

func (c *MailClient) DeleteMessage(messageID string) error {
	payload := map[string]interface{}{
		"method":    "DeleteMessage",
		"messageid": messageID,
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}

func (c *MailClient) AppendMessage(folderID, subject, body, from, to string) error {
	payload := map[string]interface{}{
		"method":   "AppendMessage",
		"folderid": folderID,
		"message": map[string]interface{}{
			"subject": subject,
			"body":    body,
			"from":    from,
			"to":      to,
		},
	}

	_, err := c.client.postJSON("/mailapi/", payload)
	return err
}
