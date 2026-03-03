package filters

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FiltersCmd = &cobra.Command{
		Use:   "filters",
		Short: "Filter/Rule operations",
		Long:  `Commands for IceWarp Filters/Rules API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Filters API",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid, err := client.Authenticate(username, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Printf("Login successful!\n")
		fmt.Printf("SID: %s\n", sid)
		fmt.Printf("Run: warpctl config set filters.sid %s\n", sid)
		return nil
	},
}

var listRulesCmd = &cobra.Command{
	Use:   "list",
	Short: "List filters/rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("filters.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl filters login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")
		ruleType, _ := cmd.Flags().GetString("type")

		rules, err := client.ListRules(account, ruleType)
		if err != nil {
			return fmt.Errorf("failed to list rules: %w", err)
		}

		fmt.Println("Filters/Rules:")
		for _, r := range rules {
			fmt.Printf("  %s (ID: %s) - Enabled: %v, Priority: %d\n", r.Name, r.ID, r.Enabled, r.Priority)
		}
		return nil
	},
}

var getRuleCmd = &cobra.Command{
	Use:   "info [rule-id]",
	Short: "Get rule information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("filters.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl filters login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")

		rule, err := client.GetRule(account, args[0])
		if err != nil {
			return fmt.Errorf("failed to get rule: %w", err)
		}

		fmt.Printf("Rule: %s (ID: %s)\n", rule.Name, rule.ID)
		fmt.Printf("Enabled: %v\n", rule.Enabled)
		fmt.Printf("Priority: %d\n", rule.Priority)
		fmt.Printf("Folder: %s\n", rule.FolderID)
		return nil
	},
}

var createRuleCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new filter/rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("filters.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl filters login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")
		folderID, _ := cmd.Flags().GetString("folder")
		conditionType, _ := cmd.Flags().GetString("condition-type")
		conditionField, _ := cmd.Flags().GetString("condition-field")
		conditionOperator, _ := cmd.Flags().GetString("condition-operator")
		conditionValue, _ := cmd.Flags().GetString("condition-value")
		actionType, _ := cmd.Flags().GetString("action-type")
		actionFolder, _ := cmd.Flags().GetString("action-folder")
		actionText, _ := cmd.Flags().GetString("action-text")
		enabled, _ := cmd.Flags().GetBool("enabled")
		priority, _ := cmd.Flags().GetInt("priority")

		err := client.CreateRule(account, args[0], folderID, conditionType, conditionField, conditionOperator, conditionValue, actionType, actionFolder, actionText, enabled, priority)
		if err != nil {
			return fmt.Errorf("failed to create rule: %w", err)
		}

		fmt.Printf("Rule '%s' created successfully\n", args[0])
		return nil
	},
}

var deleteRuleCmd = &cobra.Command{
	Use:   "delete [rule-id]",
	Short: "Delete a filter/rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("filters.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl filters login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")

		if err := client.DeleteRule(account, args[0]); err != nil {
			return fmt.Errorf("failed to delete rule: %w", err)
		}

		fmt.Printf("Rule deleted successfully\n")
		return nil
	},
}

var setRuleStateCmd = &cobra.Command{
	Use:   "set-state [rule-id]",
	Short: "Enable or disable a filter/rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("filters.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl filters login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")
		enabled, _ := cmd.Flags().GetBool("enabled")

		if err := client.SetRuleState(account, args[0], enabled); err != nil {
			return fmt.Errorf("failed to set rule state: %w", err)
		}

		fmt.Printf("Rule state updated successfully\n")
		return nil
	},
}

var moveRuleCmd = &cobra.Command{
	Use:   "move [rule-id]",
	Short: "Move a filter/rule up or down in priority",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFiltersClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("filters.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl filters login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")
		direction, _ := cmd.Flags().GetString("direction")

		if err := client.MoveRule(account, args[0], direction); err != nil {
			return fmt.Errorf("failed to move rule: %w", err)
		}

		fmt.Printf("Rule moved successfully\n")
		return nil
	},
}

func init() {
	FiltersCmd.AddCommand(loginCmd)
	FiltersCmd.AddCommand(listRulesCmd)
	FiltersCmd.AddCommand(getRuleCmd)
	FiltersCmd.AddCommand(createRuleCmd)
	FiltersCmd.AddCommand(deleteRuleCmd)
	FiltersCmd.AddCommand(setRuleStateCmd)
	FiltersCmd.AddCommand(moveRuleCmd)

	listRulesCmd.Flags().StringP("account", "a", "", "Account email")
	listRulesCmd.Flags().StringP("type", "t", "mail", "Rule type (mail)")

	getRuleCmd.Flags().StringP("account", "a", "", "Account email")

	createRuleCmd.Flags().StringP("account", "a", "", "Account email")
	createRuleCmd.Flags().StringP("folder", "f", "", "Folder ID")
	createRuleCmd.Flags().StringP("condition-type", "c", "header", "Condition type")
	createRuleCmd.Flags().StringP("condition-field", "i", "subject", "Condition field")
	createRuleCmd.Flags().StringP("condition-operator", "o", "contains", "Condition operator")
	createRuleCmd.Flags().StringP("condition-value", "v", "", "Condition value")
	createRuleCmd.Flags().StringP("action-type", "y", "move", "Action type")
	createRuleCmd.Flags().StringP("action-folder", "d", "", "Action folder")
	createRuleCmd.Flags().StringP("action-text", "x", "", "Action text")
	createRuleCmd.Flags().BoolP("enabled", "e", true, "Enable rule")
	createRuleCmd.Flags().IntP("priority", "p", 0, "Priority")

	deleteRuleCmd.Flags().StringP("account", "a", "", "Account email")

	setRuleStateCmd.Flags().StringP("account", "a", "", "Account email")
	setRuleStateCmd.Flags().BoolP("enabled", "e", true, "Enable/disable rule")

	moveRuleCmd.Flags().StringP("account", "a", "", "Account email")
	moveRuleCmd.Flags().StringP("direction", "d", "up", "Direction (up/down)")
}
