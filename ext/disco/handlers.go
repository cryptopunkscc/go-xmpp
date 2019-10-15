package disco

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/xep0030"
)

func (r *Disco) HandleStanza(s xmpp.Stanza) {
	if !r.session.Handle(s) {
		xmpp.HandleStanza(r, s)
	}
}

func (r *Disco) HandleIQ(iq *xmpp.IQ) {
	if iq.Type != "get" {
		return
	}
	if iq.Child(&xep0030.QueryInfo{}) != nil {
		if InfoRequestHandler == nil {
			return
		}
		InfoRequestHandler(&InfoRequest{
			iq:      iq,
			session: r.session,
		})
	}
	if qi, ok := iq.Child(&xep0030.QueryItems{}).(*xep0030.QueryItems); ok {
		if ItemsRequestHandler == nil {
			return
		}
		ItemsRequestHandler(&ItemsRequest{
			Node:    qi.Node,
			iq:      iq,
			session: r.session,
		})
	}
}
