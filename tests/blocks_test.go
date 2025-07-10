package tests

import (
	"testing"

	"gotest.tools/v3/assert"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func TestBlocksSetHash(t *testing.T) {
	block := &pb.Block{
		Height:       1,
		PreviousHash: "BPreviousHash",
		Timestamp:    1_000_000_000,
	}
	block.SetHash()
	want := "B4E6F736F46C8706786CD83871D23D2338D5BA148215EFEACE7C304C2C987CD0C96DBE8B1"
	assert.Equal(t, want, block.Hash)
}
