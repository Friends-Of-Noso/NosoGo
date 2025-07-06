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

	value := Host{
		Address: dns.dnsAddress, Port: dns.dnsPort, Key: dns.nodeId,
	}

	writeJSON(w, value)
}

func (dns *DNS) getSeedsHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	seeds := []Host{
		{Address: "10.0.0.1", Port: 8080, Key: "QmYou'reOnFire"},
		{Address: "10.0.0.2", Port: 8080, Key: "QmYou'reOnFire"},
	}

	writeJSON(w, seeds)
}

func (dns *DNS) getNodesHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	value := []Host{
		{Address: "10.0.0.1", Port: 8080, Key: "QmYou'reOnFire"},
		{Address: "10.0.0.2", Port: 8080, Key: "QmYou'reOnFire"},
	}

	writeJSON(w, value)
}

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
