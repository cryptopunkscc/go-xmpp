package xmppc

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-xmpp"
)

// Authenticator defines an interface for SASL authentication
type Authenticator interface {
	Name() string
	Data() string
	Challenge(string) string
}

type authFunc func(Credentials) Authenticator

func contains(list []string, item string) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// bestAuthenticator returns the a factory for the first supported authentication mechanism
func (c *Conn) bestAuthenticator() (authFunc, error) {
	mechanisms, ok := c.Features().Child(&xmpp.Mechanisms{}).(*xmpp.Mechanisms)
	if !ok {
		return nil, errors.New("authentication failed: stream has no auth mechanisms")
	}
	if contains(mechanisms.Mechanisms, "SCRAM-SHA-1") {
		return NewScramSHA1Authenticator, nil
	}
	if contains(mechanisms.Mechanisms, "PLAIN") {
		return NewPlainAuthenticator, nil
	}
	return nil, errors.New("authentication failed: no supported mechanisms found")
}

// authenticate tries to authenticate to the server
func (c *Conn) authenticate(username, password string) error {
	authFunc, err := c.bestAuthenticator()
	if err != nil {
		return errors.New("authentication failed: all mechanisms unsupported")
	}
	auth := authFunc(Credentials{
		Username: username,
		Password: password,
	})

	// Start authentication
	err = c.Write(&xmpp.Auth{
		Mechanism: auth.Name(),
		Data:      auth.Data(),
	})
	if err != nil {
		return err
	}

	for {
		msg, err := c.Read()

		if err != nil {
			return err
		}

		switch typed := msg.(type) {
		case *xmpp.Challenge:
			r := &xmpp.Response{Data: auth.Challenge(typed.Data)}
			err := c.Write(r)
			if err != nil {
				return err
			}
		case *xmpp.Success:
			return nil
		case *xmpp.Error:
			return fmt.Errorf("authentication failed: stream error: %s", typed.Condition)
		case *xmpp.SASLFailure:
			return errors.New("authentication failed: invalid credentials")
		default:
			return fmt.Errorf("authentication failed: unexpected response: %s", msg)
		}
	}
}
