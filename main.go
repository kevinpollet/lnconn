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

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(rw, "Hello")
	})

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go serveHTTP(conn, handler)
	}
}

func serveHTTP(conn net.Conn, handler http.Handler) {
	err := http.Serve(NewSingleConnListener(conn), handler)
	if err != nil && !errors.Is(err, ListenerClosedError{}) {
		log.Println(err)
		return
	}
}
