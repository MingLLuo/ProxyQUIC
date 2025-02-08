package h2h3_server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"golang.org/x/net/http2"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestStartH2Server(t *testing.T) {
	go func() {
		if err := StartH2Server("localhost:8080", "localhost:8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("StartH2Server() error = %v", err)
		}
	}()
	time.Sleep(1 * time.Second)
	caCert, err := os.ReadFile("cert.pem")
	if err != nil {
		t.Fatalf("Failed to load certificate: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	client := &http.Client{
		Transport: &http2.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to make HTTP/2 request: %v", err)
	}
	defer resp.Body.Close()

	if resp.Proto == "HTTP/2.0" {
		log.Println("Server responded with HTTP/2")
	} else {
		t.Fatalf("Server did not respond with HTTP/2, got: %v", resp.Proto)
	}

	h2Server := http.Server{
		Addr: "localhost:8080",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := h2Server.Shutdown(ctx); err != nil {
		t.Errorf("h2Server.Shutdown() error = %v", err)
	}

}
