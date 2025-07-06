package node

import (
	"context"
	"crypto/rand"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	"github.com/Friends-Of-Noso/NosoGo/dns"
	log "github.com/Friends-Of-Noso/NosoGo/logger"
)

const (
	cNodePortFlag = "node-port"
)

var (
	config = cfg.DefaultConfig()
)

type Node struct {
	cmd           *cobra.Command
	ctx           context.Context
	cancel        context.CancelFunc
	wg            *sync.WaitGroup
	address       multiaddr.Multiaddr
	port          int
	mode          string
	p2pHost       host.Host
	pubSub        *pubsub.PubSub
	topics        PubSubTopics
	subscriptions PubSubSubscription
	privateKey    crypto.PrivKey
	publicKey     crypto.PubKey
	db            *leveldb.DB
	peers         []peer.AddrInfo
	dns           *dns.DNS
	dnsAddress    string
	dnsPort       int
	seed          string // This needs to go away
	// dht           *dht.IpfsDHT
}

func NewNode(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	cmd *cobra.Command,
	address multiaddr.Multiaddr,
	port int,
	privKey string,
	pubKey string,
	mode string,
	dnsAddress string,
	dnsPort int,
	configPath string,
	dbPath string,
	seed string, // This needs to go away
) (*Node, error) {
	// TODO: This entire key thing needs a rethink!!
	var (
		privateKey    crypto.PrivKey
		publicKey     crypto.PubKey
		configPrivKey string
		configPubKey  string
		err           error
	)

	err = checkPort(port, cNodePortFlag, cfg.DefaultNodePort)
	if err != nil {
		log.Fatalf("%v", err)
		os.Exit(1)

	}

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(privKey, cfg.DefaultNodeKey) || strings.HasPrefix(pubKey, cfg.DefaultNodeKey) {
		// Use the port number as the randomness source.
		// This will always generate the same host ID on multiple executions, if the same port number is used.
		// Never do this in production code.
		r := rand.Reader

		// Creates a new RSA key pair for this host.
		privateKey, publicKey, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
		if err != nil {
			return nil, err
		}

		privProto, err := crypto.MarshalPrivateKey(privateKey)
		if err != nil {
			return nil, err
		}

		pubProto, err := crypto.MarshalPublicKey(publicKey)
		if err != nil {
			return nil, err
		}

		configPrivKey = crypto.ConfigEncodeKey(privProto)
		configPubKey = crypto.ConfigEncodeKey(pubProto)

		config.Node.PrivateKey = configPrivKey
		config.Node.PublicKey = configPubKey

		if err := cfg.WriteConfig(viper.ConfigFileUsed(), config); err != nil {
			log.Fatalf("could save config structure: %v", err)
		}

	} else {
		privRaw, err := crypto.ConfigDecodeKey(privKey)
		if err != nil {
			return nil, err
		}
		pubRaw, err := crypto.ConfigDecodeKey(pubKey)
		if err != nil {
			return nil, err
		}

		privateKey, err = crypto.UnmarshalPrivateKey(privRaw)
		if err != nil {
			return nil, err
		}
		publicKey, err = crypto.UnmarshalPublicKey(pubRaw)
		if err != nil {
			return nil, err
		}

	}

	host, err := libp2p.New(
		libp2p.ListenAddrs(address),
		libp2p.Identity(privateKey),
		// This as implication on DHT
		libp2p.DisableRelay(),
	)
	if err != nil {
		return nil, err
	}

	// If we don't have relaying, then the DHT is only good for Seeds/SuperNodes
	//
	// Create DHT for peer discovery
	// dht, err := dht.New(ctx, host)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create DHT: %v", err)
	// }

	// Create pubsub for block propagation
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub: %v", err)
	}

	return &Node{
		ctx:           ctx,
		cancel:        cancel,
		wg:            wg,
		cmd:           cmd,
		address:       address,
		port:          port,
		mode:          mode,
		dnsAddress:    dnsAddress,
		dnsPort:       dnsPort,
		p2pHost:       host,
		pubSub:        ps,
		topics:        make(PubSubTopics, 0),
		subscriptions: make(PubSubSubscription, 0),
		privateKey:    privateKey,
		publicKey:     publicKey,
		db:            db,
		peers:         make(Peers, 0),
		seed:          seed,
		// dht:           dht,
	}, nil
}

func (n *Node) Start() {
	defer n.wg.Done()
	log.Debugf("node.start() called with mode: %s", n.mode)

	// logLevel, err := n.cmd.Flags().GetString("log-level")
	// if err != nil {
	// 	log.Errorf("error getting flag 'log-level': %v", err)
	// 	n.Shutdown()
	// }
	// log.Debugf("log level: %s", logLevel)

	switch n.mode {
	case cfg.NodeModeDNS:
		n.runModeDNS()
	case cfg.NodeModeSeed:
		n.runModeSeed()
	case cfg.NodeModeSuperNode:
		n.runModeSuperNode()
	case cfg.NodeModeNode:
		n.runModeNode()
	}

}

func (n *Node) Shutdown() {
	log.Info("node shutting down...")

	// Call the Context cancel function
	n.cancel()

	// See if there's custom  cleanup for each mode
	switch n.mode {
	case cfg.NodeModeDNS:
		n.shutdownDNS()
	case cfg.NodeModeSeed:
		n.shutdownSeed()
	case cfg.NodeModeSuperNode:
		n.shutdownSuperNode()
	case cfg.NodeModeNode:
		n.shutdownNode()
	}

	// Wait for all goroutines to finish
	log.Info("waiting for threads to finish...")
	n.wg.Wait()

	// Close the database
	log.Info("closing database...")
	n.db.Close()

	log.Info("exiting")
}

func checkPort(port int, flag string, defaultPort int) error {
	goos := runtime.GOOS

	if goos == "linux" || goos == "darwin" {
		if port < 1024 && os.Geteuid() != 0 {
			return fmt.Errorf("port %d requires root privileges on %s; try a port >= 1024 (e.g. --%s %d)", port, goos, flag, defaultPort)
		}
	}

	return nil
}
