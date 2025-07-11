package dns

import (
	"fmt"
	"net/http"
	"reflect"

	"google.golang.org/protobuf/proto"
)

func (dns *DNS) getDNSHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dns.dns.WriteProtobuf(w)
}

func (dns *DNS) getSeedsHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dns.seeds.WriteProtobuf(w)
}

func (dns *DNS) getNodesHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dns.nodes.WriteProtobuf(w)
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
