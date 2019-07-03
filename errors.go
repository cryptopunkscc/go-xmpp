package xmpp

import "errors"

// ErrStreamError indicates invalid stream data
var ErrStreamError = errors.New("stream error")

// ErrEndOfStream indicates the XMPP stream has ended
var ErrEndOfStream = errors.New("end of stream")

var ErrNoPrototype = errors.New("prototype not found")
var ErrEndOfElement = errors.New("end of element")
