package node

import (
	"strings"

	cfg "github.com/Friends-Of-Noso/NosoGo/config"
	"github.com/Friends-Of-Noso/NosoGo/dns"
	log "github.com/Friends-Of-Noso/NosoGo/logger"
	"github.com/Friends-Of-Noso/NosoGo/utils"
)

const (
	cDNSPortFlag = "dns-port"
)

var (
	address string
)

func (n *Node) runModeDNS() {
	log.Debug("entering runModeDNS")

	log.Debug("Finding a non local IP/Address")
	for key, value := range n.p2pHost.Addrs() {
		log.Debugf("  address: %d, %s", key, value)
		address = value.String()
		if !strings.Contains(address, "127.0.0.1") && !strings.Contains(address, "localhost") {
			log.Debugf("  found a good one: '%s'", address)
			splitN := strings.Split(address, "/")
			log.Debugf("  splitN[2]: '%v'", splitN[2])
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

	log.Debugf("peer Address: %s", n.peer.Address)

	// TODO: This must call with ipv6==true in production
	if publicAddress := utils.GetMyIP(n.ctx, false); publicAddress != "" {
		n.peer.Address = publicAddress
	}
	log.Debugf("peer Public Address: %s", n.peer.Address)
	log.Debugf("peer Port: %d", n.peer.Port)
	log.Debugf("peer ID: %s", n.peer.Id)

	n.peer.Mode = "dns"
	log.Debugf("peer Mode: %s", n.peer.Mode)

	nodeID := n.p2pHost.ID().String()
	dnsServer, err := dns.NewDNS(
		n.ctx,
		n.wg,
		// n.cmd,
		n.dnsAddress,
		n.dnsPort,
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
