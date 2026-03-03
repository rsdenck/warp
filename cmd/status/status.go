package status

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/config"
	"github.com/icewarp/warpctl/internal/imap"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status [mailbox]",
	Short: "Show mailbox status",
	Long:  `Display message count and other information for a mailbox.`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mailbox := "INBOX"
		if len(args) > 0 {
			mailbox = args[0]
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		if cfg.Auth.Username == "" || cfg.Auth.Password == "" {
			return fmt.Errorf("username and password are required. Set IW_USERNAME and IW_PASSWORD environment variables")
		}

		client, err := imap.NewClient(cfg)
		if err != nil {
			return err
		}
		defer client.Logout()

		if err := client.Login(); err != nil {
			return err
		}

		mbox, err := client.SelectMailbox(mailbox)
		if err != nil {
			return err
		}

		fmt.Printf("Mailbox: %s\n", mailbox)
		fmt.Printf("Messages: %d\n", mbox.Messages)
		fmt.Printf("Unseen: %d\n", mbox.Unseen)
		fmt.Printf("Recent: %d\n", mbox.Recent)

		return nil
	},
}
