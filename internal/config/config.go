package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	Auth    AuthConfig    `mapstructure:"auth"`
	IMAP    IMAPConfig    `mapstructure:"imap"`
	General GeneralConfig `mapstructure:"general"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type AuthConfig struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type IMAPConfig struct {
	BatchSize int  `mapstructure:"batch_size"`
	TLSVerify bool `mapstructure:"tls_verify"`
}

type GeneralConfig struct {
	Debug bool `mapstructure:"debug"`
}

func Load() (*Config, error) {
	viper.SetConfigName("icwli")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.icwli")
	viper.AddConfigPath(".")
	viper.SetDefault("server.host", "icewarp.armazemdc.inf.br")
	viper.SetDefault("server.port", 993)
	viper.SetDefault("imap.batch_size", 5000)
	viper.SetDefault("imap.tls_verify", true)
	viper.SetDefault("general.debug", false)

	viper.BindEnv("server.host", "IW_IMAP_HOST")
	viper.BindEnv("server.port", "IW_IMAP_PORT")
	viper.BindEnv("auth.username", "IW_USERNAME")
	viper.BindEnv("auth.password", "IW_PASSWORD")
	viper.BindEnv("imap.batch_size", "IW_BATCH_SIZE")
	viper.BindEnv("imap.tls_verify", "IW_TLS_VERIFY")
	viper.BindEnv("general.debug", "IW_DEBUG")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Fprintf(os.Stderr, "Warning: Error reading config file: %v\n", err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
