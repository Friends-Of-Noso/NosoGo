package dns

import (
	"encoding/json"
	"net/http"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
)

func (dns *DNS) getDNSHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dns.dns.WriteJSON(w)
}

func (dns *DNS) getSeedsHandlerJSON(w http.ResponseWriter, r *http.Request) {
	log.Debug("dns serving seeds")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dns.seeds.WriteJSON(w)
}

func (dns *DNS) getNodesHandlerJSON(w http.ResponseWriter, r *http.Request) {
	log.Debug("dns serving nodes")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dns.nodes.WriteJSON(w)
}

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "failed to encode to JSON", http.StatusInternalServerError)
	}
}
