package protobuf

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// Creates Block Zero
func NewBlockZero() *Block {
	return &Block{
		Hash:         "BZERO",
		Height:       0,
		PreviousHash: "",
		Timestamp:    time.Now().Unix(), // TODO: This should be dated to genesis
		MerkleRoot:   "",
	}
}

// Crates a new block
func NewBlock(
	height uint64,
) (*Block, error) {
	// TODO: Add all the fields
	block := &Block{
		Height: height,
	}
	if err := block.SetHash(); err != nil {
		return nil, fmt.Errorf("error creating new block: %w", err)
	}
	return block, nil
}

// Sets the Hash field of the block
func (b *Block) SetHash() error {
	// TODO: Get more values in here
	if b.Height == 0 {
		b.Hash = "BZERO"
		return nil
	}
	value := fmt.Sprintf(
		"%d%s%d%s",
		b.Height,
		b.PreviousHash,
		b.Timestamp,
		b.getMerkleRoot(),
	)
	h := crypto.SHA256.New()
	_, err := h.Write([]byte(value))
	if err != nil {
		return fmt.Errorf("error writing to SHA256: %w", err)
	}
	b.Hash = "B" + strings.ToUpper(hex.EncodeToString(h.Sum([]byte(salt))))
	return nil
}

// Sets the MerkleRoot field of the block
func (b *Block) getMerkleRoot() string {
	// TODO: Implement the merkle root
	if b.MerkleRoot == "" {
		return fmt.Sprintf("MerkleFor:%d", b.Height)
	} else {
		return b.MerkleRoot
	}
}
