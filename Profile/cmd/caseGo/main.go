package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YoungFlores/Case_Go/Profile/internal/server"
)

// @title Profile Service API
// @version 1.0
// @description API сервис для управления профилями пользователей.
// @host localhost:8080
// @BasePath /profile/api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
		slog.Info("Starting server", "addr", srv.HTTP.Addr)
		if err := srv.HTTP.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
		}
	}()

	go func() {
		lis, err := net.Listen("tcp", ":50053")
		if err != nil {
			slog.Error("Failed to listen gRPC", "error", err)
			os.Exit(1)
		}
		slog.Info("Starting gRPC server", "addr", ":50053")
		if err := srv.GRPC.Serve(lis); err != nil {
			slog.Error("gRPC server error", "error", err)
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
