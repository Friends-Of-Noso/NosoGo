package commands

import (
	"github.com/spf13/cobra"
)

var blocksCmd = &cobra.Command{
	Use:   "blocks",
	Short: "Blocks related queries",
	//Long:  `Initializes the configuration file.`,
	//Run: runBlocks,
}

func init() {
	rootCmd.AddCommand(blocksCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blocksCmd.PersistentFlags().String("foo", "", "A help for foo")
	blocksCmd.PersistentFlags().BoolP("json", "j", false, "Outputs results in 'JSON'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blocksCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
