// Package utils This code references the following packages:
// https://github.com/ebi-yade/altsvc-go
package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// The Service represents HTTP Alternative Services declared in RFC 7838.
// https://datatracker.ietf.org/doc/html/rfc7838#section-3
type Service struct {
	Clear        bool
	ProtocolID   string // ALPN Protocol ID
	AltAuthority AltAuthority

	// https://datatracker.ietf.org/doc/html/rfc7838#section-3.1
	MaxAge  int
	Persist int
}

// AltAuthority represents the host and port of an alternative service.
// alt-authority [uri-host] ":" port
type AltAuthority struct {
	Host string // if null, the host is the same as the origin server
	Port string
}

// Parse Analyzes the input string and returns a list of services.
func Parse(s string) ([]Service, error) {
	if strings.TrimSpace(s) == "clear" {
		return []Service{{Clear: true}}, nil
	}

	rawServices := strings.Split(s, ",")
	services := make([]Service, 0, len(rawServices))
	for _, rawSvc := range rawServices {
		rawSvc = strings.TrimSpace(rawSvc)
		if rawSvc == "" {
			continue
		}
		svc, err := parseService(rawSvc)
		if err != nil {
			return nil, err
		}
		services = append(services, svc)
	}
	return services, nil
}

// parseService Parses a single service string.
//
//	h3=":8081";ma=2592000;persist=1
func parseService(s string) (Service, error) {
	var svc Service
	parts := strings.Split(s, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		key, value, ok := strings.Cut(part, "=")
		if !ok {
			return svc, fmt.Errorf("invalid parameter: %q", part)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		switch key {
		case "ma":
			ma, err := strconv.Atoi(value)
			if err != nil {
				return svc, fmt.Errorf("invalid value for 'ma': %q", value)
			}
			svc.MaxAge = ma
		case "persist":
			p, err := strconv.Atoi(value)
			if err != nil {
				return svc, fmt.Errorf("invalid value for 'persist': %q", value)
			}
			// Only the case where persist is 1 is defined in the specification, and other values should be ignored.
			if p == 1 {
				svc.Persist = 1
			}
		default:
			// key as ProtocolID，value = host:port
			unquoted, err := strconv.Unquote(value)
			if err != nil {
				return svc, fmt.Errorf("cannot unquote alt-authority value %q: %v", value, err)
			}
			host, port, ok := strings.Cut(unquoted, ":")
			if !ok {
				return svc, fmt.Errorf("invalid alt-authority format in value: %q", unquoted)
			}
			svc.ProtocolID = key
			svc.AltAuthority = AltAuthority{
				Host: host,
				Port: port,
			}
		}
	}
	return svc, nil
}
