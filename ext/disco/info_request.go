package disco

import (
	"github.com/cryptopunkscc/go-xmpp"
)

type InfoRequestHandler func(*InfoRequest)

type InfoRequest struct {
	iq      *xmpp.IQ
	session xmpp.Session
}

func (r *InfoRequest) JID() xmpp.JID {
	return r.iq.From
}

func (r *InfoRequest) Respond(info *Info) {
	response := r.iq.Response(info.queryInfo())
	r.session.Write(response)
}

