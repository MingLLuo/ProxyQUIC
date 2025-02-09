package h2h3_server

import (
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"quic-proxy/internal/utils"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/quic-go/qlog"
)

const (
	// certPath = "internal/h2h3-server/cert.pem"
	// keyPath  = "internal/h2h3-server/key.pem"
	certPath = "cert.pem"
	keyPath  = "key.pem"
)

func StartServer(h2Addr string, h3Addr string) error {
	// Generate cert first
	certGenerator := utils.DefaultTLSCertificateGenerator
	certGenerator.CertPath = certPath
	certGenerator.KeyPath = keyPath
	err := certGenerator.Generate()
	if err != nil {
		log.Fatalf("Failed to generate certificate: %v", err)
	}

	_, err = tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}

	// Start H2 server, empty handler, only Alt-svc header set
	go func() {
		err := StartH2Server(h2Addr, h3Addr)
		if err != nil {
			log.Fatalf("Failed to start H2 server: %v", err)
		}
	}()
	// Start H3 server
	return StartH3Server(h3Addr)
}

func StartH2Server(h2Addr, h3Addr string) error {
	_, h3PortInt, err := utils.SplitHostPort(h3Addr)
	if err != nil {
		log.Fatalf("Failed to split h3Addr: %v", err)
	}
	var altSvc []string
	altSvc = append(altSvc, fmt.Sprintf(`h3=":%d";ma=2592000`, h3PortInt))
	// current Path
	log.Printf("Current Path: %s", os.Getenv("PWD"))
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Failed to load certificate: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h2", "h3"},
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Proto == "HTTP/2.0" {
			log.Printf("[h2Server] HTTP/2 Protocol used")
		} else {
			log.Printf("[h2Server] Not using HTTP/2, Protocol: %s", r.Proto)
			http.Error(w, "This server only supports HTTP/2", http.StatusUpgradeRequired)
			return
		}

		// Add Alt-Svc header, remind Client Can use H3
		w.Header().Set("Alt-Svc", strings.Join(altSvc, ","))

		log.Printf("[h2Server] Request From: %s", r.RemoteAddr)
		responseMsg := fmt.Sprint("This is a HTTP/2 Server, try use HTTP/3")
		w.Write([]byte(responseMsg))
	})

	httpServer := &http.Server{
		Addr:      h2Addr,
		Handler:   handler,
		TLSConfig: tlsConfig,
	}
	log.Println("Starting HTTP/2 server on ", h2Addr)
	return httpServer.ListenAndServeTLS(certPath, keyPath)
}

func StartH3Server(serverAddress string) error {
	handler := setupHandler("")
	// QLOGDIR is an environment variable that specifies the directory to store qlog files
	// If QLOGDIR is not set, qlog files will not be generated
	server := http3.Server{
		Handler: handler,
		Addr:    serverAddress,
		QUICConfig: &quic.Config{
			Tracer: qlog.DefaultConnectionTracer,
		},
	}
	// notice, h3 Server will add Alt-Svc automatically
	// See http3.generateAltSvcHeader()
	log.Println("Starting HTTP/3 server on ", serverAddress)
	return server.ListenAndServeTLS(certPath, keyPath)
}

// Size is needed by the /demo/upload handler to determine the size of the uploaded file
type Size interface {
	Size() int64
}

// See https://en.wikipedia.org/wiki/Lehmer_random_number_generator
func generatePRData(l int) []byte {
	res := make([]byte, l)
	seed := uint64(1)
	for i := 0; i < l; i++ {
		seed = seed * 48271 % 2147483647
		res[i] = byte(seed)
	}
	return res
}

func setupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	if len(www) > 0 {
		mux.Handle("/", http.FileServer(http.Dir(www)))
	} else {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Printf("%#v\n", r)
			const maxSize = 1 << 30 // 1 GB
			num, err := strconv.ParseInt(strings.ReplaceAll(r.RequestURI, "/", ""), 10, 64)
			if err != nil || num <= 0 || num > maxSize {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.Write(generatePRData(int(num)))
		})
	}

	mux.HandleFunc("/demo/tile", func(w http.ResponseWriter, r *http.Request) {
		// Small 40x40 png
		w.Write([]byte{
			0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d,
			0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x28, 0x00, 0x00, 0x00, 0x28,
			0x01, 0x03, 0x00, 0x00, 0x00, 0xb6, 0x30, 0x2a, 0x2e, 0x00, 0x00, 0x00,
			0x03, 0x50, 0x4c, 0x54, 0x45, 0x5a, 0xc3, 0x5a, 0xad, 0x38, 0xaa, 0xdb,
			0x00, 0x00, 0x00, 0x0b, 0x49, 0x44, 0x41, 0x54, 0x78, 0x01, 0x63, 0x18,
			0x61, 0x00, 0x00, 0x00, 0xf0, 0x00, 0x01, 0xe2, 0xb8, 0x75, 0x22, 0x00,
			0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae, 0x42, 0x60, 0x82,
		})
	})

	mux.HandleFunc("/demo/tiles", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><head><style>img{width:40px;height:40px;}</style></head><body>")
		for i := 0; i < 200; i++ {
			fmt.Fprintf(w, `<img src="/demo/tile?cachebust=%d">`, i)
		}
		io.WriteString(w, "</body></html>")
	})

	mux.HandleFunc("/demo/echo", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("error reading body while handling /echo: %s\n", err.Error())
		}
		w.Write(body)
	})

	// accept file uploads and return the MD5 of the uploaded file
	// maximum accepted file size is 1 GB
	mux.HandleFunc("/demo/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			err := r.ParseMultipartForm(1 << 30) // 1 GB
			if err == nil {
				var file multipart.File
				file, _, err = r.FormFile("uploadfile")
				if err == nil {
					var size int64
					if sizeInterface, ok := file.(Size); ok {
						size = sizeInterface.Size()
						b := make([]byte, size)
						file.Read(b)
						md5 := md5.Sum(b)
						fmt.Fprintf(w, "%x", md5)
						return
					}
					err = errors.New("couldn't get uploaded file size")
				}
			}
			log.Printf("Error receiving upload: %#v", err)
		}
		io.WriteString(w, `<html><body><form action="/demo/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="uploadfile"><br>
				<input type="submit">
			</form></body></html>`)
	})

	return mux
}
