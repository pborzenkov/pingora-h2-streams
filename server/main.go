package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var maxStreams = flag.Uint("max-streams", 64, "")
var port = flag.Int("port", 8080, "port to listen on")

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("new req, proto = %s", r.Proto)
	var counter int
	for {
		fmt.Fprintf(w, "%d", counter)
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		counter++

		select {
		case <-r.Context().Done():
			return
		case <-time.After(time.Second * 5):
		}
	}
}

func main() {
	flag.Parse()

	h2s := &http2.Server{
		MaxConcurrentStreams: uint32(*maxStreams),
	}
	h1s := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", *port),
		Handler: h2c.NewHandler(http.HandlerFunc(handler), h2s),
		ConnState: func(conn net.Conn, state http.ConnState) {
			log.Printf("conn %q, state %q", conn.RemoteAddr().String(), state)
		},
	}
	log.Printf("HTTP server listening on 127.0.0.1:%d", *port)
	if err := h1s.ListenAndServe(); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
