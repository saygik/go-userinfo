package usecase

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var getADUsersMetrics = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Namespace:  "ad",
	Subsystem:  "users",
	Name:       "get_all",
	Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
}, []string{"domain"})

func observeGetADUsers(d time.Duration, domain string) {
	getADUsersMetrics.WithLabelValues(domain).Observe(d.Seconds())
}
