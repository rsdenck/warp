package mail

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	MailCmd = &cobra.Command{
		Use:   "mail",
		Short: "Mail API operations",
		Long:  `Commands for IceWarp Mail API`,
	}
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get Mail server version",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		resp, err := client.ShowVersion()
		if err != nil {
			return fmt.Errorf("failed to get version: %w", err)
		}

		fmt.Printf("Version: %s\n", resp.Version)
		fmt.Printf("Build: %s\n", resp.Build)
		fmt.Printf("Start Time: %s\n", resp.StartTime)
		return nil
	},
}

var listFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List mail folders",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		folders, err := client.ListFolders()
		if err != nil {
			return fmt.Errorf("failed to list folders: %w", err)
		}

		fmt.Println("Folders:")
		for _, f := range folders {
			fmt.Printf("  %s (ID: %s) - Messages: %d, Unread: %d\n", f.Name, f.ID, f.Messages, f.Unread)
		}
		return nil
	},
}

var createFolderCmd = &cobra.Command{
	Use:   "create-folder [name]",
	Short: "Create a new folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		parentID, _ := cmd.Flags().GetString("parent")

		if err := client.CreateFolder(args[0], parentID); err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}

		fmt.Printf("Folder '%s' created successfully\n", args[0])
		return nil
	},
}

var deleteFolderCmd = &cobra.Command{
	Use:   "delete-folder [folder-id]",
	Short: "Delete a folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if err := client.DeleteFolder(args[0]); err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}

		fmt.Printf("Folder '%s' deleted successfully\n", args[0])
		return nil
	},
}

var renameFolderCmd = &cobra.Command{
	Use:   "rename-folder [folder-id] [new-name]",
	Short: "Rename a folder",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		if err := client.RenameFolder(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to rename folder: %w", err)
		}

		fmt.Printf("Folder renamed to '%s'\n", args[1])
		return nil
	},
}

var folderInfoCmd = &cobra.Command{
	Use:   "folder-info [folder-id]",
	Short: "Get folder information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		info, err := client.GetFolderInfo(args[0])
		if err != nil {
			return fmt.Errorf("failed to get folder info: %w", err)
		}

		fmt.Printf("Folder: %s (ID: %s)\n", info.Name, info.ID)
		fmt.Printf("Type: %s\n", info.Type)
		fmt.Printf("Messages: %d\n", info.Messages)
		fmt.Printf("Unread: %d\n", info.Unread)
		fmt.Printf("Size: %d bytes\n", info.Size)
		return nil
	},
}

var listGroupRootsCmd = &cobra.Command{
	Use:   "groups",
	Short: "List group roots",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		groups, err := client.ListGroupRoots()
		if err != nil {
			return fmt.Errorf("failed to list groups: %w", err)
		}

		fmt.Println("Group Roots:")
		for _, g := range groups {
			fmt.Printf("  %s (ID: %s) - Type: %s\n", g.Name, g.ID, g.Type)
		}
		return nil
	},
}

var uploadedItemsCmd = &cobra.Command{
	Use:   "uploads",
	Short: "List uploaded files",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewMailClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		items, err := client.GetUploadedItems()
		if err != nil {
			return fmt.Errorf("failed to list uploads: %w", err)
		}

		fmt.Println("Uploaded Items:")
		for _, item := range items {
			fmt.Printf("  %s - Size: %d bytes\n", item.Name, item.Size)
		}
		return nil
	},
}

func init() {
	MailCmd.AddCommand(versionCmd)
	MailCmd.AddCommand(listFoldersCmd)
	MailCmd.AddCommand(createFolderCmd)
	MailCmd.AddCommand(deleteFolderCmd)
	MailCmd.AddCommand(renameFolderCmd)
	MailCmd.AddCommand(folderInfoCmd)
	MailCmd.AddCommand(listGroupRootsCmd)
	MailCmd.AddCommand(uploadedItemsCmd)

	createFolderCmd.Flags().StringP("parent", "p", "", "Parent folder ID")
}
