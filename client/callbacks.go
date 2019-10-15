package xmppc

import "github.com/cryptopunkscc/go-xmpp"

type IQCallback func(*xmpp.IQ)
type MessageCallback func(*xmpp.Message)
type PresenceCallback func(*xmpp.Presence)

type Callbacks struct {
	Session
	iqCallbacks       map[string]IQCallback
	messageCallbacks  map[string]MessageCallback
	presenceCallbacks map[string]PresenceCallback
}

func (c *Callbacks) WriteIQ(iq *xmpp.IQ, callback IQCallback) error {
	err := c.Write(iq)
	if err != nil {
		return err
	}
	if c.iqCallbacks == nil {
		c.iqCallbacks = make(map[string]IQCallback)
	}
	c.iqCallbacks[iq.GetID()] = callback
	return nil
}

func (c *Callbacks) WriteMessage(msg *xmpp.Message, callback MessageCallback) error {
	err := c.Write(msg)
	if err != nil {
		return err
	}
	if c.messageCallbacks == nil {
		c.messageCallbacks = make(map[string]MessageCallback)
	}
	c.messageCallbacks[msg.GetID()] = callback
	return nil
}

func (c *Callbacks) WritePresence(p *xmpp.Presence, callback PresenceCallback) error {
	err := c.Write(p)
	if err != nil {
		return err
	}
	if c.presenceCallbacks == nil {
		c.presenceCallbacks = make(map[string]PresenceCallback)
	}
	c.presenceCallbacks[p.GetID()] = callback
	return nil
}

func (c *Callbacks) Handle(s xmpp.Stanza) bool {
	id := s.GetID()
	switch typed := s.(type) {
	case *xmpp.IQ:
		if handler, ok := c.iqCallbacks[id]; ok {
			handler(typed)
			delete(c.iqCallbacks, id)
			return true
		}
	case *xmpp.Message:
		if handler, ok := c.messageCallbacks[id]; ok {
			handler(typed)
			delete(c.messageCallbacks, id)
			return true
		}
	case *xmpp.Presence:
		if handler, ok := c.presenceCallbacks[id]; ok {
			handler(typed)
			delete(c.presenceCallbacks, id)
			return true
		}
	}
	return false
}
