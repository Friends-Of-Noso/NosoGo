package node

import (
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"google.golang.org/protobuf/proto"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
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
				// log.Debug("received my own blocks message")
				continue
			}

			// Unmarshal the message
			networkMsg := &pb.BlocksSubscriptionMessage{}
			if err := proto.Unmarshal(msg.Data, networkMsg); err != nil {
				log.Error("error unmarshaling blocks sub message", err)
				continue
			}

			// Handle new payload
			switch payload := networkMsg.Payload.(type) {
			case *pb.BlocksSubscriptionMessage_NewBlock:
				n.handleNewBlock(payload.NewBlock)
			case *pb.BlocksSubscriptionMessage_NewTransactions:
				n.handleNewTransactions(payload.NewTransactions)
			default:
				log.Warn("sent a message we don't recognize from the blocks subscription")
				continue
			}
		}
	}
}

func (n *Node) handleNewBlock(newBlock *pb.BlocksSubscriptionNewBlock) {
	log.Infof("got new block(%d): %s, %s", newBlock.Block.Height, newBlock.Block.Hash, newBlock.Block.PreviousHash)
	blockKey := n.sm.BlockKey(newBlock.Block.Height)
	blockStorage := n.sm.BlockStorage()
	if err := blockStorage.Put(blockKey, newBlock.Block); err != nil {
		log.Errorf("could not store block %d on database", err, newBlock.Block.Height)
		return
	}
	if len(newBlock.Transactions) > 0 {
		transactionStorage := n.sm.TransactionStorage()
		for index, transaction := range newBlock.Transactions {
			log.Infof(
				"  transaction %d, '%s', %d, '%s', %d, '%s', '%s', '%s', '%s', %d",
				index,
				transaction.Hash,
				transaction.BlockHeight,
				transaction.Type,
				transaction.Timestamp,
				transaction.PubKey,
				transaction.Verify,
				transaction.Sender,
				transaction.Receiver,
				transaction.Amount,
			)
			transactionKey := n.sm.TransactionKey(transaction.BlockHeight, transaction.Hash)
			if err := transactionStorage.Put(transactionKey, transaction); err != nil {
				log.Errorf("could not store transaction '%s' on database", err, transaction.Hash)
				continue
			}

		}
	}
}

func (n *Node) handleNewTransactions(newTransactions *pb.BlocksSubscriptionNewTransactions) {
	log.Infof("got %d new transactions", len(newTransactions.Transactions))
	for index, transaction := range newTransactions.Transactions {
		log.Infof(
			"  transaction %d, '%s', %d, '%s', %d, '%s', '%s', '%s', '%s', %d",
			index,
			transaction.Hash,
			transaction.BlockHeight,
			transaction.Type,
			transaction.Timestamp,
			transaction.PubKey,
			transaction.Verify,
			transaction.Sender,
			transaction.Receiver,
			transaction.Amount,
		)
		transactionKey := n.sm.TransactionKey(transaction.BlockHeight, transaction.Hash)
		if transaction.BlockHeight == 0 {
			transactionStorage := n.sm.PendingTransactionStorage()
			if err := transactionStorage.Put(transactionKey, transaction); err != nil {
				log.Error("could not store pending transaction", err)
			}
		} else {
			transactionStorage := n.sm.TransactionStorage()
			if err := transactionStorage.Put(transactionKey, transaction); err != nil {
				log.Error("could not store pending transaction", err)
			}
		}
	}
}
