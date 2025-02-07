package utils

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

// TLSCertificateGenerator Define a struct to generate self-signed certificate
type TLSCertificateGenerator struct {
	Host       string        // Comma-separated hostnames and IPs to generate a certificate for
	ValidFrom  string        // Valid From (format: "Jan 2 15:04:05 2006")
	ValidFor   time.Duration // Validity period
	IsCA       bool          // Is Certificate Authority
	RsaBits    int           // Bits of RSA key to generate
	EcdsaCurve string        // Type of elliptic curve to use
	Ed25519Key bool          // Use Ed25519 key
	CertPath   string        // Certificate File Path
	KeyPath    string        // Private key File Path
}

func (t *TLSCertificateGenerator) Generate() error {
	if t.Host == "" {
		return errors.New("Missing required hostname")
	}

	var priv any
	var err error
	switch t.EcdsaCurve {
	case "":
		if t.Ed25519Key {
			_, priv, err = ed25519.GenerateKey(rand.Reader)
		} else {
			priv, err = rsa.GenerateKey(rand.Reader, t.RsaBits)
		}
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return errors.New("Unrecognized elliptic curve: " + t.EcdsaCurve)
	}
	if err != nil {
		return errors.New("Failed to generate private key: " + err.Error())
	}

	keyUsage := x509.KeyUsageDigitalSignature
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	var notBefore time.Time
	if t.ValidFrom == "" {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 2 15:04:05 2006", t.ValidFrom)
		if err != nil {
			return errors.New("Failed to parse creation date: " + err.Error())
		}
	}

	notAfter := notBefore.Add(t.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return errors.New("Failed to generate serial number: " + err.Error())
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Custom Self-Signed Certificate"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(t.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if t.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return errors.New("Failed to create certificate: " + err.Error())
	}

	if err := writePemFile(t.CertPath, "CERTIFICATE", derBytes); err != nil {
		return err
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return errors.New("Failed to marshal private key: " + err.Error())
	}

	if err := writePemFile(t.KeyPath, "PRIVATE KEY", privBytes); err != nil {
		return errors.New("Failed to write private key: " + err.Error())
	}

	log.Printf("✅ Certificate generated successfully: %s", t.CertPath)
	log.Printf("✅ Private key generated successfully: %s", t.KeyPath)
	return nil
}

// publicKey Get public key from a private key
func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

// writePemFile Write data to a pem file
func writePemFile(path string, pemType string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.New("Failed to open file" + path + " for writing: " + err.Error())
	}
	defer file.Close()

	if err := pem.Encode(file, &pem.Block{Type: pemType, Bytes: data}); err != nil {
		log.Fatalf("Failed to write data to %s: %v", path, err)
		return err
	}

	return nil
}
