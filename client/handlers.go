package xmppc

import xmpp "github.com/cryptopunkscc/go-xmpp"

// Handler defines an interface for XMPP event handlers
type Handler interface {
	Online(Session)
	HandleStanza(xmpp.Stanza)
	Offline(error)
}

// StanzaHandler defines an interface for specific stanza handler
type StanzaHandler interface{}

type messageHandler interface {
	HandleMessage(*xmpp.Message)
}

type iqHandler interface {
	HandleIQ(*xmpp.IQ)
}

type presenceHandler interface {
	HandlePresence(*xmpp.Presence)
}

// HandleStanza routes a stanza to a typed stanza handler
func HandleStanza(handler StanzaHandler, stanza xmpp.Stanza) bool {
	switch typed := stanza.(type) {
	case *xmpp.Message:
		if h, ok := handler.(messageHandler); ok {
			h.HandleMessage(typed)
		}
	case *xmpp.IQ:
		if h, ok := handler.(iqHandler); ok {
			h.HandleIQ(typed)
		}
	case *xmpp.Presence:
		if h, ok := handler.(presenceHandler); ok {
			h.HandlePresence(typed)
		}
	default:
		return false
	}
	return true
}
