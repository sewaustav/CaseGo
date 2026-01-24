package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/sewaustav/CaseGoProfile/internal/server"
)


func main() {
	srv, err := server.New()
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	go func() {
		log.Printf("Starting server on %s", srv.HTTP.Addr)
		if err := srv.HTTP.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.HTTP.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	if err := srv.DB.GetDB().Close(); err != nil {
		log.Printf("Failed to close DB: %v", err)
	} 

	log.Println("Server stopped gracefully")
}