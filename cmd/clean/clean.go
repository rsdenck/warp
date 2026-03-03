package clean

import (
	"fmt"
	"os"

	"github.com/icewarp/warpctl/internal/config"
	"github.com/icewarp/warpctl/internal/imap"
	"github.com/spf13/cobra"
)

var (
	mailbox string
	dryRun  bool
	confirm bool
)

var CleanCmd = &cobra.Command{
	Use:   "clean [mailbox]",
	Short: "Clean messages from a mailbox",
	Long: `Removes all messages from the specified mailbox (default: INBOX).
This operation is destructive and cannot be undone.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			mailbox = args[0]
		}
		if mailbox == "" {
			mailbox = "INBOX"
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Auth.Username == "" || cfg.Auth.Password == "" {
			return fmt.Errorf("username and password are required. Set IW_USERNAME and IW_PASSWORD environment variables or use 'icwli config set'")
		}

		client, err := imap.NewClient(cfg)
		if err != nil {
			return err
		}
		defer client.Logout()

		if err := client.Login(); err != nil {
			return err
		}

		count, err := client.GetMessageCount(mailbox)
		if err != nil {
			return err
		}

		if count == 0 {
			fmt.Printf("Mailbox %s is already empty\n", mailbox)
			return nil
		}

		fmt.Printf("Found %d messages in %s\n", count, mailbox)

		if dryRun {
			fmt.Println("Dry run mode - no messages will be deleted")
			return nil
		}

		if !confirm {
			fmt.Print("Are you sure you want to delete all messages? (yes/no): ")
			var response string
			fmt.Scanln(&response)
			if response != "yes" && response != "y" {
				fmt.Println("Operation cancelled")
				os.Exit(0)
			}
		}

		ids, err := client.SearchMessages(mailbox)
		if err != nil {
			return err
		}

		deleted, err := client.DeleteMessages(mailbox, ids)
		if err != nil {
			return err
		}

		fmt.Printf("Successfully deleted %d messages from %s\n", deleted, mailbox)

		return nil
	},
}

func init() {
	CleanCmd.Flags().StringVarP(&mailbox, "mailbox", "m", "INBOX", "mailbox to clean")
	CleanCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "show what would be deleted without actually deleting")
	CleanCmd.Flags().BoolVarP(&confirm, "yes", "y", false, "skip confirmation prompt")
}
