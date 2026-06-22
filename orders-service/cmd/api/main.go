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

	"orders-service/internal/bootstrap"
	"orders-service/internal/database"
)

func main() {
	slog.Info("Starting Go Orders Service...")

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
		port = "8082"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		slog.Info("Orders Service is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	// 4. Capture stop signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down orders service server gracefully...")

	// Allow up to 10 seconds to finish active requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Orders service server exited cleanly")
}
