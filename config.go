package xmpp

const (
	TLSRequired = iota
	TLSPreferred
	TLSDisabled
)

// Config represents XMPP client configuration
type Config struct {
	JID      JID
	Password string
	Host     string
	Logger   Logger
	TLSMode  int
}
