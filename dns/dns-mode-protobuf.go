package dns

import (
	"fmt"
	"net/http"
	"reflect"

	"google.golang.org/protobuf/proto"

	"github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func (dns *DNS) getDNSHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := &protobuf.DNSPeerResponse{
		Dns: &protobuf.PeerInfo{
			Address: dns.dnsAddress,
			Port:    int32(dns.dnsPort),
			Id:      dns.nodeId,
			Mode:    "dns",
		},
	}

	writeProtobuf(w, value)
}

func (dns *DNS) getSeedsHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := &protobuf.DNSSeedsResponse{
		Seeds: []*protobuf.PeerInfo{
			{Address: "10.0.0.1", Port: 8080, Id: "QmYou'reOnFireSeed", Mode: "seed"},
			{Address: "10.0.0.2", Port: 8080, Id: "QmYou'reOnFireSeed", Mode: "seed"},
			{Address: "10.0.0.3", Port: 8080, Id: "QmYou'reOnFireSeed", Mode: "seed"},
		},
	}

	writeProtobuf(w, value)
}

func (dns *DNS) getNodesHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := &protobuf.DNSNodesResponse{
		Nodes: []*protobuf.PeerInfo{
			{Address: "10.0.0.4", Port: 8080, Id: "QmYou'reOnFireNode", Mode: "node"},
			{Address: "10.0.0.5", Port: 8080, Id: "QmYou'reOnFireNode", Mode: "node"},
			{Address: "10.0.0.6", Port: 8080, Id: "QmYou'reOnFireNode", Mode: "node"},
		},
	}

	writeProtobuf(w, value)
}

// Helper to write ProtoBuf response
func writeProtobuf(w http.ResponseWriter, msg proto.Message) {
	// Determine the message name (e.g., protobuf.DNSSeedsResponse)
	msgType := reflect.TypeOf(msg).Elem()
	protoMime := fmt.Sprintf("application/x-protobuf; proto=%s.%s", msgType.PkgPath(), msgType.Name())
	w.Header().Set("Content-Type", protoMime)
	w.Header().Set("X-Content-Type-Options", "nosniff")

	data, err := proto.Marshal(msg)
	if err != nil {
		http.Error(w, "Failed to marshal protobuf", http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
