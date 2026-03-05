package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// PinchtabClient provides browser automation capabilities
type PinchtabClient struct {
	baseURL string
	token   string
	client  *http.Client
}

// Tab represents a browser tab
type Tab struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Title  string `json:"title"`
	Active bool   `json:"active"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
}

// SnapshotResponse represents the accessibility tree
type SnapshotResponse struct {
	Elements []Element `json:"elements"`
	URL      string    `json:"url"`
	Title    string    `json:"title"`
}

// Element represents an accessibility tree element
type Element struct {
	Ref         string `json:"ref"`
	Tag         string `json:"tag"`
	Text        string `json:"text"`
	Role        string `json:"role"`
	Clickable   bool   `json:"clickable"`
	Focusable   bool   `json:"focusable"`
	Visible     bool   `json:"visible"`
	Attributes  map[string]string `json:"attributes"`
}

// ActionRequest represents an action to perform
type ActionRequest struct {
	Action   string `json:"action"`   // click, type, fill, press, hover, select, scroll
	Ref      string `json:"ref,omitempty"`
	Selector string `json:"selector,omitempty"`
	Text     string `json:"text,omitempty"`
	Key      string `json:"key,omitempty"`
}

// NavigateRequest represents a navigation request
type NavigateRequest struct {
	URL   string `json:"url"`
	TabID string `json:"tab_id,omitempty"`
}

// NewPinchtabClient creates a new Pinchtab client
func NewPinchtabClient(baseURL, token string) *PinchtabClient {
	return &PinchtabClient{
		baseURL: baseURL,
		token:   token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Health checks if Pinchtab is running
func (p *PinchtabClient) Health() (*HealthResponse, error) {
	resp, err := p.get("/health")
	if err != nil {
		return nil, err
	}

	var health HealthResponse
	if err := json.Unmarshal(resp, &health); err != nil {
		return nil, fmt.Errorf("failed to parse health response: %w", err)
	}

	return &health, nil
}

// ListTabs returns all open tabs
func (p *PinchtabClient) ListTabs() ([]Tab, error) {
	resp, err := p.get("/tabs")
	if err != nil {
		return nil, err
	}

	var tabs []Tab
	if err := json.Unmarshal(resp, &tabs); err != nil {
		return nil, fmt.Errorf("failed to parse tabs response: %w", err)
	}

	return tabs, nil
}

// Navigate to a URL
func (p *PinchtabClient) Navigate(url string, tabID ...string) error {
	req := NavigateRequest{URL: url}
	if len(tabID) > 0 {
		req.TabID = tabID[0]
	}

	_, err := p.post("/navigate", req)
	return err
}

// GetSnapshot returns the accessibility tree
func (p *PinchtabClient) GetSnapshot(tabID ...string) (*SnapshotResponse, error) {
	endpoint := "/snapshot"
	if len(tabID) > 0 {
		endpoint += "?tab_id=" + tabID[0]
	}

	resp, err := p.get(endpoint)
	if err != nil {
		return nil, err
	}

	var snapshot SnapshotResponse
	if err := json.Unmarshal(resp, &snapshot); err != nil {
		return nil, fmt.Errorf("failed to parse snapshot response: %w", err)
	}

	return &snapshot, nil
}

// PerformAction executes an action on the page
func (p *PinchtabClient) PerformAction(action ActionRequest) error {
	_, err := p.post("/action", action)
	return err
}

// Click an element by reference or selector
func (p *PinchtabClient) Click(refOrSelector string) error {
	action := ActionRequest{Action: "click"}
	if refOrSelector[0] == 'e' {
		action.Ref = refOrSelector
	} else {
		action.Selector = refOrSelector
	}
	return p.PerformAction(action)
}

// Type text into an element
func (p *PinchtabClient) Type(refOrSelector, text string) error {
	action := ActionRequest{Action: "type", Text: text}
	if refOrSelector[0] == 'e' {
		action.Ref = refOrSelector
	} else {
		action.Selector = refOrSelector
	}
	return p.PerformAction(action)
}

// Fill a form field
func (p *PinchtabClient) Fill(refOrSelector, text string) error {
	action := ActionRequest{Action: "fill", Text: text}
	if refOrSelector[0] == 'e' {
		action.Ref = refOrSelector
	} else {
		action.Selector = refOrSelector
	}
	return p.PerformAction(action)
}

// GetText returns readable text from the page
func (p *PinchtabClient) GetText(mode string, tabID ...string) (string, error) {
	endpoint := "/text?mode=" + mode
	if len(tabID) > 0 {
		endpoint += "&tab_id=" + tabID[0]
	}

	resp, err := p.get(endpoint)
	if err != nil {
		return "", err
	}

	return string(resp), nil
}

// Screenshot takes a screenshot of the current page
func (p *PinchtabClient) Screenshot(quality int, tabID ...string) ([]byte, error) {
	endpoint := fmt.Sprintf("/screenshot?quality=%d", quality)
	if len(tabID) > 0 {
		endpoint += "&tab_id=" + tabID[0]
	}

	return p.get(endpoint)
}

// Helper methods for HTTP requests
func (p *PinchtabClient) get(endpoint string) ([]byte, error) {
	req, err := http.NewRequest("GET", p.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (p *PinchtabClient) post(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", p.baseURL+endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.token != "" {
		req.Header.Set("Authorization", "Bearer "+p.token)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}