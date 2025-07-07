package dns

import (
	"encoding/json"
	"net/http"
)

func (dns *DNS) getDNSHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := PeerInfo{
		Address: dns.dnsAddress, Port: dns.dnsPort, Id: dns.nodeId, Mode: "dns",
	}

	writeJSON(w, value)
}

func (dns *DNS) getSeedsHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seeds := []PeerInfo{
		{Address: "10.0.0.1", Port: 8080, Id: "QmYou'reOnFireSeed", Mode: "seed"},
		{Address: "10.0.0.2", Port: 8080, Id: "QmYou'reOnFireSeed", Mode: "seed"},
		{Address: "10.0.0.3", Port: 8080, Id: "QmYou'reOnFireSeed", Mode: "seed"},
	}

	writeJSON(w, seeds)
}

func (dns *DNS) getNodesHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := []PeerInfo{
		{Address: "10.0.0.4", Port: 8080, Id: "QmYou'reOnFireNode", Mode: "node"},
		{Address: "10.0.0.5", Port: 8080, Id: "QmYou'reOnFireNode", Mode: "node"},
		{Address: "10.0.0.6", Port: 8080, Id: "QmYou'reOnFireNode", Mode: "node"},
	}

	writeJSON(w, value)
}

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
