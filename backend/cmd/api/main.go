package main

import (
	"context"
	"github/socialforge/config"
	"github/socialforge/internal/app"
	"github/socialforge/internal/dependencies"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	cont, err := dependencies.NewContainer(ctxTimeout)
	if err != nil {
		config.Logger.Fatal("Failed to initialize dependencies", zap.Error(err))
	}
	defer cont.Close()

	go func() {
		app.Start(cont)
	}()

	<-ctx.Done()
	stop()
	log.Println("⚠️ Shutdown signal received")

	cleanup(ctx, cont)
}
func cleanup(ctx context.Context, cont *dependencies.Container) {

	if err := app.Shutdown(ctx); err != nil {
		config.Logger.Error("Failed to shutdown Fiber", zap.Error(err))
	}

	config.Logger.Info("Shutting down services...")

	if err := cont.Close(); err != nil {
		log.Fatalf("failed to close dependencies system %v", err)
	}

	cont.Close()
	config.Logger.Info("All services shut down gracefully ✅")
}
