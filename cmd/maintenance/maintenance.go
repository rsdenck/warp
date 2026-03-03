package maintenance

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	MaintenanceCmd = &cobra.Command{
		Use:   "maintenance",
		Short: "Maintenance API operations",
		Long:  `Commands for IceWarp Maintenance/Admin API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with Maintenance API",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid, err := client.Authenticate(username, password)
		if err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}

		fmt.Printf("Authentication successful!\n")
		fmt.Printf("SID: %s\n", sid)
		return nil
	},
}

var serverInfoCmd = &cobra.Command{
	Use:   "info",
	Short: "Get server information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		info, err := client.GetServerInfo()
		if err != nil {
			return fmt.Errorf("failed to get server info: %w", err)
		}

		fmt.Println("Server Information:")
		for k, v := range info {
			fmt.Printf("  %s: %v\n", k, v)
		}
		return nil
	},
}

var systemInfoCmd = &cobra.Command{
	Use:   "system",
	Short: "Get system information",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		info, err := client.GetSystemInfo()
		if err != nil {
			return fmt.Errorf("failed to get system info: %w", err)
		}

		fmt.Println("System Information:")
		for k, v := range info {
			fmt.Printf("  %s: %v\n", k, v)
		}
		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get server statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		stats, err := client.GetStatistics()
		if err != nil {
			return fmt.Errorf("failed to get statistics: %w", err)
		}

		fmt.Println("Server Statistics:")
		for k, v := range stats {
			fmt.Printf("  %s: %v\n", k, v)
		}
		return nil
	},
}

var domainsCmd = &cobra.Command{
	Use:   "domains",
	Short: "List domains",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		domains, err := client.GetDomainList()
		if err != nil {
			return fmt.Errorf("failed to get domain list: %w", err)
		}

		fmt.Println("Domains:")
		for _, d := range domains {
			fmt.Printf("  %s\n", d["name"])
		}
		return nil
	},
}

var createDomainCmd = &cobra.Command{
	Use:   "create-domain [domain-name]",
	Short: "Create a new domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		if err := client.CreateDomain(args[0]); err != nil {
			return fmt.Errorf("failed to create domain: %w", err)
		}

		fmt.Printf("Domain '%s' created successfully\n", args[0])
		return nil
	},
}

var deleteDomainCmd = &cobra.Command{
	Use:   "delete-domain [domain-name]",
	Short: "Delete a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		if err := client.DeleteDomain(args[0]); err != nil {
			return fmt.Errorf("failed to delete domain: %w", err)
		}

		fmt.Printf("Domain '%s' deleted successfully\n", args[0])
		return nil
	},
}

var usersCmd = &cobra.Command{
	Use:   "users [domain]",
	Short: "List users in a domain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		users, err := client.GetUserList(args[0])
		if err != nil {
			return fmt.Errorf("failed to get user list: %w", err)
		}

		fmt.Printf("Users in %s:\n", args[0])
		for _, u := range users {
			fmt.Printf("  %s\n", u["name"])
		}
		return nil
	},
}

var createUserCmd = &cobra.Command{
	Use:   "create-user [domain] [username] [password]",
	Short: "Create a new user",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		if err := client.CreateUser(args[0], args[1], args[2]); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		fmt.Printf("User '%s@%s' created successfully\n", args[1], args[0])
		return nil
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "delete-user [domain] [username]",
	Short: "Delete a user",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		if err := client.DeleteUser(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to delete user: %w", err)
		}

		fmt.Printf("User '%s@%s' deleted successfully\n", args[1], args[0])
		return nil
	},
}

var userInfoCmd = &cobra.Command{
	Use:   "user-info [domain] [username]",
	Short: "Get user information",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		info, err := client.GetUserInfo(args[0], args[1])
		if err != nil {
			return fmt.Errorf("failed to get user info: %w", err)
		}

		fmt.Printf("User Information for %s@%s:\n", args[1], args[0])
		for k, v := range info {
			fmt.Printf("  %s: %v\n", k, v)
		}
		return nil
	},
}

var setQuotaCmd = &cobra.Command{
	Use:   "set-quota [domain] [username] [quota-mb]",
	Short: "Set user quota",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			return fmt.Errorf("not authenticated. Use 'icwli maintenance login' first")
		}

		var quota int
		fmt.Sscanf(args[2], "%d", &quota)

		if err := client.SetUserQuota(args[0], args[1], quota); err != nil {
			return fmt.Errorf("failed to set quota: %w", err)
		}

		fmt.Printf("Quota set to %d MB for %s@%s\n", quota, args[1], args[0])
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Maintenance API",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMaintenanceClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if !client.IsAuthenticated() {
			fmt.Println("Not logged in")
			return nil
		}

		if err := client.Logout(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}

		fmt.Println("Logged out successfully")
		return nil
	},
}

func init() {
	MaintenanceCmd.AddCommand(loginCmd)
	MaintenanceCmd.AddCommand(serverInfoCmd)
	MaintenanceCmd.AddCommand(systemInfoCmd)
	MaintenanceCmd.AddCommand(statsCmd)
	MaintenanceCmd.AddCommand(domainsCmd)
	MaintenanceCmd.AddCommand(createDomainCmd)
	MaintenanceCmd.AddCommand(deleteDomainCmd)
	MaintenanceCmd.AddCommand(usersCmd)
	MaintenanceCmd.AddCommand(createUserCmd)
	MaintenanceCmd.AddCommand(deleteUserCmd)
	MaintenanceCmd.AddCommand(userInfoCmd)
	MaintenanceCmd.AddCommand(setQuotaCmd)
	MaintenanceCmd.AddCommand(logoutCmd)
}
