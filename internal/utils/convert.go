package utils

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// SplitHostPort Convert IP address to host(string) and port(int, if exists)
func SplitHostPort(addr string) (string, int, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", 0, fmt.Errorf("split address error: %w", err)
	}
	port_, err := strconv.Atoi(port)
	if err != nil {
		return "", 0, fmt.Errorf("convert port to int error: %w", err)
	}
	return host, port_, nil
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
