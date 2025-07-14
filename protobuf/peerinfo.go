package protobuf

import (
	"encoding/json"
	"fmt"
	"net/http"
	reflect "reflect"

	"google.golang.org/protobuf/proto"
)

// Helper to write JSON response
func (pi *PeerInfo) WriteJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pi); err != nil {
		http.Error(w, "failed to encode to JSON", http.StatusInternalServerError)
	}
}

// Helper to write ProtoBuf response
func (pi *PeerInfo) WriteProtoBuf(w http.ResponseWriter) {
	msgType := reflect.TypeOf(pi).Elem()

	protoMime := fmt.Sprintf("application/x-protobuf; proto=%s.%s", msgType.PkgPath(), msgType.Name())
	w.Header().Set("Content-Type", protoMime)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	data, err := proto.Marshal(pi)
	if err != nil {
		http.Error(w, "failed to marshal protobuf", http.StatusInternalServerError)
		return
	}
	w.Write(data)

}
