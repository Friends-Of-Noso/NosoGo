package commands

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	"github.com/Friends-Of-Noso/NosoGo/api"
)

var blocksStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Blocks Status",
	//Long:  `Initializes the configuration file.`,
	Run: runBlocksStatus,
}

func init() {
	blocksCmd.AddCommand(blocksStatusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blocksStatusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blocksStatusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runBlocksStatus(cmd *cobra.Command, args []string) {
	fmt.Println("blocks status called")
	json, err := cmd.Flags().GetBool("json")
	if err != nil {
		fmt.Printf("Error getting flag 'json': %v`n", err)
	}

	if json {
		fmt.Println("Output in 'JSON' format")
	}
	path := path.Join(api.APIBasePath, api.APIBlocksStatus)
	url := fmt.Sprintf("http://%s:%d/%s", config.API.Address, config.API.Port, path)
	fmt.Printf("Endpoint: '%s'\n", url)
}
