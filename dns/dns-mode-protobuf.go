package dns

import (
	"net/http"
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
