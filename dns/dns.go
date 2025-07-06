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
	JSON DNSMode = iota
	PROTOBUF
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
	Host struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
		Key     string `json:"key"`
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
		Addr:    address, //fmt.Sprintf("%s:%d", address, port),
		Handler: mux,
	}
	dns.server = server

	return dns, nil
}

func (dns *DNS) Start() {
	defer dns.wg.Done()

	log.Infof("dns server: Listening on %s", dns.dnsAddress)
	if err := dns.server.ListenAndServe(); err != nil {
		log.Fatalf("DNS Server crashed: %v", err)
	}
}

func (dns *DNS) ShutDown() {
	log.Info("dns server shuting down")
	if err := dns.server.Shutdown(dns.ctx); err != nil {
		log.Fatalf("DNS Shutdown failed: %v", err)
	}
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
