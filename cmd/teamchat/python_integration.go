package teamchat

import (
	"fmt"
	"strings"
	"time"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pythonStyleCmd = &cobra.Command{
	Use:   "python-style [channel] [message]",
	Short: "Send message using Python script approach (web automation)",
	Long:  `Replicates the Python Selenium script approach using Pinchtab for web automation`,
	Args:  cobra.RangeArgs(0, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default values (matching Python script)
		channel := "Monitoramento"
		message := "🧪 Test message from WARPCTL (Python-style integration)"

		// Override with command line args
		if len(args) >= 1 {
			channel = args[0]
		}
		if len(args) >= 2 {
			message = args[1]
		}

		fmt.Println("🐍 Python-Style TeamChat Integration")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("📤 Channel: %s\n", channel)
		fmt.Printf("💬 Message: %s\n", message)

		// Check Pinchtab availability
		pinchtabURL := viper.GetString("pinchtab.url")
		if pinchtabURL == "" {
			pinchtabURL = "http://localhost:9867"
		}

		pinchtabClient := sdk.NewPinchtabClient(pinchtabURL, viper.GetString("pinchtab.token"))

		// Health check
		health, err := pinchtabClient.Health()
		if err != nil {
			fmt.Printf("❌ Pinchtab not available: %v\n", err)
			fmt.Println("\n💡 Start Pinchtab first:")
			fmt.Println("   ./pinchtab.exe --port 9867")
			return err
		}

		fmt.Printf("✅ Pinchtab running: %s\n", health.Status)

		// Step 1: Navigate to TeamChat (like Python script)
		teamchatURL := viper.GetString("server.url") + "/teamchat"
		fmt.Printf("🌐 Navigating to: %s\n", teamchatURL)

		if err := pinchtabClient.Navigate(teamchatURL); err != nil {
			return fmt.Errorf("failed to navigate: %w", err)
		}

		// Step 2: Wait for page load (like Python script)
		fmt.Println("⏳ Waiting for page to load...")
		time.Sleep(3 * time.Second)

		// Step 3: Get page snapshot
		snapshot, err := pinchtabClient.GetSnapshot()
		if err != nil {
			return fmt.Errorf("failed to get snapshot: %w", err)
		}

		fmt.Printf("📸 Page loaded: %s\n", snapshot.Title)

		// Step 4: Check if login is needed (like Python script logic)
		loginElements := findLoginFormElements(snapshot.Elements)
		if len(loginElements) > 0 {
			fmt.Println("🔐 Login form detected, performing authentication...")
			
			username := viper.GetString("auth.username")
			password := viper.GetString("auth.password")
			
			if username == "" || password == "" {
				return fmt.Errorf("credentials required in config")
			}

			success := performWebLogin(pinchtabClient, loginElements, username, password)
			if !success {
				return fmt.Errorf("web login failed")
			}

			fmt.Println("✅ Login successful")
			
			// Wait after login
			time.Sleep(2 * time.Second)
			
			// Get new snapshot after login
			snapshot, err = pinchtabClient.GetSnapshot()
			if err != nil {
				return fmt.Errorf("failed to get post-login snapshot: %w", err)
			}
		}

		// Step 5: Find the channel (like Python script room finding logic)
		fmt.Printf("🔍 Looking for channel: %s\n", channel)
		channelElement := findChannelElement(snapshot.Elements, channel)
		
		if channelElement.Ref == "" {
			fmt.Printf("⚠️  Channel '%s' not found. Available channels:\n", channel)
			displayAvailableChannels(snapshot.Elements)
			
			// Try common variations
			variations := []string{
				"Horus Monitoramento",
				"Monitoramento", 
				"General",
				"Geral",
			}
			
			for _, variation := range variations {
				if variation != channel {
					fmt.Printf("🔄 Trying variation: %s\n", variation)
					channelElement = findChannelElement(snapshot.Elements, variation)
					if channelElement.Ref != "" {
						channel = variation
						break
					}
				}
			}
		}

		if channelElement.Ref == "" {
			return fmt.Errorf("channel not found: %s", channel)
		}

		fmt.Printf("✅ Found channel: %s (ref: %s)\n", channel, channelElement.Ref)

		// Step 6: Click on channel
		if err := pinchtabClient.Click(channelElement.Ref); err != nil {
			return fmt.Errorf("failed to click channel: %w", err)
		}

		time.Sleep(2 * time.Second)

		// Step 7: Send message (like Python script message sending)
		if err := sendMessageViaWeb(pinchtabClient, message); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		fmt.Printf("🎉 Message sent successfully to '%s'!\n", channel)
		fmt.Println("✅ Python-style integration completed")

		return nil
	},
}

// Helper functions (similar to Python script logic)
func findChannelElement(elements []sdk.Element, channelName string) sdk.Element {
	// Look for elements containing the channel name (case insensitive)
	channelLower := strings.ToLower(channelName)
	
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			elementText := strings.ToLower(element.Text)
			
			// Exact match
			if elementText == channelLower {
				return element
			}
			
			// Contains match
			if strings.Contains(elementText, channelLower) || strings.Contains(channelLower, elementText) {
				return element
			}
		}
	}
	
	return sdk.Element{}
}

func displayAvailableChannels(elements []sdk.Element) {
	fmt.Println("Available clickable elements (potential channels):")
	count := 0
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			// Filter out common UI elements
			text := strings.ToLower(element.Text)
			if !strings.Contains(text, "button") && 
			   !strings.Contains(text, "menu") && 
			   !strings.Contains(text, "close") &&
			   len(element.Text) > 2 {
				fmt.Printf("  - %s (ref: %s)\n", element.Text, element.Ref)
				count++
				if count >= 10 {
					fmt.Println("  ... and more")
					break
				}
			}
		}
	}
}

func sendMessageViaWeb(client *sdk.PinchtabClient, message string) error {
	// Get current page snapshot to find message input
	snapshot, err := client.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get snapshot: %w", err)
	}

	// Look for message input field (like Python script)
	var messageInput sdk.Element
	for _, element := range snapshot.Elements {
		if element.Visible && element.Focusable {
			// Check for textarea or input fields that could be message inputs
			if element.Tag == "textarea" {
				messageInput = element
				break
			}
			
			if element.Tag == "input" {
				inputType := element.Attributes["type"]
				placeholder := strings.ToLower(element.Attributes["placeholder"])
				
				if inputType == "text" && (strings.Contains(placeholder, "message") || 
					strings.Contains(placeholder, "mensagem") || 
					strings.Contains(placeholder, "type") ||
					strings.Contains(placeholder, "digite")) {
					messageInput = element
					break
				}
			}
		}
	}

	if messageInput.Ref == "" {
		return fmt.Errorf("message input field not found")
	}

	fmt.Printf("📝 Found message input: %s\n", messageInput.Tag)

	// Click to focus the input
	if err := client.Click(messageInput.Ref); err != nil {
		return fmt.Errorf("failed to focus message input: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Type the message
	if err := client.Type(messageInput.Ref, message); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	time.Sleep(1 * time.Second)

	// Look for send button or press Enter
	snapshot, err = client.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get updated snapshot: %w", err)
	}

	var sendButton sdk.Element
	for _, element := range snapshot.Elements {
		if element.Visible && element.Clickable {
			text := strings.ToLower(element.Text)
			if strings.Contains(text, "send") || strings.Contains(text, "enviar") ||
			   element.Tag == "button" && (element.Attributes["type"] == "submit" ||
			   strings.Contains(strings.ToLower(element.Attributes["title"]), "send")) {
				sendButton = element
				break
			}
		}
	}

	if sendButton.Ref != "" {
		fmt.Println("🚀 Clicking send button...")
		return client.Click(sendButton.Ref)
	} else {
		// Press Enter as fallback (like Python script)
		fmt.Println("⌨️  Pressing Enter to send...")
		return client.PerformAction(sdk.ActionRequest{
			Action: "press",
			Ref:    messageInput.Ref,
			Key:    "Enter",
		})
	}
}

func init() {
	TeamChatCmd.AddCommand(pythonStyleCmd)
}