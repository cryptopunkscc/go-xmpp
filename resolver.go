package xmpp

import (
	"fmt"
	"net"
)

const (
	defaultClientPort = 5222
	defaultServerPort = 5269
)

// DNSResolver represents a DSN resolver
type DNSResolver struct {
	domain string
}

// NewDNSResolver returns an instance of a DNSResolver
func NewDNSResolver(domain string) *DNSResolver {
	return &DNSResolver{
		domain: domain,
	}
}

// ClientAddress resolves the address XMPP clients should connect to
func (r *DNSResolver) ClientAddress() string {
	_, addrs, err := net.LookupSRV("xmpp-client", "tcp", r.domain)
	if err != nil {
		return fmt.Sprintf("%s:%d", r.domain, defaultClientPort)
	}
	return fmt.Sprintf("%s:%d", addrs[0].Target, addrs[0].Port)
}

// ServerAddress resolves the address XMPP servers should connect to
func (r *DNSResolver) ServerAddress() string {
	_, addrs, err := net.LookupSRV("xmpp-server", "tcp", r.domain)
	if err != nil {
		return fmt.Sprintf("%s:%d", r.domain, defaultServerPort)
	}
	return fmt.Sprintf("%s:%d", addrs[0].Target, addrs[0].Port)
}
