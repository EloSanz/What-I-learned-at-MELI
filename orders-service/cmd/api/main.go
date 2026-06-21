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

	"orders-service/internal/api"
	"orders-service/internal/database"
	"orders-service/internal/order"
)

func main() {
	slog.Info("Starting Go Orders Service...")

	// 1. Conectar a PostgreSQL
	db, err := database.ConnectDB()
	if err != nil {
		slog.Error("Could not initialize database connection", "error", err)
		os.Exit(1)
	}

	// 2. Auto-migración del esquema de GORM
	slog.Info("Running database schema auto-migrations...")
	if err := db.AutoMigrate(&order.Order{}); err != nil {
		slog.Error("Database migration failed", "error", err)
		os.Exit(1)
	}

	// 3. Inicializar capas de la aplicación
	orderRepo := order.NewRepository(db)
	orderHandler := order.NewHandler(orderRepo)
	router := api.InitRouter(orderHandler)

	// 4. Configurar Servidor HTTP con Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Ejecutar servidor en una goroutine
	go func() {
		slog.Info("Orders Service is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	// Capturar señales de detención para apagado seguro
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down orders service server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Orders service server exited cleanly")
}
