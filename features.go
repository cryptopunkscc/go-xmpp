package xmpp

import "encoding/xml"

type Features struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams features"`
	Container
}

func (m *Mechanisms) Include(name string) bool {
	if m.Mechanisms == nil {
		return false
	}
	for _, v := range m.Mechanisms {
		if v == name {
			return true
		}
	}
	return false
}

func (f *Features) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type Raw Features
	type comboType struct {
		Raw
		Proxies []proxy `xml:",any"`
	}
	combo := &comboType{}
	if err := dec.DecodeElement(combo, &start); err != nil {
		panic(err)
	}
	*f = Features(combo.Raw)
	f.XMLName = start.Name
	f.Children = proxyToInterface(combo.Proxies)
	return nil
}

func (f *Features) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type Raw Features
	type combo struct {
		Raw
		Children []interface{} `xml:",any"`
	}
	raw := &combo{}
	raw.Raw = Raw(*f)
	start.Name = f.XMLName
	raw.Children = f.Container.Children
	return enc.EncodeElement(raw, start)
}

func addFeaturesElements() {
	AddElement(&Features{})
}
