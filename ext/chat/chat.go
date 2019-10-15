package chat

import (
	"github.com/cryptopunkscc/go-xmpp"
)

type MessageHandler func(*Message)

type Chat struct {
	session xmpp.Session
	MessageHandler
}

type Message struct {
	From xmpp.JID
	To   xmpp.JID
	Body string
	Type string
}

func (r *Chat) Online(s xmpp.Session) {
	r.session = s
}

func (r *Chat) Offline(error) {
	r.session = nil
}

func (r *Chat) HandleStanza(s xmpp.Stanza) {
	xmpp.HandleStanza(r, s)
}

func (chat *Chat) HandleMessage(data *xmpp.Message) {
	msg := &Message{
		From: xmpp.JID(data.From),
		To:   xmpp.JID(data.To),
		Type: data.Type,
		Body: data.Body,
	}
	if chat.MessageHandler != nil {
		chat.MessageHandler(msg)
	}
}

func (chat *Chat) SendMessage(to xmpp.JID, body string) {
	msg := &xmpp.Message{
		To:   to,
		Type: "chat",
		Body: body,
	}
	chat.session.Write(msg)
}
