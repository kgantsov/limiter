package redis_server

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	reqDurations prometheus.Histogram
	reqCount     *prometheus.CounterVec
}

func NewMetrics(subsystem string) *Metrics {
	m := &Metrics{}

	m.reqDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Subsystem: subsystem,
		Name:      "req_durations_seconds",
		Help:      "Request latency distributions.",
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
	})
	prometheus.MustRegister(m.reqDurations)

	m.reqCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: subsystem,
			Name:      "requests_total",
			Help:      "How many redis requests processed, partitioned by status code.",
		},
		[]string{"code"},
	)
	prometheus.MustRegister(m.reqCount)

	return m
}
