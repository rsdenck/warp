package tasks

import (
	"fmt"
	"time"

	"github.com/icewarp/warpctl/internal/output"
	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	TasksCmd = &cobra.Command{
		Use:   "tasks",
		Short: "Tasks operations",
		Long:  `Commands for IceWarp Tasks API`,
	}
)

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		folderID, _ := cmd.Flags().GetString("folder")

		tasks, err := client.ListTasks(folderID)
		if err != nil {
			return fmt.Errorf("failed to list tasks: %w", err)
		}

		t := output.NewTable("TASKS")
		t.AppendHeader(table.Row{"Title", "ID", "Status", "Priority"})
		
		for _, task := range tasks {
			t.AppendRow(table.Row{task.Title, task.ID, task.Status, task.Priority})
		}
		
		t.Render()
		return nil
	},
}

var listFoldersCmd = &cobra.Command{
	Use:   "folders",
	Short: "List task folders",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
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

		t := output.NewTable("TASK FOLDERS")
		t.AppendHeader(table.Row{"Name", "ID"})
		
		for _, f := range folders {
			t.AppendRow(table.Row{f.Name, f.ID})
		}
		
		t.Render()
		return nil
	},
}

var createTaskCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		description, _ := cmd.Flags().GetString("description")
		folderID, _ := cmd.Flags().GetString("folder")
		assignee, _ := cmd.Flags().GetString("assignee")
		priority, _ := cmd.Flags().GetInt("priority")
		dueDateStr, _ := cmd.Flags().GetString("due")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		var dueDate *time.Time
		if dueDateStr != "" {
			d, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				return fmt.Errorf("invalid due date format: %w", err)
			}
			dueDate = &d
		}

		task, err := client.CreateTask(args[0], description, folderID, assignee, priority, dueDate, tags)
		if err != nil {
			return fmt.Errorf("failed to create task: %w", err)
		}

		fmt.Printf("Task '%s' created (ID: %s)\n", task.Title, task.ID)
		return nil
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update [task-id]",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetInt("priority")
		dueDateStr, _ := cmd.Flags().GetString("due")
		tags, _ := cmd.Flags().GetStringSlice("tag")

		var dueDate *time.Time
		if dueDateStr != "" {
			d, err := time.Parse("2006-01-02", dueDateStr)
			if err != nil {
				return fmt.Errorf("invalid due date format: %w", err)
			}
			dueDate = &d
		}

		task, err := client.UpdateTask(args[0], title, description, status, priority, dueDate, tags)
		if err != nil {
			return fmt.Errorf("failed to update task: %w", err)
		}

		fmt.Printf("Task '%s' updated (ID: %s)\n", task.Title, task.ID)
		return nil
	},
}

var completeTaskCmd = &cobra.Command{
	Use:   "complete [task-id]",
	Short: "Mark a task as completed",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.CompleteTask(args[0]); err != nil {
			return fmt.Errorf("failed to complete task: %w", err)
		}

		fmt.Printf("Task completed successfully\n")
		return nil
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "delete [task-id]",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteTask(args[0]); err != nil {
			return fmt.Errorf("failed to delete task: %w", err)
		}

		fmt.Printf("Task deleted successfully\n")
		return nil
	},
}

var createFolderCmd = &cobra.Command{
	Use:   "create-folder [name]",
	Short: "Create a new task folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
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
	Short: "Delete a task folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewTasksClient(&sdk.Config{
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
	TasksCmd.AddCommand(listTasksCmd)
	TasksCmd.AddCommand(listFoldersCmd)
	TasksCmd.AddCommand(createTaskCmd)
	TasksCmd.AddCommand(updateTaskCmd)
	TasksCmd.AddCommand(completeTaskCmd)
	TasksCmd.AddCommand(deleteTaskCmd)
	TasksCmd.AddCommand(createFolderCmd)
	TasksCmd.AddCommand(deleteFolderCmd)

	listTasksCmd.Flags().StringP("folder", "f", "", "Folder ID")

	createTaskCmd.Flags().StringP("description", "d", "", "Task description")
	createTaskCmd.Flags().StringP("folder", "f", "", "Folder ID")
	createTaskCmd.Flags().StringP("assignee", "a", "", "Assignee email")
	createTaskCmd.Flags().IntP("priority", "p", 3, "Priority (1-5)")
	createTaskCmd.Flags().StringP("due", "t", "", "Due date (YYYY-MM-DD)")
	createTaskCmd.Flags().StringSliceP("tag", "g", []string{}, "Tags")

	updateTaskCmd.Flags().StringP("title", "t", "", "Task title")
	updateTaskCmd.Flags().StringP("description", "d", "", "Task description")
	updateTaskCmd.Flags().StringP("status", "s", "", "Status (pending/in-progress/completed)")
	updateTaskCmd.Flags().IntP("priority", "p", 0, "Priority (1-5)")
	updateTaskCmd.Flags().StringP("due", "u", "", "Due date (YYYY-MM-DD)")
	updateTaskCmd.Flags().StringSliceP("tag", "g", []string{}, "Tags")

	createFolderCmd.Flags().StringP("parent", "p", "", "Parent folder ID")
}
