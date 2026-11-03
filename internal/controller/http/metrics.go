package http

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests by method, route and status",
	}, []string{"method", "route", "status"})

	httpRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration by method, route and status",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "route", "status"})

	httpRequestsInFlight = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "http_requests_in_flight",
		Help: "Current number of HTTP requests being processed",
	})

	httpResponseSize = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_size_bytes",
		Help: "HTTP response size in bytes by route",
		Buckets: []float64{
			256, 1024, 4096, 16384, 65536, 262144,
			1048576, 4194304, 16777216, 67108864,
		},
	}, []string{"route"})
)

// metricsMiddleware собирает общие метрики по всем HTTP-запросам.
// В метке route используется шаблон Gin (c.FullPath), а не реальный URL,
// чтобы не плодить кардинальность на путях с параметрами.
func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unmatched"
		}
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method
		elapsed := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(method, route, status).Inc()
		httpRequestDuration.WithLabelValues(method, route, status).Observe(elapsed)
		httpResponseSize.WithLabelValues(route).Observe(float64(c.Writer.Size()))
	}
}
