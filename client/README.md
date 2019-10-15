# XMPP Client in Go

This package provides a basic XMPP client implementation.

WARNING: This is an early stage project. Things *will* get broken.

## Quick start

```go
package main

import (
	"fmt"
	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmppc"
)

var quit chan bool

type Handler struct {
	session xmppc.Session
}

func (h *Handler) Online(s xmppc.Session) {
	fmt.Println("Connected", s.JID())
	h.session = s
}

func (h *Handler) Offline(err error) {
	fmt.Println("Disconnected", err)
	quit <- true
}

func (h *Handler) HandleStanza(s xmpp.Stanza) {
	// Handle an incoming stanza here
}

func main() {
	quit := make(chan bool)
	err := xmppc.Open(&Handler{}, &xmppc.Config{
		JID:      "user@host.com",
		Password: "password",
	})
	if err != nil {
		panic(err)
	}
	<-quit
}
```