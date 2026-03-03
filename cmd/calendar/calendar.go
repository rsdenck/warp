package calendar

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
	CalendarCmd = &cobra.Command{
		Use:   "calendar",
		Short: "Calendar operations",
		Long:  `Commands for IceWarp Calendar API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Calendar API",
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
		fmt.Printf("Run: warpctl config set calendar.token %s\n", resp.Token)
		return nil
	},
}

var listCalendarsCmd = &cobra.Command{
	Use:   "list",
	Short: "List calendars",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCalendarClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first")
		}
		client.SetToken(token)

		calendars, err := client.ListCalendars()
		if err != nil {
			return fmt.Errorf("failed to list calendars: %w", err)
		}

		t := output.NewTable("CALENDARS")
		t.AppendHeader(table.Row{"Name", "ID", "Description"})
		
		for _, c := range calendars {
			t.AppendRow(table.Row{c.Name, c.ID, c.Description})
		}
		
		t.Render()
		return nil
	},
}

var createCalendarCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new calendar",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCalendarClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first")
		}
		client.SetToken(token)

		description, _ := cmd.Flags().GetString("description")
		color, _ := cmd.Flags().GetString("color")

		cal, err := client.CreateCalendar(args[0], description, color)
		if err != nil {
			return fmt.Errorf("failed to create calendar: %w", err)
		}

		fmt.Printf("Calendar '%s' created (ID: %s)\n", cal.Name, cal.ID)
		return nil
	},
}

var deleteCalendarCmd = &cobra.Command{
	Use:   "delete [calendar-id]",
	Short: "Delete a calendar",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCalendarClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first")
		}
		client.SetToken(token)

		if err := client.DeleteCalendar(args[0]); err != nil {
			return fmt.Errorf("failed to delete calendar: %w", err)
		}

		fmt.Printf("Calendar deleted successfully\n")
		return nil
	},
}

var listEventsCmd = &cobra.Command{
	Use:   "events [calendar-id]",
	Short: "List events in a calendar",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCalendarClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first")
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

		events, err := client.ListEvents(args[0], start, end)
		if err != nil {
			return fmt.Errorf("failed to list events: %w", err)
		}

		t := output.NewTable(fmt.Sprintf("EVENTS (%s to %s)", start.Format("2006-01-02"), end.Format("2006-01-02")))
		t.AppendHeader(table.Row{"Title", "Start", "End"})
		
		for _, e := range events {
			t.AppendRow(table.Row{
				e.Title,
				e.Start.Format("2006-01-02 15:04"),
				e.End.Format("2006-01-02 15:04"),
			})
		}
		
		t.Render()
		return nil
	},
}

var createEventCmd = &cobra.Command{
	Use:   "create-event [calendar-id] [title]",
	Short: "Create a new event in a calendar",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCalendarClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first")
		}
		client.SetToken(token)

		description, _ := cmd.Flags().GetString("description")
		location, _ := cmd.Flags().GetString("location")
		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")
		allDay, _ := cmd.Flags().GetBool("all-day")
		attendees, _ := cmd.Flags().GetStringSlice("attendee")

		start, err := time.Parse("2006-01-02T15:04", startStr)
		if err != nil {
			return fmt.Errorf("invalid start time format: %w", err)
		}

		end, err := time.Parse("2006-01-02T15:04", endStr)
		if err != nil {
			return fmt.Errorf("invalid end time format: %w", err)
		}

		event, err := client.CreateEvent(args[0], args[1], description, location, start, end, allDay, attendees)
		if err != nil {
			return fmt.Errorf("failed to create event: %w", err)
		}

		fmt.Printf("Event '%s' created (ID: %s)\n", event.Title, event.ID)
		return nil
	},
}

var deleteEventCmd = &cobra.Command{
	Use:   "delete-event [event-id]",
	Short: "Delete an event",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewCalendarClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		token := viper.GetString("calendar.token")
		if token == "" {
			return fmt.Errorf("token not set. Use 'warpctl calendar login' first")
		}
		client.SetToken(token)

		if err := client.DeleteEvent(args[0]); err != nil {
			return fmt.Errorf("failed to delete event: %w", err)
		}

		fmt.Printf("Event deleted successfully\n")
		return nil
	},
}

func init() {
	CalendarCmd.AddCommand(loginCmd)
	CalendarCmd.AddCommand(listCalendarsCmd)
	CalendarCmd.AddCommand(createCalendarCmd)
	CalendarCmd.AddCommand(deleteCalendarCmd)
	CalendarCmd.AddCommand(listEventsCmd)
	CalendarCmd.AddCommand(createEventCmd)
	CalendarCmd.AddCommand(deleteEventCmd)

	createCalendarCmd.Flags().StringP("description", "d", "", "Calendar description")
	createCalendarCmd.Flags().StringP("color", "c", "#000000", "Calendar color")

	listEventsCmd.Flags().StringP("start", "s", "", "Start date (YYYY-MM-DD)")
	listEventsCmd.Flags().StringP("end", "e", "", "End date (YYYY-MM-DD)")

	createEventCmd.Flags().StringP("description", "d", "", "Event description")
	createEventCmd.Flags().StringP("location", "l", "", "Event location")
	createEventCmd.Flags().StringP("start", "s", "", "Start time (YYYY-MM-DDTHH:MM)")
	createEventCmd.Flags().StringP("end", "e", "", "End time (YYYY-MM-DDTHH:MM)")
	createEventCmd.Flags().BoolP("all-day", "a", false, "All day event")
	createEventCmd.Flags().StringSliceP("attendee", "t", []string{}, "Attendee email")
}
