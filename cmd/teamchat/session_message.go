package teamchat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var sessionMessageCmd = &cobra.Command{
	Use:   "session-message [channel] [message]",
	Short: "Send message using XML-RPC session (like Python script)",
	Long:  `Uses the session ID from XML-RPC authentication to send messages, exactly like the Python script`,
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default values (matching Python script)
		channel := "Monitoramento"
		message := "🧪 Test message from WARPCTL (Session-based integration)"

		// Override with command line args
		if len(args) >= 1 {
			channel = args[0]
		}
		if len(args) >= 2 {
			message = args[1]
		}

		fmt.Println("📨 Session-Based TeamChat Message")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("📤 Channel: %s\n", channel)
		fmt.Printf("💬 Message: %s\n", message)

		// Get session ID
		sessionID := viper.GetString("teamchat.session_id")
		if sessionID == "" {
			fmt.Println("❌ No session ID found")
			fmt.Println("💡 Run 'warpctl teamchat xml-auth' first to authenticate")
			return fmt.Errorf("session ID required")
		}

		fmt.Printf("🔑 Using Session ID: %s...\n", sessionID[:20])

		baseURL := viper.GetString("server.url")

		// Step 1: Find the room (like Python script)
		fmt.Printf("🔍 Looking for room: %s\n", channel)
		roomID, err := findRoom(baseURL, sessionID, channel)
		if err != nil {
			return fmt.Errorf("failed to find room: %w", err)
		}

		if roomID == "" {
			fmt.Printf("⚠️  Room '%s' not found. Trying to create it...\n", channel)
			roomID, err = createRoom(baseURL, sessionID, channel)
			if err != nil {
				return fmt.Errorf("failed to create room: %w", err)
			}
		}

		fmt.Printf("✅ Room found/created: %s (ID: %s)\n", channel, roomID)

		// Step 2: Send the message (like Python script)
		fmt.Println("📨 Sending message...")
		success, err := sendMessage(baseURL, sessionID, roomID, message)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		if success {
			fmt.Printf("🎉 Message sent successfully to '%s'!\n", channel)
		} else {
			fmt.Printf("❌ Message failed to send to '%s'\n", channel)
		}

		return nil
	},
}

// Find room using multiple methods (like Python script)
func findRoom(baseURL, sessionID, roomName string) (string, error) {
	endpoints := []string{
		baseURL + "/teamchatapi/",
		baseURL + "/rpc/",
		baseURL + "/icewarpapi/",
	}

	// Multiple methods to try (like Python script)
	methods := []map[string]interface{}{
		{"method": "chat.listRooms", "rpc": true},
		{"command": "chat.listRooms", "rpc": false},
		{"method": "chat.listPublicRooms", "rpc": true},
		{"command": "chat.listPublicRooms", "rpc": false},
		{"method": "teamchat.listRooms", "rpc": true},
		{"command": "teamchat.listRooms", "rpc": false},
		{"command": "listrooms", "rpc": false},
	}

	for _, endpoint := range endpoints {
		for _, method := range methods {
			fmt.Printf("🔄 Trying %s with %s\n", endpoint, getMethodName(method))

			rooms, err := tryListRooms(endpoint, sessionID, method)
			if err != nil {
				fmt.Printf("   ❌ Failed: %v\n", err)
				continue
			}

			// Look for matching room
			for _, room := range rooms {
				roomNameLower := strings.ToLower(room.Name)
				targetLower := strings.ToLower(roomName)

				if roomNameLower == targetLower || 
				   strings.Contains(roomNameLower, targetLower) || 
				   strings.Contains(targetLower, roomNameLower) {
					fmt.Printf("   ✅ Found match: %s (ID: %s)\n", room.Name, room.ID)
					return room.ID, nil
				}
			}

			if len(rooms) > 0 {
				fmt.Printf("   📋 Found %d rooms, but no match for '%s'\n", len(rooms), roomName)
				for i, room := range rooms {
					if i < 5 { // Show first 5 rooms
						fmt.Printf("      - %s (ID: %s)\n", room.Name, room.ID)
					}
				}
				if len(rooms) > 5 {
					fmt.Printf("      ... and %d more\n", len(rooms)-5)
				}
			}
		}
	}

	return "", nil
}

type Room struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func tryListRooms(endpoint, sessionID string, method map[string]interface{}) ([]Room, error) {
	var payload interface{}

	if method["rpc"].(bool) {
		payload = map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  method["method"],
			"params":  map[string]interface{}{"sid": sessionID},
			"id":      1,
		}
	} else {
		payload = map[string]interface{}{
			"command": method["command"],
			"sid":     sessionID,
			"params":  map[string]interface{}{},
		}
	}

	client := &http.Client{Timeout: 15 * time.Second}
	jsonData, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	// Extract rooms from different possible response formats
	var rooms []Room
	
	if resultData, ok := result["result"]; ok {
		if resultList, ok := resultData.([]interface{}); ok {
			for _, item := range resultList {
				if roomMap, ok := item.(map[string]interface{}); ok {
					room := Room{
						ID:   getString(roomMap, "id"),
						Name: getString(roomMap, "name"),
					}
					if room.ID != "" && room.Name != "" {
						rooms = append(rooms, room)
					}
				}
			}
		} else if resultMap, ok := resultData.(map[string]interface{}); ok {
			// Check for nested rooms/folders/items
			if roomsList, ok := resultMap["rooms"].([]interface{}); ok {
				for _, item := range roomsList {
					if roomMap, ok := item.(map[string]interface{}); ok {
						room := Room{
							ID:   getString(roomMap, "id"),
							Name: getString(roomMap, "name"),
						}
						if room.ID != "" && room.Name != "" {
							rooms = append(rooms, room)
						}
					}
				}
			}
		}
	}

	return rooms, nil
}

func createRoom(baseURL, sessionID, roomName string) (string, error) {
	endpoints := []string{
		baseURL + "/teamchatapi/",
		baseURL + "/rpc/",
		baseURL + "/icewarpapi/",
	}

	for _, endpoint := range endpoints {
		fmt.Printf("🔄 Trying to create room at: %s\n", endpoint)

		payload := map[string]interface{}{
			"jsonrpc": "2.0",
			"method":  "chat.createRoom",
			"params": map[string]interface{}{
				"sid":  sessionID,
				"name": roomName,
				"type": "public",
			},
			"id": 1,
		}

		client := &http.Client{Timeout: 15 * time.Second}
		jsonData, _ := json.Marshal(payload)

		req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonData))
		if err != nil {
			continue
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			continue
		}

		body, _ := io.ReadAll(resp.Body)

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			continue
		}

		if resultData, ok := result["result"].(map[string]interface{}); ok {
			if roomID := getString(resultData, "id"); roomID != "" {
				fmt.Printf("✅ Room created successfully: %s\n", roomID)
				return roomID, nil
			}
		}
	}

	return "", fmt.Errorf("failed to create room")
}

func sendMessage(baseURL, sessionID, roomID, message string) (bool, error) {
	endpoints := []string{
		baseURL + "/teamchatapi/",
		baseURL + "/rpc/",
		baseURL + "/icewarpapi/",
	}

	methods := []map[string]interface{}{
		{"method": "chat.sendMessage", "rpc": true},
		{"command": "chat.sendMessage", "rpc": false},
		{"command": "sendmessage", "rpc": false},
	}

	for _, endpoint := range endpoints {
		for _, method := range methods {
			fmt.Printf("🔄 Trying to send via %s with %s\n", endpoint, getMethodName(method))

			var payload interface{}
			sendParams := map[string]interface{}{
				"sid":    sessionID,
				"roomid": roomID,
				"text":   message,
			}

			if method["rpc"].(bool) {
				payload = map[string]interface{}{
					"jsonrpc": "2.0",
					"method":  method["method"],
					"params":  sendParams,
					"id":      1,
				}
			} else {
				payload = map[string]interface{}{
					"command": method["command"],
					"sid":     sessionID,
					"params":  sendParams,
				}
			}

			client := &http.Client{Timeout: 15 * time.Second}
			jsonData, _ := json.Marshal(payload)

			req, err := http.NewRequest("POST", endpoint, bytes.NewReader(jsonData))
			if err != nil {
				continue
			}

			req.Header.Set("Content-Type", "application/json")

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("   ❌ Request failed: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			fmt.Printf("   📥 Response: %s\n", string(body))

			if resp.StatusCode == 200 {
				var result map[string]interface{}
				if err := json.Unmarshal(body, &result); err == nil {
					if result["status"] == "ok" || result["ok"] == true || result["result"] != nil {
						fmt.Printf("   ✅ Message sent successfully!\n")
						return true, nil
					}
				}
			}
		}
	}

	return false, fmt.Errorf("all send attempts failed")
}

// Helper functions
func getMethodName(method map[string]interface{}) string {
	if method["rpc"].(bool) {
		return method["method"].(string)
	}
	return method["command"].(string)
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func init() {
	TeamChatCmd.AddCommand(sessionMessageCmd)
}