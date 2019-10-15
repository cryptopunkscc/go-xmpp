package xmppc

import xmpp "github.com/cryptopunkscc/go-xmpp"

// Config represents XMPP client configuration
type Config struct {
	JID      xmpp.JID
	Password string
}
