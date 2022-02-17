package controllers

import (
	"net"
	"net/http"
	"strings"
)

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	if strings.Index(IPAddress, "::1") > -1 {
		return "127.0.0.1"
	}
	i := strings.Index(IPAddress, ":")
	if i > -1 {
		return IPAddress[:i]
	} else {
		return IPAddress
	}

}

func ReadUserName(ip string) (names []string, err error) {
	return net.LookupAddr(ip)
}
