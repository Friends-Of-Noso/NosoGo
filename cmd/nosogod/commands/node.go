package commands

import (
	"github.com/spf13/cobra"
	//log "github.com/EIYARO-Project/core-stratumd/logger"
	//"github.com/EIYARO-Project/core-stratumd/server"
	//fs "github.com/EIYARO-Project/core-stratumd/utils"
)

// nodeCmd represents the node command
var (
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Starts the node",
		//Long:  `Starts the web server.`,
		Run: runNode,
	}
)

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	/*nodeCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is "+config.GetConfigFile()+")")

	nodeCmd.Flags().String(cAPIAddressFlag, config.API.Address, "API address")
	viper.BindPFlag("api.address", nodeCmd.Flags().Lookup(cAPIAddressFlag))

	nodeCmd.Flags().Int(cAPIPortFlag, config.API.Port, "API port")
	viper.BindPFlag("api.port", nodeCmd.Flags().Lookup(cAPIPortFlag))

	nodeCmd.Flags().String(cAPIAccessTokenFlag, config.API.AccessToken, "API access token")
	viper.BindPFlag("api.access_token", nodeCmd.Flags().Lookup(cAPIAccessTokenFlag))

	nodeCmd.Flags().String(cStratumAddressFlag, config.Stratum.Address, "Stratum address")
	viper.BindPFlag("stratum.address", nodeCmd.Flags().Lookup(cStratumAddressFlag))

	nodeCmd.Flags().Int(cStratumPortFlag, config.Stratum.Port, "Stratum port")
	viper.BindPFlag("stratum.port", nodeCmd.Flags().Lookup(cStratumPortFlag))

	nodeCmd.Flags().Int(cStratumMaxConnectionsFlag, config.Stratum.MaxConnections, "Stratum Max Connections")
	viper.BindPFlag("stratum.max_connections", nodeCmd.Flags().Lookup(cStratumMaxConnectionsFlag))*/

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runNode(cmd *cobra.Command, args []string) {
	/*log.Debug("serve called")
	if fs.FileExists(viper.ConfigFileUsed()) {
		log.Debugf("Stratum Address: '%s'", config.Stratum.Address)
		log.Debugf("Stratum Port: %d", config.Stratum.Port)
		log.Debugf("Stratum Max Connections: %d", config.Stratum.MaxConnections)
		log.Debugf("Stratum Difficulty: %d", config.Stratum.Difficulty)
		log.Debugf("API Address: '%s'", config.API.Address)
		log.Debugf("API Port: %d", config.API.Port)
		log.Debugf("API Access Token: '%s'", config.API.AccessToken)

		// TODO: Make sure we have sane values for address and port

		var waitGroup sync.WaitGroup

		// Stratum Server
		ctx := context.Background()
		server := server.NewServer(
			ctx, config.Stratum.Address,
			config.Stratum.Port,
			config.Stratum.MaxConnections,
			config.Stratum.Difficulty,
			config.API.Address,
			config.API.Port,
			config.API.AccessToken)

		// Start Polling the node's API
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			server.PollNode()
		}()

		// Start the server
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			if err := server.ListenAndServe(); err != nil {
				log.Fatalf("Error on ListenAndServe: %s", err)
			}
		}()

		// Wait for all servers to exit
		waitGroup.Wait()
	} else {
		log.Fatalf("Cannot find config file '%s', please run the 'init' command first", viper.ConfigFileUsed())
	}*/
}
