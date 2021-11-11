package main

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	l, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	handler := func(rw http.ResponseWriter, req *http.Request) {
		io.WriteString(rw, "Hello")
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go serveHTTP(conn, http.HandlerFunc(handler))
	}
}

func serveHTTP(conn net.Conn, handler http.Handler) {
	err := http.Serve(NewSingleConnListener(conn), handler)
	if err != nil && !errors.Is(err, ListenerClosedError{}) {
		log.Println(err)
		return
	}
}
