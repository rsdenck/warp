package sdk

import (
	"fmt"
	"time"
)

type DevicesClient struct {
	maintClient *MaintenanceClient
}

func NewDevicesClient(cfg *Config) *DevicesClient {
	return &DevicesClient{
		maintClient: NewMaintenanceClient(cfg),
	}
}

func (c *DevicesClient) Authenticate(username, password string) (string, error) {
	return c.maintClient.Authenticate(username, password)
}

func (c *DevicesClient) SetSID(sid string) {
	c.maintClient.sid = sid
}

func (c *DevicesClient) IsAuthenticated() bool {
	return c.maintClient.IsAuthenticated()
}

type MobileDevice struct {
	ID         string    `json:"id"`
	Account    string    `json:"account"`
	Type       string    `json:"type"`
	Model      string    `json:"model"`
	OS         string    `json:"os"`
	OSVersion  string    `json:"os_version"`
	LastSync   time.Time `json:"last_sync"`
	Policy     string    `json:"policy"`
	Status     string    `json:"status"`
	RemoteWipe bool      `json:"remote_wipe"`
	Approved   bool      `json:"approved"`
}

func (c *DevicesClient) ListDevices(account, filter string) ([]MobileDevice, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account": account,
		"s_filter":  filter,
	}

	result, err := c.maintClient.Call("getdevicelist", params)
	if err != nil {
		return nil, err
	}

	devicesData, ok := result["device"]
	if !ok {
		return []MobileDevice{}, nil
	}

	switch v := devicesData.(type) {
	case []interface{}:
		devices := make([]MobileDevice, 0, len(v))
		for _, d := range v {
			if devMap, ok := d.(map[string]interface{}); ok {
				device := MobileDevice{
					ID:         getDeviceString(devMap, "id"),
					Account:    getDeviceString(devMap, "account"),
					Type:       getDeviceString(devMap, "type"),
					Model:      getDeviceString(devMap, "model"),
					OS:         getDeviceString(devMap, "os"),
					OSVersion:  getDeviceString(devMap, "os_version"),
					Policy:     getDeviceString(devMap, "policy"),
					Status:     getDeviceString(devMap, "status"),
					RemoteWipe: getDeviceString(devMap, "remote_wipe") == "1",
					Approved:   getDeviceString(devMap, "approved") == "1",
				}
				devices = append(devices, device)
			}
		}
		return devices, nil
	}

	return []MobileDevice{}, nil
}

func (c *DevicesClient) GetDeviceInfo(deviceID string) (*MobileDevice, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_device_id": deviceID,
	}

	result, err := c.maintClient.Call("getdeviceinfo", params)
	if err != nil {
		return nil, err
	}

	device := &MobileDevice{
		ID:         getDeviceString(result, "id"),
		Account:    getDeviceString(result, "account"),
		Type:       getDeviceString(result, "type"),
		Model:      getDeviceString(result, "model"),
		OS:         getDeviceString(result, "os"),
		OSVersion:  getDeviceString(result, "os_version"),
		Policy:     getDeviceString(result, "policy"),
		Status:     getDeviceString(result, "status"),
		RemoteWipe: getDeviceString(result, "remote_wipe") == "1",
		Approved:   getDeviceString(result, "approved") == "1",
	}

	return device, nil
}

func (c *DevicesClient) SetDeviceProperties(deviceID, policy string, approved bool) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_device_id": deviceID,
		"s_policy":    policy,
		"s_approved":  boolToDeviceString(approved),
	}

	_, err := c.maintClient.Call("setdeviceproperties", params)
	return err
}

func (c *DevicesClient) DeleteDevice(deviceID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_device_id": deviceID,
	}

	_, err := c.maintClient.Call("deletedevice", params)
	return err
}

func (c *DevicesClient) RemoteWipe(deviceID string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_device_id": deviceID,
	}

	_, err := c.maintClient.Call("remotewipe", params)
	return err
}

func (c *DevicesClient) SetDeviceStatus(deviceID, status string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_device_id": deviceID,
		"s_status":    status,
	}

	_, err := c.maintClient.Call("setdevicestatus", params)
	return err
}

func (c *DevicesClient) DeleteDevicesByFilter(account, filter string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_account": account,
		"s_filter":  filter,
	}

	_, err := c.maintClient.Call("deletedevicesbyfilter", params)
	return err
}

func getDeviceString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func boolToDeviceString(b bool) string {
	if b {
		return "1"
	}
	return "0"
}
