package appeventhandler

// SECTION: Constants

const (
	Default EventHandlerPriority = iota
	Small
	Fast
)

var defaultEventHandlerOpts = EventHandlerOpts{
	LogFunc:   noopLogFunc,
	ErrorFunc: noopErrorFunc,
	Priority:  Default,
}

// SECTION: Types

// EventHandlerPriority indicates whether the EventHandler should prioritize
// handling messages quickly (Fast) or a small memory and goroutine footprint
// (Small)
type EventHandlerPriority uint8

// EventFunc is a function to run when an AppEvent occurs. It is registered to
// an AppEvent Name with Register
type EventFunc func(event AppEvent) error

// EventHandlerOpts provides options for creating an EventHandler
type EventHandlerOpts struct {
	LogFunc   func(event AppEvent)
	ErrorFunc func(error)
	Priority  EventHandlerPriority
}

// eventHandler is a composable type to handle most of the EventHandler
// implementation
type eventHandler struct {
	closed    bool
	done      chan struct{}
	eventMap  map[string]EventFunc
	errorFunc func(error)
	logFunc   func(event AppEvent)
}

// SECTION: Public Functions

func noopErrorFunc(error)        {}
func noopLogFunc(event AppEvent) {}

func NewEventHandler(opts EventHandlerOpts) EventHandler {
	opts = defaultEventHandlerOpts.apply(opts)
	h := &eventHandler{
		closed:    false,
		done:      make(chan struct{}),
		eventMap:  make(map[string]EventFunc),
		errorFunc: opts.ErrorFunc,
		logFunc:   opts.LogFunc,
	}
	switch opts.Priority {
	default:
		fallthrough
	case Small:
		return newReflectEventHandler(h, opts)
	case Fast:
		return newGoroutineEventHandler(h, opts)
	}
}

func (h *eventHandler) Close() error {
	h.closed = true
	close(h.done)
	return nil
}

func (h *eventHandler) Register(name string, fn EventFunc) (err error) {
	if h.closed {
		return ErrorHandlerClosed
	}
	h.eventMap[name] = fn
	return nil
}

func (h *eventHandler) Deregister(name string) bool {
	_, found := h.eventMap[name]
	delete(h.eventMap, name)
	return found
}

// SECTION: Private Functions

func (eho EventHandlerOpts) apply(opts EventHandlerOpts) EventHandlerOpts {
	if opts.LogFunc != nil {
		eho.LogFunc = opts.LogFunc
	}
	if opts.ErrorFunc != nil {
		eho.ErrorFunc = opts.ErrorFunc
	}
	if opts.Priority != Default {
		eho.Priority = opts.Priority
	}
	return eho
}

// handle looks up and runs an EventFunc for an AppEvent
func (h *eventHandler) handle(event AppEvent) {
	if fn, found := h.eventMap[event.Name()]; found {
		if err := fn(event); err != nil {
			h.errorFunc(err)
		}
	}
}
