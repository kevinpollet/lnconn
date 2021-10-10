package main

import (
	"errors"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = ln.Close() }()

	handler := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		_, _ = io.WriteString(rw, "Hello")
	})

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			err := http.Serve(NewSingleConnListener(conn), handler)
			if err != nil && !errors.Is(err, io.EOF) {
				log.Println(err)
				return
			}
		}()
	}
}
