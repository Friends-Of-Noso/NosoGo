package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the configuration folder, file and the log folder",
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
	log.Debug("init called")

	if utils.FileExists(config.GetConfigFile()) {
		fmt.Fprintf(os.Stderr, "Config already exists: %s", config.GetConfigFile())
		os.Exit(1)
	}

	// Config Folder
	config.ConfigDir = config.GetConfigFolder()

	// Logs Folder
	utils.EnsureDir(config.GetLogsFolder(), 0755)
	log.SetFileAndLevel(config.GetLogFile(), config.LogLevel)

	// Viper Config File
	viper.AddConfigPath(config.GetConfigFolder())
	viper.SetConfigType("toml")
	viper.SetConfigFile(config.GetConfigFile())

	// Write to Config File
	err := config.WriteConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not save config structure: %v", err)
		os.Exit(1)
	}

	// log.Infof("Created config file at '%s'", config.GetConfigFile())

	// Create LevelDB stuff
	db, err := leveldb.OpenFile(config.GetDatabasePath(), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	fmt.Fprintf(os.Stderr, "created database at '%s'\n", config.GetDatabasePath())
}
