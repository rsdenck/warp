package teamchat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// XML-RPC structures (based on Python script)
type AuthRequest struct {
	XMLName xml.Name `xml:"iq"`
	UID     string   `xml:"uid,attr"`
	Format  string   `xml:"format,attr"`
	Query   Query    `xml:"query"`
}

type Query struct {
	XMLName       xml.Name      `xml:"query"`
	Xmlns         string        `xml:"xmlns,attr"`
	CommandName   string        `xml:"commandname"`
	CommandParams CommandParams `xml:"commandparams"`
}

type CommandParams struct {
	XMLName           xml.Name `xml:"commandparams"`
	Email             string   `xml:"email"`
	Password          string   `xml:"password"`
	Digest            string   `xml:"digest"`
	AuthType          string   `xml:"authtype"`
	PersistentLogin   string   `xml:"persistentlogin"`
}

var xmlAuthCmd = &cobra.Command{
	Use:   "xml-auth",
	Short: "Authenticate using XML-RPC (like Python script)",
	Long:  `Uses the same XML-RPC authentication method as the Python script`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 XML-RPC Authentication (Python script method)")
		fmt.Println("=" + strings.Repeat("=", 50))

		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")
		baseURL := viper.GetString("server.url")

		if username == "" || password == "" {
			return fmt.Errorf("username and password required in config")
		}

		fmt.Printf("👤 User: %s\n", username)
		fmt.Printf("🌐 Server: %s\n", baseURL)

		// Try multiple endpoints like Python script
		endpoints := []string{
			baseURL + "/webmail/rpc/",
			baseURL + "/rpc/",
			baseURL + "/icewarpapi/",
			baseURL + "/webmail/icewarpapi/",
		}

		// Create XML payload (exactly like Python script)
		authReq := AuthRequest{
			UID:    "1",
			Format: "text/xml",
			Query: Query{
				Xmlns:       "admin:iq:rpc",
				CommandName: "authenticate",
				CommandParams: CommandParams{
					Email:           username,
					Password:        password,
					Digest:          "",
					AuthType:        "0",
					PersistentLogin: "0",
				},
			},
		}

		xmlData, err := xml.Marshal(authReq)
		if err != nil {
			return fmt.Errorf("failed to create XML payload: %w", err)
		}

		fmt.Printf("📤 XML Payload:\n%s\n\n", string(xmlData))

		// Try each endpoint
		for i, endpoint := range endpoints {
			fmt.Printf("🔄 Trying endpoint %d/%d: %s\n", i+1, len(endpoints), endpoint)

			sid, err := tryXMLAuth(endpoint, xmlData)
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
				continue
			}

			if sid != "" {
				fmt.Printf("✅ Authentication successful!\n")
				fmt.Printf("🔑 Session ID: %s\n", sid)

				// Save the session ID for later use
				viper.Set("teamchat.session_id", sid)
				if err := viper.WriteConfig(); err != nil {
					fmt.Printf("⚠️  Warning: Could not save session ID: %v\n", err)
				} else {
					fmt.Printf("💾 Session ID saved to config\n")
				}

				// Test the session by trying to list rooms
				fmt.Println("\n🧪 Testing session with room list...")
				if err := testSessionWithRoomList(baseURL, sid); err != nil {
					fmt.Printf("⚠️  Session test failed: %v\n", err)
				} else {
					fmt.Println("✅ Session test successful!")
				}

				return nil
			}
		}

		return fmt.Errorf("authentication failed on all endpoints")
	},
}

func tryXMLAuth(endpoint string, xmlData []byte) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(xmlData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/xml")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	responseText := string(body)
	fmt.Printf("📥 Response: %s\n", responseText)

	// Extract SID like Python script does
	// Try sid="..." attribute first
	sidAttr := regexp.MustCompile(`sid="([^"]*)"`)
	if match := sidAttr.FindStringSubmatch(responseText); len(match) > 1 {
		return match[1], nil
	}

	// Try <sid>...</sid> tag
	sidTag := regexp.MustCompile(`<sid>([^<]*)</sid>`)
	if match := sidTag.FindStringSubmatch(responseText); len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("no SID found in response")
}

func testSessionWithRoomList(baseURL, sessionID string) error {
	// Try to list rooms using the session (like Python script)
	endpoints := []string{
		baseURL + "/teamchatapi/",
		baseURL + "/rpc/",
		baseURL + "/icewarpapi/",
	}

	payload := map[string]interface{}{
		"command": "chat.listRooms",
		"sid":     sessionID,
		"params":  map[string]interface{}{},
	}

	for _, endpoint := range endpoints {
		if err := testRoomListEndpoint(endpoint, payload); err == nil {
			return nil
		}
	}

	return fmt.Errorf("failed to list rooms with session")
}

func testRoomListEndpoint(endpoint string, payload interface{}) error {
	client := &http.Client{Timeout: 15 * time.Second}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("🏠 Room list response: %s\n", string(body))
	return nil
}

func init() {
	TeamChatCmd.AddCommand(xmlAuthCmd)
}