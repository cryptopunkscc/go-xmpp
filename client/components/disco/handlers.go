package disco

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/xep0030"
	"github.com/cryptopunkscc/go-xmpp/client"
)

func (r *Disco) HandleStanza(s xmpp.Stanza) {
	if !r.session.Handle(s) {
		xmppc.HandleStanza(r, s)
	}
}

func (r *Disco) HandleIQ(iq *xmpp.IQ) {
	if iq.Type != "get" {
		return
	}
	if iq.Child(&xep0030.QueryInfo{}) != nil {
		if r.InfoRequestHandler == nil {
			return
		}
		r.InfoRequestHandler(&InfoRequest{
			iq:      iq,
			session: r.session,
		})
	}
	if qi, ok := iq.Child(&xep0030.QueryItems{}).(*xep0030.QueryItems); ok {
		if r.ItemsRequestHandler == nil {
			return
		}
		r.ItemsRequestHandler(&ItemsRequest{
			Node:    qi.Node,
			iq:      iq,
			session: r.session,
		})
	}
}
