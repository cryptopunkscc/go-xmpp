package xmpp

import "encoding/xml"

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	ID      string   `xml:"id,attr,omitempty"`
	To      JID      `xml:"to,attr,omitempty"`
	From    JID      `xml:"from,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Container
}

type RosterQuery struct {
	XMLName xml.Name     `xml:"jabber:iq:roster query"`
	Items   []RosterItem `xml:"item"`
}

type RosterItem struct {
	XMLName      xml.Name `xml:"item"`
	JID          JID      `xml:"jid,attr"`
	Name         string   `xml:"name,attr,omitempty"`
	Subscription string   `xml:"subscription,attr,omitempty"`
	Group        []string `xml:"group"`
}

// Result returns true if the IQ type is result
func (iq *IQ) Result() bool {
	return iq.Type == "result"
}

// GetID returns the id field
func (m *IQ) GetID() string { return m.ID }

// GetFrom returns the from field
func (m *IQ) GetFrom() JID { return m.From }

// GetTo returns the to field
func (m *IQ) GetTo() JID { return m.To }

// GetType returns the type field
func (m *IQ) GetType() string { return m.Type }

// GetLang returns the lang field
func (m *IQ) GetLang() string { return m.Lang }

// SetID sets the id field
func (m *IQ) SetID(s string) { m.ID = s }

// SetFrom sets the from field
func (m *IQ) SetFrom(s JID) { m.From = s }

// SetTo sets the to field
func (m *IQ) SetTo(s JID) { m.To = s }

// SetType sets the type field
func (m *IQ) SetType(s string) { m.Type = s }

// SetLang sets the lang field
func (m *IQ) SetLang(s string) { m.Lang = s }

func (iq *IQ) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type Raw IQ
	type comboType struct {
		Raw
		Proxies []proxy `xml:",any"`
	}
	combo := &comboType{}
	if err := dec.DecodeElement(combo, &start); err != nil {
		panic(err)
	}
	*iq = IQ(combo.Raw)
	iq.XMLName = start.Name
	iq.Children = proxyToInterface(combo.Proxies)
	return nil
}

func (iq *IQ) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type Raw IQ
	type combo struct {
		Raw
		Children []interface{} `xml:",any"`
	}
	raw := &combo{}
	raw.Raw = Raw(*iq)
	start.Name = xml.Name{Local: "iq"}
	raw.Children = iq.Container.Children
	return enc.EncodeElement(raw, start)
}

func initIQ() {
	AddElement(&Bind{})
	AddElement(&Error{})
	AddElement(&RosterQuery{})

	AddElement(&IQ{})
}
