package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// WhatsAppClient provides WhatsApp messaging capabilities via HTTP API
type WhatsAppClient struct {
	apiURL     string
	token      string
	client     *http.Client
	phoneNumber string
}

// WhatsAppConfig holds WhatsApp configuration
type WhatsAppConfig struct {
	APIURL      string
	Token       string
	PhoneNumber string
}

// NewWhatsAppClient creates a new WhatsApp client
func NewWhatsAppClient(config *WhatsAppConfig) (*WhatsAppClient, error) {
	if config.APIURL == "" {
		config.APIURL = "http://localhost:3000" // Default WhatsApp Web API
	}

	return &WhatsAppClient{
		apiURL:      config.APIURL,
		token:       config.Token,
		phoneNumber: config.PhoneNumber,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// Connect establishes connection to WhatsApp (placeholder for HTTP API)
func (w *WhatsAppClient) Connect() error {
	// For HTTP API, we just check if the service is available
	resp, err := w.client.Get(w.apiURL + "/status")
	if err != nil {
		return fmt.Errorf("WhatsApp API not available: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("WhatsApp API returned status %d", resp.StatusCode)
	}

	return nil
}

// Disconnect closes the WhatsApp connection (placeholder)
func (w *WhatsAppClient) Disconnect() {
	// Nothing to do for HTTP API
}

// SendMessage sends a text message to a group or individual
func (w *WhatsAppClient) SendMessage(recipient, message string) error {
	payload := map[string]interface{}{
		"chatId": recipient,
		"text":   message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	req, err := http.NewRequest("POST", w.apiURL+"/sendText", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if w.token != "" {
		req.Header.Set("Authorization", "Bearer "+w.token)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetGroups returns list of joined groups (simplified)
func (w *WhatsAppClient) GetGroups() ([]GroupInfo, error) {
	req, err := http.NewRequest("GET", w.apiURL+"/getChats", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if w.token != "" {
		req.Header.Set("Authorization", "Bearer "+w.token)
	}

	resp, err := w.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var chats []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&chats); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var groups []GroupInfo
	for _, chat := range chats {
		if isGroup, ok := chat["isGroup"].(bool); ok && isGroup {
			name, _ := chat["name"].(string)
			id, _ := chat["id"].(string)
			
			groups = append(groups, GroupInfo{
				JID:          id,
				Name:         name,
				Participants: 0, // Not available in simplified API
				CreatedAt:    time.Now(),
			})
		}
	}

	return groups, nil
}

// GroupInfo represents WhatsApp group information
type GroupInfo struct {
	JID          string    `json:"jid"`
	Name         string    `json:"name"`
	Topic        string    `json:"topic"`
	Participants int       `json:"participants"`
	CreatedAt    time.Time `json:"created_at"`
}

// GetConnectionStatus returns current connection status
func (w *WhatsAppClient) GetConnectionStatus() ConnectionStatus {
	err := w.Connect()
	connected := err == nil

	status := "Disconnected"
	if connected {
		status = "Connected"
	}

	return ConnectionStatus{
		Connected:   connected,
		Status:      status,
		PhoneNumber: w.phoneNumber,
		LoggedIn:    connected,
	}
}

// ConnectionStatus represents WhatsApp connection status
type ConnectionStatus struct {
	Connected   bool   `json:"connected"`
	Status      string `json:"status"`
	PhoneNumber string `json:"phone_number,omitempty"`
	LoggedIn    bool   `json:"logged_in"`
}

// SendZabbixAlert sends a formatted Zabbix alert to WhatsApp
func (w *WhatsAppClient) SendZabbixAlert(recipient string, alert ZabbixAlert) error {
	// Create WhatsApp-friendly message (no HTML, emoji-based formatting)
	message := fmt.Sprintf(`🚨 *ALERTA ZABBIX* - %s

🏢 *Host:* %s
⚠️ *Problema:* %s
📊 *Severidade:* %s
⏰ *Duração:* %s
🆔 *Event ID:* %s

📅 *Detectado em:* %s

---
🤖 *Sistema de Monitoramento WARPCTL*
🔗 *Integração Zabbix → WhatsApp*`,
		alert.Severity,
		alert.HostName,
		alert.TriggerName,
		alert.Severity,
		formatDuration(time.Since(alert.Timestamp)),
		generateEventID(alert.TriggerID),
		alert.Timestamp.Format("2006-01-02 15:04:05"))

	return w.SendMessage(recipient, message)
}

// PairWithCode generates a pairing code (simplified for HTTP API)
func (w *WhatsAppClient) PairWithCode() (string, error) {
	// For HTTP API, we simulate pairing
	return "PAIR-CODE-123", fmt.Errorf("pairing not implemented for HTTP API - please configure your WhatsApp Web API service")
}

// Helper functions
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func generateEventID(triggerID string) string {
	if triggerID == "" {
		triggerID = "0000"
	}
	
	if len(triggerID) >= 4 {
		return triggerID[len(triggerID)-4:] + "-WA9"
	}
	
	return triggerID + "-WA9"
}