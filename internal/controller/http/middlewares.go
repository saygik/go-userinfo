package http

import (
	"net"
	nethttp "net/http"
	"net/url"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (s *Server) registerMiddlewares() {
	s.Rtr.Use(cORSMiddleware())
	s.Rtr.Use(requestIDMiddleware())
	s.Rtr.Use(gzip.Gzip(gzip.DefaultCompression))
}

// CORSMiddleware ...
// Credentials разрешены только для своего домена: отражать любой Origin
// с Access-Control-Allow-Credentials нельзя - иначе cookie уходят на чужой сайт.
func cORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		w := c.Writer
		r := c.Request

		origin := r.Header.Get("Origin")
		if origin != "" && isTrustedOrigin(origin, r.Host) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			w.Header().Set("Access-Control-Allow-Headers", "X-Requested-With, Content-Type, Origin, Authorization, Accept, Client-Security-Token, Accept-Encoding, x-access-token, Access-Token")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.Header().Add("Vary", "Origin")
		}

		if r.Method == nethttp.MethodOptions {
			c.AbortWithStatus(nethttp.StatusOK)
			return
		}
		c.Next()
	}
}

// isTrustedOrigin разрешает запросы с того же хоста, из того же домена
// (второй уровень) и с localhost.
func isTrustedOrigin(origin, requestHost string) bool {
	parsed, err := url.Parse(origin)
	if err != nil || parsed.Host == "" {
		return false
	}

	originHost := strings.ToLower(parsed.Hostname())
	host := strings.ToLower(stripPort(requestHost))
	if originHost == "" || host == "" {
		return false
	}
	if originHost == host {
		return true
	}
	// Клиент для разработки может работать на localhost и обращаться к API
	// на другом хосте (например, к https://userinfoapi.brnv.rw). Это безопасно:
	// cookie SameSite=Lax на cross-site XHR всё равно не отправляются.
	if isLocalHost(originHost) {
		return true
	}

	originDomain := parentDomain(originHost)
	return originDomain != "" && originDomain == parentDomain(host)
}

func stripPort(host string) string {
	if h, _, err := net.SplitHostPort(host); err == nil {
		return h
	}
	return host
}

func isLocalHost(host string) bool {
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

// parentDomain - два последних уровня домена (например, brnv.rw).
func parentDomain(host string) string {
	if net.ParseIP(host) != nil {
		return ""
	}
	labels := strings.Split(host, ".")
	if len(labels) < 2 {
		return ""
	}
	return strings.Join(labels[len(labels)-2:], ".")
}

// RequestIDMiddleware ...
// Generate a unique ID and attach it to each request for future reference or use
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		suuid, _ := uuid.NewUUID() //uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", suuid.String())
		c.Next()
	}
}
