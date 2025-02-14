package h1h3

import (
	"crypto/tls"
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
	"net/http/httputil"
)

func HttpsProxy(verbose *bool, addr *string) {
	// 生成 CA 证书
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}
	customCaMitm := &goproxy.ConnectAction{Action: goproxy.ConnectMitm, TLSConfig: goproxy.TLSConfigFromCA(&cert)}
	customCaMitmHttp := &goproxy.ConnectAction{Action: goproxy.ConnectHTTPMitm, TLSConfig: goproxy.TLSConfigFromCA(&cert)}

	var customAlwaysMitm goproxy.FuncHttpsHandler = func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		if ctx.Req.ProtoMajor != 2 {
			log.Printf("req.ProtoMajor: %d", ctx.Req.ProtoMajor)
			return nil, host
		}
		if ctx.Req.Method == "CONNECT" {
			// CONNECT 请求，即 HTTPS 代理
			log.Printf("HTTPS CONNECT request intercepted: %s", host)
			return customCaMitm, host // 使用 MITM
		} else {
			// 非 CONNECT 请求，即 HTTP 请求
			log.Printf("HTTP request intercepted: %s, Method: %s", host, ctx.Req.Method)
			return customCaMitmHttp, host // 不使用 MITM
		}
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.AllowHTTP2 = true
	proxy.OnRequest().HandleConnect(customAlwaysMitm)
	//	⚠️ Note we returned a nil value as the response. If the returned response is not nil, goproxy will discard the request and send the specified response to the client.
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		dump, _ := httputil.DumpRequest(req, true)
		log.Println(string(dump))
		if req.URL.Scheme == "http" {
			return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "HTTPS Required")
		}
		if req.ProtoMajor != 2 {
			log.Printf("req.ProtoMajor: %d", req.ProtoMajor)
			return nil, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "HTTP/2 Required")
		}

		return req, nil
	})
	proxy.Verbose = *verbose
	log.Fatal(http.ListenAndServe(*addr, proxy))
}
