package bot

import (
	"fmt"
)
import "github.com/cryptopunkscc/go-xmpp"

// Context stores information about messages's context
type Context struct {
	Sender  xmpp.JID
	session xmpp.Session
}

// Reply sends back a reply
func (ctx *Context) Reply(format string, a ...interface{}) error {
	return ctx.session.Write(&xmpp.Message{
		To:   ctx.Sender,
		Type: "chat",
		Body: fmt.Sprintf(format, a...),
	})
}
