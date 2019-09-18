package xmpp

import "encoding/xml"

const (
	Away = "away"
	Chat = "chat"
	DND  = "dnd"
	XA   = "xa"
)

var PresenceContext = NewContext(&Generic{})

type Presence struct {
	XMLName xml.Name `xml:"presence"`
	Stanza
	Container
	*Context
}

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

func (p *Presence) Status() string {
	if status, ok := p.Child("status").(*PresenceStatus); ok {
		return status.Status
	}
	return ""
}

func (p *Presence) Show() string {
	if show, ok := p.Child("show").(*PresenceShow); ok {
		return show.Show
	}
	return ""
}

func (p *Presence) Priority() int {
	if prio, ok := p.Child("priority").(*PresencePriority); ok {
		return prio.Priority
	}
	return 0
}

func (p *Presence) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var err error

	p.copyStartElement(&start)
	p.Children, err = p.DecodeAll(dec)

	return err
}

func (p *Presence) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	s := p.startElement("presence")
	enc.EncodeToken(s)

	if err := EncodeAll(enc, p.Children); err != nil {
		return err
	}

	return enc.EncodeToken(s.End())
}

func initPresence() {
	PresenceContext.Add(&PresenceShow{})
	PresenceContext.Add(&PresenceStatus{})
	PresenceContext.Add(&PresencePriority{})

	StreamContext.Add(&Presence{Context: PresenceContext})
}
