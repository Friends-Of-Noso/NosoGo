package commands

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/node"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	cNodeAddressFlag = "node-address"
	cNodeAddress     = "node.address"
	cNodePortFlag    = "node-port"
	cNodePort        = "node.port"
	cNodeModeFlag    = "node-mode"
	cNodeMode        = "node.mode"

	cDNSAddressFlag = "dns-address"
	cDNSAddress     = "dns.address"
	cDNSPortFlag    = "dns-port"
	cDNSPort        = "dns.port"

	// cSeedFlag = "seed" // Needs removal in production
)

// nodeCmd represents the node command
var (
	nodeCmd = &cobra.Command{
		Use:   "node",
		Short: "Starts the node",
		//Long:  `Starts the web server.`,
		Example: `  # Various modes
  $ nosogod node --node-mode "dns"
  $ nosogod node --node-mode "seed"
  $ nosogod node --node-mode "superseed"
  $ nosogod node --node-mode "node" # This is the default mode

  # Using different node address/port combinations
  $ nosogod node --node-address "localhost" --node-port 1234
  $ nosogod node --node-address "127.0.0.1" --node-port 4321

  # In mode DNS using different address/port combinations
  $ nosogod node --node.mode "dns" --dns-address "localhost" --dns-port 1234
  $ nosogod node --node.mode "dns" --dns-address "127.0.0.1" --dns-port 4321`,
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

	nodeCmd.Flags().String(cNodeAddressFlag, config.Node.Address, "node address")
	viper.BindPFlag(cNodeAddress, nodeCmd.Flags().Lookup(cNodeAddressFlag))

	nodeCmd.Flags().Int(cNodePortFlag, config.Node.Port, "node port")
	viper.BindPFlag(cNodePort, nodeCmd.Flags().Lookup(cNodePortFlag))

	modeHelp := fmt.Sprintf("node mode: '%s', '%s', '%s', '%s'", cfg.NodeModeDNS, cfg.NodeModeSeed, cfg.NodeModeSuperNode, cfg.NodeModeNode)
	nodeCmd.Flags().String(cNodeModeFlag, config.Node.Mode, modeHelp)
	viper.BindPFlag(cNodeMode, nodeCmd.Flags().Lookup(cNodeModeFlag))
	err := nodeCmd.RegisterFlagCompletionFunc(cNodeModeFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{cfg.NodeModeDNS, cfg.NodeModeSeed, cfg.NodeModeSuperNode, cfg.NodeModeNode}, cobra.ShellCompDirectiveNoFileComp
		})
	if err != nil {
		log.Error("Error registering flag completion function", err)
	}

	nodeCmd.Flags().String(cDNSAddressFlag, config.DNS.Address, "dns address")
	viper.BindPFlag(cDNSAddressFlag, nodeCmd.Flags().Lookup(cDNSAddressFlag))

	nodeCmd.Flags().Int(cDNSPortFlag, config.DNS.Port, "dns port")
	viper.BindPFlag(cDNSPortFlag, nodeCmd.Flags().Lookup(cDNSPortFlag))

	nodeCmd.Flags().StringVarP(&seed, "seed", "s", "", "seed to connect")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runNode(cmd *cobra.Command, args []string) {
	log.Debug("node called")

	if !cfg.ValidModes[config.Node.Mode] {
		log.Fatalf("wrong node mode: '%s'", config.Node.Mode)
	}

	if utils.FileExists(viper.ConfigFileUsed()) {
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

		nodeAddressConfig := getFlagString(cmd, cNodeAddressFlag)
		log.Debugf("nodeAddress: %s", nodeAddressConfig)
		nodePortConfig := getFlagInt(cmd, cNodePortFlag)
		log.Debugf("nodePort: %d", nodePortConfig)
		nodeAddress, err := resolveToMultiaddr(nodeAddressConfig, nodePortConfig)
		if err != nil {
			log.Fatalf("unable to resolve to multiaddr: %v", err)
		}
		nodeMode := getFlagString(cmd, cNodeModeFlag)
		log.Debugf("nodeMode: %s", nodeMode)

		privKey := config.Node.PrivateKey
		log.Debugf("privKey: %s", privKey)

		pubKey := config.Node.PublicKey
		log.Debugf("pubKey: %s", pubKey)

		dnsAddressConfig := getFlagString(cmd, cDNSAddressFlag)
		log.Debugf("dnsAddrs: %s", dnsAddressConfig)
		dnsPortConfig := getFlagInt(cmd, cDNSPortFlag)
		log.Debugf("dnsPort: %d", dnsPortConfig)
		dnsAddress, err := resolveToString(dnsAddressConfig, dnsPortConfig)
		if err != nil {
			log.Fatalf("could not resolve to string: %v", err)
		}

		node, err := node.NewNode(
			ctx,
			cancel,
			&wg,
			cmd,
			nodeAddress,
			nodePortConfig,
			privKey,
			pubKey,
			nodeMode,
			dnsAddress,
			dnsPortConfig,
			config.GetConfigFolder(),
			config.GetDatabaseFolder(),
			seed,
		)
		if err != nil {
			log.Fatalf("error creating node: %v", err)
		}

		wg.Add(1)
		go node.Start()
		// Block here until we receive a termination signal
		sig := <-sigChan
		// Print a new line after the "^C" or "^\"
		if sig == syscall.SIGINT || sig == syscall.SIGQUIT || sig == syscall.SIGKILL {
			fmt.Println()
		}
		log.Infof("received signal '%s'", sig)

		// Node shutdown cancels the context to signal goroutines to stop
		node.Shutdown()

	} else {
		log.Fatalf("cannot find config file '%s', please run the 'init' command first", viper.ConfigFileUsed())
	}
}

func GetNodePortFlag() string {
	return cNodePortFlag
}

func getFlagInt(cmd *cobra.Command, flag string) int {
	flagValue, err := cmd.Flags().GetInt(flag)
	if err != nil {
		log.Fatalf("cannot retrieve flag '%s': %v", flag, err)
	}
	return flagValue
}

func getFlagString(cmd *cobra.Command, flag string) string {
	flagValue, err := cmd.Flags().GetString(flag)
	if err != nil {
		log.Fatalf("cannot retrieve flag '%s': %v", flag, err)
	}
	return flagValue
}

func resolveToMultiaddr(address string, port int) (multiaddr.Multiaddr, error) {
	ips, err := net.LookupIP(address)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("failed to resolve address %s: %v", address, err)
	}

	// Try to pick the first IPv4 (if any)
	var ip net.IP
	for _, candidate := range ips {
		if candidate.To4() != nil {
			ip = candidate
			break
		}
	}

	if ip == nil {
		return nil, fmt.Errorf("no IPv4 address found for %s", address)
	}

	// Now build the multiaddr using the resolved IP
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip.String(), port))
}

func resolveToString(address string, port int) (string, error) {
	ips, err := net.LookupIP(address)
	if err != nil || len(ips) == 0 {
		return "", fmt.Errorf("failed to resolve address %s: %v", address, err)
	}

	// Try to pick the first IPv4 (if any)
	var ip net.IP
	for _, candidate := range ips {
		if candidate.To4() != nil {
			ip = candidate
			break
		}
	}

	if ip == nil {
		return "", fmt.Errorf("no IPv4 address found for %s", address)
	}

	// Now build the multiaddr using the resolved IP
	return fmt.Sprintf("%s:%d", ip.String(), port), nil
}
