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

	"items-service/internal/api"
	"items-service/internal/database"
	"items-service/internal/item"

	"gorm.io/gorm"
)

func main() {
	slog.Info("Starting Go Items Service...")

	// 1. Conectar a PostgreSQL
	db, err := database.ConnectDB()
	if err != nil {
		slog.Error("Could not initialize database connection", "error", err)
		os.Exit(1)
	}

	// 2. Ejecutar Auto-migración del esquema de GORM
	slog.Info("Running database schema auto-migrations...")
	if err := db.AutoMigrate(&item.Item{}); err != nil {
		slog.Error("Database migration failed", "error", err)
		os.Exit(1)
	}

	// 3. Sembrar datos por defecto si la tabla está vacía
	seedDefaultItems(db)

	// 4. Inicializar capas de Clean Architecture
	itemRepo := item.NewRepository(db)
	itemHandler := item.NewHandler(itemRepo)
	router := api.InitRouter(itemHandler)

	// 5. Configurar Servidor HTTP con Graceful Shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Corremos el servidor en una goroutine para no bloquear el hilo principal
	go func() {
		slog.Info("Server is running", "port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("ListenAndServe failed", "error", err)
			os.Exit(1)
		}
	}()

	// Esperamos señales del sistema operativo para apagar el servidor de forma segura
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down items service server gracefully...")

	// Damos un tiempo de tolerancia de 10 segundos para finalizar peticiones en curso
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	slog.Info("Items service server exited cleanly")
}

// seedDefaultItems inserta el monitor de prueba si la base de datos no tiene registros
func seedDefaultItems(db *gorm.DB) {
	var count int64
	db.Model(&item.Item{}).Count(&count)
	if count == 0 {
		slog.Info("Seeding default item for Mercado Libre simulation...")
		defaultMonitor := item.Item{
			ID:    "MLA43960787",
			Title: "Monitor gamer curvo Xiaomi Gaming G34WQi LCD negro",
			Price: 619999.00,
			Stock: 55, // Stock inicial para simular compras
		}
		if err := db.Create(&defaultMonitor).Error; err != nil {
			slog.Warn("Failed to seed default item", "error", err)
		} else {
			slog.Info("Default item seeded successfully (ID: MLA43960787, Stock: 55)")
		}
	}
}
