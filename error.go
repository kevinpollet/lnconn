package main

// ListenerClosedError is returned by the SingleConnListener Accept method after a call to Close.
type ListenerClosedError struct{}

func (e ListenerClosedError) Error() string {
	return "listener is closed"
}
