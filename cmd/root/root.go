package root

import (
	"fmt"
	"os"

	"github.com/icewarp/warpctl/cmd/calendar"
	"github.com/icewarp/warpctl/cmd/certs"
	"github.com/icewarp/warpctl/cmd/cfgcmd"
	"github.com/icewarp/warpctl/cmd/clean"
	"github.com/icewarp/warpctl/cmd/conferences"
	"github.com/icewarp/warpctl/cmd/contacts"
	"github.com/icewarp/warpctl/cmd/devices"
	"github.com/icewarp/warpctl/cmd/files"
	"github.com/icewarp/warpctl/cmd/filters"
	"github.com/icewarp/warpctl/cmd/mail"
	"github.com/icewarp/warpctl/cmd/maintenance"
	"github.com/icewarp/warpctl/cmd/notes"
	"github.com/icewarp/warpctl/cmd/pinchtab"
	"github.com/icewarp/warpctl/cmd/services"
	"github.com/icewarp/warpctl/cmd/spam"
	"github.com/icewarp/warpctl/cmd/tasks"
	"github.com/icewarp/warpctl/cmd/teamchat"
	"github.com/icewarp/warpctl/cmd/zabbix"
	"github.com/icewarp/warpctl/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version   = "1.0.0"
	commit    = "dev"
	date      = "unknown"
	cfgFile   string
	debugMode bool
)

var rootCmd = &cobra.Command{
	Use:   "warpctl",
	Short: "warpctl - IceWarp CLI",
	Long: `warpctl is a professional CLI tool for managing IceWarp servers.

Complete documentation is available at https://github.com/icewarp/warpctl`,
	Version: fmt.Sprintf("%s (commit: %s, date: %s)", version, commit, date),
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(debugMode)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: $HOME/.warpctl/warpctl.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "enable debug mode")

	rootCmd.AddCommand(clean.CleanCmd)
	rootCmd.AddCommand(cfgcmd.ConfigCmd)
	rootCmd.AddCommand(teamchat.TeamChatCmd)
	rootCmd.AddCommand(mail.MailCmd)
	rootCmd.AddCommand(maintenance.MaintenanceCmd)
	rootCmd.AddCommand(calendar.CalendarCmd)
	rootCmd.AddCommand(notes.NotesCmd)
	rootCmd.AddCommand(tasks.TasksCmd)
	rootCmd.AddCommand(files.FilesCmd)
	rootCmd.AddCommand(contacts.ContactsCmd)
	rootCmd.AddCommand(conferences.ConferencesCmd)
	rootCmd.AddCommand(filters.FiltersCmd)
	rootCmd.AddCommand(spam.SpamCmd)
	rootCmd.AddCommand(devices.DevicesCmd)
	rootCmd.AddCommand(services.ServicesCmd)
	rootCmd.AddCommand(certs.CertsCmd)
	rootCmd.AddCommand(zabbix.ZabbixCmd)
	rootCmd.AddCommand(pinchtab.PinchtabCmd)
}

func initConfig() {
	if cfgFile != "" {
		os.Setenv("WARPCTL_CONFIG_FILE", cfgFile)
	}

	viper.SetConfigName("warpctl")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.warpctl")
	viper.AddConfigPath(".")
	viper.SetDefault("server.host", "icewarp.armazemdc.inf.br")
	viper.SetDefault("server.port", 993)
	viper.SetDefault("server.url", "https://icewarp.armazemdc.inf.br")
	viper.SetDefault("imap.batch_size", 5000)
	viper.SetDefault("imap.tls_verify", true)
	viper.SetDefault("general.debug", false)
	viper.SetDefault("teamchat.token", "")

	viper.BindEnv("server.host", "IW_IMAP_HOST")
	viper.BindEnv("server.port", "IW_IMAP_PORT")
	viper.BindEnv("server.url", "IW_SERVER_URL")
	viper.BindEnv("auth.username", "IW_USERNAME")
	viper.BindEnv("auth.password", "IW_PASSWORD")
	viper.BindEnv("imap.batch_size", "IW_BATCH_SIZE")
	viper.BindEnv("imap.tls_verify", "IW_TLS_VERIFY")
	viper.BindEnv("general.debug", "IW_DEBUG")
	viper.BindEnv("teamchat.token", "IW_TEAMCHAT_TOKEN")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Warning: Error reading config file: %v\n", err)
		}
	}

	if debugMode {
		viper.Set("general.debug", true)
	}
}
