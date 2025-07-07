package node

import (
	"time"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

// This file contains stuff that only makes sense during heavy development
func (n *Node) devPropagateData(height uint64) {
	// TODO This needs to go away on the real thing.
	if n.seed != "" {
		return
	}
	block := &pb.Block{
		Height:       height,
		PreviousHash: "BPreviousHash",
		Timestamp:    time.Now().Unix(),
	}
	block.SetHash()
	// Store new block
	blockKey := n.sm.BlockKey(block.Height)
	blockStorage := n.sm.BlockStorage()
	if err := blockStorage.Put(blockKey, block); err != nil {
		log.Error("could not store block on database", err)
		return
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
		return
	}

	newBlock := &pb.NewBlock{
		Block: block,
		Transactions: []*pb.Transaction{
			transaction,
		},
	}
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
		return
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
}
