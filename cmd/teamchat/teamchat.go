package teamchat

import (
	"fmt"

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

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to TeamChat",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		resp, err := client.LoginPlain(username, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Printf("Login successful!\n")
		fmt.Printf("Token: %s\n", resp.Token)
		fmt.Printf("User: %s\n", resp.User)
		return nil
	},
}

var authTestCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Check authentication status",
	RunE: func(cmd *cobra.Command, args []string) error {
		token := viper.GetString("teamchat.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'icwli teamchat login' first")
		}

		client := sdk.NewTeamChatClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})
		client.SetToken(token)

		resp, err := client.AuthTest()
		if err != nil {
			return fmt.Errorf("auth test failed: %w", err)
		}

		fmt.Printf("User ID: %s\n", resp.UserID)
		fmt.Printf("Username: %s\n", resp.Username)
		fmt.Printf("Email: %s\n", resp.Email)
		fmt.Printf("Valid: %v\n", resp.Valid)
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
	TeamChatCmd.AddCommand(authTestCmd)
	TeamChatCmd.AddCommand(logoutCmd)
	TeamChatCmd.AddCommand(presenceCmd)
	TeamChatCmd.AddCommand(conversationsCmd)
}
