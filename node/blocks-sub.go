package node

import (
	"fmt"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func (n *Node) handleBlocksTopic(sub *pubsub.Subscription) {
	defer n.wg.Done()
	for {
		select {
		case <-n.ctx.Done():
			log.Debug("node.handleBlockSubscription exiting")
			return
		default:
			msg, err := sub.Next(n.ctx)
			if err != nil {
				// log.Error("Error reading subscription", err)
				continue
			}

			// Skip messages from ourselves
			if msg.ReceivedFrom == n.p2pHost.ID() {
				log.Debug("received my own blocks message")
				continue
			}

			// Unmarshal the message
			networkMsg := &protobuf.BlocksSubMessages{}
			if err := proto.Unmarshal(msg.Data, networkMsg); err != nil {
				continue
			}

			// Handle new message
			if newBlock := networkMsg.GetNewBlock(); newBlock != nil {
				n.handleNewBlock(newBlock.Block)
			}
		}
	}
}

func (n *Node) handleNewBlock(block *protobuf.Block) {
	log.Infof("got block(%d): %s, %s", block.Height, block.Hash, block.PrevHash)
	if len(block.Transactions) > 0 {
		for index, transaction := range block.Transactions {
			log.Infof(
				"transaction %d, %s, %s, %s, %d",
				index,
				transaction.Hash,
				transaction.Sender,
				transaction.Receiver,
				transaction.Amount,
			)
		}
	}
}

func (n *Node) propagateBlock(block *protobuf.Block) error {
	// n.cacheMutex.Lock()
	// if n.blockCache[block.Hash] {
	// 	n.cacheMutex.Unlock()
	// 	return nil
	// }
	// n.blockCache[block.Hash] = true
	// n.cacheMutex.Unlock()

	// Create network message
	msg := &protobuf.BlocksSubMessages{
		Payload: &protobuf.BlocksSubMessages_NewBlock{
			NewBlock: &protobuf.NewBlock{
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
	return n.topics[BLOCKS_SUB].Publish(n.ctx, data)
}
