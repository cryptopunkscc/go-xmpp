package xmpp

import "encoding/xml"

// IQ represents an IQ stanza
type IQ struct {
	XMLName xml.Name `xml:"iq"`
	ID      string   `xml:"id,attr,omitempty"`
	To      JID      `xml:"to,attr,omitempty"`
	From    JID      `xml:"from,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Container
}

// RosterQuery represents a roster query element
type RosterQuery struct {
	XMLName xml.Name     `xml:"jabber:iq:roster query"`
	Items   []RosterItem `xml:"item"`
}

// RosterItem represents a roster item
type RosterItem struct {
	XMLName      xml.Name `xml:"item"`
	JID          JID      `xml:"jid,attr"`
	Name         string   `xml:"name,attr,omitempty"`
	Subscription string   `xml:"subscription,attr,omitempty"`
	Group        []string `xml:"group,omitempty"`
}

// Result returns true if the IQ type is result
func (iq *IQ) Result() bool {
	return iq.Type == "result"
}

// Response constructs an response IQ containing provided items
func (iq *IQ) Response(items ...interface{}) *IQ {
	r := &IQ{
		ID:   iq.ID,
		From: iq.To,
		To:   iq.From,
		Type: "result",
	}
	for _, i := range items {
		iq.AddChild(i)
	}
	return r
}

// GetID returns the id field
func (iq *IQ) GetID() string { return iq.ID }

// GetFrom returns the from field
func (iq *IQ) GetFrom() JID { return iq.From }

// GetTo returns the to field
func (iq *IQ) GetTo() JID { return iq.To }

// GetType returns the type field
func (iq *IQ) GetType() string { return iq.Type }

// GetLang returns the lang field
func (iq *IQ) GetLang() string { return iq.Lang }

// SetID sets the id field
func (iq *IQ) SetID(s string) { iq.ID = s }

// SetFrom sets the from field
func (iq *IQ) SetFrom(s JID) { iq.From = s }

// SetTo sets the to field
func (iq *IQ) SetTo(s JID) { iq.To = s }

// SetType sets the type field
func (iq *IQ) SetType(s string) { iq.Type = s }

// SetLang sets the lang field
func (iq *IQ) SetLang(s string) { iq.Lang = s }

// UnmarshalXML unmarshals an IQ from XML
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

// MarshalXML marshals IQ to XML
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
