package xmpp

import (
	"errors"
	"fmt"
)

// Authenticator defines an interface for SASL authentication
type Authenticator interface {
	Name() string
	Data() string
	Challenge(string) string
}

// Credentials holds authentication information
type Credentials struct {
	Username string
	Password string
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
	mechanisms, ok := c.Features().Child(&Mechanisms{}).(*Mechanisms)
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

// Authenticate authenticates to the server and restarts the stream
func (c *Conn) Authenticate(username, password string) error {
	authFunc, err := c.bestAuthenticator()
	if err != nil {
		return errors.New("authentication failed: all mechanisms unsupported")
	}
	auth := authFunc(Credentials{
		Username: username,
		Password: password,
	})

	// Start authentication
	err = c.Write(&Auth{
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
		case *Challenge:
			r := &Response{Data: auth.Challenge(typed.Data)}
			err := c.Write(r)
			if err != nil {
				return err
			}
		case *Success:
			//TODO: Verify success response
			return c.RestartStream(nil)
		case *Error:
			return fmt.Errorf("authentication failed: stream error: %s", typed.Condition)
		case *SASLFailure:
			return errors.New("authentication failed: invalid credentials")
		default:
			return fmt.Errorf("authentication failed: unexpected response: %s", msg)
		}
	}
}
