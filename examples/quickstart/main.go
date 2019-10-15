package main

import (
	"fmt"
	"github.com/cryptopunkscc/go-xmpp"
)

var quit chan bool

type Handler struct {
	session xmpp.Session
}

func (h *Handler) Online(s xmpp.Session) {
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
	err := xmpp.Open(&Handler{}, &xmpp.Config{
		JID:      "user@host.com",
		Password: "password",
	})
	if err != nil {
		panic(err)
	}
	<-quit
}
