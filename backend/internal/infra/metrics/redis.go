package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type RedisMetrics struct {
	CommandCount *prometheus.CounterVec
	ErrorCount   *prometheus.CounterVec
	Duration     *prometheus.HistogramVec
}

func initRedisMetrics() *RedisMetrics {
	RedisMetric = &RedisMetrics{
		CommandCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_commands_total",
				Help: "Total number of Redis commands executed",
			},
			[]string{"command", "environment", "instance", "status"},
		),
		ErrorCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "redis_errors_total",
				Help: "Total number of Redis errors",
			},
			[]string{"command", "instance", "error"},
		),
		Duration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "redis_command_duration_seconds",
				Help:    "Redis command execution duration",
				Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
			},
			[]string{"command", "environment", "instance", "status"},
		),
	}
	return RedisMetric
}
func GetRedisMetrics() *RedisMetrics {
	if RedisMetric == nil {
		InitMetrics()
	}
	return RedisMetric
}
func (m *RedisMetrics) Observe(ctx context.Context, cmd string, fn func() error) error {
	if m == nil {
		return fn()
	}

	start := time.Now()
	err := fn()
	duration := time.Since(start)

	status := "success"
	if err != nil {
		status = "error"
		errorType := "unknown"

		m.ErrorCount.WithLabelValues(cmd, errorType).Inc()
	}

	m.CommandCount.WithLabelValues(cmd, status).Inc()
	m.Duration.WithLabelValues(cmd, status).Observe(duration.Seconds())

	return err
}
