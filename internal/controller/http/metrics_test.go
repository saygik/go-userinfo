package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetricsMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(metricsMiddleware())
	router.GET("/test/:id", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	before := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/test/:id", "200"))

	req := httptest.NewRequest(http.MethodGet, "/test/42", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	after := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "/test/:id", "200"))
	if after-before != 1 {
		t.Fatalf("http_requests_total delta = %v, want 1", after-before)
	}
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestMetricsMiddlewareUnmatchedRoute(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(metricsMiddleware())

	before := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "unmatched", "404"))

	req := httptest.NewRequest(http.MethodGet, "/no-such-route", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	after := testutil.ToFloat64(httpRequestsTotal.WithLabelValues("GET", "unmatched", "404"))
	if after-before != 1 {
		t.Fatalf("http_requests_total{route=unmatched} delta = %v, want 1", after-before)
	}
}
