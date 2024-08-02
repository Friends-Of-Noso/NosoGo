package commands

import (
	"github.com/spf13/cobra"
	//ms "github.com/mitchellh/mapstructure"
	//toml "github.com/pelletier/go-toml/v2"
	//log "github.com/EIYARO-Project/core-stratumd/logger"
	//fs "github.com/EIYARO-Project/core-stratumd/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the configuration file",
	//Long:  `Initializes the configuration file.`,
	Run: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runInit(cmd *cobra.Command, args []string) {
	/*log.Debug("init called")
	if fs.FileExists(config.GetConfigFile()) {
		log.Fatalf("Config already exists: %s", config.GetConfigFile())
		os.Exit(1)
	}

	// Config Folder
	config.ConfigDir = config.GetConfigFolder()

	// Logs Folder
	fs.EnsureDir(config.GetLogsFolder(), 0755)
	if logLevel != "" {
		log.SetFileAndLevel(config.GetLogFile(), logLevel)
	} else {
		log.SetFileAndLevel(config.GetLogFile(), config.LogLevel)
	}

	// Viper Config File
	viper.AddConfigPath(config.GetConfigFolder())
	viper.SetConfigType("toml")
	viper.SetConfigFile(config.GetConfigFile())
	if err := viper.SafeWriteConfig(); err != nil {
		log.Fatalf("Error saving config file: %s", err)
	}

	// Write to Config File
	var outMap map[string]any
	ms.Decode(config, &outMap)
	b, err := toml.Marshal(outMap)
	cobra.CheckErr(err)
	fs.MustWriteFile(viper.ConfigFileUsed(), b, 0644)

	log.Infof("Created config file at %s", config.GetLogFile())*/

	// Create LevelDB stuff
}
