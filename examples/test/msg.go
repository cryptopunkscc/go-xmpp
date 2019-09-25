package main

import (
	"encoding/xml"
	"fmt"

	"github.com/cryptopunkscc/go-xmpp"
)

func main() {
	msg := &xmpp.Message{
		To: "user@host.com",
	}

	msg.AddChild(&xmpp.MessageBody{
		Content: "hello",
	})

	bytes, _ := xml.MarshalIndent(msg, "", "  ")
	fmt.Println(string(bytes))
}
