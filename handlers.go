package xmpp

// Handler defines an interface for XMPP event handlers
type Handler interface {
	Online(Session)
	HandleStanza(Stanza)
	Offline(error)
}

// StanzaHandler defines an interface for specific stanza handler
type StanzaHandler interface{}

type messageHandler interface {
	HandleMessage(*Message)
}

type iqHandler interface {
	HandleIQ(*IQ)
}

type presenceHandler interface {
	HandlePresence(*Presence)
}

// HandleStanza routes a stanza to a typed stanza handler
func HandleStanza(handler StanzaHandler, stanza Stanza) bool {
	switch typed := stanza.(type) {
	case *Message:
		if h, ok := handler.(messageHandler); ok {
			h.HandleMessage(typed)
		}
	case *IQ:
		if h, ok := handler.(iqHandler); ok {
			h.HandleIQ(typed)
		}
	case *Presence:
		if h, ok := handler.(presenceHandler); ok {
			h.HandlePresence(typed)
		}
	default:
		return false
	}
	return true
}
