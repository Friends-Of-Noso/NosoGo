package protobuf

import (
	"encoding/json"
	"fmt"
	"net/http"
	reflect "reflect"
	"sync"

	"google.golang.org/protobuf/proto"
)

// Peer direction constants
const (
	DirectionOutbound = "OUTBOUND"
	DirectionInbound  = "INBOUND"
)

type PeerList struct {
	mu    sync.RWMutex
	peers map[string]*PeerInfo
}

// NewPeerList constructs a new PeerList
func NewPeerList() *PeerList {
	return &PeerList{
		peers: make(map[string]*PeerInfo),
	}
}

// Add inserts or replaces a peer in the list
func (pl *PeerList) Add(peer *PeerInfo) {
	if peer == nil || peer.Id == "" {
		return
	}

	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.peers[peer.Id] = peer
}

// Remove deletes a peer by its ID
func (pl *PeerList) Remove(id string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	delete(pl.peers, id)
}

// Connected returns true if the peer with given ID is connected
func (pl *PeerList) Connected(id string) bool {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	peer, ok := pl.peers[id]
	return ok && peer.Connected
}

// Disconnect marks a peer as disconnected and clears its direction
func (pl *PeerList) Disconnect(id string) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	if peer, ok := pl.peers[id]; ok {
		peer.Connected = false
		peer.Direction = ""
	}
}

// Peers returns a snapshot map of all peers
func (pl *PeerList) Peers() map[string]*PeerInfo {
	pl.mu.RLock()
	defer pl.mu.RUnlock()

	// Shallow copy — safe because PeerInfo pointers are stable
	copied := make(map[string]*PeerInfo, len(pl.peers))
	for id, peer := range pl.peers {
		copied[id] = peer
	}
	return copied
}

// func (pl *PeerList) Peers() iter.Seq2[string, *PeerInfo] {
// 	return func(yield func(string, *PeerInfo) bool) {
// 		for k, v := range pl.peers {
// 			if !yield(k, v) {
// 				return
// 			}
// 		}
// 	}
// }

// Helper to write JSON response
func (pl *PeerList) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pl.peers); err != nil {
		http.Error(w, "failed to encode to JSON", http.StatusInternalServerError)
	}
}

// Helper to write ProtoBuf response
func (pl *PeerList) WriteProtobuf(w http.ResponseWriter) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	msg := &DNSPeerResponse{}

	for _, peer := range pl.peers {
		msg.Peers = append(msg.Peers, peer)
	}

	msgType := reflect.TypeOf(msg).Elem()

	protoMime := fmt.Sprintf("application/x-protobuf; proto=%s.%s", msgType.PkgPath(), msgType.Name())
	w.Header().Set("Content-Type", protoMime)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	data, err := proto.Marshal(msg)
	if err != nil {
		http.Error(w, "failed to marshal protobuf", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
