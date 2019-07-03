package xmpp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJIDDomain(t *testing.T) {
	assert.Equal(t, JID("host.com"), JID("user@host.com/res1").Domain())
}

func TestJIDBare(t *testing.T) {
	assert.Equal(t, JID("user@host.com"), JID("user@host.com/res1").Bare())
}

func TestJIDLocal(t *testing.T) {
	assert.Equal(t, "user", JID("user@host.com/res1").Local())
	assert.Equal(t, "", JID("host.com/res1").Local())
	assert.Equal(t, "", JID("host.com").Local())
}

func TestJIDResource(t *testing.T) {
	assert.Equal(t, "res1", JID("user@host.com/res1").Resource())
	assert.Equal(t, "", JID("user@host.com").Resource())
}
