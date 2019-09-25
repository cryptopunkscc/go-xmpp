package xmpp

import "encoding/xml"

// A Message represents a message stanza
type Message struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Body    string   `xml:"body,omitempty"`
	Container
}

// GetID returns the id field
func (m *Message) GetID() string { return m.ID }

// GetFrom returns the from field
func (m *Message) GetFrom() string { return m.From }

// GetTo returns the to field
func (m *Message) GetTo() string { return m.To }

// GetType returns the type field
func (m *Message) GetType() string { return m.Type }

// GetLang returns the lang field
func (m *Message) GetLang() string { return m.Lang }

// SetID sets the id field
func (m *Message) SetID(s string) { m.ID = s }

// SetFrom sets the from field
func (m *Message) SetFrom(s string) { m.From = s }

// SetTo sets the to field
func (m *Message) SetTo(s string) { m.To = s }

// SetType sets the type field
func (m *Message) SetType(s string) { m.Type = s }

// SetLang sets the lang field
func (m *Message) SetLang(s string) { m.Lang = s }

func (m *Message) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type Raw Message
	type comboType struct {
		Raw
		Proxies []proxy `xml:",any"`
	}
	combo := &comboType{}
	if err := dec.DecodeElement(combo, &start); err != nil {
		panic(err)
	}
	*m = Message(combo.Raw)
	m.XMLName = start.Name
	m.Children = proxyToInterface(combo.Proxies)
	return nil
}

func (m *Message) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type Raw Message
	type combo struct {
		Raw
		Children []interface{} `xml:",any"`
	}
	raw := &combo{}
	raw.Raw = Raw(*m)
	start.Name = m.XMLName
	raw.Children = m.Container.Children
	return enc.EncodeElement(raw, start)
}

func initMessage() {
	AddElement(&Message{})
}
