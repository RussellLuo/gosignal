package gosignal

import (
	"sync"
)

type Arguments map[string]interface{}

// Receiver is an entity that can connect to a signal to receive
// an event when the signal is emitted.
type Receiver interface {
	Receive(Arguments)
}

// FuncReceiver is a helper that transforms f into a Receiver.
func FuncReceiver(f func(Arguments)) Receiver {
	return &funcReceiver{f}
}

type funcReceiver struct {
	f func(Arguments)
}

func (r *funcReceiver) Receive(args Arguments) {
	r.f(args)
}

// The global registry of named signals.
var registry sync.Map

// Signal represents an event, which can be broadcasted to all
// the receivers connected to it.
type Signal struct {
	name      string
	receivers sync.Map
}

// New creates a signal if name is empty (i.e. anonymous signals). If name
// is non-empty (i.e. named signals), it will find the signal (by name)
// from the internal registry first, or creates (and registers) a new signal
// if the signal is not found there.
func New(name string) *Signal {
	if name == "" {
		// Create a unique signal each time for anonymous signals.
		return &Signal{name: name}
	}

	// Try to find the named signal from registry.
	sig, ok := registry.Load(name)
	if !ok {
		// Not found.
		//
		// Find or register the named signal atomically. We re-try to find
		// it first since it may have been set by others concurrently.
		sig, _ = registry.LoadOrStore(name, &Signal{name: name})
	}

	return sig.(*Signal)
}

// Name returns the name of the signal.
func (s *Signal) Name() string {
	return s.name
}

// Receivers returns all the registered receivers of the signal.
func (s *Signal) Receivers() (rs []Receiver) {
	s.receivers.Range(func(r, _ interface{}) bool {
		rs = append(rs, r.(Receiver))
		return true
	})
	return
}

// Send emits the signal.
func (s *Signal) Send(args Arguments) {
	for _, r := range s.Receivers() {
		r.Receive(args)
	}
}

// Connect registers one or more receivers, whose method Receive will be
// invoked each time the signal is emitted.
func (s *Signal) Connect(rs ...Receiver) {
	for _, r := range rs {
		s.receivers.Store(r, true)
	}
}

// Disconnect deregisters one or more receivers.
// Specifying no receivers indicates to disconnect all the receivers.
func (s *Signal) Disconnect(rs ...Receiver) {
	if len(rs) == 0 {
		s.receivers.Range(func(r, _ interface{}) bool {
			s.receivers.Delete(r)
			return true
		})
	} else {
		for _, r := range rs {
			s.receivers.Delete(r)
		}
	}
}
