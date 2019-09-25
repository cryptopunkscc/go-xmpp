package xmpp

import "encoding/xml"

// MessageContext is a space for all elements defined within the message stanza
var MessageContext = NewContext(&Generic{})

// A Message represents a message stanza
type Message struct {
	XMLName xml.Name `xml:"message"`
	ID      string   `xml:"id,attr,omitempty"`
	To      string   `xml:"to,attr,omitempty"`
	From    string   `xml:"from,attr,omitempty"`
	Type    string   `xml:"type,attr,omitempty"`
	Lang    string   `xml:"lang,attr,omitempty"`
	Container
	*Context
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

type MessageBody struct {
	XMLName xml.Name `xml:"body"`
	Content string   `xml:",chardata"`
}

func (m *Message) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var err error

	startElementToStanza(start, m)
	m.Children, err = m.DecodeAll(dec)

	return err
}

func (m *Message) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	s := stanzaToStartElement(m)
	enc.EncodeToken(s)

	if err := EncodeAll(enc, m.Children); err != nil {
		return err
	}

	return enc.EncodeToken(s.End())
}

func (m *Message) Body() *MessageBody {
	if b := m.Child("body"); b != nil {
		return b.(*MessageBody)
	}
	return nil
}

func initMessage() {
	MessageContext.Add(&MessageBody{})
	StreamContext.Add(&Message{Context: MessageContext})
}
