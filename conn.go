package main

import "net"

type connCloser struct {
	net.Conn

	ln net.Listener
}

func (c connCloser) Close() error {
	defer func() {
		_ = c.ln.Close() // Always return nil error.
	}()
	return c.Conn.Close()
}
