package main

import (
	"flag"
	"log"
	"quic-proxy/internal/config"
	h2h3_server "quic-proxy/internal/h2h3-server"
	simple_server "quic-proxy/internal/simple-server"
	"quic-proxy/internal/utils"
)

func main() {
	// Command line flags: -mode=simple / -mode=advanced / -mode=h2h3
	mode := flag.String("mode", "simple", "simple/advanced/h2h3")
	flag.Parse()

	cfg, err := config.LoadServerConfig(utils.ConfigPathCreate(*mode, "server", 0))
	if err != nil {
		log.Fatalf("failed to load server config: %v", err)
	}

	log.Printf(cfg.Description)
	if *mode == "simple" {
		err = simple_server.StartServer(cfg.ServerAddr)
	} else if *mode == "h2h3" {
		err = h2h3_server.StartServer(cfg.Http2Addr, cfg.Http3Addr)
	} else {
		log.Fatalf("unsupport mode: %s", *mode)
	}
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
