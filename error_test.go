package appeventhandler

import (
	"errors"
	"testing"
)

type errorWrapper struct {
	message string
	wrapped error
}

func (w errorWrapper) Error() string {
	return w.message
}

func (w errorWrapper) Wrap(e error) error {
	w.wrapped = e
	return w
}

func (w errorWrapper) Unwrap() error {
	return w.wrapped
}

func TestEventHandlerError_Error(t *testing.T) {
	var err error
	errMsg := "EventHandler test"
	err = EventHandlerError(errMsg)
	if actual := err.Error(); actual != errMsg {
		t.Errorf("Expected Error() to return %q but got %q", errMsg, actual)
	}
}

func TestEventHandlerError_Is(t *testing.T) {
	var err error

	err = ErrorHandlerClosed
	if !errors.Is(err, ErrorHandlerClosed) {
		t.Errorf("EventHandlerError %q should be %q", err, ErrorHandlerClosed)
	}

	errMsg := "EventHandler test"
	err = EventHandlerError(errMsg)
	if errors.Is(err, ErrorHandlerClosed) {
		t.Errorf("EventHandlerError %q should not be %q", err, ErrorHandlerClosed)
	}

	err = (errorWrapper{message: "error context"}).Wrap(ErrorEventFunc)
	if errors.Is(err, ErrorHandlerClosed) {
		t.Errorf("error %#v should not be %q", err, ErrorHandlerClosed)
	}
	if !errors.Is(err, ErrorEventFunc) {
		t.Errorf("error %#v should be %q", err, ErrorEventFunc)
	}
}
