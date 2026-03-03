package services

import (
	"fmt"
	"time"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ServicesCmd = &cobra.Command{
		Use:   "services",
		Short: "Server services operations",
		Long:  `Commands for IceWarp Server Services API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Services API",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid, err := client.Authenticate(username, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		fmt.Printf("Login successful!\n")
		fmt.Printf("SID: %s\n", sid)
		return nil
	},
}

var listServicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List server services",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		services, err := client.ListServices()
		if err != nil {
			return fmt.Errorf("failed to list services: %w", err)
		}

		fmt.Println("Server Services:")
		for _, s := range services {
			fmt.Printf("  %s - %s (PID: %d, Status: %s)\n", s.Name, s.DisplayName, s.PID, s.Status)
		}
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status [service-name]",
	Short: "Get service status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		service, err := client.GetServiceStatus(args[0])
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		fmt.Printf("Service: %s\n", service.DisplayName)
		fmt.Printf("Status: %s\n", service.Status)
		fmt.Printf("PID: %d\n", service.PID)
		fmt.Printf("Uptime: %d seconds\n", service.Uptime)
		return nil
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats [service-name]",
	Short: "Get service statistics",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		stats, err := client.GetServiceStats(args[0])
		if err != nil {
			return fmt.Errorf("failed to get stats: %w", err)
		}

		fmt.Printf("Statistics for %s:\n", args[0])
		for k, v := range stats {
			fmt.Printf("  %s: %v\n", k, v)
		}
		return nil
	},
}

var startCmd = &cobra.Command{
	Use:   "start [service-name]",
	Short: "Start a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.StartService(args[0]); err != nil {
			return fmt.Errorf("failed to start: %w", err)
		}

		fmt.Printf("Service %s started\n", args[0])
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop [service-name]",
	Short: "Stop a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.StopService(args[0]); err != nil {
			return fmt.Errorf("failed to stop: %w", err)
		}

		fmt.Printf("Service %s stopped\n", args[0])
		return nil
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart [service-name]",
	Short: "Restart a service",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.RestartService(args[0]); err != nil {
			return fmt.Errorf("failed to restart: %w", err)
		}

		fmt.Printf("Service %s restarted\n", args[0])
		return nil
	},
}

var trafficCmd = &cobra.Command{
	Use:   "traffic [service-name]",
	Short: "Get traffic chart data",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewServicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		startStr, _ := cmd.Flags().GetString("start")
		endStr, _ := cmd.Flags().GetString("end")

		start := time.Now().AddDate(0, 0, -7)
		end := time.Now()

		if startStr != "" {
			var err error
			start, err = time.Parse("2006-01-02", startStr)
			if err != nil {
				return fmt.Errorf("invalid start date: %w", err)
			}
		}

		if endStr != "" {
			var err error
			end, err = time.Parse("2006-01-02", endStr)
			if err != nil {
				return fmt.Errorf("invalid end date: %w", err)
			}
		}

		data, err := client.GetTrafficChart(args[0], start, end)
		if err != nil {
			return fmt.Errorf("failed to get traffic: %w", err)
		}

		fmt.Printf("Traffic data for %s (%s to %s):\n", args[0], start.Format("2006-01-02"), end.Format("2006-01-02"))
		fmt.Printf("%v\n", data)
		return nil
	},
}

func init() {
	ServicesCmd.AddCommand(loginCmd)
	ServicesCmd.AddCommand(listServicesCmd)
	ServicesCmd.AddCommand(statusCmd)
	ServicesCmd.AddCommand(statsCmd)
	ServicesCmd.AddCommand(startCmd)
	ServicesCmd.AddCommand(stopCmd)
	ServicesCmd.AddCommand(restartCmd)
	ServicesCmd.AddCommand(trafficCmd)

	trafficCmd.Flags().StringP("start", "s", "", "Start date (YYYY-MM-DD)")
	trafficCmd.Flags().StringP("end", "e", "", "End date (YYYY-MM-DD)")
}
