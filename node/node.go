package node

import (
	"bufio"
	"context"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
	"github.com/Friends-Of-Noso/NosoGo/utils"
	"github.com/libp2p/go-libp2p"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/syndtr/goleveldb/leveldb"
	"google.golang.org/protobuf/proto"
)

type Node struct {
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	address    string
	port       int
	host       host.Host
	dht        *dht.IpfsDHT
	pubsub     *pubsub.PubSub
	blockTopic *pubsub.Topic
	privateKey crypto.PrivKey
	publicKey  crypto.PubKey
	db         *leveldb.DB
	peers      []peer.AddrInfo
	seed       string
}

func NewNode(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	address string,
	port int,
	configPath string,
	dbPath string,
	seed string,
) (*Node, error) {
	var (
		privateKey    crypto.PrivKey
		publicKey     crypto.PubKey
		configPrivKey string
		configPubKey  string
		err           error
	)

	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	if utils.FileExists(path.Join(configPath, "keystore")) {
		keyFile, err := os.Open(path.Join(configPath, "keystore"))
		if err != nil {
			return nil, err
		}
		defer keyFile.Close()

		scanner := bufio.NewScanner(keyFile)
		if scanner.Scan() {
			configPrivKey = scanner.Text()
		}
		if scanner.Scan() {
			configPubKey = scanner.Text()
		}

		privRaw, err := crypto.ConfigDecodeKey(configPrivKey)
		if err != nil {
			return nil, err
		}
		pubRaw, err := crypto.ConfigDecodeKey(configPubKey)
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

	} else {
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

		configPrivKey := crypto.ConfigEncodeKey(privProto)
		configPubKey := crypto.ConfigEncodeKey(pubProto)

		keyFile, err := os.Create(path.Join(configPath, "keystore"))
		if err != nil {
			return nil, err
		}
		defer keyFile.Close()

		keyFile.WriteString(fmt.Sprintln(configPrivKey))
		keyFile.WriteString(configPubKey)
	}

	addr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", address, port))
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(
		libp2p.ListenAddrs(addr),
		libp2p.Identity(privateKey),
		libp2p.DisableRelay(),
	)
	if err != nil {
		return nil, err
	}

	// Create DHT for peer discovery
	kdht, err := dht.New(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create DHT: %v", err)
	}

	// Create pubsub for block propagation
	ps, err := pubsub.NewGossipSub(ctx, host)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub: %v", err)
	}

	return &Node{
		ctx:        ctx,
		cancel:     cancel,
		wg:         wg,
		address:    address,
		port:       port,
		host:       host,
		dht:        kdht,
		pubsub:     ps,
		privateKey: privateKey,
		publicKey:  publicKey,
		db:         db,
		peers:      make([]peer.AddrInfo, 0),
		seed:       seed,
	}, nil
}

func (n *Node) Start() {
	log.Debug("Node.Start() called")
	defer n.wg.Done()

	log.Infof("Listening on %s/p2p/%s", n.host.Addrs()[0], n.host.ID())
	log.Debugf("Node ID: %s", n.host.ID())
	for key, value := range n.host.Addrs() {
		log.Debugf("Address: %d, %s", key, value)
	}

	if n.seed != "" {
		// connect to seed
		log.Infof("Connecting to seed: %s", n.seed)
		targetAddr, err := multiaddr.NewMultiaddr(n.seed)
		// targetAddr, err := multiaddr.NewMultiaddr("/ip4/10.42.0.101/tcp/45050/p2p/12D3KooWKXcHejD288cQi32oqGR3aXEgY2sP3MAgpzwQ7V95CsNt")
		if err != nil {
			log.Errorf("invalid seed multiaddr: %v", err)
		} else {
			peerInfo, err := peer.AddrInfoFromP2pAddr(targetAddr)
			if err != nil {
				log.Errorf("failed to get peer info: %v", err)
			} else {
				if err := n.host.Connect(n.ctx, *peerInfo); err != nil {
					log.Error("failed to connect to seed", err)
				} else {
					log.Debugf("Connected to seed '%s'", peerInfo.String())
					// Should send a GetPeersRequest message now
				}
			}
		}
	}

	// Bootstrap DHT
	if err := n.dht.Bootstrap(n.ctx); err != nil {
		log.Errorf("failed to bootstrap DHT: %v", err)
		n.Shutdown()
	}

	// Join block topic
	blockTopic, err := n.pubsub.Join("blocks")
	if err != nil {
		log.Errorf("failed to join blocks topic: %v", err)
		n.Shutdown()
	}
	n.blockTopic = blockTopic

	// Subscribe to block topic
	blockSub, err := blockTopic.Subscribe()
	if err != nil {
		log.Errorf("failed to subscribe to blocks topic: %v", err)
		n.Shutdown()
	}

	// Handle incoming blocks
	n.wg.Add(1)
	go n.handleBlockSubscription(blockSub)

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	height := 0

	for {
		select {
		case <-n.ctx.Done():
			log.Debug("Node.Start() exiting")
			return
		case <-ticker.C:
			if n.seed == "" {
				md := md5.New()
				hash := md.Sum([]byte("MyBlock" + strconv.Itoa(height)))
				prevHash := md.Sum([]byte("MyBlock" + strconv.Itoa(height-1)))
				block := pb.Block{
					Hash:      hex.EncodeToString(hash),
					Height:    uint64(height),
					PrevHash:  hex.EncodeToString(prevHash),
					Timestamp: time.Now().Unix(),
				}
				height++
				log.Debugf("Propagating a block %d, %s, %s", block.Height, block.Hash, block.PrevHash)
				n.propagateBlock(&block)
			}
		default:
			continue
		}
	}
}

func (n *Node) Shutdown() {
	log.Debug("Node.Shutdown() called")
	n.db.Close()
	n.cancel()
}

func (n *Node) handleBlockSubscription(sub *pubsub.Subscription) {
	defer n.wg.Done()
	for {
		select {
		case <-n.ctx.Done():
			log.Debug("Node.handleBlockSubscription() exiting")
			return
		default:
			msg, err := sub.Next(n.ctx)
			if err != nil {
				// log.Error("Error reading subscription", err)
				continue
			}

			// Skip messages from ourselves
			if msg.ReceivedFrom == n.host.ID() {
				continue
			}

			// Unmarshal the message
			networkMsg := &pb.NetworkMessage{}
			if err := proto.Unmarshal(msg.Data, networkMsg); err != nil {
				continue
			}

			// Handle new block message
			if newBlock := networkMsg.GetNewBlock(); newBlock != nil {
				n.handleNewBlock(newBlock.Block)
			}
		}
	}
}

func (n *Node) handleNewBlock(block *pb.Block) {
	log.Infof("Got block(%d): %s, %s", block.Height, block.Hash, block.PrevHash)
}

func (n *Node) propagateBlock(block *pb.Block) error {
	// n.cacheMutex.Lock()
	// if n.blockCache[block.Hash] {
	// 	n.cacheMutex.Unlock()
	// 	return nil
	// }
	// n.blockCache[block.Hash] = true
	// n.cacheMutex.Unlock()

	// Create network message
	msg := &pb.NetworkMessage{
		Payload: &pb.NetworkMessage_NewBlock{
			NewBlock: &pb.NewBlockMessage{
				Block: block,
			},
		},
	}

	// Serialize the message
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal block: %v", err)
	}

	// Publish to the network
	return n.blockTopic.Publish(n.ctx, data)
}
