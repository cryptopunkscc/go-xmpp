package xmpp

import (
	"fmt"
	"net"
)

const (
	defaultClientPort = 5222
	defaultServerPort = 5269
)

// resolveSRV tries to resolve XMPP DNS entries for the domain. The service is "server" or "client".
func resolveSRV(domain string, service string) string {
	_, addrs, err := net.LookupSRV("xmpp-"+service, "tcp", domain)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s:%d", addrs[0].Target, addrs[0].Port)
}
