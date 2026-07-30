package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	AppMetric   *AppMetrics
	RedisMetric *RedisMetrics
	metricOnce  sync.Once
	registry    *prometheus.Registry
)

func InitMetrics() {
	metricOnce.Do(func() {
		registry = prometheus.NewRegistry()

		AppMetric = initAppMetrics()
		RedisMetric = initRedisMetrics()

		registry.MustRegister(
			AppMetric.HTTPRequestTotal,
			AppMetric.HTTPRequestDuration,
			AppMetric.JWTErrorTotal,
			AppMetric.DBQueryDuration,

			RedisMetric.CommandCount,
			RedisMetric.ErrorCount,
			RedisMetric.Duration,
		)

		// Set default registry
		prometheus.DefaultRegisterer = registry
		prometheus.DefaultGatherer = registry
	})
}
func normalizePath(path string) string {
	if path == "" {
		return "/"
	}

	// Contoh normalisasi path parameters
	// /users/123 -> /users/:id
	// /posts/456/comments -> /posts/:id/comments

	// Implementasi sederhana - bisa dikembangkan
	switch {
	case len(path) > 50: // Prevent too long paths in metrics
		return "[truncated]"
	default:
		return path
	}
}

func GetRegistry(logger *zap.Logger) prometheus.Registerer {
	if registry == nil {
		InitMetrics()
	}
	return registry
}

func GetGatherer() prometheus.Gatherer {
	if registry == nil {
		InitMetrics()
	}
	return registry
}
func CleanupMetrics() {
	// Jika perlu unregister metrics (jarang digunakan)
	registry = nil
	AppMetric = nil
	RedisMetric = nil
}
