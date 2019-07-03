package xmpp

import "encoding/xml"

var IQContext = NewContext(&Generic{})

type RosterItem struct {
	XMLName      xml.Name `xml:"item"`
	JID          string   `xml:"jid,attr"`
	Name         string   `xml:"name,attr"`
	Subscription string   `xml:"subscription,attr"`
	Group        []string `xml:"group"`
}

type RosterQuery struct {
	XMLName xml.Name     `xml:"jabber:iq:roster query"`
	Items   []RosterItem `xml:"item"`
}

func initIQContext() {
	IQContext.Add(&Bind{})
	IQContext.Add(&Error{})
	IQContext.Add(&RosterQuery{})

	StreamContext.Add(&Stanza{
		Stanza:  "iq",
		Context: IQContext,
	})
}
