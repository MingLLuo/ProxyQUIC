package main

import (
	"flag"
	"log"

	"quic-proxy/internal/config"
	h1h3server "quic-proxy/internal/h1h3-server"
	simpleserver "quic-proxy/internal/simple-server"
	"quic-proxy/internal/utils"
)

func main() {
	// Command line flags: -mode=simple / -mode=advanced / -mode=h1h3
	mode := flag.String("mode", "simple", "simple/advanced/h1h3")
	flag.Parse()

	cfg, err := config.LoadServerConfig(utils.ConfigPathCreate(*mode, "server", 0))
	if err != nil {
		log.Fatalf("failed to load server config: %v", err)
	}

	log.Printf(cfg.Description)
	if *mode == "simple" {
		err = simpleserver.StartServer(cfg.ServerAddr)
	} else if *mode == "h1h3" {
		err = h1h3server.StartServer(cfg.ServerAddr, cfg.Http3Addr)
	} else {
		log.Fatalf("unsupport mode: %s", *mode)
	}
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
