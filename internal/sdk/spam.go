package sdk

import (
	"fmt"
)

type SpamClient struct {
	maintClient *MaintenanceClient
}

func NewSpamClient(cfg *Config) *SpamClient {
	return &SpamClient{
		maintClient: NewMaintenanceClient(cfg),
	}
}

func (c *SpamClient) Authenticate(username, password string) (string, error) {
	return c.maintClient.Authenticate(username, password)
}

func (c *SpamClient) SetSID(sid string) {
	c.maintClient.sid = sid
}

func (c *SpamClient) IsAuthenticated() bool {
	return c.maintClient.IsAuthenticated()
}

type SpamItem struct {
	ID      string `json:"id"`
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Date    string `json:"date"`
	Size    int    `json:"size"`
	Score   string `json:"score"`
	Reason  string `json:"reason"`
	Action  string `json:"action"`
}

func (c *SpamClient) ListQuarantine(limit, offset int) ([]SpamItem, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_limit":  fmt.Sprintf("%d", limit),
		"s_offset": fmt.Sprintf("%d", offset),
	}

	result, err := c.maintClient.Call("getspamlist", params)
	if err != nil {
		return nil, err
	}

	itemsData, ok := result["item"]
	if !ok {
		return []SpamItem{}, nil
	}

	switch v := itemsData.(type) {
	case []interface{}:
		items := make([]SpamItem, 0, len(v))
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				spamItem := SpamItem{
					ID:      getSpamString(itemMap, "id"),
					From:    getSpamString(itemMap, "from"),
					To:      getSpamString(itemMap, "to"),
					Subject: getSpamString(itemMap, "subject"),
					Date:    getSpamString(itemMap, "date"),
					Size:    getSpamInt(itemMap, "size"),
					Score:   getSpamString(itemMap, "score"),
					Reason:  getSpamString(itemMap, "reason"),
					Action:  getSpamString(itemMap, "action"),
				}
				items = append(items, spamItem)
			}
		}
		return items, nil
	}

	return []SpamItem{}, nil
}

func (c *SpamClient) GetQuarantineItem(itemID string) (*SpamItem, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_item_id": itemID,
	}

	result, err := c.maintClient.Call("getspamitem", params)
	if err != nil {
		return nil, err
	}

	item := &SpamItem{
		ID:      getSpamString(result, "id"),
		From:    getSpamString(result, "from"),
		To:      getSpamString(result, "to"),
		Subject: getSpamString(result, "subject"),
		Date:    getSpamString(result, "date"),
		Size:    getSpamInt(result, "size"),
		Score:   getSpamString(result, "score"),
		Reason:  getSpamString(result, "reason"),
		Action:  getSpamString(result, "action"),
	}

	return item, nil
}

func (c *SpamClient) GetSpamBody(itemID string) (string, error) {
	if !c.maintClient.IsAuthenticated() {
		return "", fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_item_id": itemID,
	}

	result, err := c.maintClient.Call("getspambody", params)
	if err != nil {
		return "", err
	}

	return getSpamString(result, "body"), nil
}

func (c *SpamClient) DeliverQuarantineItem(itemID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_item_id": itemID,
	}

	_, err := c.maintClient.Call("deliverspamitem", params)
	return err
}

func (c *SpamClient) DeleteQuarantineItem(itemID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_item_id": itemID,
	}

	_, err := c.maintClient.Call("deletespamitem", params)
	return err
}

func (c *SpamClient) DeleteAllQuarantine() error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	_, err := c.maintClient.Call("deleteallspam", nil)
	return err
}

func (c *SpamClient) WhitelistSender(email string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_sender": email,
	}

	_, err := c.maintClient.Call("whitelist", params)
	return err
}

func (c *SpamClient) BlacklistSender(email string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_sender": email,
	}

	_, err := c.maintClient.Call("blacklist", params)
	return err
}

func (c *SpamClient) DeleteFromBlacklist(email string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_sender": email,
	}

	_, err := c.maintClient.Call("deleteblacklist", params)
	return err
}

func (c *SpamClient) DeleteFromWhitelist(email string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_sender": email,
	}

	_, err := c.maintClient.Call("deletewhitelist", params)
	return err
}

func getSpamString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getSpamInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}
