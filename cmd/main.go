package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/raiashpanda007/kairo/internal/server"
	"github.com/raiashpanda007/kairo/shared/config"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.MustLoad()

	srv := server.New(cfg.Server.Addr, logger)

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
