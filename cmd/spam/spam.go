package spam

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/output"
	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	SpamCmd = &cobra.Command{
		Use:   "spam",
		Short: "Spam/Quarantine operations",
		Long:  `Commands for IceWarp Spam/Quarantine API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Spam API",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid, err := client.Authenticate(username, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Printf("Login successful!\n")
		fmt.Printf("SID: %s\n", sid)
		fmt.Printf("Run: warpctl config set spam.sid %s\n", sid)
		return nil
	},
}

var listQuarantineCmd = &cobra.Command{
	Use:   "list",
	Short: "List spam quarantine items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		if limit == 0 {
			limit = 50
		}

		items, err := client.ListQuarantine(limit, offset)
		if err != nil {
			return fmt.Errorf("failed to list quarantine: %w", err)
		}

		t := output.NewTable("SPAM QUARANTINE")
		t.AppendHeader(table.Row{"From", "To", "Subject", "Score"})
		
		for _, i := range items {
			t.AppendRow(table.Row{i.From, i.To, i.Subject, i.Score})
		}
		
		t.Render()
		return nil
	},
}

var getItemCmd = &cobra.Command{
	Use:   "info [item-id]",
	Short: "Get spam quarantine item details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		item, err := client.GetQuarantineItem(args[0])
		if err != nil {
			return fmt.Errorf("failed to get item: %w", err)
		}

		t := output.NewTable("SPAM ITEM DETAILS")
		t.AppendRow(table.Row{"From", item.From})
		t.AppendRow(table.Row{"To", item.To})
		t.AppendRow(table.Row{"Subject", item.Subject})
		t.AppendRow(table.Row{"Date", item.Date})
		t.AppendRow(table.Row{"Size", fmt.Sprintf("%d bytes", item.Size)})
		t.AppendRow(table.Row{"Score", item.Score})
		t.AppendRow(table.Row{"Reason", item.Reason})
		t.AppendRow(table.Row{"Action", item.Action})
		
		t.Render()
		return nil
	},
}

var getBodyCmd = &cobra.Command{
	Use:   "body [item-id]",
	Short: "Get spam message body",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		body, err := client.GetSpamBody(args[0])
		if err != nil {
			return fmt.Errorf("failed to get body: %w", err)
		}

		fmt.Println(body)
		return nil
	},
}

var deliverCmd = &cobra.Command{
	Use:   "deliver [item-id]",
	Short: "Deliver a spam quarantine item to inbox",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		if err := client.DeliverQuarantineItem(args[0]); err != nil {
			return fmt.Errorf("failed to deliver: %w", err)
		}

		fmt.Printf("Item delivered to inbox\n")
		return nil
	},
}

var deleteItemCmd = &cobra.Command{
	Use:   "delete [item-id]",
	Short: "Delete a spam quarantine item",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		if err := client.DeleteQuarantineItem(args[0]); err != nil {
			return fmt.Errorf("failed to delete: %w", err)
		}

		fmt.Printf("Item deleted\n")
		return nil
	},
}

var deleteAllCmd = &cobra.Command{
	Use:   "delete-all",
	Short: "Delete all spam quarantine items",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		if err := client.DeleteAllQuarantine(); err != nil {
			return fmt.Errorf("failed to delete all: %w", err)
		}

		fmt.Printf("All spam items deleted\n")
		return nil
	},
}

var whitelistCmd = &cobra.Command{
	Use:   "whitelist [email]",
	Short: "Add sender to whitelist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		if err := client.WhitelistSender(args[0]); err != nil {
			return fmt.Errorf("failed to whitelist: %w", err)
		}

		fmt.Printf("Sender %s added to whitelist\n", args[0])
		return nil
	},
}

var blacklistCmd = &cobra.Command{
	Use:   "blacklist [email]",
	Short: "Add sender to blacklist",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		if err := client.BlacklistSender(args[0]); err != nil {
			return fmt.Errorf("failed to blacklist: %w", err)
		}

		fmt.Printf("Sender %s added to blacklist\n", args[0])
		return nil
	},
}

var cleanSpamCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean all spam messages from quarantine",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewSpamClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("spam.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl spam login' first")
		}
		client.SetSID(sid)

		if err := client.DeleteAllQuarantine(); err != nil {
			return fmt.Errorf("failed to clean spam: %w", err)
		}

		fmt.Printf("Spam quarantine cleaned\n")
		return nil
	},
}

func init() {
	SpamCmd.AddCommand(loginCmd)
	SpamCmd.AddCommand(listQuarantineCmd)
	SpamCmd.AddCommand(getItemCmd)
	SpamCmd.AddCommand(getBodyCmd)
	SpamCmd.AddCommand(deliverCmd)
	SpamCmd.AddCommand(deleteItemCmd)
	SpamCmd.AddCommand(deleteAllCmd)
	SpamCmd.AddCommand(whitelistCmd)
	SpamCmd.AddCommand(blacklistCmd)
	SpamCmd.AddCommand(cleanSpamCmd)

	listQuarantineCmd.Flags().IntP("limit", "l", 50, "Limit number of items")
	listQuarantineCmd.Flags().IntP("offset", "o", 0, "Offset for pagination")
}
