package simple_client

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"quic-proxy/internal/utils"
	"time"
)

// DoClientRequest 发起一次简单的 HTTP 请求，发送 message 到 server
func DoClientRequest(clientAddress, serverAddress, message string) error {
	serverAddress = utils.NormalizeAddress(serverAddress, "http")
	// 解析 clientAddress，分离 host 和 port
	host, port, err := utils.SplitHostPort(clientAddress)
	if err != nil {
		return fmt.Errorf("split client address error: %w", err)
	}
	// 创建 HTTP 客户端，绑定 clientAddress，port 为 0 表示随机端口
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(host),
			Port: port,
		},
		Timeout: 5 * time.Second,
	}

	// As a special case, if req.URL.Host is "localhost" (with or without
	// a port number), then a nil URL and nil error will be returned.
	client := &http.Client{
		Transport: &http.Transport{
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
		},
		Timeout: 10 * time.Second,
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

	// Debug: 打印 response 的完整报文
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalf("Failed to dump response: %v", err)
	}
	fmt.Println(string(dump))
	return nil
}
