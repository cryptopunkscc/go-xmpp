package xmpp

type Filter interface {
	ApplyFilter(Stanza) error
}
