package commands

import (
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var networkCmd = &cobra.Command{
	Use:   "network",
	Short: "Network related queries",
	//Long:  `Initializes the configuration file.`,
	//Run: runNetwork,
}

func init() {
	rootCmd.AddCommand(networkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// networkCmd.PersistentFlags().String("foo", "", "A help for foo")
	networkCmd.PersistentFlags().BoolP("json", "j", false, "Outputs results in 'JSON'")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// networkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
