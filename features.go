package xmpp

import "encoding/xml"

var FeaturesContext = NewContext(&Generic{})

type Features struct {
	XMLName xml.Name `xml:"http://etherx.jabber.org/streams features"`
	Container
}

type FeatureStartTLS struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-tls starttls"`
}

type FeatureRegister struct {
	XMLName xml.Name `xml:"http://jabber.org/features/iq-register register"`
}

type FeatureBind struct {
	XMLName xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-bind bind"`
}

type FeatureMechanisms struct {
	XMLName    xml.Name `xml:"urn:ietf:params:xml:ns:xmpp-sasl mechanisms"`
	Mechanisms []string `xml:"mechanism"`
}

type FeatureCompression struct {
	XMLName xml.Name `xml:"http://jabber.org/features/compress compression"`
	Methods []string `xml:"method"`
}

func (m *FeatureMechanisms) Include(name string) bool {
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

func (iq *Features) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) (err error) {
	iq.Children, err = FeaturesContext.DecodeAll(dec)
	return
}

func initFeatures() {
	FeaturesContext.Add(&FeatureStartTLS{})
	FeaturesContext.Add(&FeatureRegister{})
	FeaturesContext.Add(&FeatureBind{})
	FeaturesContext.Add(&FeatureMechanisms{})
	FeaturesContext.Add(&FeatureCompression{})

	StreamContext.Add(&Features{})
}
