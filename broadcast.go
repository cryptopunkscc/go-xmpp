package xmpp

// Broadcast routes XMPP events to mutiple handlers
type Broadcast struct {
	handlers []Handler
}

// Add adds a handler to the router
func (m *Broadcast) Add(h Handler) {
	if m.handlers == nil {
		m.handlers = []Handler{h}
	} else {
		m.handlers = append(m.handlers, h)
	}
}

// Online routes an Online event
func (m *Broadcast) Online(s Session) {
	for _, h := range m.handlers {
		h.Online(s)
	}
}

// Offline routes an Offline event
func (m *Broadcast) Offline(e error) {
	for _, h := range m.handlers {
		h.Offline(e)
	}
}

// HandleStanza routes a stanza
func (m *Broadcast) HandleStanza(stanza Stanza) {
	for _, h := range m.handlers {
		h.HandleStanza(stanza)
	}
}
