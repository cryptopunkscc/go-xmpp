package xep0030

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/cryptopunkscc/go-xmpp"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalQueryItems(t *testing.T) {
	xmlData := `
<query xmlns='http://jabber.org/protocol/disco#items'>
    <item jid='people.shakespeare.lit' name='Directory of Characters' xml:lang='en'/>
    <item jid='plays.shakespeare.lit' name='Play-Specific Chatrooms'/>
    <item jid='mim.shakespeare.lit' name='Gateway to Marlowe IM'/>
    <item jid='words.shakespeare.lit' name='Shakespearean Lexicon'/>
    <item jid='globe.shakespeare.lit' name='Calendar of Performances'/>
    <item jid='headlines.shakespeare.lit' name='Latest Shakespearean News'/>
    <item jid='catalog.shakespeare.lit' name='Buy Shakespeare Stuff!'/>
	<item jid='en2fr.shakespeare.lit' name='French Translation Service'/>
</query>`

	dec := xml.NewDecoder(strings.NewReader(xmlData))
	msg, err := xmpp.IQContext.Decode(dec)
	assert.Nil(t, err)

	query, ok := msg.(*QueryItems)
	assert.True(t, ok)
	assert.Equal(t, 8, len(query.Items))
	assert.Equal(t, "people.shakespeare.lit", query.Items[0].JID)
}

func TestUnmarshalQueryInfo(t *testing.T) {
	xmlData := `
<query xmlns='http://jabber.org/protocol/disco#info'>
	<identity category='conference' type='text' name='Play-Specific Chatrooms' xml:lang="en"/>
	<identity category='directory' type='chatroom' name='Play-Specific Chatrooms'/>
</query>`

	dec := xml.NewDecoder(strings.NewReader(xmlData))
	msg, err := xmpp.IQContext.Decode(dec)
	assert.Nil(t, err)

	query, ok := msg.(*QueryInfo)
	assert.True(t, ok)
	assert.Equal(t, 2, len(query.Identities))
	assert.Equal(t, "en", query.Identities[0].Lang)
	assert.Equal(t, "conference", query.Identities[0].Category)
	assert.Equal(t, "text", query.Identities[0].Type)
	assert.Equal(t, "Play-Specific Chatrooms", query.Identities[0].Name)
}
