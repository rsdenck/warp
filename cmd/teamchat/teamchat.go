package teamchat

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	TeamChatCmd = &cobra.Command{
		Use:   "teamchat",
		Short: "TeamChat API operations",
		Long:  `Commands for IceWarp TeamChat API`,
	}
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get TeamChat server version",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		resp, err := client.Version()
		if err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}

		fmt.Printf("Version: %s\n", resp.Version)
		fmt.Printf("Build: %s\n", resp.Build)
		return nil
	},
}

var loginQRCmd = &cobra.Command{
	Use:   "login-qr",
	Short: "Login to TeamChat using QR code",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		// Get QR code login
		loginResp, err := client.GetLogin()
		if err != nil {
			return fmt.Errorf("failed to get login QR: %w", err)
		}

		fmt.Printf("QR Code: %s\n", loginResp.QRCode)
		fmt.Printf("Secret: %s\n", loginResp.Secret)
		fmt.Printf("Login Type: %s\n", loginResp.LoginType)
		fmt.Println("\nScan the QR code with your TeamChat app, then press Enter to check for success...")
		
		// Wait for user input
		fmt.Scanln()

		// Check if login was successful
		successResp, err := client.GetLoginSuccess(loginResp.Secret)
		if err != nil {
			return fmt.Errorf("failed to check login success: %w", err)
		}

		if successResp.Success {
			fmt.Printf("Login successful!\n")
			fmt.Printf("Token: %s\n", successResp.Token)
			fmt.Printf("User: %s\n", successResp.User)
			
			// Save token to config
			viper.Set("teamchat.token", successResp.Token)
			if err := viper.WriteConfig(); err != nil {
				fmt.Printf("Warning: Could not save token to config: %v\n", err)
			}
		} else {
			fmt.Println("Login was not successful. Please try again.")
		}

		return nil
	},
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to TeamChat and save token",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required. Check your warpctl.yaml config")
		}

		fmt.Printf("🔐 Attempting TeamChat login for: %s\n", username)
		fmt.Printf("🌐 Server: %s\n", viper.GetString("server.url"))

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		resp, err := client.LoginPlain(username, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Printf("✅ Login successful!\n")
		fmt.Printf("👤 User: %s\n", resp.User)
		fmt.Printf("🔑 Token: %s\n", resp.Token)
		
		// Save token to config
		viper.Set("teamchat.token", resp.Token)
		if err := viper.WriteConfig(); err != nil {
			fmt.Printf("⚠️  Warning: Could not save token to config: %v\n", err)
			fmt.Printf("💡 You can manually add this token to your warpctl.yaml:\n")
			fmt.Printf("   teamchat:\n     token: \"%s\"\n", resp.Token)
		} else {
			fmt.Printf("💾 Token saved to configuration file\n")
		}
		
		// Test the token immediately
		fmt.Printf("\n🧪 Testing authentication...\n")
		client.SetToken(resp.Token)
		authResp, err := client.AuthTest()
		if err != nil {
			fmt.Printf("⚠️  Token test failed: %v\n", err)
		} else {
			fmt.Printf("✅ Authentication verified!\n")
			fmt.Printf("   User ID: %s\n", authResp.UserID)
			fmt.Printf("   Username: %s\n", authResp.Username)
			fmt.Printf("   Email: %s\n", authResp.Email)
		}

		return nil
	},
}

var authTestCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Check authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("teamchat.token")
		if token == "" {
			fmt.Println("❌ No token found in configuration")
			fmt.Println("💡 Run 'warpctl teamchat login' to authenticate")
			return fmt.Errorf("token not set")
		}

		fmt.Printf("🔍 Testing token: %s...\n", token[:20]+"...")

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		client.SetToken(token)

		resp, err := client.AuthTest()
		if err != nil {
			fmt.Printf("❌ Authentication failed: %v\n", err)
			fmt.Println("💡 Your token may be expired. Run 'warpctl teamchat login' to get a new one")
			return err
		}

		fmt.Printf("✅ Authentication successful!\n")
		fmt.Printf("👤 User ID: %s\n", resp.UserID)
		fmt.Printf("📧 Username: %s\n", resp.Username)
		fmt.Printf("📬 Email: %s\n", resp.Email)
		fmt.Printf("✔️  Valid: %v\n", resp.Valid)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Revoke authentication token",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("teamchat.token")
		if token == "" {
			fmt.Println("Not logged in")
			return nil
		}

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		client.SetToken(token)

		if err := client.Revoke(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}

		fmt.Println("Logged out successfully")
		return nil
	},
}

var presenceCmd = &cobra.Command{
	Use:   "presence [user-ids...]",
	Short: "Get user presence status",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("teamchat.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'icwli teamchat login' first")
		}

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		client.SetToken(token)

		presences, err := client.GetPresence(args)
		if err != nil {
			return fmt.Errorf("failed to get presence: %w", err)
		}

		for _, p := range presences {
			fmt.Printf("User: %s | Online: %v | Status: %s\n", p.UserID, p.Online, p.Status)
		}
		return nil
	},
}

var webLoginCmd = &cobra.Command{
	Use:   "web-login",
	Short: "Login to TeamChat via web automation (like Python script)",
	Long:  `Uses web automation to login to TeamChat, similar to the Python Selenium script approach`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🌐 TeamChat Web Login (Selenium-style automation)")
		fmt.Println("=" + strings.Repeat("=", 50))
		
		// Check if Pinchtab is available
		pinchtabURL := viper.GetString("pinchtab.url")
		if pinchtabURL == "" {
			pinchtabURL = "http://localhost:9867"
		}
		
		fmt.Printf("📡 Checking Pinchtab at: %s\n", pinchtabURL)
		
		// Simple HTTP check for Pinchtab
		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Get(pinchtabURL + "/health")
		if err != nil {
			fmt.Printf("❌ Pinchtab not available: %v\n", err)
			fmt.Println("\n💡 To use web automation:")
			fmt.Println("1. Download Pinchtab: https://github.com/pinchtab/pinchtab/releases")
			fmt.Println("2. Start: ./pinchtab.exe --port 9867")
			fmt.Println("3. Then run: warpctl teamchat web-login")
			return fmt.Errorf("Pinchtab required for web automation")
		}
		resp.Body.Close()
		
		fmt.Println("✅ Pinchtab is running")
		
		// Navigate to TeamChat login
		teamchatURL := viper.GetString("server.url") + "/teamchat"
		fmt.Printf("🌐 Navigating to: %s\n", teamchatURL)
		
		// Use Pinchtab client
		pinchtabClient := sdk.NewPinchtabClient(pinchtabURL, viper.GetString("pinchtab.token"))
		
		if err := pinchtabClient.Navigate(teamchatURL); err != nil {
			return fmt.Errorf("failed to navigate to TeamChat: %w", err)
		}
		
		// Wait for page load
		time.Sleep(3 * time.Second)
		
		// Get page snapshot
		snapshot, err := pinchtabClient.GetSnapshot()
		if err != nil {
			return fmt.Errorf("failed to get page snapshot: %w", err)
		}
		
		fmt.Printf("📸 Page loaded: %s\n", snapshot.Title)
		fmt.Printf("🔍 Found %d interactive elements\n", countInteractiveElements(snapshot.Elements))
		
		// Look for login form elements
		loginElements := findLoginFormElements(snapshot.Elements)
		if len(loginElements) == 0 {
			fmt.Println("⚠️  No login form detected. You may already be logged in.")
			fmt.Println("💡 Try accessing TeamChat directly in your browser to verify")
			return nil
		}
		
		fmt.Printf("🔐 Found %d login form elements\n", len(loginElements))
		
		// Get credentials
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")
		
		if username == "" || password == "" {
			return fmt.Errorf("username and password required in config")
		}
		
		// Perform login automation
		fmt.Printf("👤 Logging in as: %s\n", username)
		
		success := performWebLogin(pinchtabClient, loginElements, username, password)
		if success {
			fmt.Println("✅ Web login completed successfully!")
			fmt.Println("💡 You can now use: warpctl zabbix test-web-notify")
		} else {
			fmt.Println("❌ Web login failed")
			fmt.Println("💡 Check your credentials in warpctl.yaml")
		}
		
		return nil
	},
}

// Helper functions for web login
func countInteractiveElements(elements []sdk.Element) int {
	count := 0
	for _, element := range elements {
		if element.Clickable || element.Focusable {
			count++
		}
	}
	return count
}

func findLoginFormElements(elements []sdk.Element) []sdk.Element {
	var loginElements []sdk.Element
	for _, element := range elements {
		if element.Visible && element.Tag == "input" {
			inputType := element.Attributes["type"]
			if inputType == "text" || inputType == "email" || inputType == "password" {
				loginElements = append(loginElements, element)
			}
		}
		// Also look for login buttons
		if element.Visible && element.Clickable {
			text := strings.ToLower(element.Text)
			if strings.Contains(text, "login") || strings.Contains(text, "entrar") || strings.Contains(text, "sign in") {
				loginElements = append(loginElements, element)
			}
		}
	}
	return loginElements
}

func performWebLogin(client *sdk.PinchtabClient, elements []sdk.Element, username, password string) bool {
	var usernameField, passwordField, loginButton sdk.Element
	
	// Identify form elements
	for _, element := range elements {
		if element.Tag == "input" {
			inputType := element.Attributes["type"]
			placeholder := strings.ToLower(element.Attributes["placeholder"])
			name := strings.ToLower(element.Attributes["name"])
			
			if inputType == "text" || inputType == "email" || 
			   strings.Contains(placeholder, "email") || strings.Contains(placeholder, "user") ||
			   strings.Contains(name, "email") || strings.Contains(name, "user") {
				usernameField = element
			} else if inputType == "password" {
				passwordField = element
			}
		} else if element.Clickable {
			text := strings.ToLower(element.Text)
			if strings.Contains(text, "login") || strings.Contains(text, "entrar") {
				loginButton = element
			}
		}
	}
	
	// Fill form
	if usernameField.Ref != "" {
		fmt.Println("📝 Filling username field...")
		if err := client.Fill(usernameField.Ref, username); err != nil {
			fmt.Printf("❌ Failed to fill username: %v\n", err)
			return false
		}
	}
	
	if passwordField.Ref != "" {
		fmt.Println("🔒 Filling password field...")
		if err := client.Fill(passwordField.Ref, password); err != nil {
			fmt.Printf("❌ Failed to fill password: %v\n", err)
			return false
		}
	}
	
	// Submit form
	if loginButton.Ref != "" {
		fmt.Println("🚀 Clicking login button...")
		if err := client.Click(loginButton.Ref); err != nil {
			fmt.Printf("❌ Failed to click login: %v\n", err)
			return false
		}
	} else {
		// Try pressing Enter on password field
		fmt.Println("⌨️  Pressing Enter to submit...")
		if err := client.PerformAction(sdk.ActionRequest{
			Action: "press",
			Ref:    passwordField.Ref,
			Key:    "Enter",
		}); err != nil {
			fmt.Printf("❌ Failed to submit form: %v\n", err)
			return false
		}
	}
	
	// Wait for login to process
	time.Sleep(3 * time.Second)
	
	return true
}

var testMessageCmd = &cobra.Command{
	Use:   "test-message [channel] [message]",
	Short: "Send a test message to a TeamChat channel",
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("teamchat.token")
		if token == "" {
			fmt.Println("❌ No token found in configuration")
			fmt.Println("💡 Run 'warpctl teamchat login' to authenticate first")
			return fmt.Errorf("token not set")
		}

		// Default values
		channel := "Horus Monitoramento"
		message := "🧪 Test message from WARPCTL TeamChat integration"

		// Override with command line args
		if len(args) >= 1 {
			channel = args[0]
		}
		if len(args) >= 2 {
			message = args[1]
		}

		fmt.Printf("📤 Sending test message to channel: %s\n", channel)
		fmt.Printf("💬 Message: %s\n", message)

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		client.SetToken(token)

		// First test authentication
		fmt.Printf("🔍 Verifying authentication...\n")
		authResp, err := client.AuthTest()
		if err != nil {
			fmt.Printf("❌ Authentication failed: %v\n", err)
			fmt.Println("💡 Run 'warpctl teamchat login' to get a new token")
			return err
		}
		fmt.Printf("✅ Authenticated as: %s\n", authResp.Username)

		// Send the message
		fmt.Printf("📨 Sending message...\n")
		resp, err := client.PostMessage(channel, message)
		if err != nil {
			fmt.Printf("❌ Failed to send message: %v\n", err)
			return err
		}

		if resp.Success {
			fmt.Printf("✅ Message sent successfully!\n")
			if resp.ID != "" {
				fmt.Printf("📝 Message ID: %s\n", resp.ID)
			}
		} else {
			fmt.Printf("❌ Message failed: %s\n", resp.Error)
			return fmt.Errorf("message failed: %s", resp.Error)
		}

		return nil
	},
}

var conversationsCmd = &cobra.Command{
	Use:   "conversations",
	Short: "List conversations",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("teamchat.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'icwli teamchat login' first")
		}

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		client.SetToken(token)

		resp, err := client.ListConversations(50, 0)
		if err != nil {
			return fmt.Errorf("failed to list conversations: %w", err)
		}

		fmt.Printf("Total conversations: %d\n\n", resp.Total)
		for _, c := range resp.Conversations {
			fmt.Printf("ID: %s | Type: %s | Name: %s | Unread: %d\n", c.ID, c.Type, c.Name, c.Unread)
		}
		return nil
	},
}

func init() {
	TeamChatCmd.AddCommand(versionCmd)
	TeamChatCmd.AddCommand(loginCmd)
	TeamChatCmd.AddCommand(webLoginCmd)
	TeamChatCmd.AddCommand(loginQRCmd)
	TeamChatCmd.AddCommand(authTestCmd)
	TeamChatCmd.AddCommand(testMessageCmd)
	TeamChatCmd.AddCommand(logoutCmd)
	TeamChatCmd.AddCommand(presenceCmd)
	TeamChatCmd.AddCommand(conversationsCmd)
}
