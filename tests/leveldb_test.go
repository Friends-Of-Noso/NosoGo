package tests

import (
	"fmt"
	"os"
	"testing"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
	"github.com/Friends-Of-Noso/NosoGo/store"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	dbPath = "./test_data"
)

func TestBlock(t *testing.T) {
	// Remove test data, start fresh
	if utils.FileExists(dbPath) {
		err := os.RemoveAll(dbPath)
		if err != nil {
			t.Errorf("could not delete folder '%s'", dbPath)
		}
	}

	sm, err := store.NewStorageManager(dbPath)
	if err != nil {
		t.Errorf("could not create storage manager in '%s'", dbPath)
	}
	defer sm.Close()

	blockStorage := sm.BlockStorage()
	block := &pb.Block{
		Height:       1,
		Hash:         "Block1Hash",
		PreviousHash: "Block0Hash",
	}

	key := sm.BlockKey(block.Height)
	if err := blockStorage.Put(key, block); err != nil {
		t.Errorf("failed to store block: %v", err)
	}

	// Retrieve the block
	retrievedBlock := &pb.Block{}
	if err := blockStorage.Get(key, retrievedBlock); err != nil {
		t.Errorf("failed to retrieve block: %v", err)
	}

	if block.Height != retrievedBlock.Height {
		t.Errorf("mismatch height: wanted %d, got %d", block.Height, retrievedBlock.Height)
	}

	if block.Hash != retrievedBlock.Hash {
		t.Errorf("mismatch hash: wanted '%s', got '%s'", block.Hash, retrievedBlock.Hash)
	}

	if block.PreviousHash != retrievedBlock.PreviousHash {
		t.Errorf("mismatch previous hash: wanted '%s,' got '%s'", block.PreviousHash, retrievedBlock.PreviousHash)
	}
}

func BenchmarkBlocksCreate1_000(b *testing.B) { createBlocks(b, 1_000) }
func BenchmarkBlocksRead1_000(b *testing.B)   { readBlocks(b, 1_000) }

func BenchmarkBlocksCreate100_000(b *testing.B) { createBlocks(b, 100_000) }
func BenchmarkBlocksRead100_000(b *testing.B)   { readBlocks(b, 100_000) }

func BenchmarkBlocksCreate1_000_000(b *testing.B) { createBlocks(b, 1_000_000) }
func BenchmarkBlocksRead1_000_000(b *testing.B)   { readBlocks(b, 1_000_000) }

// Utilities
func createBlocks(b *testing.B, count uint64) {
	// Remove test data, start fresh
	if utils.FileExists(dbPath) {
		err := os.RemoveAll(dbPath)
		if err != nil {
			b.Errorf("could not delete folder '%s'", dbPath)
		}
	}

	sm, err := store.NewStorageManager(dbPath)
	if err != nil {
		b.Errorf("could not create storage manager in '%s'", dbPath)
	}
	defer sm.Close()

	blockStorage := sm.BlockStorage()
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
			b.Errorf("failed to store block: %v", err)
		}
	}
}

func readBlocks(b *testing.B, count uint64) {
	sm, err := store.NewStorageManager(dbPath)
	if err != nil {
		b.Errorf("could not create storage manager in '%s'", dbPath)
	}
	defer sm.Close()

	blockStorage := sm.BlockStorage()

	block := &pb.Block{}
	for i := range count {
		key := sm.BlockKey(i)
		if err := blockStorage.Get(key, block); err != nil {
			b.Errorf("Failed to retrieve block: %v", err)
		}
	}

}
