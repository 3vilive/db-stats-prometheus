package metrics

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type DbStatsCollector struct {
	Labels []string

	MaxOpenConnections *prometheus.GaugeVec // Maximum number of open connections to the database.

	// Pool status
	OpenConnections *prometheus.GaugeVec // The number of established connections both in use and idle.
	InUse           *prometheus.GaugeVec // The number of connections currently in use.
	Idle            *prometheus.GaugeVec // The number of idle connections.

	// Counters
	WaitCount         *prometheus.GaugeVec // The total number of connections waited for.
	WaitDuration      *prometheus.GaugeVec // The total time blocked waiting for a new connection.
	MaxIdleClosed     *prometheus.GaugeVec // The total number of connections closed due to SetMaxIdleConns.
	MaxLifetimeClosed *prometheus.GaugeVec // The total number of connections closed due to SetConnMaxLifetime.
}

func NewDbStatsCollector(labels []string) *DbStatsCollector {
	collector := &DbStatsCollector{
		Labels: labels,
		MaxOpenConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_max_open_connections",
			Help: "Maximum number of open connections to the database.",
		}, labels),
		OpenConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_open_connections",
			Help: "The number of established connections both in use and idle.",
		}, labels),
		InUse: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_in_use",
			Help: "The number of connections currently in use.",
		}, labels),
		Idle: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_idle",
			Help: "The number of idle connections.",
		}, labels),
		WaitCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_wait_count",
			Help: "The total number of connections waited for.",
		}, labels),
		WaitDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_wait_duration",
			Help: "The total time blocked waiting for a new connection.",
		}, labels),
		MaxIdleClosed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_max_idle_closed",
			Help: "The total number of connections closed due to SetMaxIdleConns.",
		}, labels),
		MaxLifetimeClosed: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "dbstats_max_lifetime_closed",
			Help: "The total number of connections closed due to SetConnMaxLifetime.",
		}, labels),
	}

	for _, item := range collector.Collectors() {
		prometheus.Register(item)
	}

	return collector
}

func (collector *DbStatsCollector) Set(dbStats sql.DBStats, labelValueMap map[string]string) {
	labelValues := make([]string, 0, len(collector.Labels))
	for _, label := range collector.Labels {
		labelValues = append(labelValues, labelValueMap[label])
	}

	collector.MaxOpenConnections.WithLabelValues(labelValues...).Set(float64(dbStats.MaxOpenConnections))
	collector.OpenConnections.WithLabelValues(labelValues...).Set(float64(dbStats.OpenConnections))
	collector.InUse.WithLabelValues(labelValues...).Set(float64(dbStats.InUse))
	collector.Idle.WithLabelValues(labelValues...).Set(float64(dbStats.Idle))
	collector.WaitCount.WithLabelValues(labelValues...).Set(float64(dbStats.WaitCount))
	collector.WaitDuration.WithLabelValues(labelValues...).Set(float64(dbStats.WaitDuration))
	collector.MaxIdleClosed.WithLabelValues(labelValues...).Set(float64(dbStats.MaxIdleClosed))
	collector.MaxLifetimeClosed.WithLabelValues(labelValues...).Set(float64(dbStats.MaxLifetimeClosed))
}

// get collector in stats
func (collector *DbStatsCollector) Collectors() []prometheus.Collector {
	return []prometheus.Collector{
		collector.MaxOpenConnections,
		collector.OpenConnections,
		collector.InUse,
		collector.Idle,
		collector.WaitCount,
		collector.WaitDuration,
		collector.MaxIdleClosed,
		collector.MaxLifetimeClosed,
	}
}
