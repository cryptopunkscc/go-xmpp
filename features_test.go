package xmpp

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFeaturesUnmarshalling(t *testing.T) {
	xmlData := `<features xmlns='http://etherx.jabber.org/streams'><starttls xmlns='urn:ietf:params:xml:ns:xmpp-tls'><required/></starttls><mechanisms xmlns='urn:ietf:params:xml:ns:xmpp-sasl'><mechanism>PLAIN</mechanism></mechanisms></features>`

	ctx := NewContext(nil)
	ctx.Add(&Features{})

	dec := xml.NewDecoder(strings.NewReader(xmlData))

	msg, err := ctx.Decode(dec)
	assert.Nil(t, err)

	f, ok := msg.(*Features)
	assert.True(t, ok)
	assert.Equal(t, 1, f.ChildCount("starttls"))
	assert.Equal(t, 1, f.ChildCount("mechanisms"))
	assert.Equal(t, 0, f.ChildCount("register"))

	m, ok := f.Child("mechanisms").(*FeatureMechanisms)
	assert.True(t, ok)
	assert.Equal(t, 1, len(m.Mechanisms))
	assert.True(t, m.Include("PLAIN"))
	assert.False(t, m.Include("FAKE"))
}
