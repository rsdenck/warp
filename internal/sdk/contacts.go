package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type ContactsClient struct {
	client *Client
	token  string
}

func NewContactsClient(cfg *Config) *ContactsClient {
	return &ContactsClient{
		client: NewClient(cfg),
	}
}

func (c *ContactsClient) SetToken(token string) {
	c.token = token
}

type Contact struct {
	ID        string    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Mobile    string    `json:"mobile"`
	Company   string    `json:"company"`
	JobTitle  string    `json:"job_title"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	Country   string    `json:"country"`
	ZipCode   string    `json:"zip_code"`
	WebSite   string    `json:"web_site"`
	Notes     string    `json:"notes"`
	Birthday  string    `json:"birthday"`
	Groups    []string  `json:"groups"`
	FolderID  string    `json:"folder_id"`
	Photo     string    `json:"photo"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
}

type ContactGroup struct {
	ID       string         `json:"id"`
	Name     string         `json:"name"`
	ParentID string         `json:"parent_id"`
	Contacts []Contact      `json:"contacts"`
	Children []ContactGroup `json:"children,omitempty"`
}

func (c *ContactsClient) ListContacts(folderID string) ([]Contact, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "contacts.list",
		"token":     c.token,
		"folder_id": folderID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Contacts []Contact `json:"contacts"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Contacts, nil
}

func (c *ContactsClient) ListGroups() ([]ContactGroup, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "contacts.groups.list",
		"token":  c.token,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Groups []ContactGroup `json:"groups"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Groups, nil
}

func (c *ContactsClient) CreateContact(contact *Contact) (*Contact, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":     "contacts.create",
		"token":      c.token,
		"first_name": contact.FirstName,
		"last_name":  contact.LastName,
		"email":      contact.Email,
		"phone":      contact.Phone,
		"mobile":     contact.Mobile,
		"company":    contact.Company,
		"job_title":  contact.JobTitle,
		"address":    contact.Address,
		"city":       contact.City,
		"state":      contact.State,
		"country":    contact.Country,
		"zip_code":   contact.ZipCode,
		"web_site":   contact.WebSite,
		"notes":      contact.Notes,
		"birthday":   contact.Birthday,
		"groups":     contact.Groups,
		"folder_id":  contact.FolderID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Contact
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ContactsClient) UpdateContact(contactID string, contact *Contact) (*Contact, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":     "contacts.update",
		"token":      c.token,
		"contact_id": contactID,
		"first_name": contact.FirstName,
		"last_name":  contact.LastName,
		"email":      contact.Email,
		"phone":      contact.Phone,
		"mobile":     contact.Mobile,
		"company":    contact.Company,
		"job_title":  contact.JobTitle,
		"address":    contact.Address,
		"city":       contact.City,
		"state":      contact.State,
		"country":    contact.Country,
		"zip_code":   contact.ZipCode,
		"web_site":   contact.WebSite,
		"notes":      contact.Notes,
		"birthday":   contact.Birthday,
		"groups":     contact.Groups,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Contact
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ContactsClient) DeleteContact(contactID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":     "contacts.delete",
		"token":      c.token,
		"contact_id": contactID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *ContactsClient) CreateGroup(name, parentID string) (*ContactGroup, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":    "contacts.groups.create",
		"token":     c.token,
		"name":      name,
		"parent_id": parentID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result ContactGroup
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ContactsClient) DeleteGroup(groupID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":   "contacts.groups.delete",
		"token":    c.token,
		"group_id": groupID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}
