package xmpp

import "encoding/xml"

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	StanzaFields
	Container
}

type RosterQuery struct {
	XMLName xml.Name     `xml:"jabber:iq:roster query"`
	Items   []RosterItem `xml:"item"`
}

type RosterItem struct {
	XMLName      xml.Name `xml:"item"`
	JID          string   `xml:"jid,attr"`
	Name         string   `xml:"name,attr,omitempty"`
	Subscription string   `xml:"subscription,attr,omitempty"`
	Group        []string `xml:"group"`
}

// Result returns true if the IQ type is result
func (iq *IQ) Result() bool {
	return iq.Type == "result"
}

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
	start.Name = iq.XMLName
	raw.Children = iq.Container.Children
	return enc.EncodeElement(raw, start)
}

func initIQ() {
	AddElement(&Bind{})
	AddElement(&Error{})
	AddElement(&RosterQuery{})

	AddElement(&IQ{})
}
