package ping

import "time"

type pending struct {
	cancel  chan bool
	Start   time.Time
	Timeout func()
}

func (p *pending) Cancel() {
	p.cancel <- true
}

func wait(d time.Duration, timeoutHandler func()) *pending {
	p := &pending{
		cancel:  make(chan bool),
		Start:   time.Now(),
		Timeout: timeoutHandler,
	}
	go p.wait(d)
	return p
}

func (p *pending) wait(duration time.Duration) {
	select {
	case <-time.After(duration):
		// ping timeout
		if p.Timeout != nil {
			p.Timeout()
		}
	case <-p.cancel:
		return
	}
}
