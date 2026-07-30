package config

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/internal/infra/metrics"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

type Database struct {
	*pgxpool.Pool
	notifier Notifier
}
type QueryContextKey struct {
	sql   string
	start time.Time
}
type PgxLogger struct {
}
type RecoveryMiddleware struct{}

func NewDatabase(ctx context.Context, cfg *DatabaseConfig, app *AppConfig, notifier Notifier) (*Database, error) {
	dsn := cfg.GetDSN()
	if dsn == "" {
		return nil, errors.New("database DSN is empty")
	}

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		Logger.Error("invalid DB config", zap.Error(err))
		return nil, fmt.Errorf("invalid DB config: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = 30 * time.Minute
	poolConfig.MaxConnIdleTime = 5 * time.Minute
	poolConfig.ConnConfig.Tracer = &PgxLogger{}

	var pool *pgxpool.Pool
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			break
		}

		if i < maxRetries-1 {
			Logger.Warn("Retrying database connection...",
				zap.Int("attempt", i+1),
				zap.Error(err),
			)
			time.Sleep(2 * time.Second)
		}
	}

	if pool == nil {
		Logger.Error("failed to connect", zap.Error(err), zap.Int("retries", maxRetries))
		return nil, fmt.Errorf("failed to connect after %d retries: %w", maxRetries, err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		Logger.Error("PGX pool ping failed", zap.Error(err))
		return nil, fmt.Errorf("DB ping failed: %w", err)
	}

	dbPool := &Database{
		Pool:     pool,
		notifier: notifier,
	}

	dbMutex.Lock()
	PGXDB = dbPool
	dbMutex.Unlock()

	Logger.Info("PGX pool connected successfully",
		zap.String("dsn", maskPasswordInDSN(dsn)), // Helper function to hide password
		zap.Int("max_conns", int(poolConfig.MaxConns)),
		zap.Bool("debug", app.Debug),
	)

	go checkPoolHealth(ctx)
	return dbPool, nil
}
func GetDBPool() *Database {
	if PGXDB == nil {
		return nil
	}
	return PGXDB
}
func (db *Database) Close() error {
	if db.Pool != nil {
		db.Pool.Close()
		return nil
	}
	return nil
}
func (db *Database) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return db.Pool.Ping(ctx)
}
func (db *Database) GetStats() map[string]interface{} {
	stats := db.Pool.Stat()
	return map[string]interface{}{
		"max_open_connections":       stats.MaxConns(),
		"total_connections":          stats.TotalConns(),
		"idle_connections":           stats.IdleConns(),
		"max_idle_destroy_count":     stats.MaxIdleDestroyCount(),
		"acquire_count":              stats.AcquireCount(),
		"acquire_duration":           stats.AcquireDuration().String(),
		"acquired_connections":       stats.AcquiredConns(),
		"max_lifetime_destroy_count": stats.MaxLifetimeDestroyCount(),
	}
}
func (log *PgxLogger) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if log != nil {
		Logger.Debug("Query started",
			zap.String("sql", data.SQL),
			zap.Any("args", data.Args),
		)
	}
	return context.WithValue(ctx, QueryContextKey{}, &QueryContextKey{
		sql:   data.SQL,
		start: time.Now(),
	})
}
func (log *PgxLogger) TraceQueryEnd(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	val := ctx.Value(QueryContextKey{})
	if val == nil {
		return
	}

	qc, ok := val.(*QueryContextKey)
	if !ok || qc == nil {
		return
	}

	duration := time.Since(qc.start)
	if data.Err != nil {
		Logger.Error("Query Trace End Failed %w", zap.Error(data.Err), zap.String("query", qc.sql), zap.Duration("duration", duration), zap.String("command_tag", data.CommandTag.String()))
	}
	if time.Since(time.Now()) > slowQueryThreshold {
		Logger.Error("Slow query detected Trace %w", zap.String("query", qc.sql), zap.Duration("duration", duration), zap.String("command_tag", data.CommandTag.String()))
	}
}
func (r *RecoveryMiddleware) Handle(next func(ctx context.Context) error) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		defer func() {
			if err := recover(); err != nil {
				Logger.Error("Recovered from panic",
					zap.Any("panic", err),
					zap.Stack("stack"),
				)
			}
		}()
		return next(ctx)
	}
}
func (ip *Database) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	label := ip.ExtractQueryLabel(sql)
	timer := prometheus.NewTimer(metrics.GetAppMetrics().DBQueryDuration.WithLabelValues(label, "success"))

	tag, err := ip.Pool.Exec(ctx, sql, args...)
	if err != nil {
		timer.ObserveDuration() // Record success first
		timer = prometheus.NewTimer(metrics.GetAppMetrics().DBQueryDuration.WithLabelValues(label, "error"))
		timer.ObserveDuration()
		return tag, err
	}

	timer.ObserveDuration()
	return tag, nil
}
func (ip *Database) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	label := ip.ExtractQueryLabel(sql)
	timer := prometheus.NewTimer(metrics.GetAppMetrics().DBQueryDuration.WithLabelValues(label, "success"))

	rows, err := ip.Pool.Query(ctx, sql, args...)
	if err != nil {
		timer.ObserveDuration()
		timer = prometheus.NewTimer(metrics.GetAppMetrics().DBQueryDuration.WithLabelValues(label, "error"))
		timer.ObserveDuration()
		return rows, err
	}

	timer.ObserveDuration()
	return rows, nil
}
func (ip *Database) Validate() error {
	if ip == nil || ip.Pool == nil {
		return fmt.Errorf("database pool is not initialized")
	}
	return nil
}
func (db *Database) HandleQueryCompletion(qc *QueryContextKey, data pgx.TraceQueryEndData) {
	elapsed := time.Since(qc.start)
	if data.Err != nil {
		db.notifier.SendAlert(
			AlertRequest{
				Subject: "Query execution failed",
				Message: data.Err.Error(),
				Metadata: map[string]interface{}{
					"query":       qc.sql,
					"duration":    elapsed,
					"command_tag": data.CommandTag.String(),
				},
			},
		)
	}
}
func (db *Database) ExtractQueryLabel(query string) string {
	// Ambil keyword pertama buat label
	if len(query) > 50 {
		query = query[:50]
	}
	return query
}
func ObserveDBDuration(label, status string, start time.Time) {
	duration := time.Since(start).Seconds()
	Logger.Info(fmt.Sprintf("Query %s executed in %v", label, duration))
	metrics.GetAppMetrics().DBQueryDuration.WithLabelValues(label, status).Observe(duration)
}
func maskPasswordInDSN(dsn string) string {
	if strings.Contains(dsn, "password=") {
		parts := strings.Split(dsn, " ")
		for i, part := range parts {
			if strings.HasPrefix(part, "password=") {
				parts[i] = "password=*****"
			}
		}
		return strings.Join(parts, " ")
	}
	return dsn
}
func checkPoolHealth(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			dbMutex.Lock()
			if PGXDB != nil && PGXDB.Pool != nil {
				if err := PGXDB.Pool.Ping(ctx); err != nil {
					Logger.Error("Database health check failed",
						zap.Error(err),
					)
				}
			}
			dbMutex.Unlock()
		}
	}
}
