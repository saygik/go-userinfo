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

// adCacheFillDuration — время заполнения кеша AD по домену.
// result: success | partial | error | panic.
var adCacheFillDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "ad_cache_fill_duration_seconds",
	Help:    "AD cache fill duration by domain and result",
	Buckets: []float64{1, 2, 5, 10, 15, 20, 30, 45, 60, 90, 120, 180, 300},
}, []string{"domain", "result"})

// adCacheFillErrors — ошибки заполнения кеша AD по домену и этапу.
var adCacheFillErrors = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ad_cache_fill_errors_total",
	Help: "AD cache fill errors by domain and stage",
}, []string{"domain", "stage"})

// adCacheFillTimeouts — количество проходов, не уложившихся в общий таймаут.
var adCacheFillTimeouts = promauto.NewCounter(prometheus.CounterOpts{
	Name: "ad_cache_fill_timeouts_total",
	Help: "Number of AD cache fills that hit the global timeout",
})

// adCacheLastSuccess — время последнего успешного заполнения домена.
var adCacheLastSuccess = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "ad_cache_last_success_timestamp_seconds",
	Help: "Unix timestamp of the last successful AD cache fill by domain",
}, []string{"domain"})

// adCacheLastComplete — время последнего полностью завершённого прохода.
var adCacheLastComplete = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "ad_cache_last_complete_timestamp_seconds",
	Help: "Unix timestamp of the last fully completed AD cache fill",
})

// adUsersCacheRequests — обращения к кешу обработанного списка пользователей.
// result: hit | miss | invalid.
var adUsersCacheRequests = promauto.NewCounterVec(prometheus.CounterOpts{
	Name: "ad_users_cache_requests_total",
	Help: "Processed AD users cache requests by result",
}, []string{"result"})

// adUsersBuildDuration — время пересборки обработанного списка пользователей.
var adUsersBuildDuration = promauto.NewHistogram(prometheus.HistogramOpts{
	Name:    "ad_users_build_duration_seconds",
	Help:    "Time to build processed AD users list on cache miss",
	Buckets: []float64{0.05, 0.1, 0.25, 0.5, 1, 2, 5, 10, 20, 30},
})

// adUsersProcessed — количество пользователей домена после обработки.
var adUsersProcessed = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Name: "ad_users_processed_users",
	Help: "Number of processed users per domain",
}, []string{"domain"})

func observeGetADUsers(d time.Duration, domain string) {
	getADUsersMetrics.WithLabelValues(domain).Observe(d.Seconds())
}
func observeUsersPerDomain(domain string, types string, subtype string, users int) {
	gaugeUsersPerDomain.WithLabelValues(domain, types, subtype).Set(float64(users))
}

func observeCountTicketsPerRegion(region string) {
	countTicketsPerRegion.WithLabelValues(region).Inc()
}

func observeADCacheFillDuration(d time.Duration, domain, result string) {
	adCacheFillDuration.WithLabelValues(domain, result).Observe(d.Seconds())
}

func observeADCacheFillError(domain, stage string) {
	adCacheFillErrors.WithLabelValues(domain, stage).Inc()
}

func observeADCacheFillTimeout() {
	adCacheFillTimeouts.Inc()
}

func observeADCacheLastSuccess(domain string) {
	adCacheLastSuccess.WithLabelValues(domain).Set(float64(time.Now().Unix()))
}

func observeADCacheLastComplete() {
	adCacheLastComplete.Set(float64(time.Now().Unix()))
}

func observeADUsersCacheRequest(result string) {
	adUsersCacheRequests.WithLabelValues(result).Inc()
}

func observeADUsersBuildDuration(d time.Duration) {
	adUsersBuildDuration.Observe(d.Seconds())
}

func observeADUsersProcessed(domain string, users int) {
	adUsersProcessed.WithLabelValues(domain).Set(float64(users))
}
