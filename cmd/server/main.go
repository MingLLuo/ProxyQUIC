package main

import (
	"flag"
	"log"
	"quic-proxy/internal/config"
	simple_server "quic-proxy/internal/simple-server"
	"quic-proxy/internal/utils"
)

func main() {
	// Command line flags: -mode=simple / -mode=advanced
	mode := flag.String("mode", "simple", "simple or advanced")
	if *mode != "simple" && *mode != "advanced" {
		log.Fatalf("invalid mode: %s", *mode)
	}

	cfg, err := config.LoadServerConfig(utils.ConfigPathCreate(*mode, "server", 0))
	if err != nil {
		log.Fatalf("failed to load server config: %v", err)
	}

	log.Printf(cfg.Description)
	if *mode == "simple" {
		err := simple_server.StartServer(cfg.ServerAddr)
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	} else {
		// TODO: Start advanced server
	}
}
