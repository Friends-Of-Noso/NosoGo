package config

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

const (
	cConfigFolderName = ".nosogod"
	cConfigFileName   = "config.toml"
	cLogsFolderName   = "logs"
	cLogLevel         = "info"
	cLogFileName      = "nosogod.log"
	cNodeAddress      = "0.0.0.0"
	cNodePort         = 45050
)

func homeFolder() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return home
}

type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`
	Database   *DatabaseConfig `mapstructure:"database"`
	Node       *NodeConfig     `mapstructure:"node"`
}

// DefaultConfig Default configurable parameters.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		Database:   DefaultDatabaseConfig(),
		Node:       DefaultStratumServerConfig(),
	}
}

func (c *Config) GetConfigFolder() string {
	if c.ConfigDir != "" {
		return c.ConfigDir
	} else {
		return path.Join(homeFolder(), cConfigFolderName)
	}
}

func (c *Config) GetConfigFile() string {
	if c.ConfigDir != "" {
		return path.Join(c.ConfigDir, cConfigFileName)
	} else {
		return path.Join(homeFolder(), cConfigFolderName, cConfigFileName)
	}
}

func (c *Config) GetLogsFolder() string {
	if c.ConfigDir != "" && c.LogFolder != "" {
		return path.Join(c.ConfigDir, c.LogFolder)
	} else {
		return path.Join(homeFolder(), cConfigFolderName, cLogsFolderName)
	}
}

func (c *Config) GetLogFile() string {
	if c.ConfigDir != "" && c.LogFolder != "" && c.LogFile != "" {
		return path.Join(c.ConfigDir, c.LogFolder, c.LogFile)
	} else {
		return path.Join(homeFolder(), cConfigFolderName, cLogsFolderName, cLogFileName)
	}
}

type BaseConfig struct {
	// The root directory for all data.
	// This should be set in viper so it can unmarshal into this struct
	ConfigDir string `mapstructure:"config_folder"`
	//log level to set
	LogLevel string `mapstructure:"log_level"`
	// log file name
	LogFolder string `mapstructure:"log_folder"`
	// log file name
	LogFile string `mapstructure:"log_file"`
}

// DefaultBaseConfig Default configurable base parameters.
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		LogLevel:  cLogLevel,
		LogFolder: cLogsFolderName,
		LogFile:   cLogFileName,
	}
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"name"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func DefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		"localhost",
		0,
		"nosogod",
		"nosogod-user",
		"Secret",
	}
}

type NodeConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func DefaultStratumServerConfig() *NodeConfig {
	return &NodeConfig{
		cNodeAddress,
		cNodePort,
	}
}
