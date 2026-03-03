package notes

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/output"
	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	NotesCmd = &cobra.Command{
		Use:   "notes",
		Short: "Notes operations",
		Long:  `Commands for IceWarp Notes API`,
	}
)

var listNotesCmd = &cobra.Command{
	Use:   "list",
	Short: "List notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		folderID, _ := cmd.Flags().GetString("folder")

		notes, err := client.ListNotes(folderID)
		if err != nil {
			return fmt.Errorf("failed to list notes: %w", err)
		}

		t := output.NewTable("NOTES")
		t.AppendHeader(table.Row{"Title", "ID"})
		
		for _, n := range notes {
			t.AppendRow(table.Row{n.Title, n.ID})
		}
		
		t.Render()
		return nil
	},
}

var listFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List note folders",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		folders, err := client.ListFolders()
		if err != nil {
			return fmt.Errorf("failed to list folders: %w", err)
		}

		t := output.NewTable("NOTE FOLDERS")
		t.AppendHeader(table.Row{"Name", "ID"})
		
		for _, f := range folders {
			t.AppendRow(table.Row{f.Name, f.ID})
		}
		
		t.Render()
		return nil
	},
}

var createNoteCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		content, _ := cmd.Flags().GetString("content")
		color, _ := cmd.Flags().GetString("color")
		folderID, _ := cmd.Flags().GetString("folder")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		note, err := client.CreateNote(args[0], content, color, folderID, tags)
		if err != nil {
			return fmt.Errorf("failed to create note: %w", err)
		}

		fmt.Printf("Note '%s' created (ID: %s)\n", note.Title, note.ID)
		return nil
	},
}

var updateNoteCmd = &cobra.Command{
	Use:   "update [note-id]",
	Short: "Update a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		color, _ := cmd.Flags().GetString("color")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		note, err := client.UpdateNote(args[0], title, content, color, tags)
		if err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}

		fmt.Printf("Note '%s' updated (ID: %s)\n", note.Title, note.ID)
		return nil
	},
}

var deleteNoteCmd = &cobra.Command{
	Use:   "delete [note-id]",
	Short: "Delete a note",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteNote(args[0]); err != nil {
			return fmt.Errorf("failed to delete note: %w", err)
		}

		fmt.Printf("Note deleted successfully\n")
		return nil
	},
}

var createFolderCmd = &cobra.Command{
	Use:   "create-folder [name]",
	Short: "Create a new note folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		parentID, _ := cmd.Flags().GetString("parent")

		folder, err := client.CreateFolder(args[0], parentID)
		if err != nil {
			return fmt.Errorf("failed to create folder: %w", err)
		}

		fmt.Printf("Folder '%s' created (ID: %s)\n", folder.Name, folder.ID)
		return nil
	},
}

var deleteFolderCmd = &cobra.Command{
	Use:   "delete-folder [folder-id]",
	Short: "Delete a note folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewNotesClient(&sdk.Config{
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

func init() {
	NotesCmd.AddCommand(listNotesCmd)
	NotesCmd.AddCommand(listFoldersCmd)
	NotesCmd.AddCommand(createNoteCmd)
	NotesCmd.AddCommand(updateNoteCmd)
	NotesCmd.AddCommand(deleteNoteCmd)
	NotesCmd.AddCommand(createFolderCmd)
	NotesCmd.AddCommand(deleteFolderCmd)

	listNotesCmd.Flags().StringP("folder", "f", "", "Folder ID")

	createNoteCmd.Flags().StringP("content", "c", "", "Note content")
	createNoteCmd.Flags().StringP("color", "l", "#FFFF00", "Note color")
	createNoteCmd.Flags().StringP("folder", "f", "", "Folder ID")
	createNoteCmd.Flags().StringSliceP("tag", "t", []string{}, "Tags")

	updateNoteCmd.Flags().StringP("title", "t", "", "Note title")
	updateNoteCmd.Flags().StringP("content", "c", "", "Note content")
	updateNoteCmd.Flags().StringP("color", "l", "", "Note color")
	updateNoteCmd.Flags().StringSliceP("tag", "g", []string{}, "Tags")

	createFolderCmd.Flags().StringP("parent", "p", "", "Parent folder ID")
}
