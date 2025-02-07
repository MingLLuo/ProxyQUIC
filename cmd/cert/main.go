package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"quic-proxy/internal/utils"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "This tool generates a self-signed TLS certificate.")
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
	}

	host := flag.String("host", "localhost,127.0.0.1", "Comma-separated hostnames and IPs for certificate generation")
	validFrom := flag.String("start-date", "", "Certificate start date (format: 'Jan 2 15:04:05 2006'), defaults to current time if empty")
	validFor := flag.Duration("duration", 365*24*time.Hour, "Certificate validity duration")
	isCA := flag.Bool("ca", false, "Whether to generate a CA certificate")
	rsaBits := flag.Int("rsa-bits", 2048, "RSA key size in bits")
	ecdsaCurve := flag.String("ecdsa-curve", "P256", "ECDSA curve (options: P224, P256, P384, P521)")
	ed25519Key := flag.Bool("ed25519", false, "Whether to generate an Ed25519 key")
	certPath := flag.String("cert", "cert.pem", "Output path for the certificate")
	keyPath := flag.String("key", "key.pem", "Output path for the private key")

	flag.Parse()

	generator := utils.TLSCertificateGenerator{
		Host:       *host,
		ValidFrom:  *validFrom,
		ValidFor:   *validFor,
		IsCA:       *isCA,
		RsaBits:    *rsaBits,
		EcdsaCurve: *ecdsaCurve,
		Ed25519Key: *ed25519Key,
		CertPath:   *certPath,
		KeyPath:    *keyPath,
	}

	if err := generator.Generate(); err != nil {
		log.Printf("❌ Certificate generation failed: %v", err)
	} else {
		log.Printf("✅ Certificate successfully generated: %s, Private Key: %s", *certPath, *keyPath)
	}
}
