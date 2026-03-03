package conferences

import (
	"fmt"
	"time"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ConferencesCmd = &cobra.Command{
		Use:   "conferences",
		Short: "Conferences/Meetings operations",
		Long:  `Commands for IceWarp Conferences API`,
	}
)

var listConferencesCmd = &cobra.Command{
	Use:   "list",
	Short: "List conferences",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")

		var start, end time.Time
		var err error

		if startStr != "" {
			start, err = time.Parse("2006-01-02", startStr)
			if err != nil {
				return fmt.Errorf("invalid start date format: %w", err)
			}
		} else {
			start = time.Now()
		}

		if endStr != "" {
			end, err = time.Parse("2006-01-02", endStr)
			if err != nil {
				return fmt.Errorf("invalid end date format: %w", err)
			}
		} else {
			end = start.AddDate(0, 1, 0)
		}

		conferences, err := client.ListConferences(start, end)
		if err != nil {
			return fmt.Errorf("failed to list conferences: %w", err)
		}

		fmt.Printf("Conferences (%s to %s):\n", start.Format("2006-01-02"), end.Format("2006-01-02"))
		for _, c := range conferences {
			fmt.Printf("  %s - %s to %s (Status: %s)\n", c.Title, c.StartTime.Format("2006-01-02 15:04"), c.EndTime.Format("2006-01-02 15:04"), c.Status)
		}
		return nil
	},
}

var listRoomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "List conference rooms",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		rooms, err := client.ListRooms()
		if err != nil {
			return fmt.Errorf("failed to list rooms: %w", err)
		}

		fmt.Println("Conference Rooms:")
		for _, r := range rooms {
			fmt.Printf("  %s (ID: %s, Capacity: %d)\n", r.Name, r.ID, r.Capacity)
		}
		return nil
	},
}

var createConferenceCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new conference",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		description, _ := cmd.Flags().GetString("description")
		startStr, _ := cmd.Flags().GetString("start")
		duration, _ := cmd.Flags().GetInt("duration")
		roomID, _ := cmd.Flags().GetString("room")
		password, _ := cmd.Flags().GetString("password")
		recording, _ := cmd.Flags().GetBool("recording")
		attendees, _ := cmd.Flags().GetStringSlice("attendee")

		startTime, err := time.Parse("2006-01-02T15:04", startStr)
		if err != nil {
			return fmt.Errorf("invalid start time format: %w", err)
		}

		if duration == 0 {
			duration = 60
		}

		conf, err := client.CreateConference(args[0], description, startTime, duration, attendees, roomID, password, recording)
		if err != nil {
			return fmt.Errorf("failed to create conference: %w", err)
		}

		fmt.Printf("Conference '%s' created (ID: %s)\n", conf.Title, conf.ID)
		fmt.Printf("Join URL: %s\n", conf.JoinURL)
		fmt.Printf("Host URL: %s\n", conf.HostURL)
		return nil
	},
}

var updateConferenceCmd = &cobra.Command{
	Use:   "update [conference-id]",
	Short: "Update a conference",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		startStr, _ := cmd.Flags().GetString("start")
		duration, _ := cmd.Flags().GetInt("duration")
		password, _ := cmd.Flags().GetString("password")
		recording, _ := cmd.Flags().GetBool("recording")
		attendees, _ := cmd.Flags().GetStringSlice("attendee")

		startTime, _ := time.Parse("2006-01-02T15:04", startStr)

		conf, err := client.UpdateConference(args[0], title, description, startTime, duration, attendees, password, recording)
		if err != nil {
			return fmt.Errorf("failed to update conference: %w", err)
		}

		fmt.Printf("Conference updated (ID: %s)\n", conf.ID)
		return nil
	},
}

var deleteConferenceCmd = &cobra.Command{
	Use:   "delete [conference-id]",
	Short: "Delete a conference",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteConference(args[0]); err != nil {
			return fmt.Errorf("failed to delete conference: %w", err)
		}

		fmt.Printf("Conference deleted successfully\n")
		return nil
	},
}

var addParticipantsCmd = &cobra.Command{
	Use:   "add-participants [conference-id]",
	Short: "Add participants to a conference",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		attendees, _ := cmd.Flags().GetStringSlice("attendee")

		if err := client.AddParticipants(args[0], attendees); err != nil {
			return fmt.Errorf("failed to add participants: %w", err)
		}

		fmt.Printf("Participants added successfully\n")
		return nil
	},
}

var removeParticipantsCmd = &cobra.Command{
	Use:   "remove-participants [conference-id]",
	Short: "Remove participants from a conference",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		attendees, _ := cmd.Flags().GetStringSlice("attendee")

		if err := client.RemoveParticipants(args[0], attendees); err != nil {
			return fmt.Errorf("failed to remove participants: %w", err)
		}

		fmt.Printf("Participants removed successfully\n")
		return nil
	},
}

var getInfoCmd = &cobra.Command{
	Use:   "info [conference-id]",
	Short: "Get conference information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		conf, err := client.GetConferenceInfo(args[0])
		if err != nil {
			return fmt.Errorf("failed to get conference info: %w", err)
		}

		fmt.Printf("Conference: %s (ID: %s)\n", conf.Title, conf.ID)
		fmt.Printf("Description: %s\n", conf.Description)
		fmt.Printf("Start: %s\n", conf.StartTime.Format("2006-01-02 15:04"))
		fmt.Printf("End: %s\n", conf.EndTime.Format("2006-01-02 15:04"))
		fmt.Printf("Organizer: %s\n", conf.Organizer)
		fmt.Printf("Status: %s\n", conf.Status)
		fmt.Printf("Join URL: %s\n", conf.JoinURL)
		return nil
	},
}

var createRoomCmd = &cobra.Command{
	Use:   "create-room [name]",
	Short: "Create a new conference room",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		description, _ := cmd.Flags().GetString("description")
		capacity, _ := cmd.Flags().GetInt("capacity")

		room, err := client.CreateRoom(args[0], description, capacity)
		if err != nil {
			return fmt.Errorf("failed to create room: %w", err)
		}

		fmt.Printf("Room '%s' created (ID: %s)\n", room.Name, room.ID)
		return nil
	},
}

var deleteRoomCmd = &cobra.Command{
	Use:   "delete-room [room-id]",
	Short: "Delete a conference room",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewConferencesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first to get token")
		}
		client.SetToken(token)

		if err := client.DeleteRoom(args[0]); err != nil {
			return fmt.Errorf("failed to delete room: %w", err)
		}

		fmt.Printf("Room deleted successfully\n")
		return nil
	},
}

func init() {
	ConferencesCmd.AddCommand(listConferencesCmd)
	ConferencesCmd.AddCommand(listRoomsCmd)
	ConferencesCmd.AddCommand(createConferenceCmd)
	ConferencesCmd.AddCommand(updateConferenceCmd)
	ConferencesCmd.AddCommand(deleteConferenceCmd)
	ConferencesCmd.AddCommand(addParticipantsCmd)
	ConferencesCmd.AddCommand(removeParticipantsCmd)
	ConferencesCmd.AddCommand(getInfoCmd)
	ConferencesCmd.AddCommand(createRoomCmd)
	ConferencesCmd.AddCommand(deleteRoomCmd)

	listConferencesCmd.Flags().StringP("start", "s", "", "Start date (YYYY-MM-DD)")
	listConferencesCmd.Flags().StringP("end", "e", "", "End date (YYYY-MM-DD)")

	createConferenceCmd.Flags().StringP("description", "d", "", "Description")
	createConferenceCmd.Flags().StringP("start", "s", "", "Start time (YYYY-MM-DDTHH:MM)")
	createConferenceCmd.Flags().IntP("duration", "t", 60, "Duration in minutes")
	createConferenceCmd.Flags().StringP("room", "r", "", "Room ID")
	createConferenceCmd.Flags().StringP("password", "p", "", "Meeting password")
	createConferenceCmd.Flags().BoolP("recording", "c", false, "Enable recording")
	createConferenceCmd.Flags().StringSliceP("attendee", "a", []string{}, "Attendee emails")

	updateConferenceCmd.Flags().StringP("title", "t", "", "Title")
	updateConferenceCmd.Flags().StringP("description", "d", "", "Description")
	updateConferenceCmd.Flags().StringP("start", "s", "", "Start time (YYYY-MM-DDTHH:MM)")
	updateConferenceCmd.Flags().IntP("duration", "u", 0, "Duration in minutes")
	updateConferenceCmd.Flags().StringP("password", "p", "", "Meeting password")
	updateConferenceCmd.Flags().BoolP("recording", "c", false, "Enable recording")
	updateConferenceCmd.Flags().StringSliceP("attendee", "a", []string{}, "Attendee emails")

	addParticipantsCmd.Flags().StringSliceP("attendee", "a", []string{}, "Attendee emails")
	removeParticipantsCmd.Flags().StringSliceP("attendee", "a", []string{}, "Attendee emails")

	createRoomCmd.Flags().StringP("description", "d", "", "Room description")
	createRoomCmd.Flags().IntP("capacity", "c", 100, "Room capacity")
}
