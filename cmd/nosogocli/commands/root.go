package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	"github.com/Friends-Of-Noso/NosoGo/utils"
	ver "github.com/Friends-Of-Noso/NosoGo/version"
)

var (
	config  = cfg.DefaultConfig()
	cfgFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: ver.Version,
	Use:     fmt.Sprintf("%scli", ver.Name),
	Short:   "The client for the NOSO crypto coin node",
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

	rootCmd.SetVersionTemplate(fmt.Sprintf("%s\n", ver.Title))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigType("toml")
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".nosogod" (without extension).
		viper.AddConfigPath(config.GetConfigFolder())
		viper.SetConfigType("toml")
		viper.SetConfigFile(config.GetConfigFile())
	}

	viper.AutomaticEnv() // read in environment variables that match

	if utils.FileExists(viper.ConfigFileUsed()) {
		// Read Config
		if err := viper.ReadInConfig(); err == nil {
			// fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())

			err := viper.Unmarshal(config)
			if err != nil {
				fmt.Printf("Could not unmarshal config: %s\n", err)
			}
		}
	}
}
