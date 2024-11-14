package node

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
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

	host, err := libp2p.New()
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
	log.Debug("Node.Start called")
	n.db.Close()
	n.cancel()
}
