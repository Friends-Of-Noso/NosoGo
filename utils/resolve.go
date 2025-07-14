package utils

import (
	"fmt"
	"net"

	"github.com/multiformats/go-multiaddr"
)

func ResolveToMultiaddr(address string, port int32) (multiaddr.Multiaddr, error) {
	ips, err := net.LookupIP(address)
	if err != nil || len(ips) == 0 {
		return nil, fmt.Errorf("failed to resolve address '%s': %w", address, err)
	}

	// Try to pick the first IPv4 (if any)
	var ip net.IP
	for _, candidate := range ips {
		if candidate.To4() != nil {
			ip = candidate
			break
		}
	}

	if ip == nil {
		return nil, fmt.Errorf("no IPv4 address found for '%s'", address)
	}

	// Now build the multiaddr using the resolved IP
	return multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", ip.String(), port))
}

func ResolveToString(address string, port int32) (string, error) {
	ips, err := net.LookupIP(address)
	if err != nil || len(ips) == 0 {
		return "", fmt.Errorf("failed to resolve address '%s': %w", address, err)
	}

	// Try to pick the first IPv4 (if any)
	var ip net.IP
	for _, candidate := range ips {
		if candidate.To4() != nil {
			ip = candidate
			break
		}
	}

	if ip == nil {
		return "", fmt.Errorf("no IPv4 address found for %s", address)
	}

	// Now build the multiaddr using the resolved IP
	return fmt.Sprintf("%s:%d", ip.String(), port), nil
}
