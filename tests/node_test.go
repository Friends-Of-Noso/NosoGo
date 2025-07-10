package tests

import (
	"errors"
	"testing"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
	"github.com/syndtr/goleveldb/leveldb"
)

// Test that mimics a first run with clean database
func TestBlockchainFirstRun(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	if err := startUp(t); err != nil {
		t.Errorf("failed the startup: %v", err)
	}

	databaseResetT(t)
	closeStorageManager()
}

// Test that mimics a normal run with correct data
func TestBlockchainWithCorrectData(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	createCorrectData(t)

	if err := startUp(t); err != nil {
		t.Errorf("failed the startup: %v", err)
	}

	databaseResetT(t)
	closeStorageManager()
}

// Test that mimics a normal run with correct data
func TestBlockchainWithIncorrectData(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	createIncorrectData(t)

	if err := startUp(t); err != nil {
		t.Errorf("failed the startup: %v", err)
	}

	databaseResetT(t)
	closeStorageManager()
}

func startUp(t *testing.T) error {
	// Check if we have any blocks
	blocksCount, err := blockStorage.Count()
	if err != nil {
		return err
	}

	// Check if we have any transactions
	transactionsCount, err := transactionStorage.Count()
	if err != nil {
		databaseResetT(t)
		closeStorageManager()
		return err
	}

	status := &pb.Status{}
	err = statusStorage.Get(statusKey, status)
	if errors.Is(err, leveldb.ErrNotFound) {

		// We have blocks or transactions but status is missing
		if blocksCount > 0 || transactionsCount > 0 {
			// Attempt to rescan the database
			if err := reScanBlockChain(); err != nil {
				// Ok, data is well corrupted
				databaseResetT(t)
				closeStorageManager()
				return err
			}
		} else {
			if err := initiateBlockChain(); err != nil {
				databaseResetT(t)
				closeStorageManager()
				return err
			}
		}
	} else {
		if status.LastBlock != blocksCount-1 {
			// Attempt to rescan the database
			if err := reScanBlockChain(); err != nil {
				// Ok, data is well corrupted
				databaseResetT(t)
				closeStorageManager()
				return err
			}
		}
	}
	return nil
}

// Initializes the block chain with block zero and sets status
func initiateBlockChain() error {
	blockZero := pb.NewBlockZero()

	blockZeroKey := sm.BlockKey(blockZero.Height)
	if err := blockStorage.Put(blockZeroKey, blockZero); err != nil {
		return err
	}
	status := &pb.Status{
		LastBlock: 0,
	}
	if err := statusStorage.Put(statusKey, status); err != nil {
		return err
	}

	return nil
}

// Re-scans the database and tries to recover status
func reScanBlockChain() error {
	blocks, err := blockStorage.ListValues(func() *pb.Block {
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
		// fmt.Printf("block: %d, '%s', '%s'\n", block.Height, block.Hash, block.PreviousHash)
		if block.Height != height {
			// return fmt.Errorf("mismatched block height, expected %d, got %d", height, block.Height)
			return nil
		}

		if block.Height == 0 && block.PreviousHash != previous.PreviousHash {
			// return fmt.Errorf("chain is broken: block %d does not have the the correct previous hash", block.Height)
			return nil
		}
		if block.Height != 0 && block.PreviousHash != previous.Hash {
			// return fmt.Errorf("chain is broken: block %d does not have the the correct previous hash", block.Height)
			return nil
		}

		previous = block
		height++
	}
	status := &pb.Status{
		LastBlock: height - 1,
	}
	if err := statusStorage.Put(statusKey, status); err != nil {
		return err
	}

	// Check for orphaned transactions
	transactions, err := transactionStorage.ListValues(func() *pb.Transaction {
		return &pb.Transaction{}
	})
	if err != nil {
		return err
	}
	for _, transaction := range transactions {
		key := sm.BlockKey(transaction.BlockHeight)
		ok, err := blockStorage.Has(key)
		if err != nil {
			return err
		}
		if !ok {
			return err
		}
	}
	return nil
}
