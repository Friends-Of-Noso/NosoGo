package protobuf

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
)

// Sets the Hash field of the block
func (b *Block) SetHash() error {
	// TODO: Get more values in here
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
		return fmt.Errorf("error writing to SHA256: %v", err)
	}
	b.Hash = "B" + strings.ToUpper(hex.EncodeToString(h.Sum([]byte(salt))))
	return nil
}

// Sets the MerkleRoot field of the block
func (b *Block) getMerkleRoot() string {
	if b.MerkleRoot == "" {
		return fmt.Sprintf("MerkleFor:%d", b.Height)
	} else {
		return b.MerkleRoot
	}
}
