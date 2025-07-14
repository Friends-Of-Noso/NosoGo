package dns

import (
	"net/http"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func (dns *DNS) getDNSHandlerJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peerList := pb.NewPeerList()
	peerList.Add(dns.nodePeer)
	peerList.WriteJSON(w)
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

func (dns *DNS) getResolveHandlerJSON(w http.ResponseWriter, r *http.Request) {
	log.Debug("dns resolving")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := r.PathValue("ip")
	if ip == "" {
		http.Error(w, "The IP is empty", http.StatusBadRequest)
		return
	}
	log.Debugf("resolving for '%s'", ip)

	if dns.nodePeer.Address == ip {
		log.Debug("dns peer found")
		dns.nodePeer.WriteJSON(w)
		return
	}

	seeds := dns.seeds.Peers()
	for _, peer := range seeds {
		if peer.Address == ip {
			log.Debug("seed peer found")
			peer.WriteJSON(w)
			return
		}
	}

	nodes := dns.nodes.Peers()
	for _, peer := range nodes {
		if peer.Address == ip {
			log.Debug("node peer found")
			peer.WriteJSON(w)
			return
		}
	}

	log.Debug("no peer found")
	http.Error(w, "No peer found", http.StatusNotFound)
}
