package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"items-service/internal/bootstrap"
	"items-service/internal/database"
)

func main() {
	slog.Info("Starting Go Items Service...")

	// 1. Connect to PostgreSQL
	db, err := database.ConnectDB()
	if err != nil {
		slog.Error("Could not initialize database connection", "error", err)
		os.Exit(1)
	}

	// 2. Initialize application dependencies and routing (DI Container)
	router := bootstrap.InitApp(db)

	// 3. Setup HTTP server with Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		slog.Info("Items Server is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down items service server gracefully...")

	// Allow up to 10 seconds to finish active requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Items service server exited cleanly")
}
