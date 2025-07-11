package dns

import (
	"context"
	"net/http"
	"sync"

	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
	pb "github.com/Friends-Of-Noso/NosoGo/protobuf"
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
		dns         *pb.PeerList
		seeds       *pb.PeerList
		nodes       *pb.PeerList
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
		mode:        mode,
		dns:         pb.NewPeerList(pb.DNS),
		seeds:       pb.NewPeerList(pb.SEEDS),
		nodes:       pb.NewPeerList(pb.NODES),
	}

	dns.dns.Put(&pb.PeerInfo{
		Address: address,
		Port:    int32(port),
		Id:      nodeId,
		Mode:    "dns",
	})

	dns.seeds.Put(&pb.PeerInfo{
		Address: "10.42.0.104", // BatchNAS
		Port:    8080,
		Id:      "QmUnknown",
		Mode:    "seed",
	})

	dns.nodes.Put(&pb.PeerInfo{
		Address: "10.42.0.101", // BatchDev
		Port:    8080,
		Id:      "QmUnknown",
		Mode:    "node",
	})

	dns.nodes.Put(&pb.PeerInfo{
		Address: "10.42.0.102", // BatchDev
		Port:    8080,
		Id:      "QmUnknown",
		Mode:    "node",
	})

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

func (dns *DNS) GetMode() DNSMode {
	return dns.mode
}

// DNS: Handle 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
