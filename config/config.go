package config

import (
	"os"
	"path"

	"github.com/spf13/cobra"
)

const (
	cConfigFolderName  = ".nosogod"
	cConfigFileName    = "config.toml"
	cLogsFolderName    = "logs"
	cLogLevel          = "info"
	cLogFileName       = "nosogod.log"
	cDatabasePath      = "data"
	DefaultNodeAddress = "0.0.0.0"
	DefaultNodePort    = 45050
	DefaultAPIAddress  = "127.0.0.1"
	DefaultAPIPort     = 45505
)

func homeFolder() string {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	return home
}

type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`
	API        *APIConfig  `mapstructure:"api"`
	Node       *NodeConfig `mapstructure:"node"`
}

// DefaultConfig Default configurable parameters.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		API:        DefaultAPIConfig(),
		Node:       DefaultNodeConfig(),
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

func (c *Config) GetDatabaseFolder() string {
	if c.ConfigDir != "" && c.DatabasePath != "" {
		return path.Join(c.ConfigDir, c.DatabasePath)
	} else {
		return path.Join(homeFolder(), cConfigFolderName, cDatabasePath)
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
	// LevelDB path
	DatabasePath string `mapstructure:"database_path"`
}

// DefaultBaseConfig Default configurable base parameters.
func DefaultBaseConfig() BaseConfig {
	return BaseConfig{
		LogLevel:     cLogLevel,
		LogFolder:    cLogsFolderName,
		LogFile:      cLogFileName,
		DatabasePath: cDatabasePath,
	}
}

type APIConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func DefaultAPIConfig() *APIConfig {
	return &APIConfig{
		DefaultAPIAddress,
		DefaultAPIPort,
	}
}

type NodeConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		DefaultNodeAddress,
		DefaultNodePort,
	}
}
