package appeventhandler

import "reflect"

// SECTION: Types

// reflectEventHandler implements EventHandler with a Small priority
type reflectEventHandler struct {
	*eventHandler
	sourceChannels []<-chan AppEvent
	selectCases    []reflect.SelectCase
}

// SECTION: Public Functions

func (h *reflectEventHandler) Close() error {
	err := h.eventHandler.Close()
	h.sourceChannels = nil
	h.selectCases = nil
	return err
}

// Connect funnels events from source channels into the select loop by
// appending to the channels and cases slices.
func (h *reflectEventHandler) Connect(events <-chan AppEvent) (done <-chan struct{}, err error) {
	if h.closed {
		return nil, ErrorHandlerClosed
	}
	h.sourceChannels = append(h.sourceChannels, events)
	h.selectCases = append(h.selectCases, reflect.SelectCase{
		Dir: reflect.SelectRecv, Chan: reflect.ValueOf(events),
	})
	return h.done, nil
}

// SECTION: Private Functions

func newReflectEventHandler(h *eventHandler, _ EventHandlerOpts) *reflectEventHandler {
	handler := &reflectEventHandler{
		eventHandler:   h,
		sourceChannels: nil,
		selectCases:    nil,
	}
	go handler.listen()
	return handler
}

// listen consumes from connected AppEvent channels and returns when all
// channels are closed. It removes AppEvent channels when they close.
func (h *reflectEventHandler) listen() {
	for {
		i, v, ok := reflect.Select(h.selectCases)
		if !ok {
			// remove from channels and cases
			oldChannels := h.sourceChannels
			h.sourceChannels = make([]<-chan AppEvent, len(h.sourceChannels)-1)
			copy(h.sourceChannels, oldChannels[:i])
			copy(h.sourceChannels[i:], oldChannels[i+1:])
			oldCases := h.selectCases
			h.selectCases = make([]reflect.SelectCase, len(h.selectCases)-1)
			copy(h.selectCases, oldCases[:i])
			copy(h.selectCases[i:], oldCases[i+1:])
			if len(h.sourceChannels) > 0 {
				continue
			}
			break
		}

		h.handle(v.Interface().(AppEvent))
	}
}
