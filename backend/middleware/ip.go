// middleware/ip.go
package middleware

import (
    "net"
    "net/http"
    "strings"
)

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
			return strings.Split(forwarded, ",")[0]
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
			return NormalizeIP(r.RemoteAddr)
	}
	return NormalizeIP(ip)
}


func NormalizeIP(ip string) string {
	if ip == "::1" {
			return "127.0.0.1" // Normalize IPv6 loopback to IPv4 loopback
	}
	return ip
}
