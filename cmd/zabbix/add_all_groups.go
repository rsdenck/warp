package zabbix

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var addAllGroupsCmd = &cobra.Command{
	Use:   "add-all-groups",
	Short: "Add all VDC groups to monitoring (FOR TESTING ONLY)",
	Long:  `Adds all VDC groups found in Zabbix to the monitoring configuration. This is for testing purposes only.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		channelName, _ := cmd.Flags().GetString("channel")
		if channelName == "" {
			channelName = "Horus Monitoramento"
		}

		fmt.Println("🔧 Adding ALL VDC Groups for Testing")
		fmt.Println("=" + strings.Repeat("=", 50))
		fmt.Printf("📤 Default channel: %s\n", channelName)

		// Connect to Zabbix
		fmt.Println("📡 Connecting to Zabbix...")
		zabbixClient, err := createZabbixClient()
		if err != nil {
			return fmt.Errorf("failed to connect to Zabbix: %w", err)
		}
		defer zabbixClient.Logout()

		// Get all host groups
		fmt.Println("🔍 Fetching all host groups...")
		groups, err := zabbixClient.GetHostGroups()
		if err != nil {
			return fmt.Errorf("failed to get host groups: %w", err)
		}

		fmt.Printf("✅ Found %d total groups\n", len(groups))

		// Filter VDC groups
		var vdcGroups []struct {
			GroupID string
			Name    string
		}

		for _, group := range groups {
			if strings.HasPrefix(group.Name, "VDC_") {
				vdcGroups = append(vdcGroups, struct {
					GroupID string
					Name    string
				}{
					GroupID: group.GroupID,
					Name:    group.Name,
				})
			}
		}

		fmt.Printf("🎯 Found %d VDC groups\n", len(vdcGroups))

		if len(vdcGroups) == 0 {
			return fmt.Errorf("no VDC groups found")
		}

		// Get current mappings
		currentMappings := viper.GetStringMapString("zabbix.group_mappings")
		if currentMappings == nil {
			currentMappings = make(map[string]string)
		}

		// Add all VDC groups
		newMappings := make(map[string]string)
		for k, v := range currentMappings {
			newMappings[k] = v
		}

		addedCount := 0
		for _, group := range vdcGroups {
			if _, exists := newMappings[group.GroupID]; !exists {
				newMappings[group.GroupID] = channelName
				addedCount++
				fmt.Printf("➕ Added: %s (ID: %s) → %s\n", group.Name, group.GroupID, channelName)
			} else {
				fmt.Printf("⏭️  Skipped: %s (already configured)\n", group.Name)
			}
		}

		if addedCount == 0 {
			fmt.Println("ℹ️  All VDC groups are already configured")
			return nil
		}

		// Save configuration
		fmt.Printf("\n💾 Saving %d new group mappings...\n", addedCount)
		viper.Set("zabbix.group_mappings", newMappings)
		
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("✅ Successfully added %d VDC groups to monitoring!\n", addedCount)
		fmt.Printf("📊 Total groups now configured: %d\n", len(newMappings))
		
		fmt.Println("\n⚠️  WARNING: This is for TESTING only!")
		fmt.Println("💡 In production, configure only specific groups you want to monitor")
		
		fmt.Println("\n🚀 You can now start monitoring with:")
		fmt.Println("   ./warpctl.exe zabbix start")

		return nil
	},
}

func init() {
	ZabbixCmd.AddCommand(addAllGroupsCmd)
	addAllGroupsCmd.Flags().StringP("channel", "c", "Horus Monitoramento", "Default TeamChat channel for all groups")
}