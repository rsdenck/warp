package pinchtab

import (
	"fmt"
	"os"
	"strings"

	"github.com/icewarp/warpctl/internal/output"
	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	PinchtabCmd = &cobra.Command{
		Use:   "pinchtab",
		Short: "Browser automation for Zabbix and TeamChat integration",
		Long:  `Advanced browser automation using Pinchtab for seamless Zabbix monitoring and TeamChat notifications`,
	}
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check Pinchtab service health",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		health, err := client.Health()
		if err != nil {
			return fmt.Errorf("failed to check health: %w", err)
		}

		fmt.Printf("Status: %s\n", health.Status)
		fmt.Printf("Version: %s\n", health.Version)
		fmt.Printf("Uptime: %s\n", health.Uptime)
		
		return nil
	},
}

var tabsCmd = &cobra.Command{
	Use:   "tabs",
	Short: "List all browser tabs",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		tabs, err := client.ListTabs()
		if err != nil {
			return fmt.Errorf("failed to list tabs: %w", err)
		}

		if len(tabs) == 0 {
			fmt.Println("No tabs open")
			return nil
		}

		t := output.NewTable("BROWSER TABS")
		t.AppendHeader(table.Row{"ID", "Title", "URL", "Active"})
		
		for _, tab := range tabs {
			active := "No"
			if tab.Active {
				active = "Yes"
			}
			t.AppendRow(table.Row{tab.ID, tab.Title, tab.URL, active})
		}
		
		t.Render()
		return nil
	},
}

var navigateCmd = &cobra.Command{
	Use:   "navigate [url]",
	Short: "Navigate to a URL",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		url := args[0]
		
		tabID, _ := cmd.Flags().GetString("tab")
		
		if err := client.Navigate(url, tabID); err != nil {
			return fmt.Errorf("failed to navigate: %w", err)
		}

		fmt.Printf("Navigated to: %s\n", url)
		return nil
	},
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Get accessibility tree snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		tabID, _ := cmd.Flags().GetString("tab")
		
		snapshot, err := client.GetSnapshot(tabID)
		if err != nil {
			return fmt.Errorf("failed to get snapshot: %w", err)
		}

		fmt.Printf("Page: %s\n", snapshot.Title)
		fmt.Printf("URL: %s\n", snapshot.URL)
		fmt.Printf("Elements: %d\n\n", len(snapshot.Elements))

		t := output.NewTable("ACCESSIBILITY TREE")
		t.AppendHeader(table.Row{"Ref", "Tag", "Text", "Role", "Clickable"})
		
		for _, element := range snapshot.Elements {
			if element.Visible {
				clickable := "No"
				if element.Clickable {
					clickable = "Yes"
				}
				
				text := element.Text
				if len(text) > 50 {
					text = text[:47] + "..."
				}
				
				t.AppendRow(table.Row{element.Ref, element.Tag, text, element.Role, clickable})
			}
		}
		
		t.Render()
		return nil
	},
}

var clickCmd = &cobra.Command{
	Use:   "click [ref-or-selector]",
	Short: "Click an element",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		if err := client.Click(args[0]); err != nil {
			return fmt.Errorf("failed to click: %w", err)
		}

		fmt.Printf("Clicked: %s\n", args[0])
		return nil
	},
}

var typeCmd = &cobra.Command{
	Use:   "type [ref-or-selector] [text]",
	Short: "Type text into an element",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		if err := client.Type(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to type: %w", err)
		}

		fmt.Printf("Typed '%s' into: %s\n", args[1], args[0])
		return nil
	},
}

var screenshotCmd = &cobra.Command{
	Use:   "screenshot [filename]",
	Short: "Take a screenshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		quality, _ := cmd.Flags().GetInt("quality")
		tabID, _ := cmd.Flags().GetString("tab")
		
		data, err := client.Screenshot(quality, tabID)
		if err != nil {
			return fmt.Errorf("failed to take screenshot: %w", err)
		}

		if err := os.WriteFile(args[0], data, 0644); err != nil {
			return fmt.Errorf("failed to save screenshot: %w", err)
		}

		fmt.Printf("Screenshot saved: %s\n", args[0])
		return nil
	},
}

// Advanced integration commands
var zabbixIntegrationCmd = &cobra.Command{
	Use:   "zabbix-integration",
	Short: "Automated Zabbix web interface integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		fmt.Println("🚀 Starting Zabbix Web Integration...")
		
		// Navigate to Zabbix
		zabbixURL := viper.GetString("zabbix.web_url")
		if zabbixURL == "" {
			zabbixURL = "https://monitoramento.armazem.cloud"
		}
		
		fmt.Printf("📡 Navigating to Zabbix: %s\n", zabbixURL)
		if err := client.Navigate(zabbixURL); err != nil {
			return fmt.Errorf("failed to navigate to Zabbix: %w", err)
		}

		// Wait and get snapshot
		fmt.Println("📸 Getting page snapshot...")
		snapshot, err := client.GetSnapshot()
		if err != nil {
			return fmt.Errorf("failed to get snapshot: %w", err)
		}

		fmt.Printf("✅ Page loaded: %s\n", snapshot.Title)
		fmt.Printf("📊 Found %d interactive elements\n", countClickableElements(snapshot.Elements))

		// Look for login elements
		loginElements := findLoginElements(snapshot.Elements)
		if len(loginElements) > 0 {
			fmt.Println("🔐 Login form detected, attempting authentication...")
			return performZabbixLogin(client, loginElements)
		}

		// Look for VDC groups
		fmt.Println("🔍 Searching for VDC host groups...")
		vdcElements := findVDCElements(snapshot.Elements)
		
		if len(vdcElements) > 0 {
			fmt.Printf("✅ Found %d VDC-related elements\n", len(vdcElements))
			displayVDCElements(vdcElements)
		} else {
			fmt.Println("⚠️  No VDC elements found on current page")
		}

		return nil
	},
}

var teamchatIntegrationCmd = &cobra.Command{
	Use:   "teamchat-integration",
	Short: "Automated TeamChat web interface integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := createPinchtabClient()
		
		fmt.Println("💬 Starting TeamChat Web Integration...")
		
		// Navigate to TeamChat
		teamchatURL := viper.GetString("server.url")
		if teamchatURL == "" {
			teamchatURL = "https://icewarp.armazemdc.inf.br"
		}
		
		teamchatURL += "/teamchat"
		
		fmt.Printf("📡 Navigating to TeamChat: %s\n", teamchatURL)
		if err := client.Navigate(teamchatURL); err != nil {
			return fmt.Errorf("failed to navigate to TeamChat: %w", err)
		}

		// Get snapshot
		fmt.Println("📸 Getting page snapshot...")
		snapshot, err := client.GetSnapshot()
		if err != nil {
			return fmt.Errorf("failed to get snapshot: %w", err)
		}

		fmt.Printf("✅ Page loaded: %s\n", snapshot.Title)
		
		// Look for channels
		fmt.Println("🔍 Searching for channels...")
		channelElements := findChannelElements(snapshot.Elements)
		
		if len(channelElements) > 0 {
			fmt.Printf("✅ Found %d channel elements\n", len(channelElements))
			displayChannelElements(channelElements)
		}

		return nil
	},
}

// Helper functions
func createPinchtabClient() *sdk.PinchtabClient {
	baseURL := viper.GetString("pinchtab.url")
	if baseURL == "" {
		baseURL = "http://localhost:9867"
	}
	
	token := viper.GetString("pinchtab.token")
	return sdk.NewPinchtabClient(baseURL, token)
}

func countClickableElements(elements []sdk.Element) int {
	count := 0
	for _, element := range elements {
		if element.Clickable && element.Visible {
			count++
		}
	}
	return count
}

func findLoginElements(elements []sdk.Element) []sdk.Element {
	var loginElements []sdk.Element
	for _, element := range elements {
		if element.Visible && (
			strings.Contains(strings.ToLower(element.Text), "login") ||
			strings.Contains(strings.ToLower(element.Text), "username") ||
			strings.Contains(strings.ToLower(element.Text), "password") ||
			element.Tag == "input" && (
				element.Attributes["type"] == "text" ||
				element.Attributes["type"] == "password" ||
				element.Attributes["type"] == "email")) {
			loginElements = append(loginElements, element)
		}
	}
	return loginElements
}

func findVDCElements(elements []sdk.Element) []sdk.Element {
	var vdcElements []sdk.Element
	for _, element := range elements {
		if element.Visible && strings.Contains(strings.ToUpper(element.Text), "VDC") {
			vdcElements = append(vdcElements, element)
		}
	}
	return vdcElements
}

func findChannelElements(elements []sdk.Element) []sdk.Element {
	var channelElements []sdk.Element
	for _, element := range elements {
		if element.Visible && (
			strings.Contains(strings.ToLower(element.Text), "channel") ||
			strings.Contains(strings.ToLower(element.Text), "horus") ||
			strings.Contains(strings.ToLower(element.Text), "monitoramento")) {
			channelElements = append(channelElements, element)
		}
	}
	return channelElements
}

func performZabbixLogin(client *sdk.PinchtabClient, loginElements []sdk.Element) error {
	username := viper.GetString("zabbix.username")
	password := viper.GetString("zabbix.password")
	
	if username == "" || password == "" {
		return fmt.Errorf("Zabbix credentials not configured")
	}

	// Find username and password fields
	var usernameField, passwordField, loginButton sdk.Element
	
	for _, element := range loginElements {
		if element.Tag == "input" {
			inputType := element.Attributes["type"]
			if inputType == "text" || inputType == "email" {
				usernameField = element
			} else if inputType == "password" {
				passwordField = element
			}
		} else if element.Clickable && strings.Contains(strings.ToLower(element.Text), "login") {
			loginButton = element
		}
	}

	// Perform login
	if usernameField.Ref != "" {
		fmt.Println("👤 Filling username...")
		if err := client.Fill(usernameField.Ref, username); err != nil {
			return fmt.Errorf("failed to fill username: %w", err)
		}
	}

	if passwordField.Ref != "" {
		fmt.Println("🔒 Filling password...")
		if err := client.Fill(passwordField.Ref, password); err != nil {
			return fmt.Errorf("failed to fill password: %w", err)
		}
	}

	if loginButton.Ref != "" {
		fmt.Println("🚀 Clicking login...")
		if err := client.Click(loginButton.Ref); err != nil {
			return fmt.Errorf("failed to click login: %w", err)
		}
	}

	fmt.Println("✅ Login attempt completed")
	return nil
}

func displayVDCElements(elements []sdk.Element) {
	t := output.NewTable("VDC ELEMENTS FOUND")
	t.AppendHeader(table.Row{"Ref", "Tag", "Text", "Role"})
	
	for _, element := range elements {
		text := element.Text
		if len(text) > 60 {
			text = text[:57] + "..."
		}
		t.AppendRow(table.Row{element.Ref, element.Tag, text, element.Role})
	}
	
	t.Render()
}

func displayChannelElements(elements []sdk.Element) {
	t := output.NewTable("CHANNEL ELEMENTS FOUND")
	t.AppendHeader(table.Row{"Ref", "Tag", "Text", "Role"})
	
	for _, element := range elements {
		text := element.Text
		if len(text) > 60 {
			text = text[:57] + "..."
		}
		t.AppendRow(table.Row{element.Ref, element.Tag, text, element.Role})
	}
	
	t.Render()
}

func init() {
	PinchtabCmd.AddCommand(healthCmd)
	PinchtabCmd.AddCommand(tabsCmd)
	PinchtabCmd.AddCommand(navigateCmd)
	PinchtabCmd.AddCommand(snapshotCmd)
	PinchtabCmd.AddCommand(clickCmd)
	PinchtabCmd.AddCommand(typeCmd)
	PinchtabCmd.AddCommand(screenshotCmd)
	PinchtabCmd.AddCommand(zabbixIntegrationCmd)
	PinchtabCmd.AddCommand(teamchatIntegrationCmd)

	navigateCmd.Flags().StringP("tab", "t", "", "Tab ID to navigate")
	snapshotCmd.Flags().StringP("tab", "t", "", "Tab ID to snapshot")
	screenshotCmd.Flags().IntP("quality", "q", 80, "JPEG quality (1-100)")
	screenshotCmd.Flags().StringP("tab", "t", "", "Tab ID to screenshot")
}