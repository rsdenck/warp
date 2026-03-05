package sdk

import (
	"encoding/json"
	"fmt"
)

type TeamChatClient struct {
	client *Client
	token  string
}

func NewTeamChatClient(cfg *Config) *TeamChatClient {
	return &TeamChatClient{
		client: NewClient(cfg),
	}
}

type VersionResponse struct {
	Version string `json:"version"`
	Build   string `json:"build"`
}

func (c *TeamChatClient) Version() (*VersionResponse, error) {
	resp, err := c.client.postJSON("/teamchatapi/version", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result VersionResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type AuthTestResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Valid    bool   `json:"valid"`
}

func (c *TeamChatClient) AuthTest() (*AuthTestResponse, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated. Call Login first")
	}

	payload := map[string]interface{}{
		"token": c.token,
	}

	resp, err := c.client.postJSON("/teamchatapi/auth.test", payload)
	if err != nil {
		return nil, err
	}

	var result AuthTestResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type LoginResponse struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

func (c *TeamChatClient) LoginPlain(username, password string) (*LoginResponse, error) {
	// Try the correct IceWarp authentication endpoint based on the Python script
	payload := map[string]interface{}{
		"email": username,
	}

	// First, get the login challenge
	resp, err := c.client.postJSON("/teamchatapi/iwauthentication.getLogin", payload)
	if err != nil {
		return nil, fmt.Errorf("failed to get login challenge: %w", err)
	}

	var loginChallenge struct {
		QRCode    string `json:"qrcode"`
		Secret    string `json:"secret"`
		LoginType string `json:"login_type"`
	}
	
	if err := json.Unmarshal(resp, &loginChallenge); err != nil {
		return nil, fmt.Errorf("failed to parse login challenge: %w", err)
	}

	// For now, return an error suggesting web automation
	return nil, fmt.Errorf("TeamChat API requires QR code authentication. Use web automation instead: 'warpctl zabbix test-web-notify'")
}

type LoginExternalResponse struct {
	SessionID string `json:"session_id"`
}

func (c *TeamChatClient) LoginExternal(username, provider string) (*LoginExternalResponse, error) {
	payload := map[string]interface{}{
		"user":     username,
		"provider": provider,
	}

	resp, err := c.client.postJSON("/teamchatapi/iwauthentication.login.external", payload)
	if err != nil {
		return nil, err
	}

	var result LoginExternalResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

type GetLoginResponse struct {
	QRCode    string `json:"qrcode"`
	Secret    string `json:"secret"`
	LoginType string `json:"login_type"`
}

func (c *TeamChatClient) GetLogin() (*GetLoginResponse, error) {
	payload := map[string]interface{}{
		"email": "ranlens.denck@armazem.cloud",
	}
	
	resp, err := c.client.postJSON("/teamchatapi/iwauthentication.getLogin", payload)
	if err != nil {
		return nil, err
	}

	var result GetLoginResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type GetLoginSuccessResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token,omitempty"`
	User    string `json:"user,omitempty"`
}

func (c *TeamChatClient) GetLoginSuccess(secret string) (*GetLoginSuccessResponse, error) {
	payload := map[string]interface{}{
		"secret": secret,
	}

	resp, err := c.client.postJSON("/teamchatapi/iwauthentication.getLoginSuccess", payload)
	if err != nil {
		return nil, err
	}

	var result GetLoginSuccessResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Success {
		c.token = result.Token
	}

	return &result, nil
}

func (c *TeamChatClient) Revoke() error {
	if c.token == "" {
		return nil
	}

	payload := map[string]interface{}{
		"token": c.token,
	}

	_, err := c.client.postJSON("/teamchatapi/auth.revoke", payload)
	if err != nil {
		return err
	}

	c.token = ""
	return nil
}

type UserPresence struct {
	UserID    string `json:"user_id"`
	Online    bool   `json:"online"`
	Status    string `json:"status"`
	Avatar    string `json:"avatar,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

func (c *TeamChatClient) GetPresence(userIDs []string) ([]UserPresence, error) {
	payload := map[string]interface{}{
		"token":   c.token,
		"user_id": userIDs,
	}

	resp, err := c.client.postJSON("/teamchatapi/users.getPresence", payload)
	if err != nil {
		return nil, err
	}

	var result []UserPresence
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

type Conversation struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Name        string                 `json:"name,omitempty"`
	Members     []string               `json:"members,omitempty"`
	LastMessage map[string]interface{} `json:"last_message,omitempty"`
	Unread      int                    `json:"unread"`
}

type ConversationsListResponse struct {
	Conversations []Conversation `json:"conversations"`
	Total         int            `json:"total"`
}

func (c *TeamChatClient) ListConversations(limit, offset int) (*ConversationsListResponse, error) {
	payload := map[string]interface{}{
		"token":  c.token,
		"limit":  limit,
		"offset": offset,
	}

	resp, err := c.client.postJSON("/teamchatapi/conversations.list", payload)
	if err != nil {
		return nil, err
	}

	var result ConversationsListResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type ConversationInfoResponse struct {
	Conversation Conversation `json:"conversation"`
}

func (c *TeamChatClient) GetConversationInfo(conversationID string) (*ConversationInfoResponse, error) {
	payload := map[string]interface{}{
		"token":        c.token,
		"conversation": conversationID,
	}

	resp, err := c.client.postJSON("/teamchatapi/conversations.info", payload)
	if err != nil {
		return nil, err
	}

	var result ConversationInfoResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

type Attachment struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	URL      string `json:"url,omitempty"`
}

func (c *TeamChatClient) GetAttachments(itemID string) ([]Attachment, error) {
	payload := map[string]interface{}{
		"token": c.token,
		"item":  itemID,
	}

	resp, err := c.client.postJSON("/teamchatapi/items.getAttachments", payload)
	if err != nil {
		return nil, err
	}

	var result []Attachment
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result, nil
}

func (c *TeamChatClient) GetToken() string {
	return c.token
}

func (c *TeamChatClient) SetToken(token string) {
	c.token = token
}

func (c *TeamChatClient) IsAuthenticated() bool {
	return c.token != ""
}

type PostMessageResponse struct {
	Success bool   `json:"success"`
	ID      string `json:"id,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (c *TeamChatClient) PostMessage(channel, text string) (*PostMessageResponse, error) {
	if c.token == "" {
		return nil, fmt.Errorf("not authenticated. Call Login first")
	}

	payload := map[string]interface{}{
		"token":   c.token,
		"channel": channel,
		"text":    text,
	}

	resp, err := c.client.postJSON("/teamchatapi/chat.postMessage", payload)
	if err != nil {
		return nil, err
	}

	var result PostMessageResponse
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !result.Success && result.Error != "" {
		return nil, fmt.Errorf("failed to post message: %s", result.Error)
	}

	return &result, nil
}
