package xmpp

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIQUnmarshalling(t *testing.T) {
	rawXML := `<iq id='yhc13a95' type='set'><bind xmlns='urn:ietf:params:xml:ns:xmpp-bind'><resource>balcony</resource></bind></iq>`
	dec := xml.NewDecoder(strings.NewReader(rawXML))
	stanza := &Stanza{Stanza: "iq", Context: IQContext}

	assert.NoError(t, dec.Decode(stanza))

	assert.Equal(t, "yhc13a95", stanza.ID)
	assert.Equal(t, "set", stanza.Type)
	assert.NotNil(t, stanza.Child("bind"))
}

func TestStanzaMarshalling(t *testing.T) {
	src := &Stanza{
		Stanza: "iq",
		ID:     "testid",
		Type:   "set",
	}

	src.Add(&Bind{})

	builder := &strings.Builder{}
	enc := xml.NewEncoder(builder)
	assert.NoError(t, enc.Encode(src))
	data := builder.String()

	dec := xml.NewDecoder(strings.NewReader(data))
	item, err := StreamContext.Decode(dec)
	dst := item.(*Stanza)

	assert.Nil(t, err)
	assert.Equal(t, "iq", dst.Stanza)
	assert.Equal(t, "testid", dst.ID)
	assert.Equal(t, "set", dst.Type)
	assert.Equal(t, 1, dst.ChildCount("bind"))
}
