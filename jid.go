package xmpp

import "strings"

// JID is described by https://tools.ietf.org/html/rfc6122
//
// jid = [ localpart "@" ] domainpart [ "/" resourcepart ]

// JID is the type used to store JIDs
type JID string

// Valid verifies the validity of the JID [STUB]
func (jid JID) Valid() bool {
	if jid == "" {
		return false
	}
	return true
}

// Domain returns the domain part of the JID
func (jid JID) Domain() JID {
	s := 0

	if i := strings.Index(string(jid), "@"); i != -1 {
		s = i + 1
	}

	e := len(jid)

	if i := strings.Index(string(jid), "/"); i != -1 {
		e = i
	}

	return jid[s:e]
}

// Bare returns the bare JID (ie. without the resource part)
func (jid JID) Bare() JID {
	if i := strings.Index(string(jid), "/"); i != -1 {
		return jid[0:i]
	}
	return jid
}

// Local returns the local part of the JID
func (jid JID) Local() string {
	if i := strings.Index(string(jid), "@"); i != -1 {
		return string(jid[0:i])
	}
	return ""
}

// Resource returns the resource part of the JID
func (jid JID) Resource() string {
	if i := strings.Index(string(jid), "/"); i != -1 {
		return string(jid[i+1:])
	}
	return ""
}

// String satisfies Stringer interface
func (jid JID) String() string {
	return string(jid)
}
