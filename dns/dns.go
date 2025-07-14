package dns

import (
	"context"
	"net/http"
	"sync"

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
		ctx  context.Context
		quit *chan struct{}
		wg   *sync.WaitGroup
		// cmd        *cobra.Command
		server     *http.Server
		dnsAddress string
		dnsPort    int32
		mode       DNSMode
		nodePeer   *pb.PeerInfo
		seeds      *pb.PeerList
		nodes      *pb.PeerList
	}
)

func NewDNS(
	ctx context.Context,
	quit *chan struct{},
	wg *sync.WaitGroup,
	// cmd *cobra.Command,
	nodePeer *pb.PeerInfo,
	address string,
	port int32,
	mode DNSMode,
) (*DNS, error) {
	// Create a new ServeMux
	mux := http.NewServeMux()

	dns := &DNS{
		ctx:  ctx,
		quit: quit,
		wg:   wg,
		// cmd:        cmd,
		dnsAddress: address,
		dnsPort:    port,
		nodePeer:   nodePeer,
		seeds:      pb.NewPeerList(),
		nodes:      pb.NewPeerList(),
		mode:       mode,
	}

	dns.seeds.Add(&pb.PeerInfo{
		Address:   "10.42.0.104", // BatchNAS
		Port:      8080,
		Id:        "QmUnknown1",
		Mode:      "seed",
		Connected: true,
		Direction: pb.DirectionInbound,
	})

	dns.nodes.Add(&pb.PeerInfo{
		Address:   "10.42.0.101", // BatchDev
		Port:      8080,
		Id:        "QmUnknown2",
		Mode:      "node",
		Connected: true,
		Direction: pb.DirectionInbound,
	})

	dns.nodes.Add(&pb.PeerInfo{
		Address:   "10.42.0.102", // BatchDesk
		Port:      8080,
		Id:        "QmUnknown3",
		Mode:      "node",
		Connected: true,
		Direction: pb.DirectionInbound,
	})

	// Register routes
	switch dns.mode {
	case JSON:
		mux.HandleFunc("/v1/dns", dns.getDNSHandlerJSON)
		mux.HandleFunc("/v1/seeds", dns.getSeedsHandlerJSON)
		mux.HandleFunc("/v1/nodes", dns.getNodesHandlerJSON)
		mux.HandleFunc("/v1/resolve/{ip}", dns.getResolveHandlerJSON)
	case PROTOBUF:
		mux.HandleFunc("/v1/dns", dns.getDNSHandlerProtoBuf)
		mux.HandleFunc("/v1/seeds", dns.getSeedsHandlerProtoBuf)
		mux.HandleFunc("/v1/nodes", dns.getNodesHandlerProtoBuf)
		mux.HandleFunc("/v1/resolve/{ip}", dns.getResolveHandlerProtoBuf)
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
		log.Error("ListenAndServe", err)
		close(*dns.quit)
	}
}

// Shuts down the DNS server
func (dns *DNS) ShutDown() {
	log.Info("dns server shuting down")
	if err := dns.server.Shutdown(dns.ctx); err != nil {
		log.Error("dns shutdown failed", err)
	}
}

// func (dns *DNS) GetMode() DNSMode {
// 	return dns.mode
// }

// DNS: Handle 404
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
