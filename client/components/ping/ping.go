package ping

import (
	"fmt"
	"sync"
	"time"

	"github.com/cryptopunkscc/go-xmpp"
	"github.com/cryptopunkscc/go-xmpp/xep0199"
	"github.com/cryptopunkscc/go-xmpp/client"
)

// Check if Ping satisfies Handler interface
var _ xmppc.Handler = &Ping{}

const defaultInterval = 60 * time.Second
const defaultTimeout = 30 * time.Second

// LatencyHandler receives latency information
type LatencyHandler func(time.Duration)

// Ping structure holds ping service state
type Ping struct {
	session  xmppc.Session
	stopCh   chan bool
	Interval time.Duration
	Timeout  time.Duration
	pending  map[string]*pending
	times    map[string]time.Time
	LatencyHandler
	mu sync.Mutex
}

// Online implements Handler interface
func (ping *Ping) Online(s xmppc.Session) {
	ping.session = s
	ping.start()
}

// Offline implements Handler interface
func (ping *Ping) Offline(err error) {
	ping.stop()
	ping.session = nil
}

// HandleStanza implements Handler interface
func (ping *Ping) HandleStanza(s xmpp.Stanza) {
	xmppc.HandleStanza(ping, s)
}

// Ping sends a ping request to the server
func (ping *Ping) Ping() {
	iq := &xmpp.IQ{Type: "get"}
	iq.AddChild(&xep0199.Ping{})
	ping.session.Write(iq)
	ping.addPing(iq.GetID())
}

// HandleIQ checks for incoming pongs
func (ping *Ping) HandleIQ(iq *xmpp.IQ) {
	ping.mu.Lock()
	defer ping.mu.Unlock()

	if iq.Type != "result" {
		return
	}
	p, ok := ping.pending[iq.ID]
	if !ok {
		return
	}
	p.Cancel()
	delete(ping.pending, iq.ID)
	latency := time.Now().Sub(p.Start)
	if ping.LatencyHandler != nil {
		ping.LatencyHandler(latency)
	}
}

func (ping *Ping) addPing(id string) {
	ping.mu.Lock()
	defer ping.mu.Unlock()

	if ping.pending == nil {
		ping.pending = make(map[string]*pending)
	}
	ping.pending[id] = wait(
		ping.timeout(),
		func() { ping.onTimeout(id) },
	)
}

func (ping *Ping) onTimeout(id string) {
	ping.mu.Lock()
	defer ping.mu.Unlock()

	delete(ping.pending, id)
	ping.session.Close(fmt.Errorf("ping timeout"))
}

func (ping *Ping) start() {
	ping.mu.Lock()
	defer ping.mu.Unlock()

	ping.stopCh = make(chan bool)
	go ping.pingLoop()
}

func (ping *Ping) stop() {
	ping.mu.Lock()
	defer ping.mu.Unlock()

	if ping.stopCh != nil {
		ping.stopCh <- true
		ping.stopCh = nil
		for _, p := range ping.pending {
			p.cancel <- true
		}
	}
}

func (ping *Ping) pingLoop() {
	for {
		select {
		case <-time.After(ping.interval()):
			ping.Ping()
		case <-ping.stopCh:
			return
		}
	}
}

func (ping *Ping) timeout() time.Duration {
	if ping.Timeout == 0 {
		return defaultTimeout
	}
	return ping.Timeout
}

func (ping *Ping) interval() time.Duration {
	if ping.Interval == 0 {
		return defaultInterval
	}
	return ping.Interval
}
