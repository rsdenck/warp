package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	BaseURL     string
	APIKey      string
	Username    string
	Password    string
	Timeout     time.Duration
	InsecureTLS bool
}

type Client struct {
	config     *Config
	httpClient *http.Client
	BaseURL    string
}

func NewClient(cfg *Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://icewarp.armazemdc.inf.br"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}

	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	return &Client{
		config:     cfg,
		httpClient: client,
		BaseURL:    cfg.BaseURL,
	}
}

func (c *Client) postJSON(path string, payload interface{}) ([]byte, error) {
	url := c.BaseURL + path

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, buf.String())
	}

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

func (c *Client) postXML(path string, payload string) ([]byte, error) {
	url := c.BaseURL + path

	req, err := http.NewRequest("POST", url, bytes.NewReader([]byte(payload)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/xml")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}

func (c *Client) postRawJSON(path string, payload interface{}) ([]byte, error) {
	url := c.BaseURL + path

	var body []byte
	var err error

	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, buf.String())
	}

	return buf.Bytes(), nil
}
