package cfgcmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View, set, or initialize configuration for icwli.`,
}

var configViewCmd = &cobra.Command{
	Use:   "view",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Current configuration:")
		fmt.Printf("  Server Host: %s\n", viper.GetString("server.host"))
		fmt.Printf("  Server Port: %d\n", viper.GetInt("server.port"))
		fmt.Printf("  Username: %s\n", viper.GetString("auth.username"))
		fmt.Printf("  Batch Size: %d\n", viper.GetInt("imap.batch_size"))
		fmt.Printf("  TLS Verify: %t\n", viper.GetBool("imap.tls_verify"))
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		configPath := getConfigPath()
		viper.SetConfigFile(configPath)

		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("error reading config: %w", err)
			}
		}

		viper.Set(key, value)

		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("error writing config: %w", err)
		}

		fmt.Printf("Set %s = %s\n", key, value)
		return nil
	},
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := getConfigPath()

		if _, err := os.Stat(configPath); err == nil {
			fmt.Println("Config file already exists")
			return nil
		}

		if err := os.MkdirAll(getConfigDir(), 0755); err != nil {
			return fmt.Errorf("error creating config directory: %w", err)
		}

		defaultConfig := `server:
  host: icewarp.armazemdc.inf.br
  port: 993

auth:
  username: ""
  password: ""

imap:
  batch_size: 5000
  tls_verify: true

general:
  debug: false
`

		if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("error writing config file: %w", err)
		}

		fmt.Printf("Created config file at %s\n", configPath)
		return nil
	},
}

func getConfigDir() string {
	home, _ := os.UserHomeDir()
	return fmt.Sprintf("%s/.icwli", home)
}

func getConfigPath() string {
	return fmt.Sprintf("%s/icwli.yaml", getConfigDir())
}

func init() {
	ConfigCmd.AddCommand(configViewCmd)
	ConfigCmd.AddCommand(configSetCmd)
	ConfigCmd.AddCommand(configInitCmd)
}
