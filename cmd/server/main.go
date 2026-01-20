package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pupervisor/internal/api"
	"pupervisor/internal/config"
	"pupervisor/internal/service"
	"pupervisor/web"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize services
	processSvc := service.NewProcessService()

	// Get embedded filesystems
	templatesFS := web.GetTemplatesFS()
	staticFS := web.GetStaticFS()

	// Create router
	router, err := api.NewRouter(processSvc, templatesFS, staticFS)
	if err != nil {
		log.Fatalf("Failed to create router: %v", err)
	}

	// Create server
	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting Pupervisor Web UI server on %s", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
