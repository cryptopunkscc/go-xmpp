package ping

import (
	"encoding/xml"

	"github.com/cryptopunkscc/go-xmpp"
)

type XMPPPing struct {
	XMLName xml.Name `xml:"urn:xmpp:ping ping"`
}

func init() {
	xmpp.AddElement(&XMPPPing{})
}
