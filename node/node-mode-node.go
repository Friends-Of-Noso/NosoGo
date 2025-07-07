package node

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func (n *Node) runModeNode() {
	log.Debug("entering runModeNode")

	log.Infof("node(node): Listening on %s/p2p/%s", n.p2pHost.Addrs()[0], n.p2pHost.ID())
	log.Debugf("node ID: %s", n.p2pHost.ID())
	for key, value := range n.p2pHost.Addrs() {
		log.Debugf("address: %d, %s", key, value)
	}

	// TODO: This needs to be changed to connect to a list of seeds in production
	if n.seed != "" {
		// connect to seed
		log.Infof("connecting to seed: %s", n.seed)

		// targetAddr, err := multiaddr.NewMultiaddr("/ip4/10.42.0.101/tcp/45050/p2p/12D3KooWKXcHejD288cQi32oqGR3aXEgY2sP3MAgpzwQ7V95CsNt")
		if targetAddr, err := multiaddr.NewMultiaddr(n.seed); err != nil {
			log.Error("invalid seed multiaddr", err)
		} else {
			if peerInfo, err := peer.AddrInfoFromP2pAddr(targetAddr); err != nil {
				log.Error("failed to get peer info", err)
			} else {
				if err := n.p2pHost.Connect(n.ctx, *peerInfo); err != nil {
					log.Error("failed to connect to seed", err)
				} else {
					log.Debugf("connected to seed '%s'", peerInfo.String())
					// Should send a GetPeersRequest message now
				}
			}
		}
	}

	// Bootstrap DHT
	// if err := n.dht.Bootstrap(n.ctx); err != nil {
	// 	log.Error("failed to bootstrap DHT", err)
	// 	n.Shutdown()
	// }

	// Join blocks topic
	blockTopic, err := n.pubSub.Join(BLOCKS_SUB)
	if err != nil {
		log.Error("failed to join blocks topic", err)
		n.Shutdown()
	}
	n.topics[BLOCKS_SUB] = blockTopic

	// Subscribe to block topic
	blockSub, err := blockTopic.Subscribe()
	if err != nil {
		log.Error("failed to subscribe to blocks topic", err)
		n.Shutdown()
	}

	// Handle incoming blocks
	n.wg.Add(1)
	go n.handleBlocksTopic(blockSub)

	n.subscriptions[BLOCKS_SUB] = blockSub

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	height := 0

	for {
		select {
		case <-n.ctx.Done():
			log.Debug("node.start() exiting")
			return
		case <-ticker.C:
			// if config.LogLevel != "debug" {
			// 	continue
			// }
			// TODO This needs to go away on the real thing.
			if n.seed != "" {
				continue
			}
			block := &pb.Block{
				Height:       uint64(height),
				PreviousHash: "BPreviousHash",
				Timestamp:    time.Now().Unix(),
			}
			block.SetHash()
			// Store new block
			blockKey := n.sm.BlockKey(block.Height)
			blockStorage := n.sm.BlockStorage()
			if err := blockStorage.Put(blockKey, block); err != nil {
				log.Error("could not store block on database", err)
				continue
			}

			transaction := &pb.Transaction{
				BlockHeight: block.Height,
				Type:        "COINBASE",
				Timestamp:   time.Now().Unix(),
				PubKey:      "badbeef",
				Verify:      "badbeef",
				Sender:      "COINBASE",
				Receiver:    "NReceiver",
				Amount:      100_000_000, // Coin has 8 decimals
			}
			transaction.SetHash()

			// Store new transaction
			transactionKey := n.sm.TransactionKey(transaction.BlockHeight, transaction.Hash)
			transactionStorage := n.sm.TransactionStorage()
			if err := transactionStorage.Put(transactionKey, transaction); err != nil {
				log.Error("could not store transaction on database", err)
				continue
			}

			newBlock := &pb.NewBlock{
				Block: block,
				Transactions: []*pb.Transaction{
					transaction,
				},
			}
			height++
			log.Debugf(
				"propagating a block %d, %s, %s",
				block.Height,
				block.Hash,
				block.PreviousHash,
			)
			n.propagateNewBlock(newBlock)

			transaction.BlockHeight = 0
			transaction.Type = "spend"
			transaction.Timestamp = time.Now().Unix()
			transaction.Sender = "NSender"
			transaction.Amount = 10_000_000_000
			transaction.SetHash()
			transactionKey = n.sm.TransactionKey(transaction.BlockHeight, transaction.Hash)
			pendingTransactionStorage := n.sm.PendingTransactionStorage()
			if err := pendingTransactionStorage.Put(transactionKey, transaction); err != nil {
				log.Error("could not store transaction on database", err)
				continue
			}
			newTransactions := &pb.NewTransactions{
				Transactions: []*pb.Transaction{
					transaction,
				},
			}

			log.Debugf(
				"propagating a transaction %d, %s, %s",
				transaction.BlockHeight,
				transaction.Hash,
				transaction.Type,
			)
			n.propagateNewTransactions(newTransactions)

		default:
			continue
		}
	}
}

func (n *Node) shutdownNode() {
	// NOTE: Does this have specific shutdown stuff?
}
