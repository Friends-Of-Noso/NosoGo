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

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/node"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	cConfigFlag = "config"

	cLogLevelFlag = "log-level"
	cLogLevel     = "log_level"

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
	cfgFile string

	config = cfg.DefaultConfig()

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
	// cobra.OnInitialize(nodeInitConfig)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	nodeCmd.Flags().StringVarP(&cfgFile, cConfigFlag, "c", config.GetConfigFile(), "config file")
	// if err := viper.BindPFlag(cConfigFlag, nodeCmd.Flags().Lookup(cConfigFlag)); err != nil {
	// 	fmt.Fprintf(os.Stderr, "error binding flag '%s': %v", cConfigFlag, err)
	// 	os.Exit(1)
	// }

	logLevelHelp := fmt.Sprintf(
		// "log level: '%s', '%s', '%s', '%s'",
		"log level: '%s', '%s'",
		cfg.LogLevelInfo,
		// cfg.LogLevelWarn,
		// cfg.LogLevelError,
		cfg.LogLevelDebug)

	nodeCmd.Flags().StringP(cLogLevelFlag, "l", config.LogLevel, logLevelHelp)

	if err := viper.BindPFlag(cLogLevel, nodeCmd.Flags().Lookup(cLogLevelFlag)); err != nil {
		fmt.Fprintf(os.Stderr, "error binding flag '%s': %v", cLogLevelFlag, err)
		os.Exit(1)
	}

	err := nodeCmd.RegisterFlagCompletionFunc(cLogLevelFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				cfg.LogLevelInfo,
				// cfg.LogLevelWarn,
				// cfg.LogLevelError,
				cfg.LogLevelDebug,
			}, cobra.ShellCompDirectiveNoFileComp
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error registering flag completion function: %v", err)
		os.Exit(1)
	}

	nodeCmd.Flags().String(cNodeAddressFlag, config.Node.Address, "node address")
	viper.BindPFlag(cNodeAddress, nodeCmd.Flags().Lookup(cNodeAddressFlag))

	nodeCmd.Flags().Int32(cNodePortFlag, config.Node.Port, "node port")
	viper.BindPFlag(cNodePort, nodeCmd.Flags().Lookup(cNodePortFlag))

	modeHelp := fmt.Sprintf("node mode: '%s', '%s', '%s', '%s'", cfg.NodeModeDNS, cfg.NodeModeSeed, cfg.NodeModeSuperNode, cfg.NodeModeNode)

	nodeCmd.Flags().String(cNodeModeFlag, config.Node.Mode, modeHelp)

	viper.BindPFlag(cNodeMode, nodeCmd.Flags().Lookup(cNodeModeFlag))
	err = nodeCmd.RegisterFlagCompletionFunc(cNodeModeFlag,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{cfg.NodeModeDNS, cfg.NodeModeSeed, cfg.NodeModeSuperNode, cfg.NodeModeNode}, cobra.ShellCompDirectiveNoFileComp
		})
	if err != nil {
		log.Error("Error registering flag completion function", err)
	}

	nodeCmd.Flags().String(cDNSAddressFlag, config.DNS.Address, "dns address")
	viper.BindPFlag(cDNSAddressFlag, nodeCmd.Flags().Lookup(cDNSAddressFlag))

	nodeCmd.Flags().Int32(cDNSPortFlag, config.DNS.Port, "dns port")
	viper.BindPFlag(cDNSPortFlag, nodeCmd.Flags().Lookup(cDNSPortFlag))

	nodeCmd.Flags().StringVarP(&seed, "seed", "s", "", "seed to connect")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runNode(cmd *cobra.Command, args []string) {
	// log.Debug("node called")

	nodeInitConfigAndLogs(cmd)

	// Create a cancellable context.
	ctx, cancel := context.WithCancel(context.Background())

	// Create a channel to receive OS signals.
	sigChan := make(chan os.Signal, 1)
	// Manual quit channel
	quit := make(chan struct{})

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

	nodeAddressConfig := getFlagString(cmd, cNodeAddressFlag)
	log.Debugf("nodeAddress: %s", nodeAddressConfig)
	nodePortConfig := getFlagInt32(cmd, cNodePortFlag)
	log.Debugf("nodePort: %d", nodePortConfig)
	// nodeAddress, err := resolveToMultiaddr(nodeAddressConfig, nodePortConfig)
	// if err != nil {
	// 	log.Fatalf("unable to resolve to multiaddr: %v", err)
	// }
	nodeMode := getFlagString(cmd, cNodeModeFlag)
	log.Debugf("nodeMode: %s", nodeMode)

	privKey := config.Node.PrivateKey
	log.Debugf("privKey: %s", privKey)

	pubKey := config.Node.PublicKey
	log.Debugf("pubKey: %s", pubKey)

	dnsAddressConfig := getFlagString(cmd, cDNSAddressFlag)
	log.Debugf("dnsAddrs: %s", dnsAddressConfig)
	dnsPortConfig := getFlagInt32(cmd, cDNSPortFlag)
	log.Debugf("dnsPort: %d", dnsPortConfig)
	dnsAddress, err := utils.ResolveToString(dnsAddressConfig, dnsPortConfig)
	if err != nil {
		log.Fatalf("could not resolve to string: %v", err)
	}

	node, err := node.NewNode(
		// cmd,
		ctx,
		&quit,
		&wg,
		nodeAddressConfig,
		nodePortConfig,
		privKey,
		pubKey,
		nodeMode,
		dnsAddress,
		dnsPortConfig,
		config,
		seed,
	)
	if err != nil {
		log.Fatalf("error creating node: %v", err)
	}

	wg.Add(1)
	go node.Start()

	// Block here until we receive a termination signal
	select {
	case sig := <-sigChan:
		// Print a new line after the "^C" or "^\"
		if sig == syscall.SIGINT || sig == syscall.SIGQUIT || sig == syscall.SIGKILL {
			fmt.Println()
		}
		log.Debugf("received signal '%s'", sig)
		// Node shutdown cleans up it's dependencies
		node.Shutdown()
	case <-quit:
		log.Debugf("received internal shutdown")
	}

	// cancel
	cancel()

	wg.Wait()
}

func nodeInitConfigAndLogs(cmd *cobra.Command) {
	if cfgFile != "" && !utils.FileExists(cfgFile) {
		// TODO: Output usage
		fmt.Fprintf(os.Stderr, "could not find config file '%s'\n", cfgFile)
		os.Exit(1)
	}

	if cfgFile == "" && !utils.FileExists(config.GetConfigFile()) {
		// TODO: Output usage
		fmt.Fprintf(os.Stderr, "could not find config file '%s'\n", config.GetConfigFile())
		os.Exit(1)
	}

	if cfgFile == "" {
		cfgFile = config.GetConfigFile()
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(cfgFile)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())

		err := viper.Unmarshal(config)
		if err != nil {
			fmt.Printf("could not unmarshal config: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("viper could not read config: %v\n", err)
		os.Exit(1)
	}
	//viper.AutomaticEnv() // read in environment variables that match

	if !cfg.ValidLogLevels[config.LogLevel] {
		fmt.Printf("\nError: wrong log level: '%s'.\nPlease check your config file or the usage below.\n", config.LogLevel)
		fmt.Println(cmd.UsageString())
		os.Exit(1)
	}

	if !cfg.ValidModes[config.Node.Mode] {
		fmt.Printf("\nError: wrong node mode: '%s'.\nPlease check your config file or the usage below.\n", config.Node.Mode)
		fmt.Println(cmd.UsageString())
		os.Exit(1)
	}

	log.SetFileAndLevel(config.GetLogFile(), config.LogLevel)
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

func getFlagInt32(cmd *cobra.Command, flag string) int32 {
	flagValue, err := cmd.Flags().GetInt32(flag)
	if err != nil {
		log.Fatalf("cannot retrieve flag '%s': %v", flag, err)
	}
	return flagValue
}

func getFlagInt64(cmd *cobra.Command, flag string) int64 {
	flagValue, err := cmd.Flags().GetInt64(flag)
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
