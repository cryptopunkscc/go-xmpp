package xmpp

import "encoding/xml"

const (
	Away = "away"
	Chat = "chat"
	DND  = "dnd"
	XA   = "xa"
)

var PresenceContext = NewContext(&Generic{})

type PresenceShow struct {
	XMLName xml.Name `xml:"show"`
	Show    string   `xml:",chardata"`
}

type PresenceStatus struct {
	XMLName xml.Name `xml:"status"`
	Status  string   `xml:",chardata"`
}

type PresencePriority struct {
	XMLName  xml.Name `xml:"priority"`
	Priority int      `xml:",chardata"`
}

func initPresenceContext() {
	PresenceContext.Add(&PresenceShow{})
	PresenceContext.Add(&PresenceStatus{})
	PresenceContext.Add(&PresencePriority{})

	StreamContext.Add(&Stanza{
		Stanza:  "presence",
		Context: PresenceContext,
	})
}
