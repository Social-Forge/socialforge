package main

import (
	"context"
	"github/socialforge/config"
	"github/socialforge/internal/dependencies"
	"github/socialforge/internal/infra/rabbitmq"
	"github/socialforge/internal/services"
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

	// Inbound ingestion consumer: persists messages + publishes realtime.
	outboundForIngest := services.NewOutboundService(
		cont.ConversationRepo, cont.MessageRepo, cont.MessageOutboxRepo,
		cont.ChannelRepo, cont.ContactRepo, cont.CentrifugoClient, cont.RabbitMQ,
		services.BuildSenders(cont.Config), cont.Logger,
	)
	ingestionService := services.NewIngestionService(
		cont.ChannelRepo,
		cont.ContactRepo,
		cont.ConversationRepo,
		cont.MessageRepo,
		cont.WebhookEventRepo,
		cont.AutoResponseRepo,
		outboundForIngest,
		cont.CentrifugoClient,
		cont.RabbitMQ,
		cont.Logger,
	)
	if cont.RabbitMQ != nil {
		if err := cont.RabbitMQ.Consume(rabbitmq.QueueIngestInbound, "worker-ingest", 20, ingestionService.ProcessInbound); err != nil {
			cont.Logger.Fatal("Failed to start ingest consumer", zap.Error(err))
		}
	}

	// Outbound dispatch consumer (Fase 2E plugs in real provider send).
	outboundService := services.NewOutboundService(
		cont.ConversationRepo,
		cont.MessageRepo,
		cont.MessageOutboxRepo,
		cont.ChannelRepo,
		cont.ContactRepo,
		cont.CentrifugoClient,
		cont.RabbitMQ,
		services.BuildSenders(cont.Config),
		cont.Logger,
	)
	if cont.RabbitMQ != nil {
		if err := cont.RabbitMQ.Consume(rabbitmq.QueueDispatchOutbound, "worker-dispatch", 20, outboundService.ProcessDispatch); err != nil {
			cont.Logger.Fatal("Failed to start dispatch consumer", zap.Error(err))
		}
	}

	cont.Logger.Info("🛠️  Worker started — consuming jobs")

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
