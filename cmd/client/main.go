package main

import (
	"flag"
	"log"
	"quic-proxy/internal/config"
	simple_client "quic-proxy/internal/simple-client"
	"quic-proxy/internal/utils"
)

func main() {
	mode := flag.String("mode", "simple", "simple/h2h3/advanced")
	flag.Parse()

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
	} else if *mode == "h2h3" {
		log.Fatalf("h2h3 mode is not supported in client")
	} else {
		log.Fatalf("unsupported mode: %s", *mode)
	}

}
