package roster

import (
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
)

type RosterItem struct {
	JID          xmpp.JID
	Name         string
	Subscription string
}

type Roster struct {
	session *xmppc.Callbacks
}

func (r *Roster) Online(s xmppc.Session) {
	r.session = &xmppc.Callbacks{Session: s}
}

func (r *Roster) Offline(error) {
	r.session = nil
}

func (r *Roster) HandleStanza(s xmpp.Stanza) {
	if !r.session.Handle(s) {
		xmppc.HandleStanza(r, s)
	}
}

// FetchRoster fetches the roster from the server
func (r *Roster) FetchRoster(res func([]*RosterItem)) {
	request := &xmpp.IQ{Type: "get"}
	request.AddChild(&xmpp.RosterQuery{})
	r.session.WriteIQ(request, func(iq *xmpp.IQ) {
		res(iqToList(iq))
	})
}

func (r *Roster) Add(jid xmpp.JID, name string) {
	req := &xmpp.IQ{Type: "set"}
	q := &xmpp.RosterQuery{
		Items: make([]xmpp.RosterItem, 0),
	}
	q.Items = []xmpp.RosterItem{
		xmpp.RosterItem{
			Name: name,
			JID:  jid,
		},
	}
	req.AddChild(q)
	r.session.Write(req)
}

func (r *Roster) Remove(jid xmpp.JID) {
	req := &xmpp.IQ{Type: "set"}
	q := &xmpp.RosterQuery{
		Items: make([]xmpp.RosterItem, 0),
	}
	q.Items = []xmpp.RosterItem{
		xmpp.RosterItem{
			JID:          jid,
			Subscription: "remove",
		},
	}
	req.AddChild(q)
	r.session.Write(req)
}

func iqToList(iq *xmpp.IQ) []*RosterItem {
	list := make([]*RosterItem, 0)
	if query, ok := iq.Child(&xmpp.RosterQuery{}).(*xmpp.RosterQuery); ok {
		for _, i := range query.Items {
			list = append(list, &RosterItem{
				JID:          i.JID,
				Name:         i.Name,
				Subscription: i.Subscription,
			})
		}
	}
	return list
}
