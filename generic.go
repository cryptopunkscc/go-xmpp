package xmpp

import "encoding/xml"

// A Generic represents an unknown XML element
type Generic struct {
	XMLName xml.Name
	Attrs   map[string]string `xml:"-"`
	Text    string            `xml:",chardata"`
	Container
}

// UnmarshalXML implements XML decoding
func (g *Generic) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type RawGeneric Generic
	type comboType struct {
		RawGeneric
		Proxies []proxy `xml:",any"`
	}
	combo := &comboType{}
	if err := dec.DecodeElement(combo, &start); err != nil {
		panic(err)
	}
	*g = Generic(combo.RawGeneric)
	g.XMLName = start.Name
	g.Children = proxyToInterface(combo.Proxies)
	g.Attrs = xmlAttrToMap(start.Attr)
	return nil
}

// MarshalXML implements XML encoding
func (g *Generic) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type RawGeneric Generic
	type combo struct {
		RawGeneric
		Children []interface{} `xml:",any"`
	}
	raw := &combo{}
	raw.RawGeneric = RawGeneric(*g)
	start.Name = g.XMLName
	start.Attr = mapToXMLAttr(g.Attrs)
	raw.Children = g.Container.Children
	return enc.EncodeElement(raw, start)
}

func proxyToInterface(p []proxy) []interface{} {
	list := make([]interface{}, 0)
	for _, i := range p {
		list = append(list, i.Object)
	}
	return list
}

func xmlAttrToMap(attrs []xml.Attr) (res map[string]string) {
	res = make(map[string]string)
	for _, attr := range attrs {
		res[attr.Name.Local] = attr.Value
	}
	return
}

func mapToXMLAttr(m map[string]string) (attrs []xml.Attr) {
	attrs = make([]xml.Attr, 0)
	for k, v := range m {
		attrs = append(attrs, xml.Attr{
			Name:  xml.Name{Local: k},
			Value: v,
		})
	}
	return
}
