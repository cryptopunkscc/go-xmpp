package xmpp

import "io"

// Config represents XMPP client configuration
type Config struct {
	JID      JID
	Password string
	Log      io.Writer
}
