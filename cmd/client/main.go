package main

import (
	"flag"
	"log"
	"quic-proxy/internal/config"
	simple_client "quic-proxy/internal/simple-client"
	"quic-proxy/internal/utils"
)

func main() {
	mode := flag.String("mode", "simple", "simple or advanced")
	if *mode != "simple" && *mode != "advanced" {
		log.Fatalf("invalid mode: %s", *mode)
	}
	cfg, err := config.LoadClientConfig(utils.ConfigPathCreate(*mode, "client", 0))
	if err != nil {
		log.Fatalf("failed to load client config: %v", err)
	}

	log.Printf(cfg.Description)
	if *mode == "simple" {
		err := simple_client.DoClientRequest(cfg.ClientAddr, cfg.ServerAddr, cfg.ClientMessage)
		if err != nil {
			log.Fatalf("failed to start client: %v", err)
		}
	} else {
		// TODO: Start advanced client
	}

}
