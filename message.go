package xmpp

import (
	"encoding/xml"
	"fmt"
)

// A Message represents a message stanza
type Message struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr,omitempty"`
	To      JID      `xml:"to,attr,omitempty"`
	From    JID      `xml:"from,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Subject string   `xml:"subject,omitempty"`
	Body    string   `xml:"body,omitempty"`
	Thread  string   `xml:"thread,omitempty"`
	Container
}

// Reply builds a reply struct to the message
func (m *Message) Reply(format string, a ...interface{}) *Message {
	return &Message{
		To:   m.From,
		Type: m.Type,
		Body: fmt.Sprintf(format, a...),
	}
}

// GetID returns the id field
func (m *Message) GetID() string { return m.ID }

// GetFrom returns the from field
func (m *Message) GetFrom() JID { return m.From }

// GetTo returns the to field
func (m *Message) GetTo() JID { return m.To }

// GetType returns the type field
func (m *Message) GetType() string { return m.Type }

// GetLang returns the lang field
func (m *Message) GetLang() string { return m.Lang }

// SetID sets the id field
func (m *Message) SetID(s string) { m.ID = s }

// SetFrom sets the from field
func (m *Message) SetFrom(s JID) { m.From = s }

// SetTo sets the to field
func (m *Message) SetTo(s JID) { m.To = s }

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
	start.Name = xml.Name{Local: "message"}
	raw.Children = m.Container.Children
	return enc.EncodeElement(raw, start)
}

func initMessage() {
	AddElement(&Message{})
}
