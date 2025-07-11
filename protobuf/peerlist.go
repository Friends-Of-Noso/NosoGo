package protobuf

import (
	"encoding/json"
	"fmt"
	"iter"
	"net/http"
	reflect "reflect"

	"google.golang.org/protobuf/proto"
)

type PeerListType int

const (
	DNS PeerListType = iota
	SEEDS
	NODES
)

type PeerList struct {
	peerType PeerListType
	peers    []*PeerInfo
}

func NewPeerList(plt PeerListType) *PeerList {
	return &PeerList{
		peerType: plt,
		peers:    []*PeerInfo{},
	}
}

func (pl *PeerList) Get(index int) *PeerInfo {
	return pl.peers[index]
}

func (pl *PeerList) Put(peer *PeerInfo) {
	pl.peers = append(pl.peers, peer)
}

func (pl *PeerList) Peers() iter.Seq2[int, *PeerInfo] {
	return func(yield func(int, *PeerInfo) bool) {
		for k, v := range pl.peers {
			if !yield(k, v) {
				return
			}
		}
	}
}

// Helper to write JSON response
func (pl *PeerList) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pl.peers); err != nil {
		http.Error(w, "failed to encode to JSON", http.StatusInternalServerError)
	}
}

// Helper to write ProtoBuf response
func (pl *PeerList) WriteProtobuf(w http.ResponseWriter) {
	var msgDNS *DNSPeerResponse
	var msgSeeds *DNSSeedsResponse
	var msgNodes *DNSNodesResponse
	var msgType reflect.Type
	switch pl.peerType {
	case DNS:
		msgDNS = &DNSPeerResponse{
			Dns: pl.peers,
		}
		// Determine the message name
		msgType = reflect.TypeOf(msgSeeds).Elem()
	case SEEDS:
		msgSeeds = &DNSSeedsResponse{
			Seeds: pl.peers,
		}
		// Determine the message name
		msgType = reflect.TypeOf(msgSeeds).Elem()
	case NODES:
		msgNodes = &DNSNodesResponse{
			Nodes: pl.peers,
		}
		// Determine the message name
		msgType = reflect.TypeOf(msgNodes).Elem()
	}
	protoMime := fmt.Sprintf("application/x-protobuf; proto=%s.%s", msgType.PkgPath(), msgType.Name())
	w.Header().Set("Content-Type", protoMime)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	var data []byte
	var err error
	switch pl.peerType {
	case DNS:
		data, err = proto.Marshal(msgDNS)
	case SEEDS:
		data, err = proto.Marshal(msgSeeds)
	case NODES:
		data, err = proto.Marshal(msgNodes)
	}
	if err != nil {
		http.Error(w, "failed to marshal protobuf", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
