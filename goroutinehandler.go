package appeventhandler

// SECTION: Types

// goroutineEventHandler implements EventHandler with a Fast priority
type goroutineEventHandler struct {
	*eventHandler
	all chan AppEvent
}

// SECTION: Public Functions

func (h *goroutineEventHandler) Close() error {
	err := h.eventHandler.Close()
	h.eventHandler = nil
	close(h.all)
	return err
}

// Connect funnels events from source channels into the main channel by
// reading from each in its own goroutine.
func (h *goroutineEventHandler) Connect(events <-chan AppEvent) (done <-chan struct{}, err error) {
	if h.closed {
		return nil, ErrorHandlerClosed
	}
	go func() {
		for {
			select {
			case <-h.done:
				return
			case event, open := <-events:
				if !open {
					return
				}
				h.all <- event
			}
		}
	}()
	return h.done, nil
}

// SECTION: Private Functions

func newGoroutineEventHandler(h *eventHandler, _ EventHandlerOpts) *goroutineEventHandler {
	handler := &goroutineEventHandler{
		eventHandler: h,
		all:          make(chan AppEvent),
	}
	go handler.listen()
	return handler
}

// listen consumes from a combined AppEvent channel and returns when the
// channel is closed
func (h *goroutineEventHandler) listen() {
	for event := range h.all {
		h.handle(event)
	}
}
