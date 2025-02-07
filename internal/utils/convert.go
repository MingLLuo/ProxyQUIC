package utils

import (
	"fmt"
	"net"
	"strings"
)

// SplitHostPort Convert IP address to host(string) and port(int, if exists)
func SplitHostPort(addr string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, fmt.Errorf("split address error: %w", err)
	}
	portInt, err := net.LookupPort("tcp", portStr)
	if err != nil {
		return "", 0, fmt.Errorf("lookup port error: %w", err)
	}
	return host, portInt, nil
}

// NormalizeAddress Normalize address with head
func NormalizeAddress(addr string, head string) string {
	if !strings.HasPrefix(addr, head+"://") {
		addr = head + "://" + addr
	}
	return addr
}

// ConfigPathCreate Create config path with mode and name
func ConfigPathCreate(mode string, name string, index int) string {
	return fmt.Sprintf("./config/%s/%s_%d.json", mode, name, index)
}
