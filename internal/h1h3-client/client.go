package h1h3_client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ebi-yade/altsvc-go"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
	"golang.org/x/net/context"
	"quic-proxy/internal/utils"
)

func DoClientRequest(clientAddress, serverAddress, message string) error {
	serverAddress = utils.NormalizeAddress(serverAddress, "https")
	host, port, err := utils.SplitHostPort(clientAddress)
	if err != nil {
		return fmt.Errorf("split client address error: %w", err)
	}
	// Create Https client, bind clientAddress, port 0 means random port
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		},
		Timeout: 5 * time.Second,
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			proxyURL, err := http.ProxyFromEnvironment(&http.Request{
				URL: &url.URL{Host: addr},
			})
			if err != nil {
				return nil, err
			}
			if proxyURL != nil {
				log.Printf("[DEBUG] Using proxy: %s", proxyURL.String())
				return dialer.DialContext(ctx, network, proxyURL.Host)
			}
			log.Printf("[DEBUG] Connecting directly to: %s", addr)
			return dialer.DialContext(ctx, network, addr)
		},
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}
	log.Printf("[Client] Created client with address: %s", clientAddress)

	req, err := http.NewRequest("POST", serverAddress, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request error: %w", err)
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalf("Failed to dump response: %v", err)
	}
	fmt.Println("HTTP Response:")
	fmt.Println(string(dump))

	altSvc := resp.Header.Get("Alt-Svc")
	// simple check of h3 support
	if altSvc != "" {
		services, err := altsvc.Parse(altSvc)
		if err != nil {
			return fmt.Errorf("failed to parse Alt-Svc header: %w", err)
		}
		for _, service := range services {
			if service.ProtocolID == "h3" {
				log.Printf("[Client] Found h3 service: %v", service)
				h3ServerAddr := fmt.Sprintf("https://%s:%s", host, service.AltAuthority.Port)
				err = RetryClientRequestInH3(h3ServerAddr, message)
				if err != nil {
					return fmt.Errorf("failed to retry request in h3: %w", err)
				}
				break
			}
		}
	}
	return nil
}

func RetryClientRequestInH3(h3ServerAddr, message string) error {
	// Certain HTTP implementations use the client address for logging or
	// access-control purposes. Since a QUIC client's address might change during a
	// connection (and future versions might support simultaneous use of multiple
	// addresses), such implementations will need to either actively retrieve the
	// client's current address or addresses when they are relevant or explicitly
	// accept that the original address might change.

	caCert, err := os.ReadFile("cert.pem")
	if err != nil {
		return fmt.Errorf("failed to load certificate: %w", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	roundTripper := &http3.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:            caCertPool,
			InsecureSkipVerify: true,
		},
		QUICConfig: &quic.Config{
			Tracer: qlog.DefaultConnectionTracer,
		},
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}
	resp, err := hclient.Post(h3ServerAddr+"/demo/echo", "text/plain", bytes.NewBuffer([]byte(message)))
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalf("Failed to dump response: %v", err)
	}
	fmt.Println("HTTP/3 Response:")
	fmt.Println(string(dump))
	return nil
}
