package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sewaustav/CaseGoProfile/internal/server"
)

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	srv, err := server.New()
	if err != nil {
		slog.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := srv.Run(); err != nil {
			slog.Error("Server Run error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.HTTP.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Stopping gRPC server...")
	srv.GRPC.GracefulStop()

	if err := srv.DB.GetDB().Close(); err != nil {
		slog.Error("Failed to close DB", "error", err)
	}

	slog.Info("Server stopped gracefully")
}
