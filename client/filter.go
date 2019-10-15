package xmppc

import xmpp "github.com/cryptopunkscc/go-xmpp"

type Filter interface {
	ApplyFilter(xmpp.Stanza) error
}
