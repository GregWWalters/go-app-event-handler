package appeventhandler

import "time"

// AppEvent is the interface type passed through EventHandler channels. It
// must provide a Name() string method to identify it and match it to an
// EventFunc.
type AppEvent interface {
	Name() string
}

// TimestampAppEvent is an AppEvent that provides a Timestamp method to
// indicate when the event occurred
type TimestampAppEvent interface {
	AppEvent
	Timestamp() time.Time
}

// PayloadAppEvent is an AppEvent that provides a Payload method that returns
// additional data for an EventFunc or for logging
type PayloadAppEvent interface {
	AppEvent
	Payload() any
}

// EventHandler reads from AppEvent channels and performs an EventFunc, if
// registered
type EventHandler interface {
	Connect(<-chan AppEvent) (<-chan struct{}, error)
	Register(string, EventFunc) error
	Deregister(string) bool
	Close() error
}
