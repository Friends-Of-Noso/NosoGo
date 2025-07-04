package commands

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/node"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	cNodeAddressFlag = "node-address"
	cNodeAddress     = "node.address"
	cNodePortFlag    = "node-port"
	cNodePort        = "node.port"
	cSeedFlag        = "seed"
)

// nodeCmd represents the node command
var (
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Starts the node",
		//Long:  `Starts the web server.`,
		Run: runNode,
	}
	seed string
)

func init() {
	rootCmd.AddCommand(nodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	nodeCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is "+config.GetConfigFile()+")")

	nodeCmd.Flags().String(cNodeAddressFlag, config.Node.Address, "Node address")
	viper.BindPFlag(cNodeAddress, nodeCmd.Flags().Lookup(cNodeAddressFlag))

	nodeCmd.Flags().Int(cNodePortFlag, config.Node.Port, "Node port")
	viper.BindPFlag(cNodePort, nodeCmd.Flags().Lookup(cNodePortFlag))

	nodeCmd.Flags().StringVarP(&seed, "seed", "s", "", "seed to connect")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runNode(cmd *cobra.Command, args []string) {
	log.Debug("node called")
	if utils.FileExists(viper.ConfigFileUsed()) {
		log.Debugf("Node Address: '%s'", config.Node.Address)
		log.Debugf("Node Port: %d", config.Node.Port)

		// Create a cancellable context.
		ctx, cancel := context.WithCancel(context.Background())

		// Create a channel to receive OS signals.
		sigChan := make(chan os.Signal, 1)

		// Notify on all relevant Windows and Unix signals.
		signal.Notify(sigChan,
			// Windows signals
			os.Interrupt,    // Ctrl+C
			syscall.SIGTERM, // Termination signal
			syscall.SIGABRT, // Abort signal (Windows and Unix)

			// Unix/Linux signals
			syscall.SIGHUP,  // Hangup detected (terminal or process dies)
			syscall.SIGQUIT, // Quit from keyboard (Ctrl+\ on Unix)
			syscall.SIGINT,  // Interrupt from keyboard (Ctrl+C on Unix)
			// syscall.SIGTSTP, // Stop typed at terminal (Ctrl+Z on Unix)
			// syscall.SIGUSR1, // User-defined signal 1
			// syscall.SIGUSR2, // User-defined signal 2
		)

		var wg sync.WaitGroup

		// if seed != "" {
		// 	log.Debugf("Got seed: %s", seed)
		// 	if strings.Contains(seed, ":") {
		// 		bits := strings.Split(seed, ":")
		// 		seed = fmt.Sprintf("/ip4/%s/tcp/%s", bits[0], bits[1])
		// 	} else {
		// 		seed = fmt.Sprintf("/ip4/%s/tcp/%d", seed, cfg.DefaultNodePort)
		// 	}
		// 	log.Debugf("Seed multiaddr: %s", seed)
		// }

		node, err := node.NewNode(
			ctx,
			cancel,
			&wg,
			config.Node.Address,
			config.Node.Port,
			config.GetConfigFolder(),
			config.GetDatabaseFolder(),
			seed,
		)
		if err != nil {
			log.Fatalf("Error creating node: %v", err)
		}

		wg.Add(1)
		go node.Start()
		// Block here until we receive a termination signal
		sig := <-sigChan
		// Print a new line after the "^C" or "^\"
		if sig == syscall.SIGINT || sig == syscall.SIGQUIT || sig == syscall.SIGKILL {
			fmt.Println()
		}
		log.Infof("Received signal '%s'", sig)

		// Node shutdown cancels the context to signal goroutines to stop
		log.Debug("Shutting down the node...")
		node.Shutdown()

	} else {
		log.Fatalf("Cannot find config file '%s', please run the 'init' command first", viper.ConfigFileUsed())
	}
}
