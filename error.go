package appeventhandler

// EventHandlerError is a kind of error encountered by an EventHandler. It
// may be wrapped in another error.
type EventHandlerError string

const (
	// ErrorHandlerClosed is returned when an operation is called on an
	// EventHandler that has been closed
	ErrorHandlerClosed EventHandlerError = "EventHandler closed"

	// ErrorEventFunc is returned when an EventFunc run by the EventHandler
	// returns an error. It will be wrapped or joined with the EventFunc error.
	ErrorEventFunc EventHandlerError = "EventFunc error"
)

// Error implements error for EventHandlerError
func (e EventHandlerError) Error() string {
	return string(e)
}

func (e EventHandlerError) Is(target error) bool {
	if err, ok := target.(EventHandlerError); ok {
		return string(err) == string(e)
	}
	return string(e) == target.Error()
}
