package xmpp

import "encoding/xml"

const (
	Away = "away"
	Chat = "chat"
	DND  = "dnd"
	XA   = "xa"
)

type Presence struct {
	XMLName xml.Name `xml:"presence"`
	StanzaFields
	Container
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
	if status, ok := p.Child(xml.Name{Local: "status"}).(*PresenceStatus); ok {
		return status.Status
	}
	return ""
}

func (p *Presence) Show() string {
	if show, ok := p.Child(xml.Name{Local: "show"}).(*PresenceShow); ok {
		return show.Show
	}
	return ""
}

func (p *Presence) Priority() int {
	if prio, ok := p.Child(xml.Name{Local: "priority"}).(*PresencePriority); ok {
		return prio.Priority
	}
	return 0
}

func (p *Presence) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type Raw Presence
	type comboType struct {
		Raw
		Proxies []proxy `xml:",any"`
	}
	combo := &comboType{}
	if err := dec.DecodeElement(combo, &start); err != nil {
		panic(err)
	}
	*p = Presence(combo.Raw)
	p.XMLName = start.Name
	p.Children = proxyToInterface(combo.Proxies)
	return nil
}

func (p *Presence) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type Raw Presence
	type combo struct {
		Raw
		Children []interface{} `xml:",any"`
	}
	raw := &combo{}
	raw.Raw = Raw(*p)
	start.Name = p.XMLName
	raw.Children = p.Container.Children
	return enc.EncodeElement(raw, start)
}

func initPresence() {
	AddElement(&PresenceShow{})
	AddElement(&PresenceStatus{})
	AddElement(&PresencePriority{})

	AddElement(&Presence{})
}
