package xep0199

import (
	"encoding/xml"

	"github.com/cryptopunkscc/go-xmpp"
)

type Ping struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

func init() {
	xmpp.AddElement(&Ping{})
}
