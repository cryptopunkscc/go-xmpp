package xmpp

import (
	"bytes"
	"encoding/base64"
)

// PlainAuth implements a plain SASL authentication
type PlainAuth struct {
	Username string
	Password string
}

// NewPlainAuthenticator instantiates a PLAIN authenticator using provided credentials
func NewPlainAuthenticator(creds Credentials) Authenticator {
	return &PlainAuth{
		Username: creds.Username,
		Password: creds.Password,
	}
}

// Name returns the name of the authenticator
func (auth *PlainAuth) Name() string {
	return "PLAIN"
}

// Data returns authentication data encoded for PLAIN mechanism
func (auth *PlainAuth) Data() string {
	dataBytes := bytes.Join([][]byte{
		{},
		[]byte(auth.Username),
		[]byte(auth.Password),
	}, []byte{0})

	return base64.StdEncoding.EncodeToString(dataBytes)
}

// Challenge satisfies Authenticator interface. It returns an empty string since
// PLAIN mechanism doesn't support challenges.
func (auth *PlainAuth) Challenge(string) string {
	return ""
}
