package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ZabbixClient struct {
	url   string
	token string
}

type ZabbixAlert struct {
	EventID     string    `json:"eventid"`
	TriggerID   string    `json:"triggerid"`
	HostName    string    `json:"hostname"`
	TriggerName string    `json:"trigger_name"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Message     string    `json:"message"`
	GroupName   string    `json:"group_name"`
}

type ZabbixGroup struct {
	GroupID string `json:"groupid"`
	Name    string `json:"name"`
}

type ZabbixRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
	Auth    string      `json:"auth,omitempty"`
}

type ZabbixResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *ZabbixError `json:"error"`
	ID      int         `json:"id"`
}

type ZabbixError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewZabbixClient(url, username, password string) (*ZabbixClient, error) {
	client := &ZabbixClient{url: url}
	
	// Login to get auth token - try different parameter names for different Zabbix versions
	loginParams := map[string]interface{}{
		"username": username,
		"password": password,
	}
	
	resp, err := client.call("user.login", loginParams, "")
	if err != nil {
		// Try with "user" parameter for older versions
		loginParams = map[string]interface{}{
			"user":     username,
			"password": password,
		}
		resp, err = client.call("user.login", loginParams, "")
		if err != nil {
			return nil, fmt.Errorf("failed to login to Zabbix: %w", err)
		}
	}
	
	if token, ok := resp.Result.(string); ok {
		client.token = token
	} else {
		return nil, fmt.Errorf("invalid login response")
	}

	return client, nil
}

func (c *ZabbixClient) call(method string, params interface{}, auth string) (*ZabbixResponse, error) {
	req := ZabbixRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
	
	// Only add auth if we have a token and it's not the login method
	if auth != "" && method != "user.login" {
		req.Auth = auth
	}
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	
	httpReq, err := http.NewRequest("POST", c.url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	
	httpReq.Header.Set("Content-Type", "application/json")
	
	client := &http.Client{Timeout: 30 * time.Second}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	
	var resp ZabbixResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w (body: %s)", err, string(body))
	}
	
	if resp.Error != nil {
		return nil, fmt.Errorf("Zabbix API error: %s (data: %s)", resp.Error.Message, resp.Error.Data)
	}
	
	return &resp, nil
}

func (c *ZabbixClient) GetHostGroups() ([]ZabbixGroup, error) {
	params := map[string]interface{}{
		"output": []string{"groupid", "name"},
	}

	resp, err := c.call("hostgroup.get", params, c.token)
	if err != nil {
		return nil, fmt.Errorf("failed to get host groups: %w", err)
	}

	var groups []ZabbixGroup
	if resultArray, ok := resp.Result.([]interface{}); ok {
		for _, item := range resultArray {
			if group, ok := item.(map[string]interface{}); ok {
				groupName := getString(group, "name")
				// Filter groups that start with VDC
				if strings.HasPrefix(groupName, "VDC") {
					groups = append(groups, ZabbixGroup{
						GroupID: getString(group, "groupid"),
						Name:    groupName,
					})
				}
			}
		}
	}

	return groups, nil
}

func (c *ZabbixClient) GetActiveProblems(groupID string) ([]ZabbixAlert, error) {
	// Get active triggers (problems) directly - more reliable approach
	params := map[string]interface{}{
		"output":         []string{"triggerid", "description", "priority", "value", "lastchange"},
		"selectHosts":    []string{"host"},
		"filter":         map[string]interface{}{"value": 1},
		"sortfield":      "priority",
		"sortorder":      "DESC",
		"expandComment":  true,
		"expandDescription": true,
	}

	if groupID != "" {
		params["groupids"] = []string{groupID}
	}

	resp, err := c.call("trigger.get", params, c.token)
	if err != nil {
		return nil, fmt.Errorf("failed to get triggers: %w", err)
	}

	var alerts []ZabbixAlert
	if resultArray, ok := resp.Result.([]interface{}); ok {
		for _, item := range resultArray {
			if trigger, ok := item.(map[string]interface{}); ok {
				severity := c.getSeverityName(getString(trigger, "priority"))
				
				var hostname string
				if hosts, ok := trigger["hosts"].([]interface{}); ok && len(hosts) > 0 {
					if host, ok := hosts[0].(map[string]interface{}); ok {
						hostname = getString(host, "host")
					}
				}

				// Filter only VDC hosts
				if strings.HasPrefix(hostname, "VDC") {
					// Convert lastchange timestamp
					lastChangeStr := getString(trigger, "lastchange")
					var timestamp time.Time
					if lastChangeStr != "" {
						if lastChangeInt, err := strconv.ParseInt(lastChangeStr, 10, 64); err == nil {
							timestamp = time.Unix(lastChangeInt, 0)
						} else {
							timestamp = time.Now()
						}
					} else {
						timestamp = time.Now()
					}

					alerts = append(alerts, ZabbixAlert{
						EventID:     "", // Will be generated
						TriggerID:   getString(trigger, "triggerid"),
						HostName:    hostname,
						TriggerName: getString(trigger, "description"),
						Severity:    severity,
						Status:      "PROBLEM",
						Timestamp:   timestamp,
						Message:     fmt.Sprintf("%s on %s", getString(trigger, "description"), hostname),
					})
				}
			}
		}
	}

	return alerts, nil
}

func (c *ZabbixClient) getSeverityName(severity string) string {
	switch severity {
	case "0":
		return "Not classified"
	case "1":
		return "Information"
	case "2":
		return "Warning"
	case "3":
		return "Average"
	case "4":
		return "High"
	case "5":
		return "Disaster"
	default:
		return "Unknown"
	}
}

func (c *ZabbixClient) Logout() error {
	if c.token == "" {
		return nil
	}
	
	_, err := c.call("user.logout", []interface{}{}, c.token)
	c.token = ""
	return err
}

func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}