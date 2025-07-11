package node

import (
	"fmt"
	"strings"

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	"github.com/Friends-Of-Noso/NosoGo/dns"
	log "github.com/Friends-Of-Noso/NosoGo/logger"
)

const (
	cDNSPortFlag = "dns-port"
)

var (
	address string
)

func (n *Node) runModeDNS() {
	log.Debug("entering runModeDNS")
	log.Debugf("node ID: %s", n.p2pHost.ID())

	for key, value := range n.p2pHost.Addrs() {
		log.Debugf("address: %d, %s", key, value)
		address = fmt.Sprintf("%s", value)
		if !strings.Contains(address, "127.0.0.1") && !strings.Contains(address, "localhost") {
			log.Debugf("found a good one: '%s'", address)
			splitN := strings.Split(address, "/")
			log.Debugf("splitN[2]: '%v'", splitN[2])
			address = splitN[2]
		}
	}

	log.Infof("node(dns): Listening on %s/p2p/%s", n.p2pHost.Addrs()[0], n.p2pHost.ID())
	// if config.LogLevel == "debug" {
	// 	//
	// }

	err := checkPort(n.dnsPort, cDNSPortFlag, cfg.DefaultDNSPort)
	if err != nil {
		log.Error("error checking port", err)
		n.Shutdown()
		return
	}

	nodeID := fmt.Sprintf("%s", n.p2pHost.ID())
	dnsServer, err := dns.NewDNS(
		n.ctx,
		n.wg,
		n.cmd,
		n.dnsAddress,
		n.dnsPort,
		n.address,
		n.port,
		nodeID,
		dns.JSON,
	)
	if err != nil {
		log.Error("could not create DNS server", err)
		n.Shutdown()
		return
	}

	n.dns = dnsServer
	log.Debug("starting DNS server")
	n.wg.Add(1)
	go n.dns.Start()
}

func (n *Node) shutdownDNS() {
	n.dns.ShutDown()
}
