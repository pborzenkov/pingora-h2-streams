package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

var maxStreams = flag.Int("max-streams", 100, "")
var url = flag.String("url", "http://127.0.0.1:8081", "")

func main() {
	flag.Parse()

	client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLSContext: func(ctx context.Context, network, addr string, _ *tls.Config) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, network, addr)
			},
		},
	}

	for i := 0; i < *maxStreams; i++ {
		fmt.Printf("Opening stream %d\n", i)
		resp, err := client.Get(*url)
		if err != nil {
			log.Fatalf("failed to GET: %v", err)
		}
		if resp.StatusCode/100 != 2 {
			log.Fatalf("got non-2xx status code: %d/%v", resp.StatusCode, resp.Status)
		}

		go readResp(i, resp)
		time.Sleep(time.Millisecond * 10)
	}

	select {}
}

func readResp(idx int, resp *http.Response) {
	defer resp.Body.Close()

	var buf [128]byte
	for {
		_, err := resp.Body.Read(buf[:])
		if err != nil {
			log.Fatalf("stream %d, failed to read response body: %v", idx, err)
		}
	}
}
