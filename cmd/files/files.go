package files

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	FilesCmd = &cobra.Command{
		Use:   "files",
		Short: "Files operations",
		Long:  `Commands for IceWarp Files API`,
	}
)

var listFilesCmd = &cobra.Command{
	Use:   "list",
	Short: "List files in a path",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		path := "/"
		if len(args) > 0 {
			path = args[0]
		}

		files, err := client.ListFiles(path)
		if err != nil {
			return fmt.Errorf("failed to list files: %w", err)
		}

		fmt.Printf("Files in %s:\n", path)
		for _, f := range files {
			fmt.Printf("  %s (ID: %s, Size: %d bytes, Type: %s)\n", f.Name, f.ID, f.Size, f.Type)
		}
		return nil
	},
}

var listFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List folders in a path",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		path := "/"
		if len(args) > 0 {
			path = args[0]
		}

		folders, err := client.ListFolders(path)
		if err != nil {
			return fmt.Errorf("failed to list folders: %w", err)
		}

		fmt.Printf("Folders in %s:\n", path)
		for _, f := range folders {
			fmt.Printf("  %s (ID: %s)\n", f.Name, f.ID)
		}
		return nil
	},
}

var createFolderCmd = &cobra.Command{
	Use:   "create-folder [name]",
	Short: "Create a new folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		parentPath, _ := cmd.Flags().GetString("parent")

		folder, err := client.CreateFolder(args[0], parentPath)
		if err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}

		fmt.Printf("Folder '%s' created (ID: %s)\n", folder.Name, folder.ID)
		return nil
	},
}

var deleteFileCmd = &cobra.Command{
	Use:   "delete [file-id]",
	Short: "Delete a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteFile(args[0]); err != nil {
			return fmt.Errorf("failed to delete file: %w", err)
		}

		fmt.Printf("File deleted successfully\n")
		return nil
	},
}

var deleteFolderCmd = &cobra.Command{
	Use:   "delete-folder [folder-id]",
	Short: "Delete a folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteFolder(args[0]); err != nil {
			return fmt.Errorf("failed to delete folder: %w", err)
		}

		fmt.Printf("Folder deleted successfully\n")
		return nil
	},
}

var shareFileCmd = &cobra.Command{
	Use:   "share [file-id]",
	Short: "Share a file and get shareable link",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		link, err := client.ShareFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to share file: %w", err)
		}

		fmt.Printf("Shareable link: %s\n", link)
		return nil
	},
}

var moveFileCmd = &cobra.Command{
	Use:   "move [file-id] [new-path]",
	Short: "Move a file to a new path",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.MoveFile(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to move file: %w", err)
		}

		fmt.Printf("File moved to %s\n", args[1])
		return nil
	},
}

var copyFileCmd = &cobra.Command{
	Use:   "copy [file-id] [new-path]",
	Short: "Copy a file to a new path",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewFilesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.CopyFile(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		fmt.Printf("File copied to %s\n", args[1])
		return nil
	},
}

func init() {
	FilesCmd.AddCommand(listFilesCmd)
	FilesCmd.AddCommand(listFoldersCmd)
	FilesCmd.AddCommand(createFolderCmd)
	FilesCmd.AddCommand(deleteFileCmd)
	FilesCmd.AddCommand(deleteFolderCmd)
	FilesCmd.AddCommand(shareFileCmd)
	FilesCmd.AddCommand(moveFileCmd)
	FilesCmd.AddCommand(copyFileCmd)

	createFolderCmd.Flags().StringP("parent", "p", "/", "Parent path")
}
