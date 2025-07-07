package node

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
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

	var height uint64 = 0

	for {
		select {
		case <-n.ctx.Done():
			log.Debug("node.start() exiting")
			return
		case <-ticker.C:
			// if config.LogLevel != "debug" {
			// 	continue
			// }
			n.devPropagateData(height)
			height++

		default:
			continue
		}
	}
}

func (n *Node) shutdownNode() {
	// NOTE: Does this have specific shutdown stuff?
}
