package main

import (
	"net"
	"sync"
)

// SingleConnListener is a net.Listener implementation returning the provided net.Conn.
type SingleConnListener struct {
	addr   net.Addr
	once   sync.Once
	connCh chan net.Conn
}

// NewSingleConnListener creates a new SingleConnListener instance.
func NewSingleConnListener(conn net.Conn) *SingleConnListener {
	listener := &SingleConnListener{
		addr:   conn.LocalAddr(),
		connCh: make(chan net.Conn, 1),
	}

	listener.connCh <- connCloser{
		Conn:     conn,
		listener: listener,
	}

	return listener
}

// Accept waits for and returns the next connection to the listener.
func (s *SingleConnListener) Accept() (net.Conn, error) {
	conn, ok := <-s.connCh
	if !ok {
		return nil, ListenerClosedError{}
	}

	return conn, nil
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (s *SingleConnListener) Close() error {
	s.once.Do(func() {
		close(s.connCh)
	})
	return nil
}

// Addr returns the listener's network address.
func (s *SingleConnListener) Addr() net.Addr {
	return s.addr
}

// ConnListener is a net.Listener implementation returning the provided net.Conn.
type ConnListener struct {
	addr   net.Addr
	once   sync.Once
	connCh chan net.Conn
}

// NewConnListener creates a new ConnListener instance.
func NewConnListener(addr net.Addr) *ConnListener {
	return &ConnListener{
		addr:   addr,
		connCh: make(chan net.Conn),
	}
}

// Accept waits for and returns the next connection to the listener.
func (c *ConnListener) Accept() (net.Conn, error) {
	conn, ok := <-c.connCh
	if !ok {
		return nil, ListenerClosedError{}
	}

	return conn, nil
}

// ServeConn serves the given connection and blocks until it is accepted
// through the Accept method.
func (c *ConnListener) ServeConn(conn net.Conn) {
	c.connCh <- conn
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (c *ConnListener) Close() error {
	c.once.Do(func() {
		close(c.connCh)
	})
	return nil
}

// Addr returns the listener's network address.
func (c *ConnListener) Addr() net.Addr {
	return c.addr
}
