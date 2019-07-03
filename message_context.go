package xmpp

import "encoding/xml"

// MessageContext is a space for all elements defined within the message stanza
var MessageContext = NewContext(&Generic{})

type MessageBody struct {
	XMLName xml.Name `xml:"body"`
	Content string   `xml:",chardata"`
}

func initMessageContext() {
	MessageContext.Add(&MessageBody{})

	StreamContext.Add(&Stanza{
		Stanza:  "message",
		Context: MessageContext,
	})
}
