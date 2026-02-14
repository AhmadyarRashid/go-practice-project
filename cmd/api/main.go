package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yourusername/go-enterprise-api/internal/config"
	"github.com/yourusername/go-enterprise-api/internal/database"
	"github.com/yourusername/go-enterprise-api/internal/models"
	"github.com/yourusername/go-enterprise-api/internal/routes"
	"github.com/yourusername/go-enterprise-api/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger.Init(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
		Debug:  cfg.App.Debug,
	})
	defer logger.Sync()

	logger.Info("Starting application",
		logger.String("name", cfg.App.Name),
		logger.String("env", cfg.App.Env),
		logger.String("port", cfg.App.Port),
	)

	// Connect to database
	db, err := database.New(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", logger.Err(err))
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("Failed to close database connection", logger.Err(err))
		}
	}()

	// Run migrations
	logger.Info("Running database migrations...")
	if err := db.Migrate(
		&models.User{},
		&models.Post{},
		&models.Tag{},
	); err != nil {
		logger.Fatal("Failed to run migrations", logger.Err(err))
	}
	logger.Info("Database migrations completed")

	// Setup routes
	router := routes.Setup(cfg, db)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Server starting",
			logger.String("address", srv.Addr),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", logger.Err(err))
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", logger.Err(err))
	}

	logger.Info("Server exited properly")
}
