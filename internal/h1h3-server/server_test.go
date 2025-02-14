package h1h3_server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

func TestStartH1Server(t *testing.T) {
	go func() {
		if err := StartH1Server("localhost:8080", "localhost:8081"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("StartH1Server() error = %v", err)
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
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
	resp, err := client.Get("https://localhost:8080")
	if err != nil {
		log.Fatalf("Failed to make HTTP/1.1 request: %v", err)
	}
	defer resp.Body.Close()

	if resp.Proto == "HTTP/1.1" {
		log.Println("Server responded with HTTP/1.1")
	} else {
		t.Fatalf("Server did not respond with HTTP/1.1, got: %v", resp.Proto)
	}

	h1Server := http.Server{
		Addr: "localhost:8080",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := h1Server.Shutdown(ctx); err != nil {
		t.Errorf("h1Server.Shutdown() error = %v", err)
	}
}
