package config

import (
	"fmt"
	"os"
	"path"

	ms "github.com/mitchellh/mapstructure"
	toml "github.com/pelletier/go-toml/v2"
	"github.com/spf13/viper"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelDebug = "debug"

	NodeModeDNS       = "dns"
	NodeModeSeed      = "seed"
	NodeModeSuperNode = "supernode"
	NodeModeNode      = "node"

	cConfigFolderName  = ".nosogod"
	cConfigFileName    = "config.toml"
	cLogsFolderName    = "logs"
	cLogLevel          = "info"
	cLogFileName       = "nosogod.log"
	cDatabasePath      = "data"
	DefaultNodeAddress = "0.0.0.0"
	DefaultNodePort    = 45050
	DefaultNodeMode    = NodeModeNode
	DefaultNodeKey     = "Will be changed upon first run"
	DefaultAPIAddress  = "127.0.0.1"
	DefaultAPIPort     = 45505
	DefaultDNSAddress  = "127.0.0.1"
	DefaultDNSPort     = 8080
)

var (
	ValidModes = map[string]bool{
		NodeModeDNS:       true,
		NodeModeSeed:      true,
		NodeModeSuperNode: true,
		NodeModeNode:      true,
	}

	ValidLogLevels = map[string]bool{
		LogLevelInfo: true,
		// LogLevelWarn:  true,
		// LogLevelError: true,
		LogLevelDebug: true,
	}
)

type Config struct {
	// Top level options use an anonymous struct
	BaseConfig `mapstructure:",squash"`
	API        *APIConfig  `mapstructure:"api"`
	Node       *NodeConfig `mapstructure:"node"`
	DNS        *DNSConfig  `mapstructure:"dns"`
}

// DefaultConfig Default configurable parameters.
func DefaultConfig() *Config {
	return &Config{
		BaseConfig: DefaultBaseConfig(),
		API:        DefaultAPIConfig(),
		Node:       DefaultNodeConfig(),
		DNS:        DefaultDNSConfig(),
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

func homeFolder() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not find homedir: %v", err)
		os.Exit(1)
	}

	return home
}

func WriteConfig(file string, config *Config) error {
	log.Infof("Writing config file at: %s", file)
	var outMap map[string]any
	ms.Decode(config, &outMap)
	b, err := toml.Marshal(outMap)
	if err != nil {
		return fmt.Errorf("could not marshal config structure: %v", err)
	}
	if err := utils.MustWriteFile(viper.ConfigFileUsed(), b, 0644); err != nil {
		return err
	}
	return nil
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
	Address    string `mapstructure:"address"`
	Port       int    `mapstructure:"port"`
	Mode       string `mapstructure:"mode"`
	PrivateKey string `mapstructure:"private-key"`
	PublicKey  string `mapstructure:"public-key"`
}

func DefaultNodeConfig() *NodeConfig {
	return &NodeConfig{
		Address:    DefaultNodeAddress,
		Port:       DefaultNodePort,
		Mode:       DefaultNodeMode,
		PrivateKey: DefaultNodeKey,
		PublicKey:  DefaultNodeKey,
	}
}

type DNSConfig struct {
	Address string `mapstructure:"address"`
	Port    int    `mapstructure:"port"`
}

func DefaultDNSConfig() *DNSConfig {
	return &DNSConfig{
		DefaultDNSAddress,
		DefaultDNSPort,
		// "",
	}
}
