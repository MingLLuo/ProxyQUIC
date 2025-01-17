package simple_server

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

// StartServer 启动HTTP服务
func StartServer(serverAddress string) error {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqData, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		log.Printf("[Server] Received: %s", string(reqData))
		responseMsg := fmt.Sprintf("Hello, client! I got your message: %s", string(reqData))
		_, _ = w.Write([]byte(responseMsg))
	})

	srv := &http.Server{
		Addr:    serverAddress,
		Handler: handler,
	}

	log.Printf("Server is listening on %s...", serverAddress)
	return srv.ListenAndServe()
}
