package redis_server

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	requestDuration *prometheus.HistogramVec
	requestsTotal   *prometheus.CounterVec
	requestInFlight prometheus.Gauge
	connections     prometheus.Gauge
}

func NewMetrics(subsystem string) *Metrics {
	m := &Metrics{}

	m.requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: subsystem,
		Name:      "request_duration_seconds_bucket",
		Help:      "Request latency distributions, partitioned by status code.",
		Buckets: []float64{
			0.000000001, // 1ns
			0.000000002,
			0.000000005,
			0.00000001, // 10ns
			0.00000002,
			0.00000005,
			0.0000001, // 100ns
			0.0000002,
			0.0000005,
			0.000001, // 1µs
			0.000002,
			0.000005,
			0.00001, // 10µs
			0.00002,
			0.00005,
			0.0001, // 100µs
			0.0002,
			0.0005,
			0.001, // 1ms
			0.002,
			0.005,
			0.01, // 10ms
			0.02,
			0.05,
			0.1, // 100 ms
			0.2,
			0.5,
			1.0, // 1s
			2.0,
			5.0,
			10.0, // 10s
			15.0,
			20.0,
			30.0,
		},
	}, []string{"status_code"})
	prometheus.MustRegister(m.requestDuration)

	m.requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "How many redis requests processed, partitioned by status code.",
		},
		[]string{"status_code"},
	)
	prometheus.MustRegister(m.requestsTotal)

	m.requestInFlight = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "requests_in_progress_total",
		Help:      "How many redis requests are currently being processed.",
	})
	prometheus.MustRegister(m.requestInFlight)

	m.connections = prometheus.NewGauge(prometheus.GaugeOpts{
		Subsystem: subsystem,
		Name:      "connections_total",
		Help:      "How many connections are currently open.",
	})
	prometheus.MustRegister(m.connections)

	return m
}
