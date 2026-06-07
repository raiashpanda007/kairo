package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	core_queue "github.com/raiashpanda007/kairo/core/queue"
	"github.com/raiashpanda007/kairo/internal/server"
	"github.com/raiashpanda007/kairo/shared/config"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.MustLoad()

	ctx := context.Background()

	var QueueManager = core_queue.NewQueueManager(logger)
	QueueManager.Start(ctx)

	srv := server.New(cfg.Server.Addr, logger, &QueueManager)

	go func() {
		if err := srv.Run(); err != nil {
			logger.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")
	srv.Stop()
}
