package dns

import (
	"net/http"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
)

func (dns *DNS) getDNSHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	peerList := &pb.PeerList{}
	peerList.Add(dns.nodePeer)
	peerList.WriteProtobuf(w)
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

func (dns *DNS) getResolveHandlerProtoBuf(w http.ResponseWriter, r *http.Request) {
	log.Debug("dns resolving")
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ip := r.PathValue("ip")
	log.Debugf("Resolving for '%s'", ip)

	if dns.nodePeer.Address == ip {
		log.Debug("dns peer found")
		dns.nodePeer.WriteProtoBuf(w)
		return
	}

	seeds := dns.seeds.Peers()
	for _, peer := range seeds {
		if peer.Address == ip {
			log.Debug("seed peer found")
			peer.WriteProtoBuf(w)
			return
		}
	}

	nodes := dns.nodes.Peers()
	for _, peer := range nodes {
		if peer.Address == ip {
			log.Debug("node peer found")
			peer.WriteProtoBuf(w)
			return
		}
	}

	log.Debug("no peer found")
	http.Error(w, "No peer found", http.StatusNotFound)
}
