package http

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// HandleRequestAndRedirect 处理客户端请求并转发
func HandleRequestAndRedirect(w http.ResponseWriter, req *http.Request) {
	// TODO: req.Host 可能为空吗？
	log.Printf("[PROXY] Received request: %s %s", req.Method, req.URL.String())

	// 1. 解析目标 URL
	targetURL := req.URL
	if !targetURL.IsAbs() {
		// 处理相对 URL（例如：/path）
		targetURL, _ = url.Parse("http://" + req.Host + req.RequestURI)
	}

	proxyReq, err := http.NewRequest(req.Method, targetURL.String(), req.Body)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Printf("Error creating request: %v", err)
		return
	}

	// 3. 拷贝请求头
	copyHeader(proxyReq.Header, req.Header)

	// 4. 发送请求到目标服务器
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Failed to reach target server", http.StatusBadGateway)
		log.Printf("Error forwarding request: %v", err)
		return
	}
	defer resp.Body.Close()

	// 5. 拷贝响应头和响应体，返回给客户端
	copyHeader(w.Header(), resp.Header)
	// 测试一下 Alt-svc的标头自定义
	w.Header().Add("Alt-Svc", "quic=\":443\"")
	// 关闭 HTTP 的长连接
	// w.Header().Set("Connection", "close")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	log.Printf("[PROXY] Response sent back to client with status: %d", resp.StatusCode)
}

// copyHeader 拷贝 HTTP 头信息
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
