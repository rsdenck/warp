package sdk

import (
	"fmt"
)

type FiltersClient struct {
	maintClient *MaintenanceClient
}

func NewFiltersClient(cfg *Config) *FiltersClient {
	return &FiltersClient{
		maintClient: NewMaintenanceClient(cfg),
	}
}

func (c *FiltersClient) Authenticate(username, password string) (string, error) {
	return c.maintClient.Authenticate(username, password)
}

func (c *FiltersClient) SetSID(sid string) {
	c.maintClient.sid = sid
}

func (c *FiltersClient) IsAuthenticated() bool {
	return c.maintClient.IsAuthenticated()
}

type Filter struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FolderID    string `json:"folder_id"`
	Account     string `json:"account"`
	RuleType    string `json:"rule_type"`
	Enabled     bool   `json:"enabled"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
}

func (c *FiltersClient) ListRules(account, ruleType string) ([]Filter, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account": account,
		"s_type":    ruleType,
	}

	result, err := c.maintClient.Call("getrules", params)
	if err != nil {
		return nil, err
	}

	rulesData, ok := result["rule"]
	if !ok {
		return []Filter{}, nil
	}

	switch v := rulesData.(type) {
	case []interface{}:
		rules := make([]Filter, 0, len(v))
		for _, r := range v {
			if ruleMap, ok := r.(map[string]interface{}); ok {
				filter := Filter{
					ID:       getStringValue(ruleMap, "id"),
					Name:     getStringValue(ruleMap, "name"),
					FolderID: getStringValue(ruleMap, "folder_id"),
					Enabled:  getStringValue(ruleMap, "enabled") == "1",
					Priority: getIntValue(ruleMap, "priority"),
				}
				rules = append(rules, filter)
			}
		}
		return rules, nil
	case string:
		if v == "" {
			return []Filter{}, nil
		}
		return []Filter{{ID: v}}, nil
	}

	return []Filter{}, nil
}

func (c *FiltersClient) GetRule(account, ruleID string) (*Filter, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account": account,
		"s_rule_id": ruleID,
	}

	result, err := c.maintClient.Call("getruleinfo", params)
	if err != nil {
		return nil, err
	}

	filter := &Filter{
		ID:       getStringValue(result, "id"),
		Name:     getStringValue(result, "name"),
		FolderID: getStringValue(result, "folder_id"),
		Enabled:  getStringValue(result, "enabled") == "1",
		Priority: getIntValue(result, "priority"),
	}

	return filter, nil
}

func (c *FiltersClient) CreateRule(account, name, folderID, conditionType, conditionField, conditionOperator, conditionValue, actionType, actionFolderID, actionText string, enabled bool, priority int) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	conditionJSON := fmt.Sprintf(`{"type":"%s","field":"%s","operator":"%s","value":"%s"}`,
		conditionType, conditionField, conditionOperator, conditionValue)

	actionJSON := fmt.Sprintf(`{"type":"%s","folder_id":"%s","text":"%s"}`,
		actionType, actionFolderID, actionText)

	params := map[string]string{
		"s_account":   account,
		"s_name":      name,
		"s_folder_id": folderID,
		"s_condition": conditionJSON,
		"s_action":    actionJSON,
		"s_enabled":   boolToStringValue(enabled),
		"s_priority":  fmt.Sprintf("%d", priority),
	}

	_, err := c.maintClient.Call("createrule", params)
	return err
}

func (c *FiltersClient) UpdateRule(ruleID, account, name, conditionJSON, actionJSON string, enabled bool, priority int) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_rule_id":   ruleID,
		"s_account":   account,
		"s_name":      name,
		"s_condition": conditionJSON,
		"s_action":    actionJSON,
		"s_enabled":   boolToStringValue(enabled),
		"s_priority":  fmt.Sprintf("%d", priority),
	}

	_, err := c.maintClient.Call("editrule", params)
	return err
}

func (c *FiltersClient) DeleteRule(account, ruleID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account": account,
		"s_rule_id": ruleID,
	}

	_, err := c.maintClient.Call("deleterule", params)
	return err
}

func (c *FiltersClient) SetRuleState(account, ruleID string, enabled bool) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account": account,
		"s_rule_id": ruleID,
		"s_enabled": boolToStringValue(enabled),
	}

	_, err := c.maintClient.Call("setrulestate", params)
	return err
}

func (c *FiltersClient) MoveRule(account, ruleID, direction string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account":   account,
		"s_rule_id":   ruleID,
		"s_direction": direction,
	}

	_, err := c.maintClient.Call("moverule", params)
	return err
}

func getStringValue(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getIntValue(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}

func boolToStringValue(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
