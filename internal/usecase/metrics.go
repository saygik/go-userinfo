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

var gaugeUsersPerDomain = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "ad_users_per_domain",
	Help: "Number of users per domain by types",
}, []string{"domain", "type", "subtype"})

var countTicketsPerRegion = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "glpi_hrp_tickets_per_region",
	Help: "Number of tickets per region",
}, []string{"region"})

func observeGetADUsers(d time.Duration, domain string) {
	getADUsersMetrics.WithLabelValues(domain).Observe(d.Seconds())
}
func observeUsersPerDomain(domain string, types string, subtype string, users int) {
	gaugeUsersPerDomain.WithLabelValues(domain, types, subtype).Set(float64(users))
}

func observeCountTicketsPerRegion(region string) {
	countTicketsPerRegion.WithLabelValues(region).Inc()
}
