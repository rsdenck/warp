package sdk

import (
	"fmt"
	"time"
)

type ServicesClient struct {
	maintClient *MaintenanceClient
}

func NewServicesClient(cfg *Config) *ServicesClient {
	return &ServicesClient{
		maintClient: NewMaintenanceClient(cfg),
	}
}

func (c *ServicesClient) Authenticate(username, password string) (string, error) {
	return c.maintClient.Authenticate(username, password)
}

func (c *ServicesClient) SetSID(sid string) {
	c.maintClient.sid = sid
}

func (c *ServicesClient) IsAuthenticated() bool {
	return c.maintClient.IsAuthenticated()
}

type Service struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	PID         int    `json:"pid"`
	Uptime      int    `json:"uptime"`
	Type        string `json:"type"`
}

func (c *ServicesClient) ListServices() ([]Service, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	result, err := c.maintClient.Call("getserviceslist", nil)
	if err != nil {
		return nil, err
	}

	servicesData, ok := result["service"]
	if !ok {
		return []Service{}, nil
	}

	switch v := servicesData.(type) {
	case []interface{}:
		services := make([]Service, 0, len(v))
		for _, s := range v {
			if svcMap, ok := s.(map[string]interface{}); ok {
				service := Service{
					Name:        getServiceString(svcMap, "name"),
					DisplayName: getServiceString(svcMap, "display_name"),
					Description: getServiceString(svcMap, "description"),
					Status:      getServiceString(svcMap, "status"),
					PID:         getServiceInt(svcMap, "pid"),
					Uptime:      getServiceInt(svcMap, "uptime"),
					Type:        getServiceString(svcMap, "type"),
				}
				services = append(services, service)
			}
		}
		return services, nil
	}

	return []Service{}, nil
}

func (c *ServicesClient) GetServiceStatus(serviceName string) (*Service, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_service": serviceName,
	}

	result, err := c.maintClient.Call("getservicestatus", params)
	if err != nil {
		return nil, err
	}

	service := &Service{
		Name:        getServiceString(result, "name"),
		DisplayName: getServiceString(result, "display_name"),
		Description: getServiceString(result, "description"),
		Status:      getServiceString(result, "status"),
		PID:         getServiceInt(result, "pid"),
		Uptime:      getServiceInt(result, "uptime"),
		Type:        getServiceString(result, "type"),
	}

	return service, nil
}

func (c *ServicesClient) GetServiceStats(serviceName string) (map[string]interface{}, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_service": serviceName,
	}

	return c.maintClient.Call("getservicestats", params)
}

func (c *ServicesClient) StartService(serviceName string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_service": serviceName,
	}

	_, err := c.maintClient.Call("startservcie", params)
	return err
}

func (c *ServicesClient) StopService(serviceName string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_service": serviceName,
	}

	_, err := c.maintClient.Call("stopservice", params)
	return err
}

func (c *ServicesClient) RestartService(serviceName string) error {
	if !c.maintClient.IsAuthenticated() {
		return fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_service": serviceName,
	}

	_, err := c.maintClient.Call("restartservice", params)
	return err
}

func (c *ServicesClient) GetTrafficChart(serviceName string, startTime, endTime time.Time) ([]map[string]interface{}, error) {
	if !c.maintClient.IsAuthenticated() {
		return nil, fmt.Errorf("not authenticated")
	}

	params := map[string]string{
		"s_service":    serviceName,
		"s_start_time": fmt.Sprintf("%d", startTime.Unix()),
		"s_end_time":   fmt.Sprintf("%d", endTime.Unix()),
	}

	result, err := c.maintClient.Call("gettrafficchart", params)
	if err != nil {
		return nil, err
	}

	data, ok := result["data"]
	if !ok {
		return []map[string]interface{}{}, nil
	}

	switch v := data.(type) {
	case []interface{}:
		chartData := make([]map[string]interface{}, 0, len(v))
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				chartData = append(chartData, itemMap)
			}
		}
		return chartData, nil
	}

	return []map[string]interface{}{}, nil
}

func getServiceString(m map[string]interface{}, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func getServiceInt(m map[string]interface{}, key string) int {
	if v, ok := m[key].(float64); ok {
		return int(v)
	}
	return 0
}
