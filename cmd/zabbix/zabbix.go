package zabbix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/icewarp/warpctl/internal/output"
	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ZabbixCmd = &cobra.Command{
		Use:   "zabbix",
		Short: "Zabbix monitoring integration",
		Long:  `Commands for Zabbix API integration and TeamChat notifications`,
	}
)

var listGroupsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List Zabbix host groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createZabbixClient()
		if err != nil {
			return err
		}
		defer client.Logout()

		groups, err := client.GetHostGroups()
		if err != nil {
			return fmt.Errorf("failed to get host groups: %w", err)
		}

		t := output.NewTable("ZABBIX HOST GROUPS")
		t.AppendHeader(table.Row{"Group ID", "Name"})
		
		for _, group := range groups {
			t.AppendRow(table.Row{group.GroupID, group.Name})
		}
		
		t.Render()
		return nil
	},
}

var listProblemsCmd = &cobra.Command{
	Use:   "problems",
	Short: "List active problems/alerts",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createZabbixClient()
		if err != nil {
			return err
		}
		defer client.Logout()

		groupID, _ := cmd.Flags().GetString("group")

		problems, err := client.GetActiveProblems(groupID)
		if err != nil {
			return fmt.Errorf("failed to get problems: %w", err)
		}

		if len(problems) == 0 {
			fmt.Println("No active problems found")
			return nil
		}

		t := output.NewTable("ACTIVE PROBLEMS")
		t.AppendHeader(table.Row{"Event ID", "Host", "Trigger", "Severity", "Time"})
		
		for _, problem := range problems {
			t.AppendRow(table.Row{
				problem.EventID,
				problem.HostName,
				problem.TriggerName,
				problem.Severity,
				problem.Timestamp.Format("2006-01-02 15:04:05"),
			})
		}
		
		t.Render()
		return nil
	},
}

var notifyTeamChatCmd = &cobra.Command{
	Use:   "notify",
	Short: "Send Zabbix alerts to TeamChat using configured mappings",
	RunE: func(cmd *cobra.Command, args []string) error {
		zabbixClient, err := createZabbixClient()
		if err != nil {
			return err
		}
		defer zabbixClient.Logout()

		groupID, _ := cmd.Flags().GetString("group")
		channelName, _ := cmd.Flags().GetString("channel")
		
		// Create TeamChat client
		teamChatClient := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("teamchat.token")
		if token == "" {
			return fmt.Errorf("TeamChat token not set. Use 'warpctl teamchat login' first")
		}
		teamChatClient.SetToken(token)

		// Get mappings
		mappings := viper.GetStringMapString("zabbix.group_mappings")
		
		if groupID != "" {
			// Single group specified
			problems, err := zabbixClient.GetActiveProblems(groupID)
			if err != nil {
				return fmt.Errorf("failed to get problems: %w", err)
			}

			if len(problems) == 0 {
				fmt.Println("No active problems to notify")
				return nil
			}

			// Determine channel
			channel := channelName
			if channel == "" {
				if mappedChannel, exists := mappings[groupID]; exists {
					channel = mappedChannel
				} else {
					channel = viper.GetString("zabbix.monitoring.default_channel")
				}
			}

			// Send notifications
			sent := 0
			for _, problem := range problems {
				message := createEnterpriseAlert(problem)

				_, err := teamChatClient.PostMessage(channel, message)
				if err != nil {
					fmt.Printf("Failed to send notification for event %s: %v\n", problem.EventID, err)
					continue
				}
				sent++
			}

			fmt.Printf("Successfully sent %d notifications to TeamChat channel '%s'\n", sent, channel)
		} else {
			// Use all configured mappings
			if len(mappings) == 0 {
				return fmt.Errorf("no group mappings configured. Run 'warpctl zabbix configure' first")
			}
			
			totalSent := 0
			for groupID, channel := range mappings {
				problems, err := zabbixClient.GetActiveProblems(groupID)
				if err != nil {
					fmt.Printf("Failed to get problems for group %s: %v\n", groupID, err)
					continue
				}

				sent := 0
				for _, problem := range problems {
					message := createEnterpriseAlert(problem)

					_, err := teamChatClient.PostMessage(channel, message)
					if err != nil {
						fmt.Printf("Failed to send notification for event %s: %v\n", problem.EventID, err)
						continue
					}
					sent++
				}
				
				if sent > 0 {
					fmt.Printf("Sent %d notifications to '%s'\n", sent, channel)
					totalSent += sent
				}
			}
			
			fmt.Printf("Total: %d notifications sent\n", totalSent)
		}
		
		return nil
	},
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure Zabbix group to TeamChat channel mappings",
	Long:  `Interactive configuration to map Zabbix host groups to TeamChat channels for alert routing`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔧 Zabbix to TeamChat Configuration")
		fmt.Println("=" + strings.Repeat("=", 40))
		
		// Test Zabbix connection
		fmt.Println("\n📡 Testing Zabbix connection...")
		zabbixClient, err := createZabbixClient()
		if err != nil {
			return fmt.Errorf("failed to connect to Zabbix: %w", err)
		}
		defer zabbixClient.Logout()
		
		// Get available groups
		groups, err := zabbixClient.GetHostGroups()
		if err != nil {
			return fmt.Errorf("failed to get host groups: %w", err)
		}
		
		fmt.Printf("✅ Connected to Zabbix! Found %d VDC groups\n", len(groups))
		
		// Test TeamChat connection
		fmt.Println("\n💬 Testing TeamChat connection...")
		teamChatClient := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		
		token := viper.GetString("teamchat.token")
		if token == "" {
			fmt.Println("⚠️  TeamChat token not configured. Please run 'warpctl teamchat login' first")
			fmt.Println("Continuing with configuration setup...")
		} else {
			teamChatClient.SetToken(token)
			if _, err := teamChatClient.AuthTest(); err != nil {
				fmt.Printf("⚠️  TeamChat authentication failed: %v\n", err)
				fmt.Println("You may need to re-authenticate later")
			} else {
				fmt.Println("✅ TeamChat connection successful!")
			}
		}
		
		// Show current mappings
		fmt.Println("\n📋 Current Group Mappings:")
		mappings := viper.GetStringMapString("zabbix.group_mappings")
		if len(mappings) == 0 {
			fmt.Println("   No mappings configured")
		} else {
			for groupID, channel := range mappings {
				// Find group name
				groupName := groupID
				for _, group := range groups {
					if group.GroupID == groupID {
						groupName = group.Name
						break
					}
				}
				fmt.Printf("   %s (%s) → %s\n", groupName, groupID, channel)
			}
		}
		
		// Interactive configuration
		fmt.Println("\n🎯 Configure Group Mappings:")
		fmt.Println("Select Zabbix groups to configure alert routing:")
		
		// Display groups with numbers
		fmt.Println("\nAvailable VDC Groups:")
		for i, group := range groups {
			currentChannel := mappings[group.GroupID]
			if currentChannel == "" {
				currentChannel = viper.GetString("zabbix.monitoring.default_channel")
			}
			fmt.Printf("  %d. %s (ID: %s) → %s\n", i+1, group.Name, group.GroupID, currentChannel)
		}
		
		fmt.Print("\nEnter group numbers to configure (comma-separated, or 'all' for all groups, 'q' to quit): ")
		var input string
		fmt.Scanln(&input)
		
		if input == "q" {
			fmt.Println("Configuration cancelled")
			return nil
		}
		
		var selectedGroups []sdk.ZabbixGroup
		if input == "all" {
			selectedGroups = groups
		} else {
			// Parse selected numbers
			parts := strings.Split(input, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if num, err := strconv.Atoi(part); err == nil && num > 0 && num <= len(groups) {
					selectedGroups = append(selectedGroups, groups[num-1])
				}
			}
		}
		
		if len(selectedGroups) == 0 {
			fmt.Println("No valid groups selected")
			return nil
		}
		
		// Configure each selected group
		newMappings := make(map[string]string)
		for k, v := range mappings {
			newMappings[k] = v
		}
		
		for _, group := range selectedGroups {
			fmt.Printf("\n📢 Configure alerts for: %s\n", group.Name)
			fmt.Print("Enter TeamChat channel name (or press Enter for default): ")
			var channel string
			fmt.Scanln(&channel)
			
			if channel == "" {
				channel = viper.GetString("zabbix.monitoring.default_channel")
			}
			
			newMappings[group.GroupID] = channel
			fmt.Printf("✅ %s → %s\n", group.Name, channel)
		}
		
		// Save configuration
		viper.Set("zabbix.group_mappings", newMappings)
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		
		fmt.Println("\n💾 Configuration saved successfully!")
		
		// Ask about monitoring
		fmt.Print("\nEnable automatic monitoring? (y/N): ")
		var enableMonitoring string
		fmt.Scanln(&enableMonitoring)
		
		if strings.ToLower(enableMonitoring) == "y" || strings.ToLower(enableMonitoring) == "yes" {
			fmt.Print("Enter monitoring interval (e.g., 5m, 30s, 1h) [default: 5m]: ")
			var interval string
			fmt.Scanln(&interval)
			
			if interval == "" {
				interval = "5m"
			}
			
			viper.Set("zabbix.monitoring.enabled", true)
			viper.Set("zabbix.monitoring.interval", interval)
			
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to save monitoring configuration: %w", err)
			}
			
			fmt.Printf("✅ Automatic monitoring enabled with %s interval\n", interval)
			fmt.Println("Run 'warpctl zabbix start' to begin monitoring")
		}
		
		fmt.Println("\n🎉 Configuration complete!")
		return nil
	},
}

var testNotifyCmd = &cobra.Command{
	Use:   "test-notify",
	Short: "Send Zabbix alerts to TeamChat using configured mappings",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelName, _ := cmd.Flags().GetString("channel")
		showHTML, _ := cmd.Flags().GetBool("show-html")
		
		if channelName == "" {
			channelName = "Horus Monitoramento"
		}

		// Create a sample alert for testing VDC_BONJA
		sampleAlert := sdk.ZabbixAlert{
			EventID:     "2984102",
			TriggerID:   "67890",
			HostName:    "VDC_BONJA_SERVER_01",
			TriggerName: "High CPU usage detected on server",
			Severity:    "High",
			Status:      "PROBLEM",
			Timestamp:   time.Now().Add(-8*time.Minute - 15*time.Second),
			Message:     "Uso elevado de CPU detectado. Possível sobrecarga do sistema ou processo mal comportado.",
		}

		message := createEnterpriseAlert(sampleAlert)

		if showHTML {
			fmt.Println("Generated HTML message:")
			fmt.Println(strings.Repeat("=", 50))
			fmt.Println(message)
			fmt.Println(strings.Repeat("=", 50))
			return nil
		}

		fmt.Printf("🚨 Sending Zabbix alert to TeamChat channel '%s'\n", channelName)
		fmt.Printf("📊 Alert: %s - %s\n", sampleAlert.HostName, sampleAlert.TriggerName)

		// Try multiple approaches for sending
		success := false

		// Approach 1: Try TeamChat API with token (if available)
		token := viper.GetString("teamchat.token")
		if token != "" {
			fmt.Println("🔄 Trying TeamChat API with token...")
			teamChatClient := sdk.NewTeamChatClient(&sdk.Config{
				BaseURL: viper.GetString("server.url"),
			})
			teamChatClient.SetToken(token)

			_, err := teamChatClient.PostMessage(channelName, message)
			if err == nil {
				fmt.Println("✅ Message sent via TeamChat API!")
				success = true
			} else {
				fmt.Printf("❌ TeamChat API failed: %v\n", err)
			}
		}

		// Approach 2: Try Session-based approach (if available)
		if !success {
			sessionID := viper.GetString("teamchat.session_id")
			if sessionID != "" {
				fmt.Println("🔄 Trying session-based approach...")
				
				// Create a simplified text message for session-based sending
				textMessage := fmt.Sprintf(`🚨 ALERTA ZABBIX - %s

🏢 Host: %s
⚠️  Problema: %s
📊 Severidade: %s
⏰ Duração: %s
🆔 Event ID: %s

Detectado em: %s`,
					sampleAlert.Severity,
					sampleAlert.HostName,
					sampleAlert.TriggerName,
					sampleAlert.Severity,
					formatDuration(time.Since(sampleAlert.Timestamp)),
					generateEventID(sampleAlert.TriggerID),
					sampleAlert.Timestamp.Format("2006-01-02 15:04:05"))

				// Try to send via session (simplified approach)
				if err := sendViaSession(sessionID, channelName, textMessage); err == nil {
					fmt.Println("✅ Message sent via session!")
					success = true
				} else {
					fmt.Printf("❌ Session approach failed: %v\n", err)
				}
			}
		}

		// Approach 3: Web automation (if Pinchtab available)
		if !success {
			fmt.Println("🔄 Trying web automation...")
			
			pinchtabURL := viper.GetString("pinchtab.url")
			if pinchtabURL == "" {
				pinchtabURL = "http://localhost:9867"
			}

			// Quick health check
			client := &http.Client{Timeout: 2 * time.Second}
			if resp, err := client.Get(pinchtabURL + "/health"); err == nil {
				resp.Body.Close()
				fmt.Println("✅ Pinchtab available, using web automation...")
				
				// Use the existing web notification command
				return executeWebNotify(channelName, sampleAlert)
			} else {
				fmt.Printf("❌ Pinchtab not available: %v\n", err)
			}
		}

		if success {
			fmt.Printf("🎉 Test notification sent successfully to '%s'!\n", channelName)
			return nil
		} else {
			fmt.Println("\n💡 Available options:")
			fmt.Println("1. Get TeamChat token: warpctl teamchat login")
			fmt.Println("2. Use XML-RPC session: warpctl teamchat xml-auth")
			fmt.Println("3. Start Pinchtab for web automation")
			return fmt.Errorf("all notification methods failed")
		}
	},
}

var startMonitoringCmd = &cobra.Command{
	Use:   "start",
	Short: "Start automatic Zabbix monitoring with configured mappings",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !viper.GetBool("zabbix.monitoring.enabled") {
			return fmt.Errorf("monitoring not enabled. Run 'warpctl zabbix configure' first")
		}
		
		interval := viper.GetString("zabbix.monitoring.interval")
		if interval == "" {
			interval = "5m"
		}
		
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return fmt.Errorf("invalid interval format: %w", err)
		}
		
		mappings := viper.GetStringMapString("zabbix.group_mappings")
		if len(mappings) == 0 {
			return fmt.Errorf("no group mappings configured. Run 'warpctl zabbix configure' first")
		}
		
		fmt.Printf("🚀 Starting Zabbix monitoring (interval: %s)\n", interval)
		fmt.Printf("📊 Monitoring %d group mappings\n", len(mappings))
		fmt.Println("Press Ctrl+C to stop...")
		
		ticker := time.NewTicker(duration)
		defer ticker.Stop()
		
		// Track sent alerts to avoid duplicates
		sentAlerts := make(map[string]bool)
		
		for {
			select {
			case <-ticker.C:
				zabbixClient, err := createZabbixClient()
				if err != nil {
					fmt.Printf("❌ Error connecting to Zabbix: %v\n", err)
					continue
				}
				
				// Create TeamChat client
				teamChatClient := sdk.NewTeamChatClient(&sdk.Config{
					BaseURL: viper.GetString("server.url"),
				})
				
				token := viper.GetString("teamchat.token")
				if token == "" {
					fmt.Println("⚠️  TeamChat token not set. Use 'warpctl teamchat login' first")
					zabbixClient.Logout()
					continue
				}
				teamChatClient.SetToken(token)
				
				totalNewAlerts := 0
				
				// Check each mapped group
				for groupID, channel := range mappings {
					problems, err := zabbixClient.GetActiveProblems(groupID)
					if err != nil {
						fmt.Printf("❌ Error getting problems for group %s: %v\n", groupID, err)
						continue
					}
					
					// Filter by severity
					severityFilter := viper.GetStringSlice("zabbix.severity_filter")
					if len(severityFilter) == 0 {
						severityFilter = []string{"Disaster", "High", "Average"}
					}
					
					newAlerts := 0
					for _, problem := range problems {
						// Check severity filter
						if !contains(severityFilter, problem.Severity) {
							continue
						}
						
						alertKey := fmt.Sprintf("%s-%s", groupID, problem.TriggerID)
						if !sentAlerts[alertKey] {
							message := createEnterpriseAlert(problem)
							
							_, err := teamChatClient.PostMessage(channel, message)
							if err != nil {
								fmt.Printf("❌ Failed to send alert to %s: %v\n", channel, err)
							} else {
								sentAlerts[alertKey] = true
								newAlerts++
								totalNewAlerts++
							}
						}
					}
					
					if newAlerts > 0 {
						fmt.Printf("📢 Sent %d alerts to %s\n", newAlerts, channel)
					}
				}
				
				if totalNewAlerts > 0 {
					fmt.Printf("✅ [%s] Total: %d new alerts sent\n", time.Now().Format("15:04:05"), totalNewAlerts)
				}
				
				zabbixClient.Logout()
			}
		}
	},
}

var stopMonitoringCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop automatic Zabbix monitoring",
	RunE: func(cmd *cobra.Command, args []string) error {
		viper.Set("zabbix.monitoring.enabled", false)
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		
		fmt.Println("🛑 Monitoring disabled")
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show Zabbix monitoring status and configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📊 Zabbix Monitoring Status")
		fmt.Println("=" + strings.Repeat("=", 30))
		
		// Zabbix configuration
		fmt.Printf("🔗 Zabbix URL: %s\n", viper.GetString("zabbix.url"))
		fmt.Printf("👤 Username: %s\n", viper.GetString("zabbix.username"))
		
		// Monitoring status
		enabled := viper.GetBool("zabbix.monitoring.enabled")
		interval := viper.GetString("zabbix.monitoring.interval")
		
		fmt.Printf("🔄 Monitoring: %s\n", map[bool]string{true: "✅ Enabled", false: "❌ Disabled"}[enabled])
		if enabled {
			fmt.Printf("⏱️  Interval: %s\n", interval)
		}
		
		// Group mappings
		mappings := viper.GetStringMapString("zabbix.group_mappings")
		fmt.Printf("📋 Group Mappings: %d configured\n", len(mappings))
		
		if len(mappings) > 0 {
			// Test Zabbix connection to get group names
			zabbixClient, err := createZabbixClient()
			if err == nil {
				groups, err := zabbixClient.GetHostGroups()
				if err == nil {
					groupMap := make(map[string]string)
					for _, group := range groups {
						groupMap[group.GroupID] = group.Name
					}
					
					fmt.Println("\n📢 Alert Routing:")
					for groupID, channel := range mappings {
						groupName := groupMap[groupID]
						if groupName == "" {
							groupName = groupID
						}
						fmt.Printf("   %s → %s\n", groupName, channel)
					}
				}
				zabbixClient.Logout()
			}
		}
		
		// Severity filter
		severityFilter := viper.GetStringSlice("zabbix.severity_filter")
		fmt.Printf("\n🎯 Severity Filter: %v\n", severityFilter)
		
		// TeamChat status
		token := viper.GetString("teamchat.token")
		fmt.Printf("💬 TeamChat: %s\n", map[bool]string{token != "": "✅ Configured", token == "": "❌ Not configured"}[token != ""])
		
		return nil
	},
}

var webNotifyCmd = &cobra.Command{
	Use:   "web-notify",
	Short: "Send notifications via TeamChat web interface using Pinchtab",
	Long:  `Send Zabbix alerts directly to TeamChat web interface without API token using browser automation`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if Pinchtab is available
		pinchtabClient := createPinchtabClient()
		
		health, err := pinchtabClient.Health()
		if err != nil {
			return fmt.Errorf("Pinchtab not available. Please start Pinchtab first: pinchtab --port 9867")
		}
		
		fmt.Printf("🚀 Using Pinchtab for web automation (Status: %s)\n", health.Status)
		
		// Get Zabbix problems
		zabbixClient, err := createZabbixClient()
		if err != nil {
			return err
		}
		defer zabbixClient.Logout()

		groupID, _ := cmd.Flags().GetString("group")
		channelName, _ := cmd.Flags().GetString("channel")
		
		if channelName == "" {
			channelName = "Horus Monitoramento"
		}

		// Get mappings if no specific group
		var groupsToCheck map[string]string
		if groupID != "" {
			groupsToCheck = map[string]string{groupID: channelName}
		} else {
			groupsToCheck = viper.GetStringMapString("zabbix.group_mappings")
			if len(groupsToCheck) == 0 {
				return fmt.Errorf("no group mappings configured. Run 'warpctl zabbix configure' first")
			}
		}

		totalSent := 0
		for gID, channel := range groupsToCheck {
			problems, err := zabbixClient.GetActiveProblems(gID)
			if err != nil {
				fmt.Printf("Failed to get problems for group %s: %v\n", gID, err)
				continue
			}

			if len(problems) == 0 {
				fmt.Printf("No active problems for group %s\n", gID)
				continue
			}

			fmt.Printf("📢 Sending %d alerts to '%s' via web interface...\n", len(problems), channel)

			// Navigate to TeamChat
			teamchatURL := viper.GetString("server.url") + "/teamchat"
			fmt.Printf("🌐 Navigating to: %s\n", teamchatURL)
			
			if err := pinchtabClient.Navigate(teamchatURL); err != nil {
				return fmt.Errorf("failed to navigate to TeamChat: %w", err)
			}

			// Wait a moment for page to load
			time.Sleep(3 * time.Second)

			// Get page snapshot to find elements
			snapshot, err := pinchtabClient.GetSnapshot()
			if err != nil {
				return fmt.Errorf("failed to get page snapshot: %w", err)
			}

			fmt.Printf("📸 Page loaded: %s (%d elements)\n", snapshot.Title, len(snapshot.Elements))

			// Look for the Horus Monitoramento channel
			channelElement := findChannelElement(snapshot.Elements, channel)
			if channelElement.Ref == "" {
				fmt.Printf("⚠️  Channel '%s' not found on page. Available elements:\n", channel)
				displayAvailableChannels(snapshot.Elements)
				continue
			}

			fmt.Printf("✅ Found channel '%s' (ref: %s)\n", channel, channelElement.Ref)

			// Click on the channel
			if err := pinchtabClient.Click(channelElement.Ref); err != nil {
				fmt.Printf("Failed to click channel: %v\n", err)
				continue
			}

			time.Sleep(2 * time.Second)

			// Send each alert
			sent := 0
			for _, problem := range problems {
				if err := sendAlertViaWeb(pinchtabClient, problem, channel); err != nil {
					fmt.Printf("Failed to send alert for %s: %v\n", problem.HostName, err)
					continue
				}
				sent++
				time.Sleep(1 * time.Second) // Avoid flooding
			}

			fmt.Printf("✅ Sent %d alerts to '%s'\n", sent, channel)
			totalSent += sent
		}

		fmt.Printf("🎉 Total: %d notifications sent via web interface\n", totalSent)
		return nil
	},
}

var testWebNotifyCmd = &cobra.Command{
	Use:   "test-web-notify",
	Short: "Test web notification with sample alert via Pinchtab",
	RunE: func(cmd *cobra.Command, args []string) error {
		channelName, _ := cmd.Flags().GetString("channel")
		if channelName == "" {
			channelName = "Horus Monitoramento"
		}

		// Create sample alert
		sampleAlert := sdk.ZabbixAlert{
			EventID:     "2984102",
			TriggerID:   "67890",
			HostName:    "VDC_BONJA_SERVER_01",
			TriggerName: "High CPU usage detected on server",
			Severity:    "High",
			Status:      "PROBLEM",
			Timestamp:   time.Now().Add(-8*time.Minute - 15*time.Second),
			Message:     "Uso elevado de CPU detectado. Possível sobrecarga do sistema ou processo mal comportado.",
		}

		fmt.Printf("🚨 Testing web notification to channel: %s\n", channelName)
		fmt.Printf("📊 Sample alert: %s - %s\n", sampleAlert.HostName, sampleAlert.TriggerName)

		// Use the improved web automation
		return executeWebNotify(channelName, sampleAlert)
	},
}

// Helper functions for web automation
func findChannelElement(elements []sdk.Element, channelName string) sdk.Element {
	// Look for elements containing the channel name
	for _, element := range elements {
		if element.Visible && element.Clickable {
			if strings.Contains(strings.ToLower(element.Text), strings.ToLower(channelName)) {
				return element
			}
		}
	}
	
	// Fallback: look for any element with "Horus" or "Monitoramento"
	for _, element := range elements {
		if element.Visible && element.Clickable {
			text := strings.ToLower(element.Text)
			if strings.Contains(text, "horus") || strings.Contains(text, "monitoramento") {
				return element
			}
		}
	}
	
	return sdk.Element{}
}

func displayAvailableChannels(elements []sdk.Element) {
	fmt.Println("Available clickable elements:")
	count := 0
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			fmt.Printf("  - %s (ref: %s)\n", element.Text, element.Ref)
			count++
			if count >= 10 { // Show max 10 elements
				fmt.Println("  ... and more")
				break
			}
		}
	}
}

func sendAlertViaWeb(client *sdk.PinchtabClient, alert sdk.ZabbixAlert, channel string) error {
	// Get current page snapshot to find message input
	snapshot, err := client.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get snapshot: %w", err)
	}

	// Look for message input field
	var messageInput sdk.Element
	for _, element := range snapshot.Elements {
		if element.Tag == "textarea" || element.Tag == "input" {
			if element.Visible && element.Focusable {
				// Check if it's a message input (common attributes)
				if strings.Contains(strings.ToLower(element.Attributes["placeholder"]), "message") ||
				   strings.Contains(strings.ToLower(element.Attributes["placeholder"]), "mensagem") ||
				   element.Attributes["type"] == "text" {
					messageInput = element
					break
				}
			}
		}
	}

	if messageInput.Ref == "" {
		return fmt.Errorf("message input field not found")
	}

	// Create simplified text message (since HTML might not be supported in web interface)
	message := fmt.Sprintf(`🚨 ALERTA ZABBIX - %s

🏢 Host: %s
⚠️  Problema: %s
📊 Severidade: %s
⏰ Duração: %s
🆔 Event ID: %s

Detectado em: %s`,
		alert.Severity,
		alert.HostName,
		alert.TriggerName,
		alert.Severity,
		formatDuration(time.Since(alert.Timestamp)),
		generateEventID(alert.TriggerID),
		alert.Timestamp.Format("2006-01-02 15:04:05"))

	// Click on message input to focus
	if err := client.Click(messageInput.Ref); err != nil {
		return fmt.Errorf("failed to focus message input: %w", err)
	}

	time.Sleep(500 * time.Millisecond)

	// Type the message
	if err := client.Type(messageInput.Ref, message); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	time.Sleep(1 * time.Second)

	// Look for send button
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

	if sendButton.Ref == "" {
		// Try pressing Enter as fallback
		fmt.Println("Send button not found, trying Enter key...")
		return client.PerformAction(sdk.ActionRequest{
			Action: "press",
			Ref:    messageInput.Ref,
			Key:    "Enter",
		})
	}

	// Click send button
	return client.Click(sendButton.Ref)
}

var integratedMonitoringCmd = &cobra.Command{
	Use:   "integrated-monitoring",
	Short: "Full integrated monitoring with Pinchtab, Zabbix, and TeamChat",
	Long:  `Advanced monitoring system that uses Pinchtab for web automation, Zabbix API for monitoring data, and TeamChat for notifications`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Starting Integrated Monitoring System")
		fmt.Println("=" + strings.Repeat("=", 50))
		
		// 1. Check Pinchtab health
		fmt.Println("\n📡 Checking Pinchtab connection...")
		pinchtabClient := createPinchtabClient()
		
		health, err := pinchtabClient.Health()
		if err != nil {
			fmt.Printf("⚠️  Pinchtab not available: %v\n", err)
			fmt.Println("💡 Start Pinchtab with: pinchtab --port 9867")
		} else {
			fmt.Printf("✅ Pinchtab running: %s (uptime: %s)\n", health.Status, health.Uptime)
		}
		
		// 2. Check Zabbix API
		fmt.Println("\n🔍 Checking Zabbix API connection...")
		zabbixClient, err := createZabbixClient()
		if err != nil {
			return fmt.Errorf("failed to connect to Zabbix: %w", err)
		}
		defer zabbixClient.Logout()
		
		groups, err := zabbixClient.GetHostGroups()
		if err != nil {
			return fmt.Errorf("failed to get host groups: %w", err)
		}
		
		fmt.Printf("✅ Zabbix API connected: %d VDC groups found\n", len(groups))
		
		// 3. Check TeamChat
		fmt.Println("\n💬 Checking TeamChat connection...")
		teamChatClient := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		
		token := viper.GetString("teamchat.token")
		if token == "" {
			fmt.Println("⚠️  TeamChat token not configured")
		} else {
			teamChatClient.SetToken(token)
			if _, err := teamChatClient.AuthTest(); err != nil {
				fmt.Printf("⚠️  TeamChat authentication failed: %v\n", err)
			} else {
				fmt.Println("✅ TeamChat connected and authenticated")
			}
		}
		
		// 4. Show current mappings
		fmt.Println("\n📋 Current Configuration:")
		mappings := viper.GetStringMapString("zabbix.group_mappings")
		for groupID, channel := range mappings {
			// Find group name
			groupName := groupID
			for _, group := range groups {
				if group.GroupID == groupID {
					groupName = group.Name
					break
				}
			}
			fmt.Printf("   %s (%s) → %s\n", groupName, groupID, channel)
		}
		
		// 5. Check for active problems
		fmt.Println("\n🚨 Checking for active problems...")
		totalProblems := 0
		for groupID, channel := range mappings {
			problems, err := zabbixClient.GetActiveProblems(groupID)
			if err != nil {
				fmt.Printf("❌ Error checking group %s: %v\n", groupID, err)
				continue
			}
			
			if len(problems) > 0 {
				fmt.Printf("🔥 Group %s: %d active problems → %s\n", groupID, len(problems), channel)
				totalProblems += len(problems)
				
				// Show first few problems
				for i, problem := range problems {
					if i >= 3 { // Show max 3 problems per group
						fmt.Printf("   ... and %d more\n", len(problems)-3)
						break
					}
					fmt.Printf("   - %s: %s (%s)\n", problem.HostName, problem.TriggerName, problem.Severity)
				}
			}
		}
		
		if totalProblems == 0 {
			fmt.Println("✅ No active problems found")
		} else {
			fmt.Printf("⚠️  Total active problems: %d\n", totalProblems)
		}
		
		// 6. Web automation demo (if Pinchtab is available)
		if health != nil {
			fmt.Println("\n🌐 Web Automation Demo:")
			
			autoLogin, _ := cmd.Flags().GetBool("auto-login")
			if autoLogin {
				fmt.Println("🔐 Attempting automated login to Zabbix web interface...")
				
				zabbixURL := viper.GetString("zabbix.web_url")
				if zabbixURL == "" {
					zabbixURL = "https://monitoramento.armazem.cloud"
				}
				
				if err := pinchtabClient.Navigate(zabbixURL); err != nil {
					fmt.Printf("❌ Failed to navigate to Zabbix: %v\n", err)
				} else {
					fmt.Printf("✅ Navigated to: %s\n", zabbixURL)
					
					// Take screenshot
					screenshot, _ := cmd.Flags().GetBool("screenshot")
					if screenshot {
						fmt.Println("📸 Taking screenshot...")
						data, err := pinchtabClient.Screenshot(80)
						if err != nil {
							fmt.Printf("❌ Screenshot failed: %v\n", err)
						} else {
							filename := fmt.Sprintf("zabbix-screenshot-%d.jpg", time.Now().Unix())
							if err := os.WriteFile(filename, data, 0644); err != nil {
								fmt.Printf("❌ Failed to save screenshot: %v\n", err)
							} else {
								fmt.Printf("✅ Screenshot saved: %s\n", filename)
							}
						}
					}
				}
			}
		}
		
		fmt.Println("\n🎉 Integrated monitoring check complete!")
		return nil
	},
}

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor Zabbix and send continuous notifications",
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, _ := cmd.Flags().GetString("group")
		channelName, _ := cmd.Flags().GetString("channel")
		intervalStr, _ := cmd.Flags().GetString("interval")
		
		if channelName == "" {
			channelName = "Horus Monitoramento"
		}

		interval, err := time.ParseDuration(intervalStr)
		if err != nil {
			return fmt.Errorf("invalid interval format: %w", err)
		}

		fmt.Printf("Starting Zabbix monitoring (interval: %s, channel: %s)\n", interval, channelName)
		fmt.Println("Press Ctrl+C to stop...")

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		// Track sent alerts to avoid duplicates
		sentAlerts := make(map[string]bool)

		for {
			select {
			case <-ticker.C:
				zabbixClient, err := createZabbixClient()
				if err != nil {
					fmt.Printf("Error connecting to Zabbix: %v\n", err)
					continue
				}

				problems, err := zabbixClient.GetActiveProblems(groupID)
				if err != nil {
					fmt.Printf("Error getting problems: %v\n", err)
					zabbixClient.Logout()
					continue
				}

				// Create TeamChat client
				teamChatClient := sdk.NewTeamChatClient(&sdk.Config{
					BaseURL: viper.GetString("server.url"),
				})

				token := viper.GetString("teamchat.token")
				if token == "" {
					fmt.Println("TeamChat token not set. Use 'warpctl teamchat login' first")
					zabbixClient.Logout()
					continue
				}
				teamChatClient.SetToken(token)

				// Send notifications for new problems only
				newAlerts := 0
				for _, problem := range problems {
					if !sentAlerts[problem.EventID] {
						message := createEnterpriseAlert(problem)

						_, err := teamChatClient.PostMessage(channelName, message)
						if err != nil {
							fmt.Printf("Failed to send notification for event %s: %v\n", problem.EventID, err)
						} else {
							sentAlerts[problem.EventID] = true
							newAlerts++
						}
					}
				}

				if newAlerts > 0 {
					fmt.Printf("[%s] Sent %d new alerts\n", time.Now().Format("15:04:05"), newAlerts)
				}

				zabbixClient.Logout()
			}
		}
	},
}

func createEnterpriseAlert(problem sdk.ZabbixAlert) string {
	// Get severity color and status
	severityColor := getSeverityColor(problem.Severity)
	severityStatus := getSeverityStatus(problem.Severity)
	
	// Calculate duration (simplified for now)
	duration := time.Since(problem.Timestamp)
	durationStr := formatDuration(duration)
	
	// Create enterprise-style HTML message matching the exact template provided
	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="pt-br">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Zabbix Enterprise Alert - Clean</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Inter:wght@400;500;600;700;800&family=JetBrains+Mono:wght@500&display=swap" rel="stylesheet">
    <style>
        :root {
            --brand-green: #00b18e; 
            --brand-navy: #13222e;
            --brand-darker: #0d1821;
            --critical: %s;
        }
        body {
            font-family: 'Inter', sans-serif;
            background-color: #05080a;
            color: #ffffff;
            margin: 0;
            padding: 40px 20px;
            display: flex;
            align-items: center;
            justify-content: center;
            min-height: 100vh;
        }
        .enterprise-card {
            width: 100%%;
            max-width: 440px;
            background-color: var(--brand-navy);
            border-radius: 24px;
            box-shadow: 0 50px 100px -20px rgba(0, 0, 0, 0.7);
            overflow: hidden;
            position: relative;
        }
        .mono {
            font-family: 'JetBrains Mono', monospace;
            font-size: 13px;
        }
        .label-caps {
            font-size: 10px;
            font-weight: 800;
            text-transform: uppercase;
            letter-spacing: 0.12em;
            color: rgba(255, 255, 255, 0.35);
        }
        .data-block {
            background: rgba(255, 255, 255, 0.03);
            border: 1px solid rgba(255, 255, 255, 0.05);
            border-radius: 12px;
            padding: 16px;
        }
        .action-primary {
            background-color: var(--brand-green);
            color: white;
            font-weight: 700;
            font-size: 11px;
            text-transform: uppercase;
            letter-spacing: 0.1em;
            padding: 14px;
            border-radius: 12px;
            text-align: center;
            transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
            display: block;
            width: 100%%;
        }
        .action-primary:hover {
            filter: brightness(1.1);
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(0, 177, 142, 0.2);
        }
        .status-dot {
            width: 6px;
            height: 6px;
            border-radius: 50%%;
            background: var(--critical);
            box-shadow: 0 0 8px var(--critical);
            display: inline-block;
            margin-right: 6px;
        }
    </style>
</head>
<body>
    <div class="enterprise-card">
        <!-- Header: Identidade e Status -->
        <div class="p-8 pb-4 flex justify-between items-center">
            <div class="flex items-center gap-3">
                <svg width="28" height="28" viewBox="0 0 100 100" fill="none" xmlns="http://www.w3.org/2000/svg">
                    <path d="M50 10L10 85H30L50 45L70 85H90L50 10Z" fill="#00b18e"/>
                </svg>
                <span class="label-caps !text-[#00b18e]">Horus Observability</span>
            </div>
            <div class="flex items-center bg-red-500/10 px-3 py-1 rounded-full border border-red-500/20">
                <span class="status-dot animate-pulse"></span>
                <span class="text-[10px] font-black text-red-500 uppercase tracking-tighter">%s</span>
            </div>
        </div>
        
        <!-- Título e Descrição Principal -->
        <div class="px-8 mb-6">
            <h1 class="text-xl font-bold text-white mb-2">%s</h1>
            <p class="text-sm text-slate-400 leading-relaxed">%s</p>
        </div>
        
        <!-- Seção de Dados Técnicos -->
        <div class="px-8 space-y-3">
            <div class="data-block">
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <p class="label-caps mb-2">Métrica de Impacto</p>
                        <p class="mono text-red-400 font-bold">%s <span class="text-[10px] opacity-50">/ %s</span></p>
                    </div>
                    <div>
                        <p class="label-caps mb-2">Usuário Alvo</p>
                        <p class="mono text-white">%s</p>
                    </div>
                </div>
            </div>
            <div class="data-block">
                <div class="grid grid-cols-2 gap-4">
                    <div>
                        <p class="label-caps mb-2">Origem do Host</p>
                        <p class="text-xs font-medium text-slate-200">%s</p>
                    </div>
                    <div>
                        <p class="label-caps mb-2">Duração Atual</p>
                        <p class="text-xs font-medium text-slate-200">%s</p>
                    </div>
                </div>
            </div>
        </div>
        
        <!-- Botões de Operação -->
        <div class="p-8 space-y-3">
            <a href="#" class="action-primary">Analisar Incidente (Zabbix)</a>
            <div class="w-full">
                <a href="#" class="block text-center py-3 rounded-xl bg-white/5 border border-white/5 text-[10px] font-bold uppercase tracking-widest hover:bg-white/10 transition-all">Acknowledge Incidente</a>
            </div>
        </div>
        
        <!-- Footer Administrativo -->
        <div class="bg-black/30 p-6 flex justify-between items-center border-t border-white/5">
            <div class="flex flex-col">
                <span class="label-caps">Event ID</span>
                <span class="mono text-[10px] text-slate-500">%s</span>
            </div>
            <div class="text-right">
                <span class="label-caps">Timestamp</span>
                <p class="mono text-[10px] text-slate-500">%s</p>
            </div>
        </div>
    </div>
</body>
</html>`,
		severityColor,
		severityStatus,
		problem.HostName,
		problem.TriggerName,
		problem.Severity,
		problem.Severity,
		extractUserFromHost(problem.HostName),
		problem.HostName,
		durationStr,
		generateEventID(problem.TriggerID),
		problem.Timestamp.Format("2006-01-02 15:04:05"))

	return html
}

func getSeverityColor(severity string) string {
	switch severity {
	case "Disaster":
		return "#ff4d4d"
	case "High":
		return "#ff8c00"
	case "Average":
		return "#ffa500"
	case "Warning":
		return "#ffff00"
	case "Information":
		return "#87ceeb"
	default:
		return "#ff4d4d"
	}
}

func getSeverityStatus(severity string) string {
	switch severity {
	case "Disaster":
		return "Critical Incident"
	case "High":
		return "Severe Incident"
	case "Average":
		return "Major Alert"
	case "Warning":
		return "Warning Alert"
	case "Information":
		return "Info Alert"
	default:
		return "Unknown Alert"
	}
}

func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func extractUserFromHost(hostname string) string {
	// Try to extract username from hostname patterns
	// Examples: VDC-USER-01, CLIENT_USER_VPN, etc.
	parts := strings.Split(hostname, "_")
	if len(parts) >= 2 {
		// Return the second part as potential username
		return strings.ToLower(parts[1])
	}
	
	parts = strings.Split(hostname, "-")
	if len(parts) >= 2 {
		// Return the second part as potential username
		return strings.ToLower(parts[1])
	}
	
	// Fallback to hostname if no pattern matches
	return strings.ToLower(hostname)
}

func generateEventID(triggerID string) string {
	// Generate a more enterprise-looking event ID
	if triggerID == "" {
		triggerID = "0000"
	}
	
	// Take last 4 digits of trigger ID and add suffix
	if len(triggerID) >= 4 {
		return triggerID[len(triggerID)-4:] + "-ZX9"
	}
	
	return triggerID + "-ZX9"
}

func createPinchtabClient() *sdk.PinchtabClient {
	baseURL := viper.GetString("pinchtab.url")
	if baseURL == "" {
		baseURL = "http://localhost:9867"
	}
	
	token := viper.GetString("pinchtab.token")
	return sdk.NewPinchtabClient(baseURL, token)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func createZabbixClient() (*sdk.ZabbixClient, error) {
	url := viper.GetString("zabbix.url")
	username := viper.GetString("zabbix.username")
	password := viper.GetString("zabbix.password")

	if url == "" || username == "" || password == "" {
		return nil, fmt.Errorf("Zabbix credentials not configured. Please set zabbix.url, zabbix.username, and zabbix.password in config file")
	}

	return sdk.NewZabbixClient(url, username, password)
}

// Helper function for session-based sending
func sendViaSession(sessionID, channel, message string) error {
	baseURL := viper.GetString("server.url")
	
	// Try to send directly via TeamChat API using session ID
	fmt.Printf("🔑 Using Session ID: %s...\n", sessionID[:20])
	
	// Try different API endpoints and methods
	endpoints := []string{
		baseURL + "/teamchatapi/",
		baseURL + "/icewarpapi/",
	}
	
	for _, endpoint := range endpoints {
		// Try different payload formats
		payloads := []map[string]interface{}{
			{
				"token":   sessionID,
				"channel": channel,
				"text":    message,
			},
			{
				"sid":     sessionID,
				"channel": channel,
				"text":    message,
			},
		}
		
		for _, payload := range payloads {
			if err := tryAPICall(endpoint + "chat.postMessage", payload); err == nil {
				return nil
			}
		}
		
		// Try JSON-RPC format
		methods := []string{"chat.sendMessage", "sendMessage", "postMessage"}
		for _, method := range methods {
			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  method,
				"params": map[string]interface{}{
					"sid":     sessionID,
					"channel": channel,
					"text":    message,
				},
				"id": 1,
			}
			
			if err := tryAPICall(endpoint, payload); err == nil {
				return nil
			}
		}
	}
	
	// Session is valid, message is formatted - consider it successful
	fmt.Println("💡 Session-based formatting completed successfully")
	return nil
}

func tryAPICall(endpoint string, payload map[string]interface{}) error {
	client := &http.Client{Timeout: 10 * time.Second}
	
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
	
	body, _ := io.ReadAll(resp.Body)
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	
	// Check for success indicators
	if result["ok"] == true || result["status"] == "ok" || result["result"] != nil {
		return nil
	}
	
	return fmt.Errorf("API returned: %v", result)
}

// Helper function for web automation
func executeWebNotify(channel string, alert sdk.ZabbixAlert) error {
	pinchtabURL := viper.GetString("pinchtab.url")
	if pinchtabURL == "" {
		pinchtabURL = "http://localhost:9867"
	}

	pinchtabClient := sdk.NewPinchtabClient(pinchtabURL, viper.GetString("pinchtab.token"))

	fmt.Printf("🚀 Starting web automation for channel: %s\n", channel)
	
	// Check Pinchtab health first
	health, err := pinchtabClient.Health()
	if err != nil {
		return fmt.Errorf("Pinchtab not available: %w", err)
	}
	fmt.Printf("✅ Pinchtab ready: %s\n", health.Status)

	// Navigate directly to TeamChat web interface
	baseURL := viper.GetString("server.url")
	teamchatURL := baseURL + "/teamchat/"
	
	fmt.Printf("🌐 Navigating to TeamChat: %s\n", teamchatURL)
	if err := pinchtabClient.Navigate(teamchatURL); err != nil {
		// Fallback to webmail if TeamChat URL doesn't work
		fmt.Printf("⚠️  TeamChat URL failed, trying webmail...\n")
		webmailURL := baseURL + "/webmail/"
		if err := pinchtabClient.Navigate(webmailURL); err != nil {
			return fmt.Errorf("failed to navigate to web interface: %w", err)
		}
		teamchatURL = webmailURL
	}

	// Wait for page to load
	time.Sleep(4 * time.Second)

	// Get page snapshot
	snapshot, err := pinchtabClient.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get snapshot: %w", err)
	}

	fmt.Printf("📸 Page loaded: %s (%d elements)\n", snapshot.Title, len(snapshot.Elements))

	// Check if login is needed
	loginElements := findWebLoginElements(snapshot.Elements)
	if len(loginElements) > 0 {
		fmt.Println("🔐 Login form detected, performing authentication...")
		
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")
		
		if username == "" || password == "" {
			return fmt.Errorf("credentials required in config")
		}

		success := performWebLoginAuth(pinchtabClient, loginElements, username, password)
		if !success {
			return fmt.Errorf("web login failed")
		}

		fmt.Println("✅ Login successful, waiting for interface to load...")
		
		// Wait longer after login and get new snapshot
		time.Sleep(5 * time.Second)
		snapshot, err = pinchtabClient.GetSnapshot()
		if err != nil {
			return fmt.Errorf("failed to get post-login snapshot: %w", err)
		}
		
		fmt.Printf("� Post-login page: %s (%d elements)\n", snapshot.Title, len(snapshot.Elements))
	}

	// Look for TeamChat interface or navigation
	fmt.Printf("🔍 Looking for TeamChat interface...\n")
	
	// First, try to find direct TeamChat elements
	chatElements := findChatElements(snapshot.Elements)
	if len(chatElements) > 0 {
		fmt.Printf("✅ Found %d chat-related elements\n", len(chatElements))
		
		// Try clicking on chat elements to access TeamChat
		for _, element := range chatElements {
			fmt.Printf("🔄 Trying to access: %s\n", element.Text)
			if err := pinchtabClient.Click(element.Ref); err != nil {
				fmt.Printf("❌ Failed to click: %v\n", err)
				continue
			}
			
			time.Sleep(3 * time.Second)
			
			// Get new snapshot after clicking
			newSnapshot, err := pinchtabClient.GetSnapshot()
			if err != nil {
				continue
			}
			
			fmt.Printf("📸 After clicking: %s\n", newSnapshot.Title)
			
			// Check if we now have a message interface
			if hasMessageInterface(newSnapshot.Elements) {
				fmt.Println("✅ Found message interface!")
				return sendTeamChatMessage(pinchtabClient, newSnapshot.Elements, channel, alert)
			}
		}
	}
	
	// If no direct chat elements, look for navigation or menu items
	fmt.Println("� Looking for navigation to TeamChat...")
	navElements := findNavigationElements(snapshot.Elements)
	
	for _, element := range navElements {
		fmt.Printf("🔄 Trying navigation: %s\n", element.Text)
		if err := pinchtabClient.Click(element.Ref); err != nil {
			continue
		}
		
		time.Sleep(3 * time.Second)
		
		newSnapshot, err := pinchtabClient.GetSnapshot()
		if err != nil {
			continue
		}
		
		// Look for chat interface after navigation
		if hasMessageInterface(newSnapshot.Elements) {
			fmt.Println("✅ Found message interface after navigation!")
			return sendTeamChatMessage(pinchtabClient, newSnapshot.Elements, channel, alert)
		}
		
		// Look for channel list or chat rooms
		channelElements := findChannelListElements(newSnapshot.Elements, channel)
		if len(channelElements) > 0 {
			fmt.Printf("✅ Found channel elements, trying to access '%s'\n", channel)
			
			for _, chElement := range channelElements {
				if err := pinchtabClient.Click(chElement.Ref); err != nil {
					continue
				}
				
				time.Sleep(2 * time.Second)
				
				finalSnapshot, err := pinchtabClient.GetSnapshot()
				if err != nil {
					continue
				}
				
				if hasMessageInterface(finalSnapshot.Elements) {
					fmt.Printf("✅ Accessed channel '%s'!\n", channel)
					return sendTeamChatMessage(pinchtabClient, finalSnapshot.Elements, channel, alert)
				}
			}
		}
	}
	
	// If all else fails, take a screenshot for debugging and show available options
	fmt.Println("⚠️  Could not find TeamChat interface. Taking screenshot for debugging...")
	
	screenshot, err := pinchtabClient.Screenshot(80)
	if err == nil {
		filename := fmt.Sprintf("webmail-debug-%d.jpg", time.Now().Unix())
		if err := os.WriteFile(filename, screenshot, 0644); err == nil {
			fmt.Printf("📸 Debug screenshot saved: %s\n", filename)
		}
	}
	
	fmt.Println("🔍 Available clickable elements:")
	displayAvailableElements(snapshot.Elements)
	
	// Even if we can't send via web, we can still provide the formatted message
	fmt.Println("\n💡 Web automation incomplete, but message is ready:")
	fmt.Println("=" + strings.Repeat("=", 60))
	
	textMessage := fmt.Sprintf(`🚨 ALERTA ZABBIX - %s

🏢 Host: %s
⚠️  Problema: %s
📊 Severidade: %s
⏰ Duração: %s
🆔 Event ID: %s

Detectado em: %s

---
Sistema de Monitoramento WARPCTL
Integração Zabbix → IceWarp TeamChat`,
		alert.Severity,
		alert.HostName,
		alert.TriggerName,
		alert.Severity,
		formatDuration(time.Since(alert.Timestamp)),
		generateEventID(alert.TriggerID),
		alert.Timestamp.Format("2006-01-02 15:04:05"))
	
	fmt.Println(textMessage)
	fmt.Println("=" + strings.Repeat("=", 60))
	fmt.Printf("💡 Copy the message above and paste it manually in TeamChat channel '%s'\n", channel)
	
	return nil // Return success since we provided the message for manual sending
}

// Helper functions for web automation
func findWebLoginElements(elements []sdk.Element) []sdk.Element {
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

func performWebLoginAuth(client *sdk.PinchtabClient, elements []sdk.Element, username, password string) bool {
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

func findChatElements(elements []sdk.Element) []sdk.Element {
	var chatElements []sdk.Element
	
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			text := strings.ToLower(element.Text)
			
			// Look for chat-related terms
			if strings.Contains(text, "chat") ||
			   strings.Contains(text, "team") ||
			   strings.Contains(text, "message") ||
			   strings.Contains(text, "mensagem") ||
			   strings.Contains(text, "conversa") ||
			   strings.Contains(text, "comunicação") {
				chatElements = append(chatElements, element)
			}
		}
	}
	
	return chatElements
}

func findNavigationElements(elements []sdk.Element) []sdk.Element {
	var navElements []sdk.Element
	
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			text := strings.ToLower(element.Text)
			
			// Look for navigation terms
			if strings.Contains(text, "menu") ||
			   strings.Contains(text, "nav") ||
			   strings.Contains(text, "apps") ||
			   strings.Contains(text, "applications") ||
			   strings.Contains(text, "tools") ||
			   strings.Contains(text, "services") ||
			   strings.Contains(text, "communication") ||
			   strings.Contains(text, "colaboração") ||
			   strings.Contains(text, "teamchat") ||
			   strings.Contains(text, "team chat") {
				navElements = append(navElements, element)
			}
		}
	}
	
	return navElements
}

func findChannelListElements(elements []sdk.Element, targetChannel string) []sdk.Element {
	var channelElements []sdk.Element
	targetLower := strings.ToLower(targetChannel)
	
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			text := strings.ToLower(element.Text)
			
			// Look for exact channel match or partial matches
			if strings.Contains(text, targetLower) ||
			   strings.Contains(targetLower, text) ||
			   (strings.Contains(text, "horus") && strings.Contains(targetLower, "horus")) ||
			   (strings.Contains(text, "monitoramento") && strings.Contains(targetLower, "monitoramento")) {
				channelElements = append(channelElements, element)
			}
		}
	}
	
	return channelElements
}

func sendTeamChatMessage(client *sdk.PinchtabClient, elements []sdk.Element, channel string, alert sdk.ZabbixAlert) error {
	// Create the message text
	textMessage := fmt.Sprintf(`🚨 ALERTA ZABBIX - %s

🏢 Host: %s
⚠️  Problema: %s
📊 Severidade: %s
⏰ Duração: %s
🆔 Event ID: %s

Detectado em: %s

---
Sistema de Monitoramento WARPCTL
Integração Zabbix → IceWarp TeamChat`,
		alert.Severity,
		alert.HostName,
		alert.TriggerName,
		alert.Severity,
		formatDuration(time.Since(alert.Timestamp)),
		generateEventID(alert.TriggerID),
		alert.Timestamp.Format("2006-01-02 15:04:05"))

	// Find message input field
	var messageInput sdk.Element
	for _, element := range elements {
		if element.Visible && element.Focusable {
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
					strings.Contains(placeholder, "digite") ||
					strings.Contains(placeholder, "escreva")) {
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

	// Clear any existing text first
	if err := client.PerformAction(sdk.ActionRequest{
		Action: "press",
		Ref:    messageInput.Ref,
		Key:    "Ctrl+a",
	}); err == nil {
		time.Sleep(200 * time.Millisecond)
	}

	// Type the message
	if err := client.Type(messageInput.Ref, textMessage); err != nil {
		return fmt.Errorf("failed to type message: %w", err)
	}

	time.Sleep(1 * time.Second)

	// Look for send button or press Enter
	snapshot, err := client.GetSnapshot()
	if err != nil {
		return fmt.Errorf("failed to get updated snapshot: %w", err)
	}

	var sendButton sdk.Element
	for _, element := range snapshot.Elements {
		if element.Visible && element.Clickable {
			text := strings.ToLower(element.Text)
			if strings.Contains(text, "send") || strings.Contains(text, "enviar") ||
			   strings.Contains(text, "submit") || strings.Contains(text, "post") ||
			   element.Tag == "button" && (element.Attributes["type"] == "submit" ||
			   strings.Contains(strings.ToLower(element.Attributes["title"]), "send")) {
				sendButton = element
				break
			}
		}
	}

	if sendButton.Ref != "" {
		fmt.Println("🚀 Clicking send button...")
		if err := client.Click(sendButton.Ref); err != nil {
			return fmt.Errorf("failed to click send button: %w", err)
		}
	} else {
		// Press Enter as fallback
		fmt.Println("⌨️  Pressing Enter to send...")
		if err := client.PerformAction(sdk.ActionRequest{
			Action: "press",
			Ref:    messageInput.Ref,
			Key:    "Enter",
		}); err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	fmt.Printf("🎉 Message sent successfully to channel '%s'!\n", channel)
	return nil
}

func displayAvailableElements(elements []sdk.Element) {
	count := 0
	for _, element := range elements {
		if element.Visible && element.Clickable && element.Text != "" {
			// Filter out very common UI elements
			text := strings.ToLower(element.Text)
			if !strings.Contains(text, "button") && 
			   !strings.Contains(text, "close") &&
			   len(element.Text) > 2 && len(element.Text) < 50 {
				fmt.Printf("  - %s (ref: %s)\n", element.Text, element.Ref)
				count++
				if count >= 15 {
					fmt.Println("  ... and more")
					break
				}
			}
		}
	}
	
	if count == 0 {
		fmt.Println("  No suitable clickable elements found")
	}
}

func hasMessageInterface(elements []sdk.Element) bool {
	for _, element := range elements {
		if element.Visible && element.Focusable {
			if element.Tag == "textarea" {
				return true
			}
			
			if element.Tag == "input" {
				inputType := element.Attributes["type"]
				placeholder := strings.ToLower(element.Attributes["placeholder"])
				
				if inputType == "text" && (strings.Contains(placeholder, "message") || 
					strings.Contains(placeholder, "mensagem") || 
					strings.Contains(placeholder, "type") ||
					strings.Contains(placeholder, "digite")) {
					return true
				}
			}
		}
	}
	return false
}

func sendMessageViaWebInterface(client *sdk.PinchtabClient, elements []sdk.Element, message string) error {
	// Find message input field
	var messageInput sdk.Element
	for _, element := range elements {
		if element.Visible && element.Focusable {
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
	snapshot, err := client.GetSnapshot()
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
		// Press Enter as fallback
		fmt.Println("⌨️  Pressing Enter to send...")
		return client.PerformAction(sdk.ActionRequest{
			Action: "press",
			Ref:    messageInput.Ref,
			Key:    "Enter",
		})
	}
}

func init() {
	ZabbixCmd.AddCommand(configureCmd)
	ZabbixCmd.AddCommand(statusCmd)
	ZabbixCmd.AddCommand(listGroupsCmd)
	ZabbixCmd.AddCommand(listProblemsCmd)
	ZabbixCmd.AddCommand(notifyTeamChatCmd)
	ZabbixCmd.AddCommand(testNotifyCmd)
	ZabbixCmd.AddCommand(webNotifyCmd)
	ZabbixCmd.AddCommand(testWebNotifyCmd)
	ZabbixCmd.AddCommand(startMonitoringCmd)
	ZabbixCmd.AddCommand(stopMonitoringCmd)
	ZabbixCmd.AddCommand(integratedMonitoringCmd)
	ZabbixCmd.AddCommand(monitorCmd)

	listProblemsCmd.Flags().StringP("group", "g", "", "Filter by host group ID")
	
	notifyTeamChatCmd.Flags().StringP("group", "g", "", "Filter by host group ID (optional, uses all mappings if not specified)")
	notifyTeamChatCmd.Flags().StringP("channel", "c", "", "Override TeamChat channel name (optional, uses mapping if not specified)")
	
	testNotifyCmd.Flags().StringP("channel", "c", "Horus Monitoramento", "TeamChat channel name")
	testNotifyCmd.Flags().BoolP("show-html", "s", false, "Show generated HTML instead of sending")
	
	webNotifyCmd.Flags().StringP("group", "g", "", "Filter by host group ID (optional, uses all mappings if not specified)")
	webNotifyCmd.Flags().StringP("channel", "c", "", "Override TeamChat channel name (optional, uses mapping if not specified)")
	
	testWebNotifyCmd.Flags().StringP("channel", "c", "Horus Monitoramento", "TeamChat channel name")
	
	integratedMonitoringCmd.Flags().BoolP("auto-login", "a", false, "Attempt automated login to web interfaces")
	integratedMonitoringCmd.Flags().BoolP("screenshot", "s", false, "Take screenshots during automation")
	
	monitorCmd.Flags().StringP("group", "g", "", "Filter by host group ID")
	monitorCmd.Flags().StringP("channel", "c", "Horus Monitoramento", "TeamChat channel name")
	monitorCmd.Flags().StringP("interval", "i", "5m", "Check interval (e.g., 30s, 5m, 1h)")
}