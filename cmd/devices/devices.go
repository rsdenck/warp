package devices

import (
	"fmt"

	"github.com/icewarp/warpctl/internal/sdk"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	DevicesCmd = &cobra.Command{
		Use:   "devices",
		Short: "Mobile devices operations",
		Long:  `Commands for IceWarp Mobile Devices API`,
	}
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Devices API",
	RunE: func(cmd *cobra.Command, args []string) error {
		username := viper.GetString("auth.username")
		password := viper.GetString("auth.password")

		if username == "" || password == "" {
			return fmt.Errorf("username and password are required")
		}

		client := sdk.NewDevicesClient(&sdk.Config{
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

var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "List mobile devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewDevicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		account, _ := cmd.Flags().GetString("account")
		filter, _ := cmd.Flags().GetString("filter")

		devices, err := client.ListDevices(account, filter)
		if err != nil {
			return fmt.Errorf("failed to list devices: %w", err)
		}

		fmt.Println("Mobile Devices:")
		for _, d := range devices {
			fmt.Printf("  %s - %s (%s) - Status: %s\n", d.Model, d.OS, d.OSVersion, d.Status)
		}
		return nil
	},
}

var getDeviceInfoCmd = &cobra.Command{
	Use:   "info [device-id]",
	Short: "Get device information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewDevicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		device, err := client.GetDeviceInfo(args[0])
		if err != nil {
			return fmt.Errorf("failed to get device info: %w", err)
		}

		fmt.Printf("Device: %s\n", device.Model)
		fmt.Printf("Type: %s\n", device.Type)
		fmt.Printf("OS: %s %s\n", device.OS, device.OSVersion)
		fmt.Printf("Account: %s\n", device.Account)
		fmt.Printf("Status: %s\n", device.Status)
		fmt.Printf("Approved: %v\n", device.Approved)
		fmt.Printf("Remote Wipe: %v\n", device.RemoteWipe)
		return nil
	},
}

var deleteDeviceCmd = &cobra.Command{
	Use:   "delete [device-id]",
	Short: "Delete a mobile device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewDevicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.DeleteDevice(args[0]); err != nil {
			return fmt.Errorf("failed to delete device: %w", err)
		}

		fmt.Printf("Device deleted\n")
		return nil
	},
}

var remoteWipeCmd = &cobra.Command{
	Use:   "wipe [device-id]",
	Short: "Remote wipe a mobile device",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewDevicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		if err := client.RemoteWipe(args[0]); err != nil {
			return fmt.Errorf("failed to remote wipe: %w", err)
		}

		fmt.Printf("Remote wipe initiated\n")
		return nil
	},
}

var setDeviceStatusCmd = &cobra.Command{
	Use:   "set-status [device-id]",
	Short: "Set device status",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := sdk.NewDevicesClient(&sdk.Config{
			BaseURL: viper.GetString("server.url"),
		})

		sid := viper.GetString("maintenance.sid")
		if sid == "" {
			return fmt.Errorf("SID not set. Use 'warpctl maintenance login' first")
		}
		client.SetSID(sid)

		status, _ := cmd.Flags().GetString("status")

		if err := client.SetDeviceStatus(args[0], status); err != nil {
			return fmt.Errorf("failed to set status: %w", err)
		}

		fmt.Printf("Device status updated\n")
		return nil
	},
}

func init() {
	DevicesCmd.AddCommand(loginCmd)
	DevicesCmd.AddCommand(listDevicesCmd)
	DevicesCmd.AddCommand(getDeviceInfoCmd)
	DevicesCmd.AddCommand(deleteDeviceCmd)
	DevicesCmd.AddCommand(remoteWipeCmd)
	DevicesCmd.AddCommand(setDeviceStatusCmd)

	listDevicesCmd.Flags().StringP("account", "a", "", "Account email")
	listDevicesCmd.Flags().StringP("filter", "f", "", "Filter")

	setDeviceStatusCmd.Flags().StringP("status", "s", "active", "Device status")
}
