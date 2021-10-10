package main

import (
	"io"
	"net"
	"sync"
)

// SingleConnListener is a net.Listener implementation returning the provided net.Conn.
type SingleConnListener struct {
	addr   net.Addr
	once   sync.Once
	connCh chan net.Conn
}

// NewSingleConnListener creates a new SingleConnListener instance with the provided net.Conn.
func NewSingleConnListener(conn net.Conn) *SingleConnListener {
	listener := &SingleConnListener{
		addr:   conn.LocalAddr(),
		connCh: make(chan net.Conn, 1),
	}

	listener.connCh <- connCloser{
		Conn: conn,
		ln:   listener,
	}

	return listener
}

// Accept waits for and returns the next connection to the ln.
func (s *SingleConnListener) Accept() (net.Conn, error) {
	conn, ok := <-s.connCh
	if !ok {
		return nil, io.EOF
	}

	return conn, nil
}

// Close closes the ln.
// Any blocked Accept operations will be unblocked and return errors.
func (s *SingleConnListener) Close() error {
	s.once.Do(func() {
		close(s.connCh)
	})
	return nil
}

// Addr returns the ln's network address.
func (s *SingleConnListener) Addr() net.Addr {
	return s.addr
}
