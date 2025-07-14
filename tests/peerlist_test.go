package tests

import (
	"testing"

	"gotest.tools/v3/assert"

	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

// Test Adding a peer to the list
func TestPeerListAdd(t *testing.T) {
	t.Parallel()

	peerList := pb.NewPeerList()
	assert.Equal(t, 0, len(peerList.Peers()))

	peer := &pb.PeerInfo{
		Address:   "0.0.0.0",
		Port:      8080,
		Id:        "QmTesting",
		Mode:      "node",
		Connected: true,
		Direction: pb.DirectionInbound,
	}
	peerList.Add(peer)
	assert.Equal(t, 1, len(peerList.Peers()))
}

// Test Adding a peer to the list
func TestPeerListRemove(t *testing.T) {
	t.Parallel()

	peerList := pb.NewPeerList()
	assert.Equal(t, 0, len(peerList.Peers()))

	peer := &pb.PeerInfo{
		Address:   "0.0.0.0",
		Port:      8080,
		Id:        "QmTesting",
		Mode:      "node",
		Connected: true,
		Direction: pb.DirectionInbound,
	}
	peerList.Add(peer)
	assert.Equal(t, 1, len(peerList.Peers()))
	peerList.Remove(peer.Id)
	assert.Equal(t, 0, len(peerList.Peers()))
}

// Test if a peer is connected
func TestPeerListConnected(t *testing.T) {
	t.Parallel()

	peerList := pb.NewPeerList()
	peer := &pb.PeerInfo{
		Address:   "0.0.0.0",
		Port:      8080,
		Id:        "QmTesting",
		Mode:      "node",
		Connected: true,
		Direction: pb.DirectionInbound,
	}
	peerList.Add(peer)
	assert.Equal(t, true, peerList.Connected(peer.Id))
}

// Test disconnecting a peer
func TestPeerListDisconnect(t *testing.T) {
	t.Parallel()

	peerList := pb.NewPeerList()
	peer := &pb.PeerInfo{
		Address:   "0.0.0.0",
		Port:      8080,
		Id:        "QmTesting",
		Mode:      "node",
		Connected: true,
		Direction: pb.DirectionInbound,
	}
	peerList.Add(peer)
	peerList.Disconnect(peer.Id)
	assert.Equal(t, false, peerList.Connected(peer.Id))
	assert.Equal(t, "", peerList.Peers()[peer.Id].Direction)
	// This also works due to shallow copy
	assert.Equal(t, false, peer.Connected)
	assert.Equal(t, "", peer.Direction)
}
