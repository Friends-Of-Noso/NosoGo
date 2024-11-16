package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var networkStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Network Status",
	//Long:  `Initializes the configuration file.`,
	Run: runNetworkStatus,
}

func init() {
	networkCmd.AddCommand(networkStatusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// networkStatusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// networkStatusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runNetworkStatus(cmd *cobra.Command, args []string) {
	fmt.Println("network status called")
	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		fmt.Printf("Error getting flag 'json': %v`n", err)
	}

	if json {
		fmt.Println("Output in 'JSON' format")
	}

	fmt.Printf("API: %s:%d\n", config.API.Address, config.API.Port)
}
