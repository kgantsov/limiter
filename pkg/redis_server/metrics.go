package redis_server

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	reqDurations prometheus.Summary
	reqCount     *prometheus.CounterVec
}

func NewMetrics(subsystem string) *Metrics {
	m := &Metrics{}

	m.reqDurations = prometheus.NewSummary(prometheus.SummaryOpts{
		Subsystem:  subsystem,
		Name:       "req_durations_seconds",
		Help:       "Request latency distributions.",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
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
