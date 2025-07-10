package tests

import (
	"fmt"
	"testing"

	"gotest.tools/v3/assert"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
	"github.com/Friends-Of-Noso/NosoGo/store"
)

const (
	dbPath    = "test_data"
	statusKey = "status-main"
)

var (
	sm *store.StorageManager

	statusStorage      *store.Storage[*pb.Status]
	blockStorage       *store.Storage[*pb.Block]
	transactionStorage *store.Storage[*pb.Transaction]
	// peerInfoStorage    *store.Storage[*pb.PeerInfo]
)

// Test storing and retrieving a block
func TestBlockPutGet(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	block := &pb.Block{
		Height:       1,
		Hash:         "Block1Hash",
		PreviousHash: "Block0Hash",
	}

	key := sm.BlockKey(block.Height)
	err = blockStorage.Put(key, block)
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	// Retrieve the block
	retrievedBlock := &pb.Block{}
	err = blockStorage.Get(key, retrievedBlock)
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	assert.Equal(t, block.Height, retrievedBlock.Height)

	assert.Equal(t, block.Hash, retrievedBlock.Hash)

	assert.Equal(t, block.PreviousHash, retrievedBlock.PreviousHash)

	databaseResetT(t)
	closeStorageManager()
}

// Test retrieving a list of blocks with correct data
func TestBlocksStorageListCorrect(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	createCorrectData(t)

	blocks, err := blockStorage.ListKeys()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for index, value := range blocks {
		key := sm.BlockKey(uint64(index))
		if key != value {
			databaseResetT(t)
			closeStorageManager()
			t.Fatalf("list out of order: wanted '%s', got '%s'", key, value)
		}
		height++
	}

	databaseResetT(t)
	closeStorageManager()
}

// Test retrieving a list of blocks with incorrect data
func TestBlocksStorageListIncorrect(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	createIncorrectData(t)

	blocks, err := blockStorage.ListKeys()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for index, value := range blocks {
		key := sm.BlockKey(uint64(index))
		if key != value {
			break
		}
		height++
	}

	databaseResetT(t)
	closeStorageManager()
}

// Test retrieving a list of blocks with values with correct data
func TestBlocksStorageListWithValuesCorrect(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	createCorrectData(t)

	blocks, err := blockStorage.ListValues(func() *pb.Block {
		return &pb.Block{}
	})
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for _, block := range blocks {
		if height != block.Height {
			t.Errorf("list out of order: wanted '%d', got '%d'", height, block.Height)
			break
		}
		height++
	}

	databaseResetT(t)
	closeStorageManager()
}

// Test retrieving a list of blocks with values with incorrect data
func TestBlocksStorageListWithValuesIncorrect(t *testing.T) {
	// t.SkipNow()

	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetT(t)

	createIncorrectData(t)

	blocks, err := blockStorage.ListValues(func() *pb.Block {
		return &pb.Block{}
	})
	if err != nil {
		closeStorageManager()
		t.Fatalf("could not retrieve list of blocks: %v", err)
	}
	var height uint64 = 0
	for _, block := range blocks {
		if height != block.Height {
			break
		}
		height++
	}

	databaseResetT(t)
	closeStorageManager()
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

// Reset the database
func databaseReset() error {
	// Remove test data, start fresh

	// Delete Status
	if err := statusStorage.Delete(statusKey); err != nil {
		return err
	}

	// Delete blocks
	blocks, err := blockStorage.ListValues(func() *pb.Block { return &pb.Block{} })
	if err != nil {
		return err
	}
	for _, block := range blocks {
		key := sm.BlockKey(block.Height)
		if err := blockStorage.Delete(key); err != nil {
			return err
		}
	}

	// Delete transactions
	transactions, err := transactionStorage.ListKeys()
	if err != nil {
		return err
	}
	for _, key := range transactions {
		if err := transactionStorage.Delete(key); err != nil {
			return err
		}
	}

	return nil
}

// Reset the database with testing.T
func databaseResetT(t *testing.T) {
	if err := databaseReset(); err != nil {
		closeStorageManager()
		t.Fatalf("could not delete data: %v", err)
	}
}

// Reset the database with `testing.B`
func databaseResetB(b *testing.B) {
	if err := databaseReset(); err != nil {
		closeStorageManager()
		b.Fatalf("could not delete data: %v", err)
	}
}

// Tests helper that creates a correct sequence of data
func createCorrectData(t *testing.T) {
	var limit uint64 = 10

	status := &pb.Status{
		LastBlock: limit,
	}

	if err := statusStorage.Put(statusKey, status); err != nil {
		closeStorageManager()
		t.Fatalf("could not store status: %v", err)
	}

	var (
		block    *pb.Block
		previous *pb.Block
	)

	for height := range limit {
		// Block
		if height == 0 {
			block = pb.NewBlockZero()
			previous = pb.NewBlockZero()
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
			closeStorageManager()
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
				closeStorageManager()
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
		closeStorageManager()
		t.Fatalf("could not store status: %v", err)
	}

	var (
		block    *pb.Block
		previous *pb.Block
	)
	for height := range limit - 1 {
		// Block
		if height == 0 {
			block = pb.NewBlockZero()
			previous = pb.NewBlockZero()
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
			closeStorageManager()
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
		closeStorageManager()
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
			closeStorageManager()
			t.Fatalf("could not store transaction: %v", err)
		}
	}

}

// Benchmark helper to create many Blocks
func createBlocks(b *testing.B, count uint64) {
	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		b.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	databaseResetB(b)

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
			closeStorageManager()
			b.Fatalf("failed to store block: %v", err)
		}
	}
	closeStorageManager()
}

// Benchmark helper to read many Blocks
func readBlocks(b *testing.B, count uint64) {
	err := initializeStorage()
	if err != nil {
		closeStorageManager()
		b.Fatalf("could not open storage manager in '%s': %v", dbPath, err)
	}

	var block *pb.Block
	for i := range count {
		key := sm.BlockKey(i)
		if err := blockStorage.Get(key, block); err != nil {
			closeStorageManager()
			b.Fatalf("Failed to retrieve block: %v", err)
		}
	}
	databaseResetB(b)
	closeStorageManager()
}
