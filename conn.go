package main

import "net"

type connCloser struct {
	net.Conn

	listener net.Listener
}

func (c connCloser) Close() error {
	// Always return a nil error.
	defer c.listener.Close()

	return c.Conn.Close()
}
