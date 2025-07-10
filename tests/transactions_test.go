package tests

import (
	"testing"

	"gotest.tools/v3/assert"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func TestTransactionsSetHash(t *testing.T) {
	transaction := &pb.Transaction{
		BlockHeight: 1,
		Type:        "COINBASE",
		Timestamp:   1_000_000_000,
		PubKey:      "badbeef",
		Verify:      "badbeef",
		Sender:      "NSender",
		Receiver:    "NReceiver",
	}
	transaction.SetHash()
	want := "T4E6F736F2997166BEB866E8F02A81AA7F0D5C70038257C31ADC5988DE2D59D16A240D45D"
	assert.Equal(t, want, transaction.Hash)
}
