package core

import (
	"net"
	"net/http"
	"strings"
)

func ReadClientIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}

	if len(IPAddress) > 0 {
		IPAddress, _, _ = net.SplitHostPort(IPAddress)
	}

	if IPAddress == "::1" {
		IPAddress = "127.0.0.1"
	} else {
		if strings.Contains(IPAddress, ":") {
			tmpSlice := strings.Split(IPAddress, ":")
			if len(tmpSlice) > 0 {
				IPAddress = tmpSlice[0]
			}
		}
	}

	return IPAddress
}
