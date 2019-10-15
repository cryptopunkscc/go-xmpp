package disco

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/xep0030"
	"github.com/cryptopunkscc/go-xmpp/client"
)

type Disco struct {
	session *xmppc.Callbacks
	InfoRequestHandler
	ItemsRequestHandler
}

func (disco *Disco) RequestInfo(jid xmpp.JID, handler func(*Info)) error {
	stanza := &xmpp.IQ{Type: "get", To: jid}
	stanza.AddChild(&xep0030.QueryInfo{})
	disco.session.WriteIQ(stanza, func(iq *xmpp.IQ) {
		if q, ok := iq.Child(&xep0030.QueryInfo{}).(*xep0030.QueryInfo); ok {
			handler(queryInfoToInfo(q))
		}
	})
	return nil
}

func (disco *Disco) RequestItems(jid xmpp.JID, handler func(*Items)) error {
	iq := &xmpp.IQ{Type: "get", To: jid}
	iq.AddChild(&xep0030.QueryItems{})
	disco.session.WriteIQ(iq, func(iq *xmpp.IQ) {
		if q, ok := iq.Child(&xep0030.QueryInfo{}).(*xep0030.QueryItems); ok {
			handler(queryItemsToItems(q))
		}
	})
	return nil
}

func (r *Disco) Online(s xmppc.Session) {
	r.session = &xmppc.Callbacks{Session: s}
}

func (r *Disco) Offline(error) {
	r.session = nil
}
