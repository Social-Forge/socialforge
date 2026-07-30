package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/prometheus/client_golang/prometheus"
)

type AppMetrics struct {
	HTTPRequestTotal    *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
	JWTErrorTotal       *prometheus.CounterVec
	DBQueryDuration     *prometheus.HistogramVec
}

func initAppMetrics() *AppMetrics {
	AppMetric = &AppMetrics{
		HTTPRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total HTTP requests received",
			},
			[]string{"method", "path", "status"},
		),
		HTTPRequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		JWTErrorTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "jwt_errors_total",
				Help: "Total JWT validation errors",
			},
			[]string{"error_type"},
		),
		DBQueryDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "db_query_duration_seconds",
				Help: "DB query duration",
			},
			[]string{"query_label", "status"},
		),
	}

	return AppMetric
}
func GetAppMetrics() *AppMetrics {
	if AppMetric == nil {
		InitMetrics()
	}
	return AppMetric
}
func HTTPMetrics(m *AppMetrics) fiber.Handler {
	return func(c fiber.Ctx) error {

		start := time.Now()
		err := c.Next()

		duration := time.Since(start)
		status := c.Response().StatusCode()
		if status == 0 {
			status = http.StatusOK
		}

		path := normalizePath(c.Path())
		method := c.Method()
		statusStr := strconv.Itoa(status)

		m.HTTPRequestTotal.WithLabelValues(method, path, statusStr).Inc()
		m.HTTPRequestDuration.WithLabelValues(method, path, statusStr).Observe(duration.Seconds())

		return err
	}
}
