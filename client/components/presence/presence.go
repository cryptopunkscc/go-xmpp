package presence

import (
	"sync"

	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/client"
)

// Presence holds the state of the presence service
type Presence struct {
	RequestHandler
	UpdateHandler
	session  xmppc.Session
	show     string
	status   string
	priority int
	mu       sync.Mutex
}

type Update struct {
	JID      xmpp.JID
	Online   bool
	Status   string
	Show     string
	Priority int
}

type Request struct {
	p   *Presence
	JID xmpp.JID
}

type RequestHandler func(*Request)
type UpdateHandler func(*Update)

func (req *Request) Allow() {
	req.p.Allow(req.JID)
}

func (req *Request) Deny() {
	req.p.Deny(req.JID)
}

func (p *Presence) Online(s xmppc.Session) {
	p.session = s
	p.broadcastPresence()
}

func (p *Presence) Offline(error) {}

func (p *Presence) HandleStanza(s xmpp.Stanza) {
	xmppc.HandleStanza(p, s)
}

// Subscribe to a user's presence
func (p *Presence) Subscribe(jid xmpp.JID) {
	p.session.Write(&xmpp.Presence{
		Type: "subscribe",
		To:   jid,
	})
}

// Unsubscribe from a user's presence
func (p *Presence) Unsubscribe(jid xmpp.JID) {
	p.session.Write(&xmpp.Presence{
		Type: "unsubscribe",
		To:   jid,
	})
}

// Allow user to subscribe
func (p *Presence) Allow(jid xmpp.JID) {
	p.session.Write(&xmpp.Presence{
		Type: "subscribed",
		To:   jid,
	})
}

// Deny user to subscribe
func (p *Presence) Deny(jid xmpp.JID) {
	p.session.Write(&xmpp.Presence{
		Type: "unsubscribed",
		To:   jid,
	})
}

// SetStatus sets the status text
func (p *Presence) SetStatus(status string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.status = status
	p.broadcastPresence()
}

// SetPriority sets resource priority
func (p *Presence) SetPriority(priority int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.priority = priority
	p.broadcastPresence()
}

// SetShow sets presence type (away, xa, dnd, chat or an empty string)
func (p *Presence) SetShow(show string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.show = show
	p.broadcastPresence()
}

func (p *Presence) broadcastPresence() {
	if p.session == nil {
		return
	}
	p.session.Write(&xmpp.Presence{
		Status:   p.status,
		Show:     p.show,
		Priority: p.priority,
	})
}

func (p *Presence) onSubscriptionRequest(pres *xmpp.Presence) {
	if p.RequestHandler != nil {
		p.RequestHandler(&Request{
			p:   p,
			JID: xmpp.JID(pres.From),
		})
	}
}

func (p *Presence) onStatus(pres *xmpp.Presence) {
	if p.UpdateHandler != nil {
		p.UpdateHandler(&Update{
			JID:      pres.From,
			Online:   (pres.Type == ""),
			Status:   pres.Status,
			Show:     pres.Show,
			Priority: pres.Priority,
		})

	}
}

func (p *Presence) HandlePresence(pres *xmpp.Presence) {
	switch pres.Type {
	case "", "unavailable":
		// Sender is (no longer) available for communication
		p.onStatus(pres)

	case "subscribed":
		// Sender has allowed the recipient to receive their presence

	case "subscribe":
		// Sender wishes to subscribe to the recipient's presence
		p.onSubscriptionRequest(pres)

	case "unsubscribe":
		// Sender is unsubscribing from the receiver's presence

	case "unsubscribed":
		// Subscription request has been denied or a previously granted subscription has been canceled

	case "error":
		// An error has occurred regarding processing of a previously sent presence stanza
	}
}
