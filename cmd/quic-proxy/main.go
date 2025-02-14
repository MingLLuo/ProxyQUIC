package main

import (
	"flag"
	"quic-proxy/internal/proxy/h1h3"
)

func main() {
	verbose := flag.Bool("v", true, "should every proxy request be logged to stdout")
	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()
	h1h3.HttpsProxy(verbose, addr)
}
