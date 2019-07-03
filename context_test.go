package xmpp

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type basic struct {
	XMLName xml.Name
	Type    string `xml:"type,attr"`
}

func (b *basic) Name() string {
	return "basic"
}

func (b *basic) Namespace() string {
	return ""
}

type unknown struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

func (u *unknown) Name() string {
	return u.XMLName.Local
}

func (u *unknown) Namespace() string {
	return u.XMLName.Space
}

func TestDecode(t *testing.T) {
	s := basicSpace()
	dec := xmlDecoder(`<basic type="square"></basic>`)

	element, err := s.Decode(dec)
	assert.Nil(t, err)

	e, ok := element.(*basic)
	assert.True(t, ok)
	assert.Equal(t, "basic", e.XMLName.Local)
	assert.Equal(t, "square", e.Type)
}

func TestDecodeAll(t *testing.T) {
	s := basicSpace()
	dec := xmlDecoder(`<basic></basic><basic></basic><random></random>`)

	list, err := s.DecodeAll(dec)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(list))
	assert.Equal(t, "basic", ResolveName(list[0]))
	assert.Equal(t, "basic", ResolveName(list[1]))
	assert.Equal(t, "random", ResolveName(list[2]))
}

func TestDefault(t *testing.T) {
	space := basicSpace()
	dec := xmlDecoder(`<random>content</random>`)

	element, err := space.Decode(dec)
	assert.Nil(t, err)

	typed, ok := element.(*unknown)
	assert.True(t, ok)
	assert.Equal(t, "random", typed.XMLName.Local)
	assert.Equal(t, "content", typed.Content)
}

//
// Some test helpers
//

func xmlDecoder(data string) *xml.Decoder {
	return xml.NewDecoder(strings.NewReader(data))
}

func basicSpace() *Context {
	s := NewContext(&unknown{})
	s.Add(&basic{})
	return s
}
