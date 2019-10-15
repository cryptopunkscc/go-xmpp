package disco

import (
	"github.com/cryptopunkscc/go-xmpp"
)

type ItemsRequestHandler func(*ItemsRequest)

type ItemsRequest struct {
	Node    string
	iq      *xmpp.IQ
	session xmpp.Session
}

func (r *ItemsRequest) JID() xmpp.JID {
	return r.iq.From
}

func (r *ItemsRequest) Respond(items *Items) {
	response := r.iq.Response(items.queryItems())
	r.session.Write(response)
}
