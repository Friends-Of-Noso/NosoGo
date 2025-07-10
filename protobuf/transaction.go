package protobuf

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"strings"
)

// Creates a new transaction
func NewTransaction(
	height uint64,
) (*Transaction, error) {
	// TODO: Add more fields
	transaction := &Transaction{
		BlockHeight: height,
	}
	if err := transaction.SetHash(); err != nil {
		return nil, fmt.Errorf("error creating new transaction: %v", err)
	}
	return transaction, nil
}

// Sets the Hash field of the transaction
func (t *Transaction) SetHash() error {
	// TODO: Get more values in here
	value := fmt.Sprintf(
		"%d%s%d%s%s%s%s",
		t.BlockHeight,
		t.Type,
		t.Timestamp,
		t.PubKey,
		t.Verify,
		t.Sender,
		t.Receiver,
	)
	h := crypto.SHA256.New()
	_, err := h.Write([]byte(value))
	if err != nil {
		return fmt.Errorf("error writing to SHA256: %v", err)
	}
	t.Hash = "T" + strings.ToUpper(hex.EncodeToString(h.Sum([]byte(salt))))
	return nil
}
