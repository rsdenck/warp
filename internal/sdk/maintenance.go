package sdk

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type MaintenanceClient struct {
	client *Client
	sid    string
}

func NewMaintenanceClient(cfg *Config) *MaintenanceClient {
	return &MaintenanceClient{
		client: NewClient(cfg),
	}
}

type XMLRequest struct {
	MethodName string    `xml:"methodName"`
	Params     XMLParams `xml:"params"`
}

type XMLParams struct {
	XMLName xml.Name
	Items   []XMLItem `xml:",any"`
}

type XMLItem struct {
	XMLName xml.Name
	Value   string    `xml:",chardata"`
	Items   []XMLItem `xml:",any"`
}

type XMLResponse struct {
	MethodResponse XMLMethodResponse `xml:"methodResponse"`
}

type XMLMethodResponse struct {
	Params XMLResponseParams `xml:"params"`
}

type XMLResponseParams struct {
	Param XMLParam `xml:"param"`
}

type XMLParam struct {
	Value XMLValue `xml:"value"`
}

type XMLValue struct {
	String  string    `xml:"string"`
	Array   XMLArray  `xml:"array"`
	Struct  XMLStruct `xml:"struct"`
	Int     string    `xml:"int"`
	Boolean string    `xml:"boolean"`
}

type XMLArray struct {
	Data []XMLValue `xml:"data>value"`
}

type XMLStruct struct {
	Member []XMLMember `xml:"member"`
}

type XMLMember struct {
	Name  string   `xml:"name"`
	Value XMLValue `xml:"value"`
}

func (c *MaintenanceClient) Authenticate(username, password string) (string, error) {
	xmlPayload := fmt.Sprintf(`<?xml version="1.0"?>
<methodCall>
	<methodName>authenticate</methodName>
	<params>
		<param><value><string>%s</string></value></param>
		<param><value><string>%s</string></value></param>
	</params>
</methodCall>`, username, password)

	resp, err := c.client.postXML("/icewarpapi/", xmlPayload)
	if err != nil {
		return "", err
	}

	var result XMLResponse
	if err := xml.Unmarshal(resp, &result); err != nil {
		return "", fmt.Errorf("failed to parse XML response: %w", err)
	}

	if result.MethodResponse.Params.Param.Value.String != "" {
		c.sid = result.MethodResponse.Params.Param.Value.String
		return c.sid, nil
	}

	return "", fmt.Errorf("authentication failed")
}

func (c *MaintenanceClient) Call(method string, params map[string]string) (map[string]interface{}, error) {
	if c.sid == "" {
		return nil, fmt.Errorf("not authenticated. Call Authenticate first")
	}

	var paramStrings []string
	paramStrings = append(paramStrings, fmt.Sprintf("<param><value><string>%s</string></value></param>", c.sid))

	for _, value := range params {
		paramStrings = append(paramStrings, fmt.Sprintf("<param><value><string>%s</string></value></param>", value))
	}

	xmlPayload := fmt.Sprintf(`<?xml version="1.0"?>
<methodCall>
	<methodName>%s</methodName>
	<params>
		%s
	</params>
</methodCall>`, method, strings.Join(paramStrings, "\n"))

	resp, err := c.client.postXML("/icewarpapi/", xmlPayload)
	if err != nil {
		return nil, err
	}

	var result XMLResponse
	if err := xml.Unmarshal(resp, &result); err != nil {
		return nil, fmt.Errorf("failed to parse XML response: %w", err)
	}

	output := make(map[string]interface{})

	if result.MethodResponse.Params.Param.Value.Struct.Member != nil {
		for _, member := range result.MethodResponse.Params.Param.Value.Struct.Member {
			if member.Value.String != "" {
				output[member.Name] = member.Value.String
			} else if len(member.Value.Array.Data) > 0 {
				arr := make([]string, len(member.Value.Array.Data))
				for i, v := range member.Value.Array.Data {
					arr[i] = v.String
				}
				output[member.Name] = arr
			}
		}
	}

	return output, nil
}

func (c *MaintenanceClient) GetServerInfo() (map[string]interface{}, error) {
	return c.Call("getinfo", nil)
}

func (c *MaintenanceClient) GetDomainList() ([]map[string]interface{}, error) {
	result, err := c.Call("getdomainlist", nil)
	if err != nil {
		return nil, err
	}

	domains := make([]map[string]interface{}, 0)
	if domainList, ok := result["domain"].([]string); ok {
		for _, domainName := range domainList {
			domains = append(domains, map[string]interface{}{
				"name": domainName,
			})
		}
	}

	return domains, nil
}

func (c *MaintenanceClient) CreateDomain(domainName string) error {
	params := map[string]string{
		"s_domain_name": domainName,
	}
	_, err := c.Call("createdomain", params)
	return err
}

func (c *MaintenanceClient) DeleteDomain(domainName string) error {
	params := map[string]string{
		"s_domain_name": domainName,
	}
	_, err := c.Call("deletedomain", params)
	return err
}

func (c *MaintenanceClient) GetUserList(domain string) ([]map[string]interface{}, error) {
	params := map[string]string{
		"s_domain": domain,
	}
	result, err := c.Call("getuserlist", params)
	if err != nil {
		return nil, err
	}

	users := make([]map[string]interface{}, 0)
	if userList, ok := result["user"].([]string); ok {
		for _, user := range userList {
			users = append(users, map[string]interface{}{
				"name": user,
			})
		}
	}

	return users, nil
}

func (c *MaintenanceClient) CreateUser(domain, username, password string) error {
	params := map[string]string{
		"s_domain":    domain,
		"s_user_name": username,
		"s_password":  password,
	}
	_, err := c.Call("createuser", params)
	return err
}

func (c *MaintenanceClient) DeleteUser(domain, username string) error {
	params := map[string]string{
		"s_domain":    domain,
		"s_user_name": username,
	}
	_, err := c.Call("deleteuser", params)
	return err
}

func (c *MaintenanceClient) GetUserInfo(domain, username string) (map[string]interface{}, error) {
	params := map[string]string{
		"s_domain":    domain,
		"s_user_name": username,
	}
	return c.Call("getuserinfo", params)
}

func (c *MaintenanceClient) SetUserQuota(domain, username string, quotaMB int) error {
	params := map[string]string{
		"s_domain":     domain,
		"s_user_name":  username,
		"s_quota_size": fmt.Sprintf("%d", quotaMB),
	}
	_, err := c.Call("setuserquota", params)
	return err
}

func (c *MaintenanceClient) GetDomainInfo(domain string) (map[string]interface{}, error) {
	params := map[string]string{
		"s_domain_name": domain,
	}
	return c.Call("getdomaininfo", params)
}

func (c *MaintenanceClient) GetSystemInfo() (map[string]interface{}, error) {
	return c.Call("getsysteminfo", nil)
}

func (c *MaintenanceClient) GetStatistics() (map[string]interface{}, error) {
	return c.Call("getstatistics", nil)
}

func (c *MaintenanceClient) Logout() error {
	if c.sid == "" {
		return nil
	}

	_, err := c.Call("logout", nil)
	c.sid = ""
	return err
}

func (c *MaintenanceClient) GetSID() string {
	return c.sid
}

func (c *MaintenanceClient) IsAuthenticated() bool {
	return c.sid != ""
}
