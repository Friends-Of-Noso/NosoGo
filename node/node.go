package node

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/libp2p/go-libp2p"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	"github.com/Friends-Of-Noso/NosoGo/dns"
	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
	"github.com/Friends-Of-Noso/NosoGo/store"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	cNodePortFlag = "node-port"
)

type Node struct {
	// cmd                   *cobra.Command
	ctx                   context.Context
	quit                  *chan struct{}
	wg                    *sync.WaitGroup
	peer                  *pb.PeerInfo
	p2pHost               host.Host
	pubSub                *pubsub.PubSub
	topics                PubSubTopics
	subscriptions         PubSubSubscription
	privateKey            crypto.PrivKey
	publicKey             crypto.PubKey
	sm                    *store.StorageManager
	dnsPeers              *pb.PeerList
	seedPeers             *pb.PeerList
	nodePeers             *pb.PeerList
	dns                   *dns.DNS
	dnsAddress            string
	dnsPort               int32
	statusStorage         *store.Storage[*pb.Status]
	blockStorage          *store.Storage[*pb.Block]
	transactionStorage    *store.Storage[*pb.Transaction]
	bannedPeerInfoStorage *store.Storage[*pb.PeerInfo]
	status                *pb.Status
	seed                  string // This needs to go away
	// dht           *dht.IpfsDHT
}

func NewNode(
	// cmd *cobra.Command,
	ctx context.Context,
	// cancel context.CancelFunc,
	quit *chan struct{},
	wg *sync.WaitGroup,
	address string,
	port int32,
	privKey string,
	pubKey string,
	mode string,
	dnsAddress string,
	dnsPort int32,
	config *cfg.Config,
	seed string, // This needs to go away
) (*Node, error) {
	// if !utils.FileExists(configPath) {
	// 	return nil, fmt.Errorf("could not find config ", configPath)
	// }

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
	}

	sm, err := store.NewStorageManager(config.GetDatabasePath())
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

		if err := config.WriteConfig(); err != nil {
			log.Fatalf("could not save config structure: %v", err)
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

	nodeAddress, err := utils.ResolveToMultiaddr(address, port)
	if err != nil {
		log.Fatalf("unable to resolve to multiaddr: %v", err)
	}

	host, err := libp2p.New(
		libp2p.ListenAddrs(nodeAddress),
		libp2p.Identity(privateKey),
		// This as implication on DHT
		libp2p.DisableRelay(),
	)
	if err != nil {
		return nil, err
	}

	peer := &pb.PeerInfo{
		Address: address,
		Port:    port,
		Mode:    mode,
		Id:      host.ID().String(),
	}

	// If we don't have relaying, then the DHT is only good for Seeds/SuperNodes
	//
	// Create DHT for peer discovery
	// dht, err := dht.New(ctx, host)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to create DHT: %w", err)
	// }

	// Create pubsub for block propagation
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub: %w", err)
	}

	return &Node{
		// cmd:                   cmd,
		ctx:                   ctx,
		quit:                  quit,
		wg:                    wg,
		peer:                  peer,
		dnsAddress:            dnsAddress,
		dnsPort:               dnsPort,
		p2pHost:               host,
		pubSub:                ps,
		topics:                make(PubSubTopics, 0),
		subscriptions:         make(PubSubSubscription, 0),
		privateKey:            privateKey,
		publicKey:             publicKey,
		sm:                    sm,
		dnsPeers:              &pb.PeerList{},
		seedPeers:             &pb.PeerList{},
		nodePeers:             &pb.PeerList{},
		status:                &pb.Status{},
		statusStorage:         sm.StatusStorage(),
		blockStorage:          sm.BlockStorage(),
		transactionStorage:    sm.TransactionStorage(),
		bannedPeerInfoStorage: sm.PeerInfoStorage(),
		seed:                  seed,
		// dht:           dht,
	}, nil
}

func (n *Node) Start() {
	defer n.wg.Done()
	log.Infof("node starting in mode: %s", n.peer.Mode)

	if err := n.startUp(); err != nil {
		log.Errorf("failed calling startUp", err)
		n.Shutdown()
		return
	}

	switch n.peer.Mode {
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
	// n.cancel()
	close(*n.quit)

	// See if there's custom  cleanup for each mode
	switch n.peer.Mode {
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
	// log.Info("waiting for threads to finish...")
	// n.wg.Wait()

	// Close the database
	log.Info("closing database...")
	n.sm.Close()

	log.Info("exiting node")
}

// Loads the status
func (n *Node) loadStatus() error {
	if err := n.statusStorage.Get(pb.StatusKey, n.status); err != nil {
		return err
	}
	return nil
}

// Loads the status
func (n *Node) saveStatus() error {
	if err := n.statusStorage.Put(pb.StatusKey, n.status); err != nil {
		return err
	}
	return nil
}

// Propagates a new block
func (n *Node) propagateNewBlock(newblock *pb.BlocksSubscriptionNewBlock) error {
	// Create network message
	msg := &pb.BlocksSubscriptionMessage{
		Payload: &pb.BlocksSubscriptionMessage_NewBlock{
			NewBlock: newblock,
		},
	}

	// Serialize the message
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}

	// Publish to the network
	return n.topics[BLOCKS_SUB].Publish(n.ctx, data)
}

// Propagates a new block
func (n *Node) propagateNewTransactions(newTransactions *pb.BlocksSubscriptionNewTransactions) error {
	// Create network message
	msg := &pb.BlocksSubscriptionMessage{
		Payload: &pb.BlocksSubscriptionMessage_NewTransactions{
			NewTransactions: newTransactions,
		},
	}

	// Serialize the message
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Publish to the network
	return n.topics[BLOCKS_SUB].Publish(n.ctx, data)
}

func (n *Node) startUp() error {

	log.Info("checking blockchain")

	// return fmt.Errorf("test exit: %d", 1)

	// Check if we have any blocks
	blocksCount, err := n.blockStorage.Count()
	if err != nil {
		return err
	}

	// Check if we have any transactions
	transactionsCount, err := n.transactionStorage.Count()
	if err != nil {
		return err
	}

	err = n.loadStatus()
	if errors.Is(err, leveldb.ErrNotFound) {
		log.Debug("loadStatus -> ErrNotFound")
		// We have blocks or transactions but status is missing
		if blocksCount > 0 || transactionsCount > 0 {
			log.Info("blockchain seems to be out of sync: attempting a re-scan")
			// Attempt to rescan the database
			if err := n.reScanBlockChain(); err != nil {
				// Ok, data is well corrupted
				return err
			}
		} else {
			if err := n.initiateBlockChain(); err != nil {
				return err
			}
		}
	} else {
		log.Debugf("loadStatus -> status: %d, '%s'", n.status.LastBlock, n.status.LastHash)
		if n.status.LastBlock != blocksCount-1 {
			log.Info("blockchain seems to be out of sync: attempting a re-scan")
			// Attempt to rescan the database
			if err := n.reScanBlockChain(); err != nil {
				// Ok, data is well corrupted
				return err
			}
		}
	}
	return nil
}

// Initializes the block chain with block zero and sets status
func (n *Node) initiateBlockChain() error {
	log.Info("no blockchain found, creating it")
	blockZero := pb.NewBlockZero()

	blockZeroKey := n.sm.BlockKey(blockZero.Height)
	if err := n.blockStorage.Put(blockZeroKey, blockZero); err != nil {
		return err
	}
	status := &pb.Status{
		LastBlock: blockZero.Height,
		LastHash:  blockZero.Hash,
	}
	if err := n.statusStorage.Put(pb.StatusKey, status); err != nil {
		return err
	}

	return nil
}

// Re-scans the database and tries to recover status
func (n *Node) reScanBlockChain() error {
	log.Info("re-scanning the blockchain")
	blocks, err := n.blockStorage.ListValues(func() *pb.Block {
		return &pb.Block{}
	})
	if err != nil {
		return err
	}

	// Check that all blocks are sequential
	var (
		height   uint64 = 0
		previous        = pb.NewBlockZero()
	)

	for _, block := range blocks {
		if block.Height != height {
			return fmt.Errorf("mismatched block height, expected %d, got %d", height, block.Height)
		}

		if height == 0 && block.PreviousHash != previous.PreviousHash {
			return fmt.Errorf("chain is broken: block %d does not have the the correct previous hash", block.Height)
		}
		if height != 0 && block.PreviousHash != previous.Hash {
			return fmt.Errorf("chain is broken: block %d does not have the the correct previous hash", block.Height)
		}

		// ????????????????????????
		// if err := n.loadStatus(); err != nil {
		// 	log.Error("reScanBlockChain.loadStatus", err)
		// 	n.Shutdown()
		// 	return err
		// }

		n.status.LastBlock = block.Height
		n.status.LastHash = block.Hash

		if err := n.saveStatus(); err != nil {
			log.Error("reScanBlockChain.saveStatus", err)
			n.Shutdown()
			return err
		}

		previous = block
		height++
	}

	// Check for orphaned transactions
	transactions, err := n.transactionStorage.ListValues(func() *pb.Transaction {
		return &pb.Transaction{}
	})
	if err != nil {
		return fmt.Errorf("could not retrieve transactions: %w", err)
	}
	for _, transaction := range transactions {
		key := n.sm.BlockKey(transaction.BlockHeight)
		ok, err := n.blockStorage.Has(key)
		if err != nil {
			return fmt.Errorf("error querying for block: %w", err)
		}
		if !ok {
			return fmt.Errorf("found an orphan transaction: '%s'", transaction.Hash)
		}
	}
	return nil
}
