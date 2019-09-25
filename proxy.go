package xmpp

import (
	"encoding/xml"
)

type proxy struct {
	Object interface{} `xml:"-"`
}

func (c *proxy) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	e := GetElement(start.Name)
	err := dec.DecodeElement(e, &start)
	c.Object = e
	return err
}
