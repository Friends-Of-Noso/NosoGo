package tests

import (
	"errors"
	"fmt"
	"testing"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
	"github.com/Friends-Of-Noso/NosoGo/store"
	"github.com/Friends-Of-Noso/NosoGo/utils"
	"github.com/syndtr/goleveldb/leveldb"
)

const (
	dbPath    = "./test_data"
	statusKey = "status-main"
)

var (
	sm *store.StorageManager

	statusStorage      *store.Storage[*pb.Status]      = nil
	blockStorage       *store.Storage[*pb.Block]       = nil
	transactionStorage *store.Storage[*pb.Transaction] = nil
	// peerInfoStorage    *store.Storage[*pb.PeerInfo]    = nil
)

// Test storing and retrieving a block
func TestBlockPutGet(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	block := &pb.Block{
		Height:       1,
		Hash:         "Block1Hash",
		PreviousHash: "Block0Hash",
	}

	key := sm.BlockKey(block.Height)
	if err := blockStorage.Put(key, block); err != nil {
		t.Fatalf("failed to store block: %v", err)
	}

	// Retrieve the block
	retrievedBlock := &pb.Block{}
	if err := blockStorage.Get(key, retrievedBlock); err != nil {
		t.Fatalf("failed to retrieve block: %v", err)
	}

	if block.Height != retrievedBlock.Height {
		t.Fatalf("mismatch height: wanted %d, got %d", block.Height, retrievedBlock.Height)
	}

	if block.Hash != retrievedBlock.Hash {
		t.Fatalf("mismatch hash: wanted '%s', got '%s'", block.Hash, retrievedBlock.Hash)
	}

	if block.PreviousHash != retrievedBlock.PreviousHash {
		t.Fatalf("mismatch previous hash: wanted '%s,' got '%s'", block.PreviousHash, retrievedBlock.PreviousHash)
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test retrieving a list of blocks with correct data
func TestBlocksStorageListCorrect(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	createCorrectData(t)

	blocks, err := blockStorage.ListKeys()
	if err != nil {
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for index, value := range blocks {
		key := sm.BlockKey(uint64(index))
		if key != value {
			t.Fatalf("list out of order: wanted '%s', got '%s'", key, value)
		}
		height++
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test retrieving a list of blocks with incorrect data
func TestBlocksStorageListIncorrect(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	createIncorrectData(t)

	blocks, err := blockStorage.ListKeys()
	if err != nil {
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for index, value := range blocks {
		key := sm.BlockKey(uint64(index))
		if key != value {
			return
		}
		height++
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test retrieving a list of blocks with values with correct data
func TestBlocksStorageListWithValuesCorrect(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	createCorrectData(t)

	blocks, err := blockStorage.ListValues(func() *pb.Block {
		return &pb.Block{}
	})
	if err != nil {
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for _, block := range blocks {
		if height != block.Height {
			// closeStorageManager()
			// databaseResetT(t)
			t.Errorf("list out of order: wanted '%d', got '%d'", height, block.Height)
		}
		height++
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test retrieving a list of blocks with values with incorrect data
func TestBlocksStorageListWithValuesIncorrect(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	createIncorrectData(t)

	blocks, err := blockStorage.ListValues(func() *pb.Block {
		return &pb.Block{}
	})
	if err != nil {
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for _, block := range blocks {
		if height != block.Height {
			closeStorageManager()
			databaseResetT(t)
			return
		}
		height++
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test that mimics a first run with clean database
func TestStartupSequenceFirstRun(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	if err := startUp(t); err != nil {
		t.Errorf("failed the startup: %v", err)
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test that mimics a normal run with correct data
func TestCheckBlockChainWithCorrectData(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	createCorrectData(t)

	if err := startUp(t); err != nil {
		t.Errorf("failed the startup: %v", err)
	}

	closeStorageManager()
	databaseResetT(t)
}

// Test that mimics a normal run with correct data
func TestCheckBlockChainWithIncorrectData(t *testing.T) {
	// t.SkipNow()

	databaseResetT(t)

	err := initializeStorage()
	if err != nil {
		t.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	createIncorrectData(t)

	if err := startUp(t); err != nil {
		t.Errorf("failed the startup: %v", err)
	}

	closeStorageManager()
	databaseResetT(t)
}

// Benchmark that creates 1 thousand blocks
func BenchmarkBlocksCreate1_000(b *testing.B) { createBlocks(b, 1_000) }

// Benchmark that reads 1 thousand blocks
func BenchmarkBlocksRead1_000(b *testing.B) { readBlocks(b, 1_000) }

// Benchmark that creates 100 thousand blocks
func BenchmarkBlocksCreate100_000(b *testing.B) { createBlocks(b, 100_000) }

// Benchmark that reads 100 thousand blocks
func BenchmarkBlocksRead100_000(b *testing.B) { readBlocks(b, 100_000) }

// Benchmark that creates 1 million blocks
func BenchmarkBlocksCreate1_000_000(b *testing.B) { createBlocks(b, 1_000_000) }

// Benchmark that reads 1 million blocks
func BenchmarkBlocksRead1_000_000(b *testing.B) { readBlocks(b, 1_000_000) }

// Initialize storage variables
func initializeStorage() error {
	var err error
	sm, err = store.NewStorageManager(dbPath)
	if err != nil {
		return err
	}
	statusStorage = sm.StatusStorage()
	blockStorage = sm.BlockStorage()
	transactionStorage = sm.TransactionStorage()
	// peerInfoStorage = sm.PeerInfoStorage()

	return nil
}

// Close the Storage Manager
func closeStorageManager() {
	sm.Close()
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
		closeStorageManager()
		databaseResetT(t)
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
				closeStorageManager()
				databaseResetT(t)
				return err
			}
		} else {
			if err := initiateBlockChain(); err != nil {
				closeStorageManager()
				databaseResetT(t)
				return err
			}
		}
	} else {
		if status.LastBlock != blocksCount-1 {
			// Attempt to rescan the database
			if err := reScanBlockChain(); err != nil {
				// Ok, data is well corrupted
				closeStorageManager()
				databaseResetT(t)
				return err
			}
		}
	}
	return nil
}

// Initializes the block chain with block zero and sets status
func initiateBlockChain() error {
	blockZero := getBlockZero()

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
		previous        = getBlockZero()
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

func getBlockZero() *pb.Block {
	block := &pb.Block{
		Height:       0,
		PreviousHash: "BZERO",
		Timestamp:    0, // This must be the inception date
		MerkleRoot:   "MZERO",
	}
	block.SetHash()
	return block
}

// Reset the database
func databaseReset() error {
	// Remove test data, start fresh
	if utils.FileExists(dbPath) {
		err := utils.RemoveGlob(dbPath + "/*")
		if err != nil {
			return err
		}
	}
	return nil
}

// Reset the database with testing.T
func databaseResetT(t *testing.T) {
	if err := databaseReset(); err != nil {
		t.Fatalf("could not delete folder '%s': %v", dbPath, err)
	}
}

// Reset the database with `testing.B`
func databaseResetB(b *testing.B) {
	if err := databaseReset(); err != nil {
		b.Fatalf("could not delete folder '%s': %v", dbPath, err)
	}
}

// Tests helper that creates a correct sequence of data
func createCorrectData(t *testing.T) {
	var limit uint64 = 10

	status := &pb.Status{
		LastBlock: limit,
	}

	if err := statusStorage.Put(statusKey, status); err != nil {
		t.Fatalf("could not store status: %v", err)
	}

	var (
		block    *pb.Block
		previous *pb.Block
	)

	for height := range limit {
		// Block
		if height == 0 {
			block = getBlockZero()
			previous = getBlockZero()
		} else {
			block = &pb.Block{
				Height:       height,
				PreviousHash: previous.Hash,
			}
			block.SetHash()
			previous = block
		}

		key := sm.BlockKey(block.Height)
		if err := blockStorage.Put(key, block); err != nil {
			t.Fatalf("could not put block: %v", err)
		}

		// Transaction
		for range 3 {
			transaction := &pb.Transaction{
				BlockHeight: block.Height,
				Type:        "spend",
			}
			transaction.SetHash()
			key := sm.TransactionKey(block.Height, transaction.Hash)
			if err := transactionStorage.Put(key, transaction); err != nil {
				t.Fatalf("could not store transaction: %v", err)
			}
		}
	}
}

// Tests helper that creates a correct sequence of data
func createIncorrectData(t *testing.T) {
	var limit uint64 = 10

	status := &pb.Status{
		LastBlock: limit,
	}

	if err := statusStorage.Put(statusKey, status); err != nil {
		t.Fatalf("could not store status: %v", err)
	}

	var (
		block    *pb.Block
		previous *pb.Block
	)
	for height := range limit - 1 {
		// Block
		if height == 0 {
			block = getBlockZero()
			previous = getBlockZero()
		} else {
			block = &pb.Block{
				Height:       height,
				PreviousHash: previous.Hash,
			}
			block.SetHash()
			previous = block
		}

		key := sm.BlockKey(block.Height)
		if err := blockStorage.Put(key, block); err != nil {
			t.Fatalf("could not put block: %v", err)
		}
	}

	block = &pb.Block{
		Height:       limit,
		PreviousHash: previous.Hash,
	}
	block.SetHash()
	key := sm.BlockKey(block.Height)
	if err := blockStorage.Put(key, block); err != nil {
		t.Fatalf("could not put block: %v", err)
	}

	// Transaction
	for range 3 {
		transaction := &pb.Transaction{
			BlockHeight: 42,
			Type:        "spend",
		}
		transaction.SetHash()
		key := sm.TransactionKey(block.Height, transaction.Hash)
		if err := transactionStorage.Put(key, transaction); err != nil {
			t.Fatalf("could not store transaction: %v", err)
		}
	}

}

// Benchmark helper to create many Blocks
func createBlocks(b *testing.B, count uint64) {
	databaseResetB(b)

	err := initializeStorage()
	if err != nil {
		b.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	var block *pb.Block
	for i := range count {

		if i == 0 {
			block = &pb.Block{
				Height:       uint64(i),
				Hash:         fmt.Sprintf("Block%dHash", i),
				PreviousHash: "",
			}
		} else {
			block = &pb.Block{
				Height:       uint64(i),
				Hash:         fmt.Sprintf("Block%dHash", i),
				PreviousHash: fmt.Sprintf("Block%dHash", i-1),
			}
		}

		key := sm.BlockKey(block.Height)
		if err := blockStorage.Put(key, block); err != nil {
			b.Fatalf("failed to store block: %v", err)
		}
	}
	closeStorageManager()
}

// Benchmark helper to read many Blocks
func readBlocks(b *testing.B, count uint64) {
	err := initializeStorage()
	if err != nil {
		b.Fatalf("could not create storage manager in '%s'", dbPath)
	}

	block := &pb.Block{}
	for i := range count {
		key := sm.BlockKey(i)
		if err := blockStorage.Get(key, block); err != nil {
			b.Fatalf("Failed to retrieve block: %v", err)
		}
	}
	closeStorageManager()
	databaseResetB(b)
}
