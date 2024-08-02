package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	//cfg "github.com/EIYARO-Project/core-stratumd/config"
	//log "github.com/EIYARO-Project/core-stratumd/logger"
	//fs "github.com/EIYARO-Project/core-stratumd/utils"
	ver "github.com/Friends-Of-Noso/NosoGo/version"
)

const (
	cLogLevelFlag = "log-level"
)

var (
	//config   = cfg.DefaultConfig()
	cfgFile  string
	logLevel string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: ver.Version,
	Use:     ver.Name,
	Short:   "The node for the NOSO crypto coin",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVarP(&logLevel, cLogLevelFlag, "l", "", "Log Level")
	viper.BindPFlag("log_level", rootCmd.Flags().Lookup(cLogLevelFlag))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", ver.Title))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	/*if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigType("toml")
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".stratumd" (without extension).
		viper.AddConfigPath(config.GetConfigFolder())
		viper.SetConfigType("toml")
		viper.SetConfigFile(config.GetConfigFile())
	}

	viper.AutomaticEnv() // read in environment variables that match

	if fs.FileExists(viper.ConfigFileUsed()) {
		// Logger
		if logLevel != "" {
			log.SetFileAndLevel(config.GetLogFile(), logLevel)
		}
		// Read Config
		if err := viper.ReadInConfig(); err == nil {
			log.Infof("Using config file: %s", viper.ConfigFileUsed())

			err := viper.Unmarshal(config)
			if err != nil {
				log.Fatalf("Could not unmarshal config: %s", err)
			}
		}
		if logLevel == "" {
			log.SetFileAndLevel(config.GetLogFile(), config.LogLevel)
		}
	} else {
		log.SetFileAndLevel("", "info")
	}*/
}
