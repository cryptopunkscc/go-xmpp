package xmpp

// Session defines an interface for session
type Session interface {
	JID() JID
	Write(Stanza) error
	Close(error)
	AddFilter(Filter)
}

type session struct {
	conn    *Conn
	handler Handler
	id      UniqueID
	filters []Filter
	in      chan interface{}
	err     chan error
}

// Open opens a new session using provided config and routes session events
// to the provided handler
func Open(handler Handler, cfg *Config) error {
	var err error

	domain := cfg.JID.Domain().String()
	server := NewDNSResolver(domain).ClientAddress()
	conn, err := Connect(server, cfg.JID, cfg.Password)
	if err != nil {
		return err
	}
	s := &session{
		conn:    conn,
		handler: handler,
	}
	if handler != nil {
		handler.Online(s)
	}
	go s.loop()
	return nil
}

// JID returns the JID session is bound to
func (s *session) JID() JID {
	return s.conn.JID()
}

// Write writes a stanza to the XMPP stream
func (s *session) Write(stanza Stanza) error {
	s.injectID(stanza)
	if err := s.applyFilters(stanza); err != nil {
		return err
	}
	return s.conn.Write(stanza)
}

// Close closes the XMPP session with an optional error
func (s *session) Close(err error) {
	s.err <- err
}

// AddFilter adds an outgoing packet filter
func (s *session) AddFilter(f Filter) {
	if s.filters == nil {
		s.filters = make([]Filter, 0)
	}
	s.filters = append(s.filters, f)
}

func (s *session) applyFilters(stanza Stanza) error {
	for _, f := range s.filters {
		if err := f.ApplyFilter(stanza); err != nil {
			return err
		}
	}
	return nil
}

func (s *session) injectID(stanza Stanza) string {
	if stanza.GetID() == "" {
		stanza.SetID(s.id.Next())
	}
	return stanza.GetID()
}

func (s *session) handleNonza(packet interface{}) {
	//bytes, _ := xml.MarshalIndent(packet, "", "  ")
	//log.Println("Nonza received:")
	//for _, l := range strings.Split(string(bytes), "\n") {
	//	log.Println(l)
	//}
}

func (s *session) handlePacket(p interface{}) {
	if stanza, ok := p.(Stanza); ok {
		s.handler.HandleStanza(stanza)
	} else {
		s.handleNonza(p)
	}
}

func (s *session) reader() {
	for {
		msg, err := s.conn.Read()
		if err != nil {
			s.err <- err
			return
		}
		s.in <- msg
	}
}

func (s *session) loop() {
	s.err = make(chan error, 1)
	s.in = make(chan interface{})
	go s.reader()

	for {
		select {
		case err := <-s.err:
			s.close(err)
			return
		case packet := <-s.in:
			s.handlePacket(packet)
		}
	}
}

func (s *session) close(err error) {
	s.conn.Close()
	s.err = nil
	s.in = nil
	s.handler.Offline(err)
}
