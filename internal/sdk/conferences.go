package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type ConferencesClient struct {
	client *Client
	token  string
}

func NewConferencesClient(cfg *Config) *ConferencesClient {
	return &ConferencesClient{
		client: NewClient(cfg),
	}
}

func (c *ConferencesClient) SetToken(token string) {
	c.token = token
}

type Conference struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"start_time"`
	EndTime      time.Time `json:"end_time"`
	Duration     int       `json:"duration"`
	Organizer    string    `json:"organizer"`
	Participants []string  `json:"participants"`
	RoomID       string    `json:"room_id"`
	Password     string    `json:"password"`
	Recording    bool      `json:"recording"`
	RecordingURL string    `json:"recording_url,omitempty"`
	Status       string    `json:"status"`
	JoinURL      string    `json:"join_url"`
	HostURL      string    `json:"host_url"`
	Created      time.Time `json:"created"`
}

type ConferenceRoom struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Capacity    int    `json:"capacity"`
	Owner       string `json:"owner"`
}

func (c *ConferencesClient) ListConferences(start, end time.Time) ([]Conference, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "conferences.list",
		"token":  c.token,
		"start":  start.Format(time.RFC3339),
		"end":    end.Format(time.RFC3339),
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Conferences []Conference `json:"conferences"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Conferences, nil
}

func (c *ConferencesClient) ListRooms() ([]ConferenceRoom, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "conferences.rooms.list",
		"token":  c.token,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Rooms []ConferenceRoom `json:"rooms"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Rooms, nil
}

func (c *ConferencesClient) CreateConference(title, description string, startTime time.Time, duration int, participants []string, roomID, password string, recording bool) (*Conference, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":       "conferences.create",
		"token":        c.token,
		"title":        title,
		"description":  description,
		"start_time":   startTime.Format(time.RFC3339),
		"duration":     duration,
		"participants": participants,
		"room_id":      roomID,
		"password":     password,
		"recording":    recording,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Conference
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ConferencesClient) UpdateConference(conferenceID string, title, description string, startTime time.Time, duration int, participants []string, password string, recording bool) (*Conference, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":        "conferences.update",
		"token":         c.token,
		"conference_id": conferenceID,
		"title":         title,
		"description":   description,
		"start_time":    startTime.Format(time.RFC3339),
		"duration":      duration,
		"participants":  participants,
		"password":      password,
		"recording":     recording,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Conference
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ConferencesClient) DeleteConference(conferenceID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":        "conferences.delete",
		"token":         c.token,
		"conference_id": conferenceID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *ConferencesClient) AddParticipants(conferenceID string, participants []string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":        "conferences.participants.add",
		"token":         c.token,
		"conference_id": conferenceID,
		"participants":  participants,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *ConferencesClient) RemoveParticipants(conferenceID string, participants []string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":        "conferences.participants.remove",
		"token":         c.token,
		"conference_id": conferenceID,
		"participants":  participants,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *ConferencesClient) GetConferenceInfo(conferenceID string) (*Conference, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":        "conferences.info",
		"token":         c.token,
		"conference_id": conferenceID,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Conference
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ConferencesClient) CreateRoom(name, description string, capacity int) (*ConferenceRoom, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "conferences.rooms.create",
		"token":       c.token,
		"name":        name,
		"description": description,
		"capacity":    capacity,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result ConferenceRoom
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *ConferencesClient) DeleteRoom(roomID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":  "conferences.rooms.delete",
		"token":   c.token,
		"room_id": roomID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}
