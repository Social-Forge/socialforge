package main

import (
	"context"
	"github/socialforge/config"
	"github/socialforge/internal/dependencies"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// The worker process shares the same dependency container as the API but runs
// background jobs instead of serving HTTP: channel ingestion consumers,
// message dispatch (outbox), auto-assign, subscription-expiry and AI jobs.
//
// Those consumers are introduced from Fase 2 onward. For now this boots the
// container (verifying all infra connectivity) and idles until a shutdown
// signal, so the worker service is deployable and observable from Fase 0.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cont, err := dependencies.NewContainer(ctxTimeout)
	if err != nil {
		config.Logger.Fatal("Failed to initialize worker dependencies", zap.Error(err))
	}
	defer cont.Close()

	cont.Logger.Info("🛠️  Worker started — waiting for jobs")

	// TODO(Fase 2+): start RabbitMQ consumers here, e.g.
	//   registerIngestConsumers(ctx, cont)
	//   registerDispatchConsumers(ctx, cont)
	//   startCronJobs(ctx, cont) // subscription expiry, etc.

	<-ctx.Done()
	stop()
	log.Println("⚠️ Worker shutdown signal received")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := cont.Close(); err != nil {
		config.Logger.Error("Failed to close worker dependencies", zap.Error(err))
	}
	_ = shutdownCtx
	config.Logger.Info("Worker shut down gracefully ✅")
}
