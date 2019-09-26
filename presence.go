package xmpp

import "encoding/xml"

const (
	Away = "away"
	Chat = "chat"
	DND  = "dnd"
	XA   = "xa"
)

type Presence struct {
	XMLName  xml.Name `xml:"presence"`
	ID       string   `xml:"id,attr,omitempty"`
	To       string   `xml:"to,attr,omitempty"`
	From     string   `xml:"from,attr,omitempty"`
	Type     string   `xml:"type,attr,omitempty"`
	Lang     string   `xml:"lang,attr,omitempty"`
	Show     string   `xml:"show,omitempty"`
	Status   string   `xml:"status,omitempty"`
	Priority int      `xml:"priority,omitempty"`
	Container
}

// GetID returns the id field
func (m *Presence) GetID() string { return m.ID }

// GetFrom returns the from field
func (m *Presence) GetFrom() string { return m.From }

// GetTo returns the to field
func (m *Presence) GetTo() string { return m.To }

// GetType returns the type field
func (m *Presence) GetType() string { return m.Type }

// GetLang returns the lang field
func (m *Presence) GetLang() string { return m.Lang }

// SetID sets the id field
func (m *Presence) SetID(s string) { m.ID = s }

// SetFrom sets the from field
func (m *Presence) SetFrom(s string) { m.From = s }

// SetTo sets the to field
func (m *Presence) SetTo(s string) { m.To = s }

// SetType sets the type field
func (m *Presence) SetType(s string) { m.Type = s }

// SetLang sets the lang field
func (m *Presence) SetLang(s string) { m.Lang = s }

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
	start.Name = xml.Name{Local: "presence"}
	raw.Children = p.Container.Children
	return enc.EncodeElement(raw, start)
}

func initPresence() {
	AddElement(&Presence{})
}
