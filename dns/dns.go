package dns

import (
	"context"
	"net/http"
	"sync"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

type (
	DNSMode int
)

const (
	PROTOBUF DNSMode = iota
	JSON
)

type (
	DNS struct {
		ctx         context.Context
		wg          *sync.WaitGroup
		cmd         *cobra.Command
		server      *http.Server
		dnsAddress  string
		dnsPort     int
		mode        DNSMode
		nodeAddress multiaddr.Multiaddr
		nodePort    int
		nodeId      string
	}

	// Response structures
	PeerInfo struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
		Id      string `json:"id"`
		Mode    string `json:"mode"`
	}
)

func NewDNS(
	ctx context.Context,
	wg *sync.WaitGroup,
	cmd *cobra.Command,
	address string,
	port int,
	nodeAddress multiaddr.Multiaddr,
	nodePort int,
	nodeId string,
	mode DNSMode,
) (*DNS, error) {
	// Create a new ServeMux
	mux := http.NewServeMux()

	dns := &DNS{
		ctx:         ctx,
		wg:          wg,
		cmd:         cmd,
		dnsAddress:  address,
		dnsPort:     port,
		nodeAddress: nodeAddress,
		nodePort:    nodePort,
		nodeId:      nodeId,
		mode:        JSON,
	}

	// Register routes
	switch mode {
	case JSON:
		mux.HandleFunc("/v1/dns", dns.getDNSHandlerJSON)
		mux.HandleFunc("/v1/seeds", dns.getSeedsHandlerJSON)
		mux.HandleFunc("/v1/nodes", dns.getNodesHandlerJSON)
	case PROTOBUF:
		mux.HandleFunc("/v1/dns", dns.getDNSHandlerProtoBuf)
		mux.HandleFunc("/v1/seeds", dns.getSeedsHandlerProtoBuf)
		mux.HandleFunc("/v1/nodes", dns.getNodesHandlerProtoBuf)
	}
	mux.HandleFunc("/", notFoundHandler)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	dns.server = server

	return dns, nil
}

// Starts the DNS server
func (dns *DNS) Start() {
	defer dns.wg.Done()

	log.Infof("dns server: Listening on %s", dns.dnsAddress)
	if err := dns.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Debug("dns server crashed")
	}
}

// Shuts down the DNS server
func (dns *DNS) ShutDown() {
	log.Info("dns server shuting down")
	if err := dns.server.Shutdown(dns.ctx); err != nil {
		log.Error("dns shutdown failed", err)
	}
}

// DNS: Handle 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
