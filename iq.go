package xmpp

import "encoding/xml"

var IQContext = NewContext(&Generic{})

type IQ struct {
	XMLName xml.Name `xml:"iq"`
	StanzaFields
	Container
	*Context
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
	var err error

	iq.copyStartElement(&start)
	iq.Children, err = iq.DecodeAll(dec)

	return err
}

func (iq *IQ) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	s := iq.startElement("iq")
	enc.EncodeToken(s)

	if err := EncodeAll(enc, iq.Children); err != nil {
		return err
	}

	return enc.EncodeToken(s.End())
}

func initIQ() {
	IQContext.Add(&Bind{})
	IQContext.Add(&Error{})
	IQContext.Add(&RosterQuery{})

	StreamContext.Add(&IQ{Context: IQContext})
}
