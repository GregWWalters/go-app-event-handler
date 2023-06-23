package appeventhandler

import (
	"errors"
	"testing"
	"time"
)

func TestNewReflectEventHandler(t *testing.T) {
	ehSmall := NewEventHandler(EventHandlerOpts{Priority: Small})
	eh, ok := ehSmall.(*reflectEventHandler)
	if !ok {
		t.Fatalf("expected EventHandler to be type %T but got %T", eh, ehSmall)
	}
}

func Test_reflectEventHandler_Close(t *testing.T) {
	eh := NewEventHandler(EventHandlerOpts{
		Priority: Small,
		LogFunc: func(event AppEvent) {
			t.Logf("received %#v", event)
		},
		ErrorFunc: func(err error) {
			t.Logf("error handling event: %s", err)
		},
	})

	eventChan := make(chan AppEvent)
	defer close(eventChan)
	done, err := eh.Connect(eventChan)
	if err != nil {
		t.Error("failed to connect to EventHandler:", err)
	}

	select {
	case <-done:
		t.Fatal("EventHandler closed prematurely")
	default:
		eventChan <- appEvent{
			name: "test event",
			ts:   time.Now(),
			data: []byte(`{ "color": "brown" }`),
		}
	}

	if err := eh.Close(); err != nil {
		t.Fatal("Error closing EventHandler:", err)
	}

	reh := eh.(*reflectEventHandler)
	if !reh.eventHandler.closed {
		t.Error("EventHandler closed flag is not set")
	}
	if reh.sourceChannels != nil {
		t.Error("EventHandler source channels are not nil")
	}
	if reh.selectCases != nil {
		t.Error("EventHandler select cases are not nil")
	}

	select {
	case <-done:
		break
	default:
		t.Error("EventHandler done channel is not closed")
	}

	if _, err := eh.Connect(eventChan); !errors.Is(err, ErrorHandlerClosed) {
		t.Errorf("should have received error %q but got %q", ErrorHandlerClosed, err)
	}
}

func Test_reflectEventHandler_Register(t *testing.T) {
	// Create EventHandler
	eh := NewEventHandler(EventHandlerOpts{
		Priority: Small,
		LogFunc: func(event AppEvent) {
			t.Logf("received %#v", event)
		},
		ErrorFunc: func(err error) {
			t.Logf("error handling event: %s", err)
		},
	})
	defer eh.Close()

	// Connect channel to EventHandler
	eventChan := make(chan AppEvent)
	defer close(eventChan)
	done, err := eh.Connect(eventChan)
	if err != nil {
		t.Fatal("error connecting to EventHandler:", err)
	}

	// Register EventFunc
	eventName := "TEST EVENT"
	calledFunc := false
	if err := eh.Register(eventName, func(event AppEvent) error {
		calledFunc = true
		t.Log("test EventFunc called")
		return nil
	}); err != nil {
		t.Fatal("failed to register EventFunc:", err)
	}

	// Send event
	select {
	case <-done:
		t.Fatal("EventHandler closed before AppEvent could be sent")
	default:
		eventChan <- appEvent{
			name: eventName,
			ts:   time.Now(),
			data: nil,
		}
	}

	// Check if EventFunc was called
	if !calledFunc {
		t.Error("event function was not called")
	}
}

func Test_reflectEventHandler_Deregister(t *testing.T) {
	// TODO
	t.SkipNow()
}

func Test_reflectEventHandler_handle(t *testing.T) {
	// TODO
	t.SkipNow()
}
