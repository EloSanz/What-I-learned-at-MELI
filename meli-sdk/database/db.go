package database

import (
	"fmt"
	"log/slog"

	"github.com/user/meli-sdk/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect initializes the PostgreSQL database connection using the provided configuration.
func Connect(cfg *config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	slog.Info("Database connection successfully established via SDK", "dbname", cfg.Name)
	return db, nil
}
