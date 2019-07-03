package xmpp

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

type named struct{}

func (*named) Name() string {
	return "override"
}

func (*named) Namespace() string {
	return "overridens"
}

func TestResolveName(t *testing.T) {
	type reflected struct{}
	type tagged struct {
		XMLName xml.Name `xml:"taggedname"`
	}
	type namespaced struct {
		XMLName xml.Name `xml:"namespace taggedname"`
	}
	type complicated struct {
		XMLName xml.Name `xml:"namespace taggedname,option1,option2"`
	}
	type invalid struct {
		XMLName xml.Name `xml:""`
	}

	assert.Equal(t, "reflected", ResolveName(&reflected{}))
	assert.Equal(t, "reflected", ResolveName(reflected{}))
	assert.Equal(t, "override", ResolveName(&named{}))
	assert.Equal(t, "taggedname", ResolveName(&tagged{}))
	assert.Equal(t, "taggedname", ResolveName(&namespaced{}))
	assert.Equal(t, "taggedname", ResolveName(&complicated{}))
	assert.Equal(t, "invalid", ResolveName(&invalid{}))
}

func TestResolveNamespace(t *testing.T) {
	type unnamed struct{}
	type namespaced struct {
		XMLName xml.Name `xml:"xmpp:namespace taggedname"`
	}
	type complicated struct {
		XMLName xml.Name `xml:"xmpp:namespace taggedname,option1,option2"`
	}

	assert.Equal(t, "overridens", ResolveNamespace(&named{}))
	assert.Equal(t, "", ResolveNamespace(&unnamed{}))
	assert.Equal(t, "xmpp:namespace", ResolveNamespace(&namespaced{}))
	assert.Equal(t, "xmpp:namespace", ResolveNamespace(&complicated{}))
}
