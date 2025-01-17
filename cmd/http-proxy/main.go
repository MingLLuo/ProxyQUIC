package main

import (
	"flag"
	"log"
	"net/http"
	"quic-proxy/internal/config"
	http_proxy "quic-proxy/internal/proxy/http"
	"quic-proxy/internal/utils"
)

func main() {
	// Command line flags: -mode=simple / -mode=advanced
	mode := flag.String("mode", "simple", "simple or advanced")
	if *mode != "simple" && *mode != "advanced" {
		log.Fatalf("invalid mode: %s", *mode)
	}

	cfg, err := config.LoadHttpProxyConfig(utils.ConfigPathCreate(*mode, "http_proxy", 0))
	if err != nil {
		log.Fatalf("failed to load http config: %v", err)
	}

	log.Printf(cfg.Description)
	if *mode == "simple" {
		handler := http.HandlerFunc(http_proxy.HandleRequestAndRedirect)
		if err := http.ListenAndServe(cfg.ProxyAddr, handler); err != nil {
			log.Fatalf("failed to start http proxy: %v", err)
		}
	}
}
