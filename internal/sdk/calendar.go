package sdk

import (
	"encoding/json"
	"fmt"
	"time"
)

type CalendarClient struct {
	client *Client
	token  string
}

func NewCalendarClient(cfg *Config) *CalendarClient {
	return &CalendarClient{
		client: NewClient(cfg),
	}
}

func (c *CalendarClient) SetToken(token string) {
	c.token = token
}

type Calendar struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	Owner       string    `json:"owner"`
	IsShared    bool      `json:"is_shared"`
	IsPublic    bool      `json:"is_public"`
	Created     time.Time `json:"created"`
	Modified    time.Time `json:"modified"`
}

type CalendarEvent struct {
	ID          string    `json:"id"`
	CalendarID  string    `json:"calendar_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	AllDay      bool      `json:"all_day"`
	Location    string    `json:"location"`
	Attendees   []string  `json:"attendees"`
	Reminder    int       `json:"reminder"`
	Recurrence  string    `json:"recurrence"`
	Organizer   string    `json:"organizer"`
	Status      string    `json:"status"`
}

func (c *CalendarClient) ListCalendars() ([]Calendar, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method": "calendars.list",
		"token":  c.token,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Calendars []Calendar `json:"calendars"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Calendars, nil
}

func (c *CalendarClient) CreateCalendar(name, description, color string) (*Calendar, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "calendars.create",
		"token":       c.token,
		"name":        name,
		"description": description,
		"color":       color,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result Calendar
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *CalendarClient) DeleteCalendar(calendarID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "calendars.delete",
		"token":       c.token,
		"calendar_id": calendarID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}

func (c *CalendarClient) ListEvents(calendarID string, start, end time.Time) ([]CalendarEvent, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "calendars.events.list",
		"token":       c.token,
		"calendar_id": calendarID,
		"start":       start.Format(time.RFC3339),
		"end":         end.Format(time.RFC3339),
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Events []CalendarEvent `json:"events"`
	}
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return result.Events, nil
}

func (c *CalendarClient) CreateEvent(calendarID, title, description, location string, start, end time.Time, allDay bool, attendees []string) (*CalendarEvent, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":      "calendars.events.create",
		"token":       c.token,
		"calendar_id": calendarID,
		"title":       title,
		"description": description,
		"location":    location,
		"start":       start.Format(time.RFC3339),
		"end":         end.Format(time.RFC3339),
		"all_day":     allDay,
		"attendees":   attendees,
	}

	resp, err := c.client.postJSON("/cloudapi/", payload)
	if err != nil {
		return nil, err
	}

	var result CalendarEvent
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *CalendarClient) UpdateEvent(eventID string, updates map[string]interface{}) (*CalendarEvent, error) {
	if c.token == "" {
		return nil, fmt.Errorf("token required")
	}

	updates["method"] = "calendars.events.update"
	updates["token"] = c.token
	updates["event_id"] = eventID

	resp, err := c.client.postJSON("/cloudapi/", updates)
	if err != nil {
		return nil, err
	}

	var result CalendarEvent
	if err := json.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

func (c *CalendarClient) DeleteEvent(eventID string) error {
	if c.token == "" {
		return fmt.Errorf("token required")
	}

	payload := map[string]interface{}{
		"method":   "calendars.events.delete",
		"token":    c.token,
		"event_id": eventID,
	}

	_, err := c.client.postJSON("/cloudapi/", payload)
	return err
}
