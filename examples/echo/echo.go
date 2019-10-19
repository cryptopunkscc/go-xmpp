package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/ext/ping"
	"github.com/cryptopunkscc/go-xmpp/ext/presence"
)

/*
	A simple echo XMPP client. Takes config file as an argument. The config file should look like this:

	{
		"jid": "<jid>",
		"password": "<password>"
	}
*/

// App holds all application state
type App struct {
	xmpp.Broadcast
	session  xmpp.Session
	quit     chan bool
	presence presence.Presence
}

type config struct {
	JID      xmpp.JID `json:"jid"`
	Password string   `json:"password"`
}

// Online is called when an XMPP session begins (successful login and bind)
func (app *App) Online(s xmpp.Session) {
	log.Printf("Connected as %s.\n", s.JID())
	app.session = s
}

// Offline is called when an XMPP session ends
func (app *App) Offline(err error) {
	if err == nil {
		log.Println("Disconnected.")
	} else {
		log.Printf("Disconnected due to error: %s.\n", err.Error())
	}
	app.quit <- true
}

// HandleStanza is called when a stanza is received
func (app *App) HandleStanza(stanza xmpp.Stanza) {
	xmpp.HandleStanza(app, stanza)
}

// HandleMessage is called when a message packet is received
func (app *App) HandleMessage(msg *xmpp.Message) {
	if msg.Body != "" {
		// Print out the message and echo it back
		log.Printf("[%s] %s\n", msg.From.Bare(), msg.Body)
		err := app.session.Write(&xmpp.Message{
			To:   msg.From,
			Type: msg.Type,
			Body: msg.Body,
		})
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

// Run sets up and runs the app
func (app *App) Run(cfg *config) {
	app.quit = make(chan bool)

	// Add XMPP ext
	app.Add(app)
	app.Add(&ping.Ping{
		LatencyHandler: func(l time.Duration) {
			log.Println("Ping:", l)
		},
	})
	app.Add(&app.presence)

	// Automatically accept presence subscriptions
	app.presence.RequestHandler = func(r *presence.Request) {
		log.Printf("Automatically authorizing %s to see our presence.\n", r.JID)
		r.Allow()
	}

	// Open an XMPP session
	err := xmpp.Open(&app.Broadcast, &xmpp.Config{
		JID:      cfg.JID,
		Password: cfg.Password,
		TLSMode:  xmpp.TLSPreferred,
	})
	if err != nil {
		panic(err)
	}

	// Wait for quit signal
	<-app.quit
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], "<config_file>")
		return
	}
	// Load the config file and run the app
	var cfg config
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}
	app := &App{}
	app.Run(&cfg)
}
