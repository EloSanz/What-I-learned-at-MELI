package config

import (
	"encoding/json"
	"os"
)

// DBConfig holds the database connection settings.
type DBConfig struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Port     string `json:"port"`
	SSLMode  string `json:"sslmode"`
}

// LoadDBConfig loads the config from a JSON file if it exists,
// otherwise it falls back to reading environment variables.
func LoadDBConfig(configPath string) (*DBConfig, error) {
	// Try to load from JSON file first
	file, err := os.Open(configPath)
	if err == nil {
		defer file.Close()
		var cfg DBConfig
		if err := json.NewDecoder(file).Decode(&cfg); err == nil {
			return &cfg, nil
		}
	}

	// Fallback to environment variables
	return &DBConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		Name:     getEnv("DB_NAME", "meli_db"),
		Port:     getEnv("DB_PORT", "5432"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
