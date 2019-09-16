package xmpp

import "encoding/xml"

// MessageContext is a space for all elements defined within the message stanza
var MessageContext = NewContext(&Generic{})

type Message struct {
	XMLName xml.Name `xml:"message"`
	Stanza
	Container
	*Context
}

type MessageBody struct {
	XMLName xml.Name `xml:"body"`
	Content string   `xml:",chardata"`
}

func (m *Message) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var err error

	m.copyStartElement(&start)
	m.Children, err = m.DecodeAll(dec)

	return err
}

func (m *Message) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	s := m.startElement("message")
	enc.EncodeToken(s)

	if err := EncodeAll(enc, m.Children); err != nil {
		return err
	}

	return enc.EncodeToken(s.End())
}

func initMessage() {
	MessageContext.Add(&MessageBody{})
	StreamContext.Add(&Message{Context: MessageContext})
}
