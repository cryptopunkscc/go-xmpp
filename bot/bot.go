package bot

import (
	"errors"
	"fmt"
	"github.com/cryptopunkscc/go-xmpp"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Compile-time check to make sure Bot implements the Handler interface
var _ xmpp.Handler = &Bot{}

// Engine represents a bot engine
type Engine interface {
	Online(writer ChatWriter)
	Offline(error)
}

type ChatWriter interface {
	Send(xmpp.JID, string, ...interface{}) error
}

// Bot hold the bot state
type Bot struct {
	Sync     bool
	session  xmpp.Session
	commands map[string]reflect.Value
	engine   Engine
}

// Instantiate a new bot using the provided engine
func New(e Engine) *Bot {
	bot := &Bot{}
	bot.setEngine(e)
	return bot
}

// Online is invoked when a new XMPP session is established
func (bot *Bot) Online(s xmpp.Session) {
	bot.session = s
	bot.engine.Online(bot)
}

// HandleStanza handles incoming XMPP stanzas
func (bot *Bot) HandleStanza(s xmpp.Stanza) {
	xmpp.HandleStanza(bot, s)
}

// Offline is invoked when the XMPP session ends
func (bot *Bot) Offline(err error) {
	bot.session = nil
	bot.engine.Offline(err)
}

// HandleMessage handles incoming XMPP messages
func (bot *Bot) HandleMessage(m *xmpp.Message) {
	if m.Body == "" {
		return
	}
	bot.execute(m)
}

// Send a formatted message to a JID
func (bot *Bot) Send(jid xmpp.JID, format string, a ...interface{}) error {
	if bot.session == nil {
		return errors.New("offline")
	}
	msg := &xmpp.Message{
		To:   jid,
		Type: "chat",
		Body: fmt.Sprintf(format, a...),
	}
	return bot.session.Write(msg)
}

func (bot *Bot) setEngine(e Engine) {
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Ptr {
		return
	}
	if v.Elem().Kind() != reflect.Struct {
		return
	}
	bot.engine = e
	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		m := v.Method(i)
		if m.Type().NumIn() == 0 {
			continue
		}
		if m.Type().In(0) != reflect.TypeOf(&Context{}) {
			continue
		}
		name := strings.ToLower(t.Method(i).Name)
		bot.addCommand(name, m)
	}
}

func (bot *Bot) addCommand(name string, cmd reflect.Value) {
	if bot.commands == nil {
		bot.commands = make(map[string]reflect.Value)
	}
	bot.commands[name] = cmd
}

func (bot *Bot) execute(m *xmpp.Message) {
	parts := regexp.MustCompile("'.+'|\".+\"|\\S+").FindAllString(m.Body, -1)

	cmd := parts[0]
	args := parts[1:]

	ctx := &Context{Sender: m.From, session: bot.session}

	fn, ok := bot.commands[cmd]
	if !ok {
		ctx.Reply("unknown command")
		return
	}
	if fn.Type().NumIn() < len(args)+1 {
		ctx.Reply("too many arguments")
		return
	}

	a := []reflect.Value{reflect.ValueOf(ctx)}

	for i := 1; i < fn.Type().NumIn(); i++ {
		at := fn.Type().In(i)
		arg := ""
		if i-1 < len(args) {
			arg = args[i-1]
			if len(arg) > 1 && arg[0] == '"' && arg[len(arg)-1] == '"' {
				arg = arg[1 : len(arg)-1]
			}
		}

		switch at.Kind() {
		case reflect.Int:
			v, _ := strconv.Atoi(arg)
			a = append(a, reflect.ValueOf(v))
		case reflect.String:
			a = append(a, reflect.ValueOf(arg))
		}
	}

	if bot.Sync {
		fn.Call(a)
	} else {
		go fn.Call(a)
	}
}
