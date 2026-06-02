package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tendo-mulira/tnotes-teams/internal/config"
	"github.com/tendo-mulira/tnotes-teams/internal/database"
	"github.com/tendo-mulira/tnotes-teams/internal/handler"
	"github.com/tendo-mulira/tnotes-teams/internal/middleware"
	"github.com/tendo-mulira/tnotes-teams/internal/repository"
	"github.com/tendo-mulira/tnotes-teams/internal/service"
	"github.com/tendo-mulira/tnotes-teams/internal/websocket"
	"github.com/tendo-mulira/tnotes-teams/internal/worker"
	"github.com/tendo-mulira/tnotes-teams/pkg/api"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, entClient, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	defer entClient.Close()

	// Run auto-migrations in development
	if cfg.Environment == "development" {
		ctx := context.Background()
		if err := database.Migrate(ctx, entClient); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Database migrations completed")
		if err := database.SeedDevelopment(ctx, entClient); err != nil {
			log.Printf("Seed warning: %v", err)
		}
	}

	// Initialize repositories
	repos := repository.NewRepositories(entClient, db)

	// Initialize WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Initialize services
	services := service.NewServices(repos, cfg, hub)

	// Initialize handlers
	handlers := handler.NewHandlers(services)

	// Initialize middleware
	mw := middleware.NewMiddleware(cfg, repos)

	// Set up router
	router := api.SetupRoutes(handlers, mw, hub, services, cfg)

	// Start team update engine (background worker)
	teamEngine := worker.NewTeamEngine(services, hub, cfg)
	go teamEngine.Start(context.Background())

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("TNotes Teams Web server starting on port %s", cfg.Port)
		log.Printf("Environment: %s", cfg.Environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
