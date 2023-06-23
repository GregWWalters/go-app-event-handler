package appeventhandler

import (
	"testing"
	"time"
)

type appEvent struct {
	name string
	ts   time.Time
	data []byte
}

func (ae appEvent) Name() string {
	return ae.name
}
func (ae appEvent) Timestamp() time.Time {
	return ae.ts
}
func (ae appEvent) Payload() any {
	return ae.data
}

func TestNewEventHandler(t *testing.T) {
	ehDefault := NewEventHandler(EventHandlerOpts{})
	eh, ok := ehDefault.(*reflectEventHandler)
	if !ok {
		t.Fatalf("expected EventHandler to be type %T but got %T", eh, ehDefault)
	}
	if eh.closed {
		t.Error("EventHandler is closed")
	}
	if eh.done == nil {
		t.Error("EventHandler done channel is nil")
	}
	if eh.eventMap == nil {
		t.Errorf("EventHandler eventMap is nil")
	}
	if eh.logFunc == nil {
		t.Error("logFunc is nil")
	}
	if eh.errorFunc == nil {
		t.Error("errorFunc is nil")
	}
}
