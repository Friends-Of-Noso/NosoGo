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

	value := &protobuf.DNSHostResponse{
		Dns: &protobuf.DNSHost{
			Address: dns.dnsAddress,
			Port:    int32(dns.dnsPort),
			Id:      dns.nodeId,
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
		Seeds: []*protobuf.DNSHost{
			{Address: "10.0.0.1", Port: 8080, Id: "QmYou'reOnFire"},
			{Address: "10.0.0.2", Port: 8080, Id: "QmYou'reOnFire"},
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
		Nodes: []*protobuf.DNSHost{
			{Address: "10.0.0.1", Port: 8080, Id: "QmYou'reOnFire"},
			{Address: "10.0.0.2", Port: 8080, Id: "QmYou'reOnFire"},
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
