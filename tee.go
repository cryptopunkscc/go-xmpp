package xmpp

import (
	"io"
)

type tee struct {
	target io.ReadWriter
	logger Logger
}

func (t *tee) Write(p []byte) (n int, err error) {
	t.logger.Sent(p)
	return t.target.Write(p)
}

func (t *tee) Read(p []byte) (n int, err error) {
	n, err = t.target.Read(p)
	t.logger.Received(p[:n])
	return n, err
}
