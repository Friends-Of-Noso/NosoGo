package node

import (
	"context"
	"strconv"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/syndtr/goleveldb/leveldb"

	log "github.com/Friends-Of-Noso/NosoGo/logger"
)

type Node struct {
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
	address string
	port    int
	host    *host.Host
	db      *leveldb.DB
	peers   []peer.AddrInfo
}

func NewNode(
	ctx context.Context,
	cancel context.CancelFunc,
	wg *sync.WaitGroup,
	address string,
	port int,
	dbPath string,
) (*Node, error) {
	db, err := leveldb.OpenFile(dbPath, nil)
	if err != nil {
		return nil, err
	}

	host, err := libp2p.New(
		libp2p.ListenAddrStrings("/ip4/" + address + "/tcp/" + strconv.Itoa(port)),
		// libp2p.DisableRelay(),
	)
	if err != nil {
		return nil, err
	}

	return &Node{
		ctx:     ctx,
		cancel:  cancel,
		wg:      wg,
		address: address,
		port:    port,
		host:    &host,
		db:      db,
		peers:   make([]peer.AddrInfo, 0),
	}, nil
}

func (n *Node) Start() {
	log.Debug("Node.Start called")
	defer n.wg.Done()

	for {
		select {
		case <-n.ctx.Done():
			log.Debug("Node.Start() exiting")
			return
		default:
			continue
		}
	}
}

func (n *Node) Shutdown() {
	log.Debug("Node.Shutdown called")
	n.db.Close()
	n.cancel()
}
