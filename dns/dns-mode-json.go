package dns

import (
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
